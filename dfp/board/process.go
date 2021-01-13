package dfpboard

import (
	"context"
	"time"

	"github.com/disaster37/gobot-fat/dfpconfig"
	"github.com/disaster37/gobot-fat/models"
	log "github.com/sirupsen/logrus"
	"gobot.io/x/gobot/drivers/gpio"
)

// Wash it's run on routine for no blocking.
func (h *DFPBoard) wash() {

	h.state.IsWashed = true
	h.Publish(NewState, h.state)
	chFinished := make(chan bool)

	// Blink green led
	go func() {
		for {
			select {
			case <-chFinished:
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

		out := h.Subscribe()

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
		case evt := <-out:
			switch evt.Name {
			case Stop:
				h.forceStopRelais()
				chFinished <- true
				return
			}
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
		case evt := <-out:
			switch evt.Name {
			case Stop:
				h.forceStopRelais()
				chFinished <- true
				return
			}
		case <-timer.C:

		}

		// Stop pump and drum
		log.Debugf("Stop pump and drump, washing finished")
		h.forceStopRelais()
		chFinished <- true

		h.state.IsWashed = false
		h.Publish(NewState, h.state)
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
	// Handle state change

	h.On(NewState, func(s interface{}) {
		err := h.stateUsecase.Update(ctx, s.(*models.DFPState))
		if err != nil {
			log.Errorf("Error when update DFP state: %s", err.Error())
		}
	})

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

	})

	// When button set
	h.buttonEmergencyStop.On(gpio.ButtonPush, func(s interface{}) {
		if log.IsLevelEnabled(log.DebugLevel) {
			log.Debug("Button emergency stop pushed")
		}

		//@TODO Emergency func
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
			h.Publish("wash", true)
		}
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
				h.Publish(NewState, h.state)
			}
		} else {
			// Unset security mode
			if h.state.IsSecurity {
				log.Info("Unset security mode")
				h.state.IsSecurity = false
				h.turnOffRedLed()
				h.Publish(NewSecurity, false)
				h.Publish(NewState, h.state)
			}
		}
	}
	h.captorSecurityUpper.On(gpio.ButtonPush, security)
	h.captorSecurityUpper.On(gpio.ButtonRelease, security)
	h.captorSecurityUnder.On(gpio.ButtonPush, security)
	h.captorSecurityUnder.On(gpio.ButtonRelease, security)

	h.isInitialized = true

}
