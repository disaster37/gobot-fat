package dfpboard

import (
	"context"
	"time"

	"github.com/disaster37/go-arest/arest"
	"github.com/disaster37/go-arest/arest/device/gpio/button"
	"github.com/disaster37/go-arest/arest/device/gpio/led"
	"github.com/disaster37/go-arest/arest/device/gpio/relay"
	"github.com/disaster37/go-arest/arest/serial"
	"github.com/disaster37/gobot-fat/board"
	"github.com/disaster37/gobot-fat/dfp"
	dfpconfig "github.com/disaster37/gobot-fat/dfp_config"
	dfpstate "github.com/disaster37/gobot-fat/dfp_state"
	"github.com/disaster37/gobot-fat/event"
	"github.com/disaster37/gobot-fat/models"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// DFPHandler manage DFP command
type DFPHandler struct {
	state            *models.DFPState
	board            arest.Arest
	configUsecase    dfpconfig.Usecase
	eventUsecase     event.Usecase
	stateUsecase     dfpstate.Usecase
	configHandler    *viper.Viper
	isOnline         bool
	chStop           chan bool
	relayDrum        relay.Relay
	relayPump        relay.Relay
	ledGreen         led.Led
	ledRed           led.Led
	ledButtons       []led.Led
	buttonStart      button.Button
	buttonStop       button.Button
	buttonForceDrum  button.Button
	buttonForcePump  button.Button
	buttonWash       button.Button
	buttonSet        button.Button
	captorSecurities []button.Button
	captorWaters     []button.Button
	config           *models.DFPConfig
	timerLED         *time.Timer
	turnOffLED       bool
	isRunning        bool
}

// NewDFP create handler to manage FAT
func NewDFP(configHandler *viper.Viper, configUsecase dfpconfig.Usecase, eventUsecase event.Usecase, stateUsecase dfpstate.Usecase, state *models.DFPState) (dfpHandler dfp.Board) {

	dfpHandler = &DFPHandler{
		state:            state,
		configUsecase:    configUsecase,
		eventUsecase:     eventUsecase,
		stateUsecase:     stateUsecase,
		configHandler:    configHandler,
		chStop:           make(chan bool),
		isOnline:         false,
		isRunning:        false,
		captorSecurities: make([]button.Button, 0, 0),
		captorWaters:     make([]button.Button, 0, 0),
		timerLED:         time.NewTimer(60 * time.Second),
		turnOffLED:       false,
	}

	return dfpHandler
}

// State return the current state
func (h *DFPHandler) State() models.DFPState {
	return *h.state
}

// Name is the board name
func (h *DFPHandler) Name() string {
	return h.state.Name
}

// IsOnline is true if board is online
func (h *DFPHandler) IsOnline() bool {
	return h.isOnline
}

// Board get public board data
func (h *DFPHandler) Board() *models.Board {
	return &models.Board{
		Name:     h.state.Name,
		IsOnline: h.isOnline,
	}
}

// StartDFP start DFP
func (h *DFPHandler) StartDFP(ctx context.Context) error {
	if !h.state.IsRunning {
		h.state.IsRunning = true
		err := h.stateUsecase.Update(ctx, h.state)
		if err != nil {
			return err
		}

		h.sendEvent(ctx, "start_dfp", "board")

		log.Debug("Start DFP")
	} else {
		log.Debug("DFP already running")
	}

	if !h.state.Security() {
		err := h.ledGreen.TurnOn(ctx)
		if err != nil {
			return err
		}
		err = h.ledRed.TurnOff(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

// StopDFP stop dfp
func (h *DFPHandler) StopDFP(ctx context.Context) error {

	routine := board.NewRoutine(ctx, h.stopDFP)
	select {
	case err := <-routine.Error():
		return err
	case <-routine.Result():
		break
	}

	if h.state.IsRunning {
		h.state.IsRunning = false
		err := h.stateUsecase.Update(ctx, h.state)
		if err != nil {
			return err
		}

		h.sendEvent(ctx, "stop_dfp", "board")

		log.Debug("Stop DFP")
	} else {
		log.Debug("DFP Already stopped")
	}

	return nil
}

// ForceWashing run imediate washing if state permit it
func (h *DFPHandler) ForceWashing(ctx context.Context) error {

	if h.state.ShouldWash() {
		routine := board.NewRoutine(ctx, h.wash)
		select {
		case err := <-routine.Error():
			return err
		case <-routine.Result():
			break
		}
		log.Debug("Run wash successfully")
	} else {
		log.Debug("Wash can't be start because of state")
	}

	return nil
}

// StartManualDrum start drum
func (h *DFPHandler) StartManualDrum(ctx context.Context) error {
	if !h.state.IsEmergencyStopped {
		err := h.relayDrum.On(ctx)
		if err != nil {
			return err
		}

		h.sendEvent(ctx, "manual_start_drum", "motor")
		log.Debug("Start drum successfully")
	} else {
		log.Debug("Drum can't start because of state")
	}

	return nil
}

// StopManualDrum stop drum
func (h *DFPHandler) StopManualDrum(ctx context.Context) error {

	h.sendEvent(ctx, "manual_stop_drum", "motor")

	return h.relayDrum.Off(ctx)
}

// StartManualPump push startPump buton on DFP
func (h *DFPHandler) StartManualPump(ctx context.Context) error {

	if !h.state.IsEmergencyStopped {
		err := h.relayPump.On(ctx)
		if err != nil {
			return err
		}

		h.sendEvent(ctx, "manual_start_pump", "motor")

		log.Debug("Start pump successfully")
	} else {
		log.Debug("Can't start pump because of state")
	}

	return nil
}

// StopManualPump stop pump
func (h *DFPHandler) StopManualPump(ctx context.Context) error {
	err := h.relayPump.Off(ctx)
	h.sendEvent(ctx, "manual_stop_pump", "motor")
	return err
}

// Start run the main board function
func (h *DFPHandler) Start(ctx context.Context) error {

	// Put timeout to start
	ctxWithTiemout, cancel := context.WithTimeout(ctx, 10*time.Second)
	defer cancel()

	c, err := serial.NewClient(h.configHandler.GetString("url"), 1*time.Second, false)
	if err != nil {
		return err
	}
	h.board = c

	// Read arbitrary value to check if board is online
	_, err = h.board.ReadValue(ctxWithTiemout, "isRebooted")
	if err != nil {
		return err
	}

	// Load config
	config, err := h.configUsecase.Get(ctx)
	if err != nil {
		return err
	}
	h.config = config

	// Relay
	NORelay := relay.NewOutput()
	NORelay.SetOutputNO()
	highSignal := arest.NewLevel()
	highSignal.SetLevelHigh()
	stateOff := relay.NewState()
	stateOff.SetStateOff()
	relayDrum, err := relay.NewRelay(h.board, h.configHandler.GetInt("pin.relay.drum"), highSignal, NORelay, stateOff)
	if err != nil {
		return err
	}
	h.relayDrum = relayDrum

	relayPump, err := relay.NewRelay(h.board, h.configHandler.GetInt("pin.relay.pump"), highSignal, NORelay, stateOff)
	if err != nil {
		return err
	}
	h.relayPump = relayPump

	// Buttons
	buttonStart, err := button.NewButton(h.board, h.configHandler.GetInt("pin.button.start"), highSignal, true)
	if err != nil {
		return err
	}
	h.buttonStart = buttonStart

	buttonStop, err := button.NewButton(h.board, h.configHandler.GetInt("pin.button.stop"), highSignal, true)
	if err != nil {
		return err
	}
	h.buttonStop = buttonStop

	buttonWash, err := button.NewButton(h.board, h.configHandler.GetInt("pin.button.wash"), highSignal, true)
	if err != nil {
		return err
	}
	h.buttonWash = buttonWash

	buttonForceDrum, err := button.NewButton(h.board, h.configHandler.GetInt("pin.button.force_drump"), highSignal, true)
	if err != nil {
		return err
	}
	h.buttonForceDrum = buttonForceDrum

	buttonForcePump, err := button.NewButton(h.board, h.configHandler.GetInt("pin.button.force_pump"), highSignal, true)
	if err != nil {
		return err
	}
	h.buttonForcePump = buttonForcePump

	buttonSet, err := button.NewButton(h.board, h.configHandler.GetInt("pin.button.set"), highSignal, true)
	if err != nil {
		return err
	}
	h.buttonSet = buttonSet

	// Captors
	lowSignal := arest.NewLevel()
	lowSignal.SetLevelLow()
	captorWaterUpper, err := button.NewButton(h.board, h.configHandler.GetInt("pin.captor.water_upper"), lowSignal, true)
	if err != nil {
		return err
	}
	h.captorWaters = append(h.captorWaters, captorWaterUpper)

	captorWaterUnder, err := button.NewButton(h.board, h.configHandler.GetInt("pin.captor.water_under"), highSignal, true)
	if err != nil {
		return err
	}
	h.captorWaters = append(h.captorWaters, captorWaterUnder)

	captorSecurityUpper, err := button.NewButton(h.board, h.configHandler.GetInt("pin.captor.security_upper"), lowSignal, true)
	if err != nil {
		return err
	}
	h.captorSecurities = append(h.captorSecurities, captorSecurityUpper)

	captorSecurityUnder, err := button.NewButton(h.board, h.configHandler.GetInt("pin.captor.security_under"), highSignal, true)
	if err != nil {
		return err
	}
	h.captorSecurities = append(h.captorSecurities, captorSecurityUnder)

	// Leds
	var defaultLedState bool
	if h.state.IsRunning {
		defaultLedState = true
	} else {
		defaultLedState = false
	}
	ledGreen, err := led.NewLed(h.board, h.configHandler.GetInt("pin.led.green"), defaultLedState)
	if err != nil {
		return err
	}
	h.ledGreen = ledGreen

	if !h.state.IsRunning {
		defaultLedState = true
	} else {
		defaultLedState = false
	}
	ledRed, err := led.NewLed(h.board, h.configHandler.GetInt("pin.led.red"), defaultLedState)
	if err != nil {
		return err
	}
	h.ledRed = ledRed

	ledButtonAuto, err := led.NewLed(h.board, h.configHandler.GetInt("pin.led.button_auto"), false)
	if err != nil {
		return err
	}
	h.ledButtons = append(h.ledButtons, ledButtonAuto)

	ledButtonStop, err := led.NewLed(h.board, h.configHandler.GetInt("pin.led.button_stop"), false)
	if err != nil {
		return err
	}
	h.ledButtons = append(h.ledButtons, ledButtonStop)

	ledButtonForceDrum, err := led.NewLed(h.board, h.configHandler.GetInt("pin.led.button_force_drum"), false)
	if err != nil {
		return err
	}
	h.ledButtons = append(h.ledButtons, ledButtonForceDrum)

	ledButtonForcePump, err := led.NewLed(h.board, h.configHandler.GetInt("pin.led.button_force_pump"), false)
	if err != nil {
		return err
	}
	h.ledButtons = append(h.ledButtons, ledButtonForcePump)

	ledButtonWash, err := led.NewLed(h.board, h.configHandler.GetInt("pin.led.button_wash"), false)
	if err != nil {
		return err
	}
	h.ledButtons = append(h.ledButtons, ledButtonWash)

	ledButtonSet, err := led.NewLed(h.board, h.configHandler.GetInt("pin.led.button_set"), false)
	if err != nil {
		return err
	}
	h.ledButtons = append(h.ledButtons, ledButtonSet)

	ledLCDEnable, err := led.NewLed(h.board, h.configHandler.GetInt("pin.led.lcd_enable"), false)
	if err != nil {
		return err
	}
	h.ledButtons = append(h.ledButtons, ledLCDEnable)

	ledLCDLED, err := led.NewLed(h.board, h.configHandler.GetInt("pin.led.lcd_led"), false)
	if err != nil {
		return err
	}
	h.ledButtons = append(h.ledButtons, ledLCDLED)

	// Handle reboot
	h.handleReboot(ctx)
	board.NewHandler(ctx, 10*time.Second, h.chStop, h.handleReboot)

	// Handle config
	h.handleConfig(ctx)
	board.NewHandler(ctx, 10*time.Second, h.chStop, h.handleConfig)

	// Handle state
	h.isRunning = true
	board.NewHandler(ctx, 1*time.Nanosecond, h.chStop, h.handleState)

	h.isOnline = true

	h.sendEvent(ctx, "board_dfp_start", "board")

	log.Infof("Board %s initialized successfully", h.Name())

	return nil

}

// Stop permit to stop the board
func (h *DFPHandler) Stop(ctx context.Context) error {

	h.chStop <- true

	h.isRunning = false

	h.sendEvent(ctx, "board_dfp_stop", "board")

	log.Infof("Board %s sucessfully stoped", h.Name())

	return nil
}
