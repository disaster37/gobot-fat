package dfpboard

import (
	"context"
	"time"

	"github.com/disaster37/gobot-fat/dfpconfig"
	"github.com/disaster37/gobot-fat/models"
	log "github.com/sirupsen/logrus"
	"gobot.io/x/gobot/drivers/gpio"
)

// Wash  run on routine for no blocking.
func (h *DFPBoard) wash() {

	h.state.IsWashed = true
	if err := h.stateUsecase.Update(context.Background(), h.state); err != nil {
		log.Errorf("Error when save state in wash routine: %s", err.Error())
	}
	chFinished := make(chan bool)
	chStopped := make(chan bool)
	out := h.Subscribe()

	// Check stop event
	go func() {
		select {
		case evt := <-out:
			switch evt.Name {
			case Stop:
				h.forceStopRelais()
				chStopped <- true
				return
			}
		case <-chFinished:
			return
		}

	}()

	// Blink green led
	go func() {
		for {
			select {
			case <-chFinished:
				h.turnOnGreenLed()
				return
			case <-chStopped:
				if h.state.IsRunning {
					h.turnOnGreenLed()
				} else {
					h.turnOffGreenLed()
				}
				return
			default:
				h.ledGreen.Toggle()
				time.Sleep(1 * time.Second)
			}
		}
	}()

	// Run wash
	go func() {

		// Start pump and wait some time
		log.Debugf("Start pump and wait before continue %d s", h.config.StartWashingPumpBeforeWashing)
		timer := time.NewTimer(time.Duration(h.config.StartWashingPumpBeforeWashing) * time.Second)
		err := h.relayPump.On()
		if err != nil {
			log.Errorf("When start pump: %s", err.Error())
			h.forceStopRelais()
			chFinished <- true
			return
		}
		select {
		case <-chStopped:
			h.forceStopRelais()
			return
		case <-timer.C:
		}

		// Start drump
		log.Debugf("Start drum and wait before continue %d s", h.config.WashingDuration)
		timer = time.NewTimer(time.Duration(h.config.WashingDuration) * time.Second)
		err = h.relayDrum.On()
		if err != nil {
			log.Errorf("When start drum: %s", err.Error())
			h.forceStopRelais()
			chFinished <- true
			return
		}
		select {
		case <-chStopped:
			h.forceStopRelais()
			return

		case <-timer.C:
		}

		// Stop pump and drum
		log.Debugf("Stop pump and drump, washing finished")
		h.forceStopRelais()

		h.state.IsWashed = false
		if err := h.stateUsecase.Update(context.Background(), h.state); err != nil {
			log.Errorf("Error when save state in wash routine: %s", err.Error())
		}

		chFinished <- true
		h.Publish(NewWash, h.state)
		return
	}()

}

func (h *DFPBoard) sendEvent(ctx context.Context, kind string, name string, args ...interface{}) {
	event := &models.Event{
		SourceID:   h.state.Name,
		SourceName: h.state.Name,
		Timestamp:  time.Now(),
		EventType:  name,
		EventKind:  kind,
	}
	err := h.eventUsecase.Store(ctx, event)
	if err != nil {
		log.Errorf("Error when store new event: %s", err.Error())
	}
}

