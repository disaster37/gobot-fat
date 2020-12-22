package dfpboard

import (
	"context"

	dfpconfig "github.com/disaster37/gobot-fat/dfp_config"
	dfpstate "github.com/disaster37/gobot-fat/dfp_state"
	"github.com/disaster37/gobot-fat/event"
	"github.com/disaster37/gobot-fat/models"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/stianeikeland/go-rpio/v4"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"
)

// DFPBoard is the DFP board
type DFPBoard struct {
	state               *models.DFPState
	board               *raspi.Adaptor
	gobot               *gobot.Robot
	configUsecase       dfpconfig.Usecase
	eventUsecase        event.Usecase
	stateUsecase        dfpstate.Usecase
	configHandler       *viper.Viper
	isOnline            bool
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
	buttonSet           *gpio.ButtonDriver
	captorWaterUpper    *gpio.ButtonDriver
	captorWaterUnder    *gpio.ButtonDriver
	captorSecurityUpper *gpio.ButtonDriver
	captorSecurityUnder *gpio.ButtonDriver
	config              *models.DFPConfig
	isRunning           bool
	gobot.Eventer
}

//NewDFPBoard return the DFP board with all IO created but not started
func NewDFP(configHandler *viper.Viper, configUsecase dfpconfig.Usecase, eventUsecase event.Usecase, stateUsecase dfpstate.Usecase, state *models.DFPState) (dfpBoard *DFPBoard) {

	dfpBoard = &DFPBoard{
		configHandler: configHandler,
		configUsecase: configUsecase,
		eventUsecase:  eventUsecase,
		stateUsecase:  stateUsecase,
		state:         state,
		board:         raspi.NewAdaptor(),
		Eventer:       gobot.NewEventer(),
	}

	// Create relay
	dfpBoard.relayDrum = gpio.NewRelayDriver(dfpBoard.board, configHandler.GetString("pin.relay.drum"))
	dfpBoard.relayPump = gpio.NewRelayDriver(dfpBoard.board, configHandler.GetString("pin.relay.pump"))

	// Create LED
	dfpBoard.ledGreen = gpio.NewLedDriver(dfpBoard.board, configHandler.GetString("pin.led.green"))
	dfpBoard.ledRed = gpio.NewLedDriver(dfpBoard.board, configHandler.GetString("pin.led.red"))

	// Create button
	dfpBoard.buttonSet = gpio.NewButtonDriver(dfpBoard.board, configHandler.GetString("pin.button.set"))
	dfpBoard.buttonSet.DefaultState = 1
	dfpBoard.buttonStart = gpio.NewButtonDriver(dfpBoard.board, configHandler.GetString("pin.button.start"))
	dfpBoard.buttonStart.DefaultState = 1
	dfpBoard.buttonStop = gpio.NewButtonDriver(dfpBoard.board, configHandler.GetString("pin.button.stop"))
	dfpBoard.buttonStop.DefaultState = 1
	dfpBoard.buttonWash = gpio.NewButtonDriver(dfpBoard.board, configHandler.GetString("pin.button.wash"))
	dfpBoard.buttonWash.DefaultState = 1
	dfpBoard.buttonForceDrum = gpio.NewButtonDriver(dfpBoard.board, configHandler.GetString("pin.button.force_drum"))
	dfpBoard.buttonForceDrum.DefaultState = 1
	dfpBoard.buttonForcePump = gpio.NewButtonDriver(dfpBoard.board, configHandler.GetString("pin.button.force_pump"))
	dfpBoard.buttonForcePump.DefaultState = 1

	// Create water captors
	dfpBoard.captorSecurityUpper = gpio.NewButtonDriver(dfpBoard.board, configHandler.GetString("pin.captor.security_upper"))
	dfpBoard.captorSecurityUnder = gpio.NewButtonDriver(dfpBoard.board, configHandler.GetString("pin.captor.security_under"))
	dfpBoard.captorSecurityUnder.DefaultState = 1
	dfpBoard.captorWaterUpper = gpio.NewButtonDriver(dfpBoard.board, configHandler.GetString("pin.captor.water_upper"))
	dfpBoard.captorWaterUnder = gpio.NewButtonDriver(dfpBoard.board, configHandler.GetString("pin.captor.water_under"))
	dfpBoard.captorWaterUnder.DefaultState = 1

	dfpBoard.gobot = gobot.NewRobot(
		configHandler.GetString("name"),
		[]gobot.Connection{dfpBoard.board},
		[]gobot.Device{
			dfpBoard.relayDrum,
			dfpBoard.relayPump,
			dfpBoard.ledGreen,
			dfpBoard.ledRed,
			dfpBoard.buttonSet,
			dfpBoard.buttonStart,
			dfpBoard.buttonStop,
			dfpBoard.buttonWash,
			dfpBoard.buttonForceDrum,
			dfpBoard.buttonForcePump,
			dfpBoard.captorSecurityUpper,
			dfpBoard.captorSecurityUnder,
			dfpBoard.captorWaterUpper,
			dfpBoard.captorWaterUnder,
		},
		dfpBoard.work,
	)

	dfpBoard.AddEvent("state")

	return dfpBoard

}

