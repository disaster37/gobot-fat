package dfpboard

import (
	"context"
	"time"

	"github.com/disaster37/gobot-fat/models"
	log "github.com/sirupsen/logrus"
)

// sendEvent permit to send event on Elasticsearch
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

// Wash permit to run wash
// It consider that wash can started without more control. You need to check before call it.
func (h *DFPHandler) wash() {

	h.state.IsWashed = true
	err := h.stateUsecase.Update(context.Background(), h.state)
	if err != nil {
		log.Errorf("Error when save dfp state: %s", err.Error())
		return
	}

	// init all time to stop the if needed
	timerPump := time.NewTimer(time.Second)
	timerWashing := time.NewTimer(time.Second)
	timerPump.Stop()
	timerWashing.Stop()

	// Blink green led
	timerLed := h.ledGreen.Blink(time.Duration(h.config.StartWashingPumpBeforeWashing+h.config.WashingDuration) * time.Second)

	// run whatchdog to stop washing if needed
	go func() {
		for {
			if h.state.IsWashed != true {
				log.Debugf("Detect aborded wash")
				if !timerPump.Stop() {
					<-timerPump.C
				}
				if !timerWashing.Stop() {
					<-timerWashing.C
				}
				if !timerLed.Stop() {
					<-timerLed.C
				}

				return
			}
			time.Sleep(time.Millisecond * 1)
		}
	}()

	// Start pump and wait some time
	timerPump.Reset(time.Duration(h.config.StartWashingPumpBeforeWashing) * time.Second)
	err = h.relayPump.On()
	if err != nil {
		log.Errorf("Error when start pump: %s", err.Error())
	}
	<-timerPump.C

	// Start drum and wait some time
	timerWashing.Reset(time.Duration(h.config.WashingDuration) * time.Second)
	err = h.relayDrum.On()
	if err != nil {
		log.Errorf("Error when start drum")
	}
	<-timerWashing.C

	// Stop all
	err = h.relayDrum.Off()
	if err != nil {
		h.forceStopMotors()
		log.Errorf("Error when stop drum: %s", err.Error())
	}
	err = h.relayPump.Off()
	if err != nil {
		h.forceStopMotors()
		log.Errorf("Error when stop pump: %s", err.Error())
	}

	// Update state and save it
	h.state.IsWashed = false
	h.state.LastWashing = time.Now()
	err = h.stateUsecase.Update(context.Background(), h.state)
	if err != nil {
		log.Errorf("Error when save dfp state: %s", err.Error())
	}

	// Send event for stats
	h.sendEvent("washing", "motor")
	log.Debugf("Washing successfully finished")
	if !timerLed.Stop() {
		<-timerLed.C
	}

	time.Sleep(time.Second * 5)
	if !h.state.IsRunning || h.state.Security() {
		err := h.ledGreen.TurnOff()
		if err != nil {
			log.Errorf("Error when turn off green led: %s", err.Error())
		}
	}

}

// forceStopMotors run on routine while stop motor not completed and event if get errors.
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

// turnOnLight turn on led on buttons and on LCD
func (h *DFPHandler) turnOnLight() {

	// Buttons led
	for i, led := range h.ledButtons {
		err := led.TurnOn()
		if err != nil {
			log.Errorf("Error appear when turn on button led %d: %s", i, err.Error())
		}
	}
}

// turnOffLight turn off led on buttons and on LCD
func (h *DFPHandler) turnOffLight() {

	// Buttons led
	for i, led := range h.ledButtons {
		err := led.TurnOff()
		if err != nil {
			log.Errorf("Error appear when turn off button led %d: %s", i, err.Error())
		}
	}
}

func (h *DFPHandler) stopDFP() error {

	h.state.IsWashed = false

	err := h.relayDrum.Off()
	if err != nil {
		return err
	}
	err = h.relayPump.Off()
	if err != nil {
		return err
	}

	err = h.ledGreen.TurnOff()
	if err != nil {
		return err
	}
	err = h.ledRed.TurnOn()
	if err != nil {
		return err
	}

	return nil

}
