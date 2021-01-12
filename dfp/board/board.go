package dfpboard

import (
	"context"
	"fmt"
	"time"

	"github.com/disaster37/gobot-fat/dfp"
	"github.com/disaster37/gobot-fat/event"
	"github.com/disaster37/gobot-fat/models"
	"github.com/disaster37/gobot-fat/usecase"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/stianeikeland/go-rpio/v4"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"
)

const (
	NewConfig   = "new-config"
	NewReboot   = "new-reboot"
	NewOffline  = "new-offline"
	NewWash     = "new-wash"
	NewSecurity = "new-security"
	NewState    = "new-state"
	Stop        = "stop"
)

// DFPAdaptor is DFP board interface
type DFPAdaptor interface {
	gobot.Adaptor
	gpio.DigitalReader
	gpio.DigitalWriter
}

// DFPBoard is the DFP board
type DFPBoard struct {
	state               *models.DFPState
	config              *models.DFPConfig
	board               DFPAdaptor
	gobot               *gobot.Robot
	eventUsecase        event.Usecase
	stateUsecase        usecase.UsecaseCRUD
	configHandler       *viper.Viper
	isOnline            bool
	isInitialized       bool
	relayDrum           *gpio.RelayDriver
	relayPump           *gpio.RelayDriver
	ledGreen            *gpio.LedDriver
	ledRed              *gpio.LedDriver
	ledButtons          []*gpio.LedDriver
	buttonStart         *gpio.ButtonDriver
	buttonStop          *gpio.ButtonDriver
	buttonForceDrum     *gpio.ButtonDriver
	buttonForcePump     *gpio.ButtonDriver
	buttonWash          *gpio.ButtonDriver
	buttonEmergencyStop *gpio.ButtonDriver
	captorWaterUpper    *gpio.ButtonDriver
	captorWaterUnder    *gpio.ButtonDriver
	captorSecurityUpper *gpio.ButtonDriver
	captorSecurityUnder *gpio.ButtonDriver
	globalEventer       gobot.Eventer
	isRunning           bool
	name                string
	chStop              chan bool
	timeBetweenWash     *time.Ticker
	gobot.Eventer
}

// NewDFP create board to manage DFP
func NewDFP(configHandler *viper.Viper, config *models.DFPConfig, state *models.DFPState, eventUsecase event.Usecase, dfpStateUsecase usecase.UsecaseCRUD, eventer gobot.Eventer) (dfpBoard dfp.Board) {

	//Create client
	c := raspi.NewAdaptor()

	return newDFP(c, configHandler, config, state, eventUsecase, dfpStateUsecase, eventer)

}

func newDFP(board DFPAdaptor, configHandler *viper.Viper, config *models.DFPConfig, state *models.DFPState, eventUsecase event.Usecase, dfpStateUsecase usecase.UsecaseCRUD, eventer gobot.Eventer) dfp.Board {

	buttonPollingDuration := configHandler.GetDuration("button_polling") * time.Millisecond
	// Create struct
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
		chStop:              make(chan bool),
		timeBetweenWash:     time.NewTicker(time.Duration(1 * time.Nanosecond)),
		Eventer:             gobot.NewEventer(),
	}

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

	dfpBoard.AddEvent(NewConfig)
	dfpBoard.AddEvent(NewReboot)
	dfpBoard.AddEvent(NewOffline)
	dfpBoard.AddEvent(NewWash)
	dfpBoard.AddEvent(NewSecurity)
	dfpBoard.AddEvent(NewState)

	log.Infof("Board %s initialized successfully", dfpBoard.Name())

	return dfpBoard

}

// Start will init some item, like INPUT_PULLUP button, then start gobot
func (h *DFPBoard) Start(ctx context.Context) (err error) {

	// Start connexion to set some initial state on I/O
	if err := h.board.Connect(); err != nil {
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

	if err := rpio.Open(); err != nil {
		log.Errorf("Error when open rpio: %s", err.Error())
		return err
	}
	defer rpio.Close()

	for _, button := range listPins {

		// Need to translate pin
		translatedPin, err := translatePin(button.Pin(), "3")
		if err != nil {
			return err
		}
		pin := rpio.Pin(translatedPin)
		pin.Input()
		pin.PullUp()

		button.DefaultState = 1
	}

	log.Infof("RPIO initialized")
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

	if err := h.gobot.Start(false); err != nil {
		return err
	}

	log.Infof("Board initialized")

	h.isOnline = true

	h.sendEvent(ctx, fmt.Sprintf("start_%s", h.name), "board")

	return nil

}

// Stop permit to stop gobot.
// It send event of name `stop`. It can be used to stop routines.
func (h *DFPBoard) Stop(ctx context.Context) (err error) {

	err = h.gobot.Stop()
	if err != nil {
		return err
	}

	// Stop internal routine
	h.chStop <- true

	h.isOnline = false
	h.isInitialized = false

	h.sendEvent(ctx, fmt.Sprintf("stop_%s", h.name), "board")

	h.Publish(Stop, true)

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
