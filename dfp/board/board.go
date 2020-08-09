package dfpboard

import (
	"context"
	"time"

	"github.com/disaster37/go-arest"
	"github.com/disaster37/go-arest/device/gpio/button"
	"github.com/disaster37/go-arest/device/gpio/led"
	"github.com/disaster37/go-arest/device/gpio/relay"
	"github.com/disaster37/go-arest/serial"
	"github.com/disaster37/gobot-fat/dfp"
	dfpconfig "github.com/disaster37/gobot-fat/dfp_config"
	dfpstate "github.com/disaster37/gobot-fat/dfp_state"
	"github.com/disaster37/gobot-fat/event"
	"github.com/disaster37/gobot-fat/helper"
	"github.com/disaster37/gobot-fat/models"
	"github.com/labstack/gommon/log"
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
	routines         []*time.Ticker
	relayDrum        relay.Relay
	relayPump        relay.Relay
	ledGreen         led.Led
	ledRed           led.Led
	ledButtons       []led.Led
	buttonAuto       button.Button
	buttonStop       button.Button
	buttonForceDrum  button.Button
	buttonForcePump  button.Button
	buttonWash       button.Button
	buttonSet        button.Button
	captorSecurities []button.Button
	captorWaters     []button.Button
	config           *models.DFPConfig
}