// Start will init some item, like INPUT_PULLUP button, then start gobot
func (h *DFPBoard) Start(ctx context.Context) (err error) {

	// Load config
	config, err := h.configUsecase.Get(context.TODO())
	if err != nil {
		return err
	}
	h.config = config

	if log.IsLevelEnabled(log.DebugLevel) {
		log.Debugf("Current config: %+v", h.config)
		log.Debugf("Current state %+v", h.state)
	}

	// Start connection on board and set INPUT_PULLUP on some pins
	err = h.board.Connect()
	if err != nil {
		return err
	}

	listPins := []int{
		h.configHandler.GetInt("pin.button.set"),
		h.configHandler.GetInt("pin.button.start"),
		h.configHandler.GetInt("pin.button.stop"),
		h.configHandler.GetInt("pin.button.wash"),
		h.configHandler.GetInt("pin.button.force_drum"),
		h.configHandler.GetInt("pin.button.force_pump"),
		h.configHandler.GetInt("pin.captor.water_upper"),
		h.configHandler.GetInt("pin.captor.water_under"),
		h.configHandler.GetInt("pin.captor.security_upper"),
		h.configHandler.GetInt("pin.captor.security_under"),
	}

	err = rpio.Open()
	if err != nil {
		return err
	}
	for _, pin := range listPins {
		pin := rpio.Pin(pin)
		pin.Input()
		pin.PullUp()
	}
	rpio.Close()

	return h.gobot.Start(false)

}

// Stop permit to stop gobot.
// It send event of name `stop`. It can be used to stop routines.
func (h *DFPBoard) Stop(ctx context.Context) (err error) {

	err = h.gobot.Stop()
	if err != nil {
		return err
	}

	h.Publish("stop", true)

	return nil
}

// Name return the current board name
func (h *DFPBoard) Name() string {
	return h.gobot.Name
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

// StartDFP put dfp on auto
func (h *DFPBoard) StartDFP(ctx context.Context) (err error) {

	if !h.state.IsRunning {
		h.state.IsRunning = true
		err = h.ledGreen.On()
		if err != nil {
			return
		}
		h.Publish("state", h.state)
		h.sendEvent(ctx, "board", "dfp_start")
	}

	return
}

// StopDFP stop dfp and disable auto
func (h *DFPBoard) StopDFP(ctx context.Context) (err error) {

	if h.state.IsRunning {
		h.state.IsRunning = false
		err = h.ledGreen.Off()
		if err != nil {
			return
		}
		h.Publish("state", h.state)
		h.sendEvent(ctx, "board", "dfp_stop")
	}

	return
}

// ForceWashing start a washing cycle
func (h *DFPBoard) ForceWashing(ctx context.Context) (err error) {
	if !h.state.IsWashed && !h.state.IsEmergencyStopped {
		log.Debug("Run force wash")
		h.wash()
	}

	return
}

// StartManualDrum force start drum motor
// Only if not already wash and is not on emergency stopped
func (h *DFPBoard) StartManualDrum(ctx context.Context) (err error) {

	if !h.state.IsWashed && !h.state.IsEmergencyStopped {
		if log.IsLevelEnabled(log.DebugLevel) {
			log.Debug("Run force drum")
		}

		err = h.relayDrum.On()
		if err != nil {
			return
		}

	}
	return
}

// StopManualDrum force stop drum motor
// Only if not current washing
func (h *DFPBoard) StopManualDrum(ctx context.Context) (err error) {

	if !h.state.IsWashed {
		if log.IsLevelEnabled(log.DebugLevel) {
			log.Debug("Stop force drum")
		}

		err = h.relayDrum.Off()
		if err != nil {
			return
		}
	}
	return
}

// StartManualPump force start pump
// Only if not already wash and is not on emergency stopped
func (h *DFPBoard) StartManualPump(ctx context.Context) (err error) {

	if !h.state.IsWashed && !h.state.IsEmergencyStopped {
		if log.IsLevelEnabled(log.DebugLevel) {
			log.Debug("Run force pump")
		}

		err = h.relayPump.On()
		if err != nil {
			return
		}
	}

	return
}

// StopManualPump force stop pump
// Only if not already wash
func (h *DFPBoard) StopManualPump(ctx context.Context) (err error) {

	// Stop force pump
	if !h.state.IsWashed {
		if log.IsLevelEnabled(log.DebugLevel) {
			log.Debug("Stop force pump")
		}

		err = h.relayPump.Off()
		if err != nil {
			return
		}
	}

	return
}

// State return copy of current state
func (h *DFPBoard) State() (state models.DFPState) {
	return *h.state
}
