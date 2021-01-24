package dfpboard

import (
	"context"
	"sync"
	"time"

	"github.com/disaster37/gobot-fat/dfpconfig"
	"github.com/disaster37/gobot-fat/dfpstate"
	"github.com/disaster37/gobot-fat/models"
	log "github.com/sirupsen/logrus"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
)

// Wash  run on routine for no blocking.
// All routines are stopped when receive EventStopDFP,  EventBoardStop, EventSetSecurity or EventSetEmergencyStop  internal event
func (h *DFPBoard) wash() {
	// Only one watch
	defer h.Unlock()
	h.Lock()

	h.state.IsWashed = true
	if err := h.stateUsecase.Update(context.Background(), h.state); err != nil {
		log.Errorf("Error when save state in wash routine: %s", err.Error())
	}
	var wg sync.WaitGroup
	wg.Add(1)
	chFinishedStopEvent := make(chan bool, 1)
	chStoppedWash := make(chan bool, 1)

	ledControl := h.blinkGreenLed()

	handleError := func() {
		h.forceStopRelais()
		chFinishedStopEvent <- true
		ledControl.Stop()
		wg.Wait()
		h.state.IsWashed = false
		if err := h.stateUsecase.Update(context.Background(), h.state); err != nil {
			log.Errorf("Error when save state in wash routine: %s", err.Error())
		}
		ledControl.Wait()
		if h.state.IsRunning {
			h.turnOnGreenLed()
		} else {
			h.turnOffGreenLed()
		}
		return
	}

	// Check stop event
	go func() {
		out := h.Subscribe()
		for {
			select {
			case evt := <-out:
				if evt.Name == EventStopDFP || evt.Name == EventBoardStop || evt.Name == EventSetSecurity || evt.Name == EventSetEmergencyStop {

					// Not stop when security event and security is disabled
					if evt.Name == EventSetSecurity && h.state.IsDisableSecurity {
						break
					}
					h.forceStopRelais()
					chStoppedWash <- true
					ledControl.Stop()
					ledControl.Wait()
					h.Unsubscribe(out)
					wg.Done()
					return
				}
			case <-chFinishedStopEvent:
				h.Unsubscribe(out)
				wg.Done()
				return
			}
		}

	}()

	// Run wash
	go func() {

		var err error

		// Start pump and wait some time
		log.Debugf("Start pump and wait before continue %d s", h.config.StartWashingPumpBeforeWashing)
		timer := time.NewTimer(time.Duration(h.config.StartWashingPumpBeforeWashing) * time.Second)
		if err = h.relayPump.On(); err != nil {
			log.Errorf("When start pump: %s", err.Error())
			handleError()
			return

		}
		select {
		case <-chStoppedWash:
			h.forceStopRelais()
			wg.Wait()
			h.state.IsWashed = false
			if err = h.stateUsecase.Update(context.Background(), h.state); err != nil {
				log.Errorf("Error when save state in wash routine: %s", err.Error())
			}
			return
		case <-timer.C:
		}

		// Start drump
		log.Debugf("Start drum and wait before continue %d s", h.config.WashingDuration)
		timer = time.NewTimer(time.Duration(h.config.WashingDuration) * time.Second)
		if err = h.relayDrum.On(); err != nil {
			log.Errorf("When start drum: %s", err.Error())
			handleError()
			return
		}
		select {
		case <-chStoppedWash:
			h.forceStopRelais()
			wg.Wait()
			h.state.IsWashed = false
			if err = h.stateUsecase.Update(context.Background(), h.state); err != nil {
				log.Errorf("Error when save state in wash routine: %s", err.Error())
			}
			return

		case <-timer.C:
		}

		// Stop pump and drum
		log.Debugf("Stop pump and drump, washing finished")
		h.forceStopRelais()

		chFinishedStopEvent <- true
		ledControl.Stop()

		h.state.IsWashed = false
		h.state.LastWashing = time.Now()
		if err = h.stateUsecase.Update(context.Background(), h.state); err != nil {
			log.Errorf("Error when save state in wash routine: %s", err.Error())
		}

		// Reinit timer
		h.waitTimeForceWash = time.NewTicker(time.Duration(h.config.ForceWashingDuration) * time.Second)
		h.waitTimeForceWashFrozen = time.NewTicker(time.Duration(h.config.ForceWashingDurationWhenFrozen) * time.Second)

		ledControl.Wait()
		wg.Wait()

		h.Publish(EventWash, h.state)
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
	if err := h.eventUsecase.Create(ctx, event); err != nil {
		log.Errorf("Error when store new event: %s", err.Error())
	}
}

func (h *DFPBoard) work() {

	ctx := context.TODO()

	/*******
	 * Process external events
	 */
	// Handle config
	h.on(h.globalEventer, dfpconfig.NewDFPConfig, func(s interface{}) {
		//panic("plop2")

		dfpConfig := s.(*models.DFPConfig)
		log.Debugf("New config received for board %s, we update it", h.name)

		h.config = dfpConfig

		// Publish internal event
		h.Publish(EventNewConfig, dfpConfig)

	})

	// Handle state
	h.on(h.globalEventer, dfpstate.NewDFPState, func(s interface{}) {

		//panic("plop")
		dfpState := s.(*models.DFPState)
		log.Debugf("New state received for board %s, we update it", h.name)

		h.state.IsDisableSecurity = dfpState.IsDisableSecurity

		// Publish internal event
		h.Publish(EventNewState, dfpState)

	})

	/*******
	 * Process on button events
	 */
	// When button start

	h.on(h.buttonStart, gpio.ButtonPush, func(s interface{}) {
		log.Debug("Button start pushed")

		if err := h.StartDFP(ctx); err != nil {
			log.Errorf("When start DFP: %s", err.Error())
		}
	})

	// When button stop
	h.on(h.buttonStop, gpio.ButtonPush, func(s interface{}) {

		log.Debug("Button stop pushed")

		if err := h.StopDFP(ctx); err != nil {
			log.Errorf("When stop DFP: %s", err.Error())
		}

	})

	// When button wash
	h.on(h.buttonWash, gpio.ButtonPush, func(s interface{}) {

		log.Debug("Button wash pushed")

		if err := h.ForceWashing(ctx); err != nil {
			log.Errorf("When force washing: %s", err.Error())
		}
	})

	// Manual drum
	h.on(h.buttonForceDrum, gpio.ButtonPush, func(s interface{}) {
		// Start
		log.Debug("Button force drum pushed")

		if err := h.StartManualDrum(ctx); err != nil {
			log.Errorf("When start manual drum: %s", err.Error())
		}

		h.Publish(EventNewInput, "button_drum_pushed")

	})
	h.on(h.buttonForceDrum, gpio.ButtonRelease, func(s interface{}) {
		// Stop
		log.Debug("Button force drum released")

		if err := h.StopManualDrum(ctx); err != nil {
			log.Errorf("When stop manual drum: %s", err.Error())
		}

		h.Publish(EventNewInput, "button_drum_released")

	})

	// Manual pump
	h.on(h.buttonForcePump, gpio.ButtonPush, func(s interface{}) {
		// Start
		log.Debug("Button force pump pushed")

		if err := h.StartManualPump(ctx); err != nil {
			log.Errorf("When start manual pump: %s", err.Error())
		}

		h.Publish(EventNewInput, "button_pomp_pushed")

	})
	h.on(h.buttonForcePump, gpio.ButtonRelease, func(s interface{}) {
		// Stop
		log.Debug("Button force pump released")

		if err := h.StopManualPump(ctx); err != nil {
			log.Errorf("When stop manual pump: %s", err.Error())
		}

		h.Publish(EventNewInput, "button_pomp_released")

	})

	// When button emergency stop
	h.on(h.buttonEmergencyStop, gpio.ButtonPush, func(s interface{}) {
		log.Debug("Button emergency stop pushed")

		if err := h.SetEmergencyStop(ctx); err != nil {
			log.Errorf("When set emergency stop for DFP: %s", err.Error())
		}
	})
	h.on(h.buttonEmergencyStop, gpio.ButtonRelease, func(s interface{}) {
		log.Debug("Button emergency stop released")

		if err := h.UnsetEmergencyStop(ctx); err != nil {
			log.Errorf("When unset emergency stop for DFP: %s", err.Error())
		}
	})

	/*******
	 * Process on Captor event
	 */

	// When water captor ask wash
	wash := func(s interface{}) {
		log.Debug("Water captor pushed")

		select {
		case <-h.timeBetweenWash.C:
			// Timer finished
			if h.state.ShouldWash() {
				h.wash()
			}
			h.timeBetweenWash = time.NewTicker(time.Duration(h.config.WaitTimeBetweenWashing) * time.Second)
			break

		default:
			log.Debug("Wash not lauched because of need to wait some time before run again")
			break

		}
	}
	h.on(h.captorWaterUpper, gpio.ButtonPush, wash)
	h.on(h.captorWaterUnder, gpio.ButtonPush, wash)

	// When water captor ask security
	security := func(s interface{}) {

		if h.captorSecurityUpper.Active || h.captorSecurityUnder.Active {

			if err := h.SetSecurity(ctx); err != nil {
				log.Errorf("When set security for DFP: %s", err.Error())
			}
		} else {

			if err := h.UnsetSecurity(ctx); err != nil {
				log.Errorf("When unset security for DFP: %s", err.Error())
			}
		}

	}
	h.on(h.captorSecurityUpper, gpio.ButtonPush, security)
	h.on(h.captorSecurityUpper, gpio.ButtonRelease, security)
	h.on(h.captorSecurityUnder, gpio.ButtonPush, security)
	h.on(h.captorSecurityUnder, gpio.ButtonRelease, security)

	/*********
	 * Scheduling routines
	 */

	// Read temperature sensor
	ticker := gobot.Every(30*time.Minute, func() {
		//@TODO read sensors and update state.
		//No need to save state for that.
	})
	h.schedulingRoutines = append(h.schedulingRoutines, ticker)

	// Force washing when inactivity
	h.runWashInactivity()

	h.isInitialized = true

}

// runWashInactivity force wash if not running from ForceWashingDuration and from ForceWashingDurationWhenFrozen
// All routines are stopped when receive EventBoardStop internal event
func (h *DFPBoard) runWashInactivity() {

	chStop := make(chan bool)
	h.waitTimeForceWash = time.NewTicker(time.Duration(h.config.ForceWashingDuration) * time.Second)
	h.waitTimeForceWashFrozen = time.NewTicker(time.Duration(h.config.ForceWashingDurationWhenFrozen) * time.Second)

	// Check if stop event
	go func() {
		out := h.Subscribe()

		for {
			select {
			case evt := <-out:
				if evt.Name == EventBoardStop {
					chStop <- true
					h.Unsubscribe(out)
					return
				}
			}
		}
	}()

	// Force wash if needed
	go func() {
		for {
			select {
			case <-chStop:
				return
			case <-h.waitTimeForceWash.C:
				h.waitTimeForceWash = time.NewTicker(time.Duration(h.config.ForceWashingDuration) * time.Second)
				if int(h.state.AmbientTemperature) > h.config.TemperatureThresholdWhenFrozen {
					if h.state.ShouldWash() {
						h.wash()
					}
					break
				}
				continue
			case <-h.waitTimeForceWashFrozen.C:
				h.waitTimeForceWashFrozen = time.NewTicker(time.Duration(h.config.ForceWashingDurationWhenFrozen) * time.Second)
				if int(h.state.AmbientTemperature) <= h.config.TemperatureThresholdWhenFrozen {
					if h.state.ShouldWash() {
						h.wash()
					}
					break
				}
				continue
			}
		}
	}()

}

// Use on instead gobot.Eventer.On because of it not close routine at board is stopped.
// So, if you start / stop / start board, you have so many routine
func (h *DFPBoard) on(driver gobot.Eventer, event string, f func(data interface{})) {

	halt := make(chan bool)

	// Detect stop board
	go func() {
		out := h.Subscribe()

		for {
			select {
			case evt := <-out:
				if evt.Name == EventBoardStop {
					halt <- true
					h.Unsubscribe(out)
					return
				}
			}
		}
	}()

	// Handle on event
	go func() {
		out := driver.Subscribe()
		for {
			select {
			case <-halt:
				driver.Unsubscribe(out)
				return
			case evt := <-out:
				if evt.Name == event {
					f(evt.Data)
				}
			}
		}

	}()
}