// NewDFP create handler to manage FAT
func NewDFP(configHandler *viper.Viper, configUsecase dfpconfig.Usecase, eventUsecase event.Usecase, stateUsecase dfpstate.Usecase, state *models.DFPState) (dfpHandler dfp.Board) {

	dfpHandler = &DFPHandler{
		state:            state,
		configUsecase:    configUsecase,
		eventUsecase:     eventUsecase,
		stateUsecase:     stateUsecase,
		configHandler:    configHandler,
		routines:         make([]*time.Ticker, 0, 0),
		isOnline:         false,
		captorSecurities: make([]button.Button, 0, 0),
		captorWaters:     make([]button.Button, 0, 0),
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

// Auto put DFP on auto mode
func (h *DFPHandler) Auto(ctx context.Context) error {
	if !h.state.IsRunning {
		h.state.IsRunning = true
		err := h.stateUsecase.Update(context.Background(), h.state)
		if err != nil {
			return err
		}

		log.Debug("Put DFP on auto mode")
	} else {
		log.Debug("DFP already on auto mode")
	}

	return nil
}

// StopDFP stop dfp
func (h *DFPHandler) StopDFP(ctx context.Context) error {

	err := h.relayDrum.Off()
	if err != nil {
		return err
	}
	err = h.relayPump.Off()
	if err != nil {
		return err
	}

	if h.state.IsRunning {
		h.state.IsRunning = false
		err := h.stateUsecase.Update(context.Background(), h.state)
		if err != nil {
			return err
		}
	}

	return err
}

// ForceWashing push washing buton on DFP
func (h *DFPHandler) ForceWashing(ctx context.Context) error {

	if h.state.ShouldWash() {
		err := h.wash()
		if err != nil {
			return err
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
		err := h.relayDrum.On()
		if err != nil {
			return err
		}
		log.Debug("Start drum successfully")
	} else {
		log.Debug("Drum can't start because of state")
	}

	return nil
}

// StopManualDrum stop drum
func (h *DFPHandler) StopManualDrum(ctx context.Context) error {
	return h.relayDrum.Off()
}

// StartManualPump push startPump buton on DFP
func (h *DFPHandler) StartManualPump(ctx context.Context) error {

	if !h.state.IsEmergencyStopped {
		err := h.relayPump.On()
		if err != nil {
			return err
		}
		log.Debug("Start pump successfully")
	} else {
		log.Debug("Can't start pump because of state")
	}

	return nil
}

// StopManualPump stop pump
func (h *DFPHandler) StopManualPump(ctx context.Context) error {
	return h.relayPump.Off()
}

// Start run the main board function
func (h *DFPHandler) Start() error {

	c, err := serial.NewClient(h.configHandler.GetString("url"))
	if err != nil {
		return err
	}
	h.board = c

	// Read arbitrary value to check if board is online
	_, err = h.board.ReadValue("isRebooted")
	if err != nil {
		return err
	}

	// Load config
	config, err := h.configUsecase.Get(context.Background())
	if err != nil {
		return err
	}
	h.config = config

	// Init I/O

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
	buttonAuto, err := button.NewButton(h.board, h.configHandler.GetInt("pin.button.auto"), highSignal, true)
	if err != nil {
		return err
	}
	h.buttonAuto = buttonAuto

	buttonStop, err := button.NewButton(h.board, h.configHandler.GetInt("pin.button.stop"), highSignal, true)
	if err != nil {
		return err
	}
	h.buttonStop = buttonStop

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
	ledGreen, err := led.NewLed(h.board, h.configHandler.GetInt("led.pin.green"), defaultLedState)
	if err != nil {
		return err
	}
	h.ledGreen = ledGreen

	if !h.state.IsRunning {
		defaultLedState = true
	} else {
		defaultLedState = false
	}
	ledRed, err := led.NewLed(h.board, h.configHandler.GetInt("led.pin.red"), defaultLedState)
	if err != nil {
		return err
	}
	h.ledRed = ledRed

	ledButtonAuto, err := led.NewLed(h.board, h.configHandler.GetInt("led.pin.button_auto"), false)
	if err != nil {
		return err
	}
	h.ledButtons = append(h.ledButtons, ledButtonAuto)

	ledButtonStop, err := led.NewLed(h.board, h.configHandler.GetInt("led.pin.button_stop"), false)
	if err != nil {
		return err
	}
	h.ledButtons = append(h.ledButtons, ledButtonStop)

	ledButtonForceDrum, err := led.NewLed(h.board, h.configHandler.GetInt("led.pin.button_force_drum"), false)
	if err != nil {
		return err
	}
	h.ledButtons = append(h.ledButtons, ledButtonForceDrum)

	ledButtonForcePump, err := led.NewLed(h.board, h.configHandler.GetInt("led.pin.button_force_pump"), false)
	if err != nil {
		return err
	}
	h.ledButtons = append(h.ledButtons, ledButtonForcePump)

	ledButtonWash, err := led.NewLed(h.board, h.configHandler.GetInt("led.pin.button_wash"), false)
	if err != nil {
		return err
	}
	h.ledButtons = append(h.ledButtons, ledButtonWash)

	ledButtonSet, err := led.NewLed(h.board, h.configHandler.GetInt("led.pin.button_set"), false)
	if err != nil {
		return err
	}
	h.ledButtons = append(h.ledButtons, ledButtonSet)

	// Handle reboot
	h.routines = append(h.routines, helper.Every(10*time.Second, handleReboot(h)))

	// Handle state
	//h.routines = append(h.routines, helper.Every(1*time.Second, handleState(h)))

	h.isOnline = true

	log.Infof("Board %s initialized successfully", h.Name())

	return nil

}

// Stop permit to stop the board
func (h *DFPHandler) Stop() error {
	for _, routine := range h.routines {
		routine.Stop()
	}

	log.Infof("Board %s sucessfully stoped", h.Name())

	return nil
}

// handleReboot permit to check on background if board is rebooted
// If board is rebooted, it wil reset with current state
func handleReboot(handler *DFPHandler) func() {
	return func() {

		data, err := handler.board.ReadValue("isRebooted")
		if err != nil {
			log.Errorf("Error when read value isRebooted: %s", err.Error())
			handler.isOnline = false
			return
		}

		if data.(bool) {
			log.Info("Board %s has been rebooted, reset state", handler.Name())

			// Auto mode
			if handler.state.IsRunning && !handler.state.IsEmergencyStopped {
				err := handler.Auto(context.Background())
				if err != nil {
					log.Errorf("Error when reset auto mode: %s", err.Error())
				} else {
					log.Info("Successfully reset auto mode")
				}

			} else {
				// Stop / Ermergency mode
				err := handler.StopDFP(context.Background())
				if err != nil {
					log.Errorf("Error when reset stop mode: %s", err.Error())
				} else {
					log.Info("Seccessfully reset stop mode")
				}
			}

			// Washing
			if handler.state.IsWashed {
				err := handler.ForceWashing(context.Background())
				if err != nil {
					log.Errorf("Error when reset wash: %s", err.Error())
				} else {
					log.Info("Successfully reset wash")
				}
			}

			// Acknolege reboot
			_, err := handler.board.CallFunction("acknoledgeRebooted", "")
			if err != nil {
				log.Errorf("Error when aknoledge reboot on board %s: %s", handler.Name(), err.Error())
			}

			handler.isOnline = true

		}
	}
}

func handleState(h *DFPHandler) {
	for {

		// Read all values
		err := h.buttonAuto.Read()
		if err != nil {
			log.Errorf("Error when read button auto: %s", err.Error())
		}
		err = h.buttonForceDrum.Read()
		if err != nil {
			log.Errorf("Error when read button force drum: %s", err.Error())
		}
		err = h.buttonForcePump.Read()
		if err != nil {
			log.Errorf("Error when read button force pump: %s", err.Error())
		}
		err = h.buttonSet.Read()
		if err != nil {
			log.Errorf("Error when read button set: %s", err.Error())
		}
		err = h.buttonStop.Read()
		if err != nil {
			log.Errorf("Error when read button stop: %s", err.Error())
		}
		err = h.buttonWash.Read()
		if err != nil {
			log.Errorf("Error when read button wash; %s", err.Error())
		}
		for i, captor := range h.captorWaters {
			err := captor.Read()
			if err != nil {
				log.Errorf("Error when read water captor %d: %s", i, err.Error())
			}
		}
		for i, captor := range h.captorSecurities {
			err := captor.Read()
			if err != nil {
				log.Errorf("Error when read security captor %d: %s", i, err.Error())
			}
		}

		// Manage Security captor first
		for _, captor := range h.captorSecurities {
			if captor.IsPushed() && !h.state.IsDisableSecurity {
				h.state.IsSecurity = true
				err := h.StopDFP(context.Background())
				if err != nil {
					h.forceStopMotors()
					log.Errorf("Error when stop motor because of security state")
				}
				break
			}
		}

		buttonPushed := false

		// Stop / Auto button pushed
		if h.buttonStop.IsPushed() {
			buttonPushed = true
			err = h.StopDFP(context.Background())
			if err != nil {
				log.Errorf("Error when stop DFP: %s", err.Error())
			}
		} else if h.buttonAuto.IsPushed() {
			buttonPushed = true
			err = h.Auto(context.Background())
			if err != nil {
				log.Errorf("Error when set auto mode: %s", err.Error())
			}
		}

		// Force drum button
		if h.buttonForceDrum.IsPushed() {
			buttonPushed = true
			err = h.StartManualDrum(context.Background())
			if err != nil {
				log.Errorf("Error when force drum motor: %s", err.Error())
			}
		} else if h.buttonForceDrum.IsReleazed() {
			err = h.StartManualDrum(context.Background())
			if err != nil {
				log.Errorf("Error when stop drum motor: %s", err.Error())
			}
		}

		// Force pump
		if h.buttonForcePump.IsPushed() {
			buttonPushed = true
			err := h.StartManualPump(context.Background())
			if err != nil {
				log.Errorf("Error when force pump: %s", err.Error())
			}
		} else if h.buttonForcePump.IsReleazed() {
			err := h.StopManualPump(context.Background())
			if err != nil {
				log.Errorf("Error when stop pump: %s", err.Error)
			}
		}

		// Force wash
		if h.buttonWash.IsPushed() {
			buttonPushed = true
			err := h.ForceWashing(context.Background())
			if err != nil {
				log.Errorf("Error when force washing: %s", err.Error)
			}
		}

		//Set button
		if h.buttonSet.IsPushed() {
			buttonPushed = true
		}

		// Manage button led and screen
		//@TODO
		if buttonPushed {

		}

		// Manage captor state
		for _, captor := range h.captorWaters {
			if captor.IsPushed() && time.Now().After(h.state.LastWashing.Add(time.Duration(h.config.WaitTimeBetweenWashing)*time.Second)) {
				err := h.ForceWashing(context.Background())
				if err != nil {
					log.Errorf("Error when run wash: %s", err.Error())
				}
				break
			}
		}

	}
}

func (h *DFPHandler) sendEvent(eventType string, eventKind string) {
	event := &models.Event{
		SourceID:   h.state.Name,
		SourceName: h.state.Name,
		Timestamp:  time.Now(),
		EventType:  eventType,
		EventKind:  eventKind,
	}
	err := h.eventUsecase.Store(context.Background(), event)
	if err != nil {
		log.Errorf("Error when store new event: %s", err.Error())
	}
}

func (h *DFPHandler) wash() error {

	h.state.IsWashed = true

	// Blink led
	timerLed := h.ledGreen.Blink(time.Duration(h.config.StartWashingPumpBeforeWashing+h.config.WashingDuration) * time.Second)
	go func() {
		isTimerLedFinished := false
		go func() {
			for !isTimerLedFinished {
				if !h.state.IsWashed {
					timerLed.Stop()
					break
				}
				time.Sleep(1 * time.Millisecond)
			}
		}()
		<-timerLed.C
		isTimerLedFinished = true
	}()

	// Start pump and wait some time
	timerPump := time.NewTimer(time.Duration(h.config.StartWashingPumpBeforeWashing) * time.Second)
	isTimerPumpFinished := true
	go func() {
		for !isTimerPumpFinished {
			if !h.state.ShouldMotorStart() {
				timerPump.Stop()
				break
			}
			time.Sleep(1 * time.Millisecond)
		}
	}()
	err := h.relayPump.On()
	if err != nil {
		isTimerPumpFinished = true
		return err
	}
	<-timerPump.C
	isTimerPumpFinished = true

	// Start drum and wait some time
	timerWashing := time.NewTimer(time.Duration(h.config.WashingDuration) * time.Second)
	isTimerWashingFinished := false
	go func() {
		for !isTimerWashingFinished {
			if !h.state.ShouldMotorStart() {
				timerWashing.Stop()
				break
			}
			time.Sleep(1 * time.Millisecond)
		}
	}()
	err = h.relayDrum.On()
	if err != nil {
		isTimerWashingFinished = true
		h.forceStopMotors()
		return err
	}
	<-timerWashing.C
	isTimerWashingFinished = true

	// Stop all
	err = h.relayDrum.Off()
	if err != nil {
		h.forceStopMotors()
		return err
	}

	err = h.relayPump.Off()
	if err != nil {
		h.forceStopMotors()
		return err
	}

	h.state.IsWashed = false
	h.state.LastWashing = time.Now()
	err = h.stateUsecase.Update(context.Background(), h.state)
	if err != nil {
		return err
	}

	h.sendEvent("washing", "motor")

	log.Debugf("Washing successfully finished")
	return nil

}

func (h *DFPHandler) forceStopMotors() {
	go func() {
		isOk := false
		for !isOk {
			isOk = true
			err := h.relayDrum.Off()
			if err != nil {
				log.Errorf("Error appear when try to stop drum: %s", err.Error())
				isOk = false
			}

			err = h.relayPump.Off()
			if err != nil {
				log.Errorf("Error appear when try to stop pump: %s", err.Error())
				isOk = false
			}

		}

		h.state.IsWashed = false
	}()
}
