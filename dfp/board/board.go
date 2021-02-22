package dfpboard

import (
	"context"
	"sync"
	"time"

	"github.com/disaster37/gobot-fat/dfp"
	"github.com/disaster37/gobot-fat/helper"
	"github.com/disaster37/gobot-fat/mail"
	"github.com/disaster37/gobot-fat/models"
	"github.com/disaster37/gobot-fat/usecase"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
)

const (
	// EventNewConfig  receive new DFPConfig
	EventNewConfig = "dfp-new-config"

	// EventNewState receive New DFPState
	EventNewState = "dfp-new-state"

	// EventBoardStop board stopped
	EventBoardStop = "dfp-board-stop"

	// EventWash wash
	EventWash = "dfp-new-wash"

	// EventStopDFP DFP stopped
	EventStopDFP = "dfp-stop-dfp"

	// EventStartDFP DFP started
	EventStartDFP = "dfp-start-dfp"

	// EventSetSecurity set security
	EventSetSecurity = "dfp-set-security"

	// EventUnsetSecurity unset security
	EventUnsetSecurity = "dfp-unset-security"

	// EventSetEmergencyStop set emergency stop
	EventSetEmergencyStop = "dfp-set-emergency-stop"

	// EventUnsetEmergencyStop unset emergency stop
	EventUnsetEmergencyStop = "dfp-unset-emergency-stop"

	// EventNewInput Permit to test work function
	EventNewInput = "new-input"
)

// DFPAdaptor is DFP board interface
type DFPAdaptor interface {
	gobot.Adaptor
	gpio.DigitalReader
	gpio.DigitalWriter
	SetInputPullup(listPins []*gpio.ButtonDriver) (err error)
}

// DFPBoard is the DFP board
type DFPBoard struct {
	state                   *models.DFPState
	config                  *models.DFPConfig
	board                   DFPAdaptor
	gobot                   *gobot.Robot
	eventUsecase            usecase.UsecaseCRUD
	stateUsecase            usecase.UsecaseCRUD
	configHandler           *viper.Viper
	mailClient              mail.Mail
	isOnline                bool
	isInitialized           bool
	relayDrum               *gpio.RelayDriver
	relayPump               *gpio.RelayDriver
	ledGreen                *gpio.LedDriver
	ledRed                  *gpio.LedDriver
	buttonStart             *gpio.ButtonDriver
	buttonStop              *gpio.ButtonDriver
	buttonForceDrum         *gpio.ButtonDriver
	buttonForcePump         *gpio.ButtonDriver
	buttonWash              *gpio.ButtonDriver
	buttonEmergencyStop     *gpio.ButtonDriver
	captorWaterUpper        *gpio.ButtonDriver
	captorWaterUnder        *gpio.ButtonDriver
	captorSecurityUpper     *gpio.ButtonDriver
	captorSecurityUnder     *gpio.ButtonDriver
	globalEventer           gobot.Eventer
	isRunning               bool
	name                    string
	timeBetweenWash         *time.Ticker
	waitTimeForceWash       *time.Ticker
	waitTimeForceWashFrozen *time.Ticker
	waitTimeUnsetSecurity   *time.Ticker
	schedulingRoutines      []*time.Ticker
	gobot.Eventer
	sync.Mutex
}

// NewDFP create board to manage DFP
func NewDFP(configHandler *viper.Viper, config *models.DFPConfig, state *models.DFPState, eventUsecase usecase.UsecaseCRUD, dfpStateUsecase usecase.UsecaseCRUD, eventer gobot.Eventer, mailClient mail.Mail) (dfpBoard dfp.Board) {

	//Create client
	c := NewRaspiAdaptor()

	return newDFP(c, configHandler, config, state, eventUsecase, dfpStateUsecase, eventer, mailClient)

}

