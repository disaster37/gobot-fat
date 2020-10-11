package dfpboard

import (
	"context"
	"time"

	"github.com/disaster37/gobot-fat/models"
	log "github.com/sirupsen/logrus"
	"gobot.io/x/gobot/drivers/gpio"
)

// Wash it's run on routine for no blocking.
func (h *DFPBoard) wash() {

	h.state.IsWashed = true
	h.Publish("state", h.state)
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
			case "stop":
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
			case "stop":
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
		h.Publish("state", h.state)
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

	/****************
	 * Init state
	 */

	// If run normally
	if h.state.IsRunning && !h.state.IsSecurity && !h.state.IsEmergencyStopped {
		if log.IsLevelEnabled(log.DebugLevel) {
			log.Debug("DFP run...")
		}
		h.turnOnGreenLed()
		h.turnOffRedLed()
	} else {
		// It stopped or security
		if log.IsLevelEnabled(log.DebugLevel) {
			log.Debug("DFP stopped or in security")
		}
		h.forceStopRelais()
		h.turnOffGreenLed()
		h.turnOnRedLed()
	}

	// If on current wash
	if h.state.IsWashed {
		h.wash()
	}

	/*******
	 * Routines on backgroup
	 */
	// Update state
	go func() {
		out := h.Subscribe()
		for {
			select {
			case evt := <-out:
				switch evt.Name {
				case "state":
					err := h.stateUsecase.Update(ctx, evt.Data.(*models.DFPState))
					if err != nil {
						log.Errorf("Error when update DFP state: %s", err.Error())
					}
				case "stop":
					return
				}
			}
		}
	}()

	// Load config
	go func() {
		out := h.Subscribe()
		duration := 1 * time.Minute
		timer := time.NewTicker(duration)
		for {
			select {
			case evt := <-out:
				switch evt.Name {
				case "stop":
					return
				}
			case <-timer.C:
				timer = time.NewTicker(duration)
				config, err := h.configUsecase.Get(ctx)
				if err != nil {
					log.Errorf("Error when load DFP config: %s", err.Error())
					continue
				}

				h.config = config
			}
		}
	}()

	// Handle security captor
	go func() {
		out := h.Subscribe()
		for {
			select {
			case evt := <-out:
				switch evt.Name {
				case "stop":
					return
				}
			default:
				if h.captorSecurityUpper.Active || h.captorSecurityUnder.Active {
					// Set security mode
					if !h.state.IsSecurity {
						log.Info("Set security mode")
						h.state.IsSecurity = true
						h.turnOnRedLed()
						h.forceStopRelais()
						h.Publish("security", true)
						h.Publish("state", h.state)
					}
				} else {
					// Unset security mode
					if h.state.IsSecurity {
						log.Info("Unset security mode")
						h.state.IsSecurity = false
						h.turnOffRedLed()
						h.Publish("security", false)
						h.Publish("state", h.state)
					}
				}
			}

			time.Sleep(1 * time.Millisecond)
		}
	}()

	// Auto wash events
	go func() {
		timer := time.NewTicker(time.Duration(h.config.WaitTimeBetweenWashing) * time.Second)
		out := h.Subscribe()
		for {
			select {
			case evt := <-out:
				switch evt.Name {
				case "stop":
					return
				case "wash":
					select {
					case <-timer.C:
						// Timer finished
						if h.state.ShouldWash() {
							h.wash()
						}

						timer = time.NewTicker(time.Duration(h.config.WaitTimeBetweenWashing) * time.Second)
					default:
						if log.IsLevelEnabled(log.DebugLevel) {
							log.Debug("Wash not lauched because of need to wait some time before run again")
						}
					}
				}
			}
		}
	}()

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
	h.buttonSet.On(gpio.ButtonPush, func(s interface{}) {
		if log.IsLevelEnabled(log.DebugLevel) {
			log.Debug("Button set pushed")
		}
	})

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

}