func (h *DFPBoard) work() {

	ctx := context.TODO()

	/*******
	 * Process external events
	 */
	// Handle config

	h.globalEventer.On(dfpconfig.NewDFPConfig, func(s interface{}) {
		dfpConfig := s.(*models.DFPConfig)
		log.Debugf("New config received for board %s, we update it", h.name)

		h.config = dfpConfig

		// Publish internal event
		h.Publish(NewConfig, dfpConfig)
	})

	/******
	 * Process internal events
	 */

	// Handle wash

	h.On(NewWash, func(s interface{}) {
		select {
		case <-h.timeBetweenWash.C:
			// Timer finished
			if h.state.ShouldWash() {
				h.wash()
			}

			h.timeBetweenWash = time.NewTicker(time.Duration(h.config.WaitTimeBetweenWashing) * time.Second)
		default:
			if log.IsLevelEnabled(log.DebugLevel) {
				log.Debug("Wash not lauched because of need to wait some time before run again")
			}
		}

	})

	/*******
	 * Process on button events
	 */
	// When button start

	h.buttonStart.On(gpio.ButtonPush, func(s interface{}) {
		if log.IsLevelEnabled(log.DebugLevel) {
			log.Debug("Button start pushed")
		}

		err := h.StartDFP(ctx)
		if err != nil {
			log.Errorf("When start DFP: %s", err.Error())
		}

		h.Publish(NewInput, "button_start_pushed")

	})

	// When button stop
	h.buttonStop.On(gpio.ButtonPush, func(s interface{}) {
		if log.IsLevelEnabled(log.DebugLevel) {
			log.Debug("Button stop pushed")
		}

		err := h.StopDFP(ctx)
		if err != nil {
			log.Errorf("When stop DFP: %s", err.Error())
		}

		h.Publish(NewInput, "button_stop_pushed")

	})

	// When button wash
	h.buttonWash.On(gpio.ButtonPush, func(s interface{}) {
		if log.IsLevelEnabled(log.DebugLevel) {
			log.Debug("Button wash pushed")
		}

		// Run force wash if not already wash, or is not on emergency stopped
		if log.IsLevelEnabled(log.DebugLevel) {

			err := h.ForceWashing(ctx)
			if err != nil {
				log.Errorf("When force washing: %s", err.Error())
			}
		}

		h.Publish(NewInput, "button_wash_pushed")

	})

	// Manual drum
	h.buttonForceDrum.On(gpio.ButtonPush, func(s interface{}) {
		// Start
		if log.IsLevelEnabled(log.DebugLevel) {
			log.Debug("Button force drum pushed")
		}

		err := h.StartManualDrum(ctx)
		if err != nil {
			log.Errorf("When start manual drum: %s", err.Error())
		}

		h.Publish(NewInput, "button_drum_pushed")

	})
	h.buttonForceDrum.On(gpio.ButtonRelease, func(s interface{}) {
		// Stop
		if log.IsLevelEnabled(log.DebugLevel) {
			log.Debug("Button force drum released")
		}

		err := h.StopManualDrum(ctx)
		if err != nil {
			log.Errorf("When stop manual drum: %s", err.Error())
		}

		h.Publish(NewInput, "button_drum_released")

	})

	// Manual pump
	h.buttonForcePump.On(gpio.ButtonPush, func(s interface{}) {
		// Start
		if log.IsLevelEnabled(log.DebugLevel) {
			log.Debug("Button force pump pushed")
		}

		err := h.StartManualPump(ctx)
		if err != nil {
			log.Errorf("When start manual pump: %s", err.Error())
		}

		h.Publish(NewInput, "button_pomp_pushed")

	})
	h.buttonForcePump.On(gpio.ButtonRelease, func(s interface{}) {
		// Stop
		if log.IsLevelEnabled(log.DebugLevel) {
			log.Debug("Button force pump released")
		}

		err := h.StopManualPump(ctx)
		if err != nil {
			log.Errorf("When stop manual pump: %s", err.Error())
		}

		h.Publish(NewInput, "button_pomp_released")

	})

	// When button emergency stop
	h.buttonEmergencyStop.On(gpio.ButtonPush, func(s interface{}) {
		if log.IsLevelEnabled(log.DebugLevel) {
			log.Debug("Button emergency stop pushed")
		}

		// Publish event to stop current wash
		h.Publish(Stop, nil)

		// Stop all relay
		h.forceStopRelais()

		// Publish even for external board
		h.globalEventer.Publish(EmergencyStopOn, nil)

		// Set red led
		h.turnOnRedLed()

		// Update state
		h.state.IsEmergencyStopped = true
		if err := h.stateUsecase.Update(ctx, h.state); err != nil {
			log.Errorf("Error when save state after emergency stop pushed: %s", err.Error())
		}

		h.Publish(NewInput, "button_emergency_stop_pushed")
	})
	h.buttonEmergencyStop.On(gpio.ButtonRelease, func(s interface{}) {
		if log.IsLevelEnabled(log.DebugLevel) {
			log.Debug("Button emergency stop released")
		}

		// Publish even for external board
		h.globalEventer.Publish(EmergencyStopOff, nil)

		// Turn off red label if DFP is running
		if h.state.IsRunning {
			h.turnOffRedLed()
		}

		// Update state
		h.state.IsEmergencyStopped = false
		if err := h.stateUsecase.Update(ctx, h.state); err != nil {
			log.Errorf("Error when save state after button emergency stop released: %s", err.Error())
		}

		h.Publish(NewInput, "button_emergency_stop_released")
	})

	/*******
	 * Process on Captor event
	 */

	// When water captor ask wash
	wash := func(s interface{}) {
		if log.IsLevelEnabled(log.DebugLevel) {
			log.Debug("Water captor pushed")
		}

		// Lauch event only if can wash
		if h.state.ShouldWash() {
			h.wash()
		}

		h.Publish(NewInput, "captor_water_pushed")
	}
	h.captorWaterUpper.On(gpio.ButtonPush, wash)
	h.captorWaterUnder.On(gpio.ButtonPush, wash)

	// When water captor ask security
	security := func(s interface{}) {
		if h.captorSecurityUpper.Active || h.captorSecurityUnder.Active {
			// Set security mode
			if !h.state.IsSecurity {
				log.Info("Set security mode")
				h.state.IsSecurity = true
				h.turnOnRedLed()
				h.forceStopRelais()
				h.Publish(NewSecurity, true)

				if err := h.stateUsecase.Update(ctx, h.state); err != nil {
					log.Errorf("Error when save state after detect security: %s", err.Error())
				}

				h.Publish(NewInput, "captor_security_pushed")
			}
		} else {
			// Unset security mode
			if h.state.IsSecurity {
				log.Info("Unset security mode")
				h.state.IsSecurity = false
				h.turnOffRedLed()
				h.Publish(NewSecurity, false)

				if err := h.stateUsecase.Update(ctx, h.state); err != nil {
					log.Errorf("Error when save state after unset security: %s", err.Error())
				}

				h.Publish(NewInput, "captor_security_released")
			}
		}
	}
	h.captorSecurityUpper.On(gpio.ButtonPush, security)
	h.captorSecurityUpper.On(gpio.ButtonRelease, security)
	h.captorSecurityUnder.On(gpio.ButtonPush, security)
	h.captorSecurityUnder.On(gpio.ButtonRelease, security)

	h.isInitialized = true

}