// newDFP create board to manage DFP
func newDFP(board DFPAdaptor, configHandler *viper.Viper, config *models.DFPConfig, state *models.DFPState, eventUsecase usecase.UsecaseCRUD, dfpStateUsecase usecase.UsecaseCRUD, eventer gobot.Eventer, mailClient mail.Mail) dfp.Board {

	buttonPollingDuration := configHandler.GetDuration("button_polling") * time.Millisecond

	// Init board
	dfpBoard := &DFPBoard{
		board:               board,
		eventUsecase:        eventUsecase,
		stateUsecase:        dfpStateUsecase,
		configHandler:       configHandler,
		name:                configHandler.GetString("name"),
		config:              config,
		state:               state,
		isOnline:            false,
		isInitialized:       false,
		globalEventer:       eventer,
		relayDrum:           gpio.NewRelayDriver(board, configHandler.GetString("pin.relay.drum")),
		relayPump:           gpio.NewRelayDriver(board, configHandler.GetString("pin.relay.pomp")),
		ledGreen:            gpio.NewLedDriver(board, configHandler.GetString("pin.led.green")),
		ledRed:              gpio.NewLedDriver(board, configHandler.GetString("pin.led.red")),
		buttonEmergencyStop: gpio.NewButtonDriver(board, configHandler.GetString("pin.button.emergency_stop"), buttonPollingDuration),
		buttonStart:         gpio.NewButtonDriver(board, configHandler.GetString("pin.button.start"), buttonPollingDuration),
		buttonStop:          gpio.NewButtonDriver(board, configHandler.GetString("pin.button.stop"), buttonPollingDuration),
		buttonWash:          gpio.NewButtonDriver(board, configHandler.GetString("pin.button.wash"), buttonPollingDuration),
		buttonForceDrum:     gpio.NewButtonDriver(board, configHandler.GetString("pin.button.force_drum"), buttonPollingDuration),
		buttonForcePump:     gpio.NewButtonDriver(board, configHandler.GetString("pin.button.force_pump"), buttonPollingDuration),
		captorSecurityUpper: gpio.NewButtonDriver(board, configHandler.GetString("pin.captor.security_upper"), buttonPollingDuration),
		captorSecurityUnder: gpio.NewButtonDriver(board, configHandler.GetString("pin.captor.security_under"), buttonPollingDuration),
		captorWaterUpper:    gpio.NewButtonDriver(board, configHandler.GetString("pin.captor.water_upper"), buttonPollingDuration),
		captorWaterUnder:    gpio.NewButtonDriver(board, configHandler.GetString("pin.captor.water_under"), buttonPollingDuration),
		timeBetweenWash:     time.NewTicker(time.Duration(1 * time.Nanosecond)),
		Eventer:             gobot.NewEventer(),
		schedulingRoutines:  make([]*time.Ticker, 0),
		mailClient:          mailClient,
	}

	// Create gobot robot
	dfpBoard.gobot = gobot.NewRobot(
		dfpBoard.Name(),
		[]gobot.Connection{dfpBoard.board},
		[]gobot.Device{
			dfpBoard.relayDrum,
			dfpBoard.relayPump,
			dfpBoard.ledGreen,
			dfpBoard.ledRed,
			dfpBoard.buttonEmergencyStop,
			dfpBoard.buttonForceDrum,
			dfpBoard.buttonForcePump,
			dfpBoard.buttonStart,
			dfpBoard.buttonStop,
			dfpBoard.buttonWash,
			dfpBoard.captorSecurityUnder,
			dfpBoard.captorSecurityUpper,
			dfpBoard.captorWaterUnder,
			dfpBoard.captorWaterUpper,
		},
		dfpBoard.work,
	)

	// Add events on eventer
	dfpBoard.AddEvent(EventNewConfig)
	dfpBoard.AddEvent(EventNewState)
	dfpBoard.AddEvent(EventWash)
	dfpBoard.AddEvent(EventStopDFP)
	dfpBoard.AddEvent(EventStartDFP)
	dfpBoard.AddEvent(EventSetSecurity)
	dfpBoard.AddEvent(EventUnsetSecurity)
	dfpBoard.AddEvent(EventSetEmergencyStop)
	dfpBoard.AddEvent(EventUnsetEmergencyStop)
	dfpBoard.AddEvent(EventBoardStop)

	log.Infof("Board %s initialized successfully", dfpBoard.Name())

	return dfpBoard
}

// Start will init some item, like INPUT_PULLUP button, then start gobot
func (h *DFPBoard) Start(ctx context.Context) (err error) {

	// Start connexion to set some initial state on I/O
	if err = h.board.Connect(); err != nil {
		return err
	}

	// Set all input as INPUT_PULLUP and set default state as 1
	listPins := []*gpio.ButtonDriver{
		h.buttonEmergencyStop,
		h.buttonForceDrum,
		h.buttonForcePump,
		h.buttonStart,
		h.buttonStop,
		h.buttonWash,
		h.captorSecurityUnder,
		h.captorSecurityUpper,
		h.captorWaterUnder,
		h.captorWaterUpper,
	}
	if err = h.board.SetInputPullup(listPins); err != nil {
		return err
	}
	h.captorSecurityUpper.DefaultState = 0
	h.captorWaterUpper.DefaultState = 0

	// Init state
	if h.state.IsRunning && !h.state.IsSecurity && !h.state.IsEmergencyStopped {
		h.turnOnGreenLed()
		h.turnOffRedLed()
	} else {
		// It stopped or security
		h.forceStopRelais()
		h.turnOffGreenLed()
		h.turnOnRedLed()
	}

	// If on current wash
	if h.state.IsWashed {
		h.wash()
	}

	if err = h.gobot.Start(false); err != nil {
		return err
	}

	log.Infof("Board initialized")

	h.isOnline = true

	helper.SendEvent(ctx, h.eventUsecase, h.name, helper.KindEventStartBoard, h.name)

	return nil

}

// Stop permit to stop gobot.
// It send event of name `stop`. It can be used to stop routines.
func (h *DFPBoard) Stop(ctx context.Context) (err error) {

	// Internal event
	h.Publish(EventBoardStop, nil)

	// Stop outputs
	h.forceStopRelais()
	h.turnOffGreenLed()
	h.turnOffRedLed()

	// Not publish on global event to avoid stop pump and uvc each time we restart program
	// It can be dangerous to stop board.

	// Then stop board
	if h.isOnline {
		if err = h.gobot.Stop(); err != nil {
			return err
		}
	}

	// Stop scheduling routines
	for _, ticker := range h.schedulingRoutines {
		ticker.Stop()
	}
	h.schedulingRoutines = make([]*time.Ticker, 0)

	h.isOnline = false
	h.isInitialized = false

	helper.SendEvent(ctx, h.eventUsecase, h.name, helper.KindEventStopBoard, h.name)

	return nil
}

// Name return the current board name
func (h *DFPBoard) Name() string {
	return h.name
}

// Board return public board data
func (h *DFPBoard) Board() *models.Board {
	return &models.Board{
		Name:     h.Name(),
		IsOnline: h.isOnline,
	}
}

// IsOnline return is board is online
func (h *DFPBoard) IsOnline() bool {
	return h.isOnline
}

// State return copy of current state
func (h *DFPBoard) State() (state models.DFPState) {
	return *h.state
}

func (h *DFPBoard) IO() models.DFPIO {
	io := models.DFPIO{}

	// Led state
	if h.ledGreen.State() {
		io.GreenLed = "on"
	} else {
		io.GreenLed = "off"
	}
	if h.ledRed.State() {
		io.RedLed = "on"
	} else {
		io.RedLed = "off"
	}

	// Relais state
	if h.relayDrum.State() {
		io.DrumRelay = "on"
	} else {
		io.DrumRelay = "off"
	}
	if h.relayPump.State() {
		io.PumpRelay = "on"
	} else {
		io.PumpRelay = "off"
	}

	return io
}
