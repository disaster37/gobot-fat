package dfpboard

import (
	"context"
	"time"

	"github.com/disaster37/gobot-fat/board"
	"github.com/disaster37/gobot-fat/models"
	log "github.com/sirupsen/logrus"
)

// sendEvent permit to send event on Elasticsearch
func (h *DFPHandler) sendEvent(ctx context.Context, eventType string, eventKind string) {

	event := &models.Event{
		SourceID:   h.state.Name,
		SourceName: h.state.Name,
		Timestamp:  time.Now(),
		EventType:  eventType,
		EventKind:  eventKind,
	}
	err := h.eventUsecase.Store(ctx, event)
	if err != nil {
		log.Errorf("Error when store new event: %s", err.Error())
	}
}

// Wash permit to run wash
// It consider that wash can started without more control. You need to check before call it.
func (h *DFPHandler) wash(mainCtx context.Context, chCancel chan bool) (res interface{}, err error) {

	ctx, cancel := context.WithCancel(mainCtx)
	defer cancel()

	h.state.IsWashed = true
	err = h.stateUsecase.Update(ctx, h.state)
	if err != nil {
		return nil, err
	}

	// Blink green led
	chCancelLed, chErrLed := h.ledGreen.BlinkCanceleable(ctx)
	go func() {
		select {
		case err := <-chErrLed:
			if err != nil {
				log.Errorf("Error when cancel blink green led: %s", err.Error())
			}
		}
	}()

	// run whatchdog to stop washing if needed
	go func() {
		for {
			if h.state.IsWashed != true {
				log.Debugf("Detect aborded wash")

				// Cancel routine and timers
				cancel()
				chCancel <- true
				chCancelLed <- true
			}
			time.Sleep(1 * time.Nanosecond)
		}
	}()

	// Start pump and wait some time
	log.Debugf("Start pump and wait before continue %d s", h.config.StartWashingPumpBeforeWashing)
	timerPump := time.NewTimer(time.Duration(h.config.StartWashingPumpBeforeWashing) * time.Second)
	err = h.relayPump.On(ctx)
	if err != nil {
		routine := board.NewRoutine(mainCtx, h.forceStopMotors)
		select {
		case err := <-routine.Error():
			return nil, err
		case <-routine.Result():
			log.Debug("Motors force stopped")
			return nil, err
		}
	}
	select {
	case <-chCancel:
		err = h.relayPump.Off(mainCtx)
		if err != nil {
			routine := board.NewRoutine(mainCtx, h.forceStopMotors)
			select {
			case err := <-routine.Error():
				return nil, err
			case <-routine.Result():
				log.Debug("Motors force stopped")
				return nil, nil
			}
		}
		return nil, nil
	case <-timerPump.C:
		break
	}

	// Start drum and wait some time
	log.Debugf("Start drum and wait %d s", h.config.WashingDuration)
	timerWashing := time.NewTimer(time.Duration(h.config.WashingDuration) * time.Second)
	err = h.relayDrum.On(ctx)
	if err != nil {

		routine := board.NewRoutine(mainCtx, h.forceStopMotors)
		select {
		case err := <-routine.Error():
			return nil, err
		case <-routine.Result():
			log.Debug("Motors force stopped")
			return nil, err
		}

	}
	select {
	case <-chCancel:
		err = h.relayPump.Off(mainCtx)
		if err != nil {
			routine := board.NewRoutine(mainCtx, h.forceStopMotors)
			select {
			case err := <-routine.Error():
				return nil, err
			case <-routine.Result():
				log.Debug("Motors force stopped")
				return nil, err
			}
		}
		err = h.relayDrum.Off(mainCtx)
		if err != nil {
			routine := board.NewRoutine(mainCtx, h.forceStopMotors)
			select {
			case err := <-routine.Error():
				return nil, err
			case <-routine.Result():
				log.Debug("Motors force stopped")
				return nil, err
			}
		}
		return nil, nil
	case <-timerWashing.C:
		break
	}

	// Stop all
	err = h.relayDrum.Off(ctx)
	if err != nil {
		routine := board.NewRoutine(mainCtx, h.forceStopMotors)
		select {
		case err := <-routine.Error():
			return nil, err
		case <-routine.Result():
			log.Debug("Motors force stopped")
			return nil, err
		}
	}
	err = h.relayPump.Off(ctx)
	if err != nil {
		routine := board.NewRoutine(mainCtx, h.forceStopMotors)
		select {
		case err := <-routine.Error():
			return nil, err
		case <-routine.Result():
			log.Debug("Motors force stopped")
			return nil, err
		}

	}

	// Update state and save it
	h.state.IsWashed = false
	h.state.LastWashing = time.Now()
	err = h.stateUsecase.Update(ctx, h.state)
	if err != nil {
		return nil, err
	}

	// Send event for stats
	h.sendEvent(ctx, "washing", "motor")
	log.Debugf("Washing successfully finished")

	time.Sleep(time.Second * 5)
	if !h.state.IsRunning || h.state.Security() {
		err := h.ledGreen.TurnOff(ctx)
		if err != nil {
			return nil, err
		}
	}

	return nil, nil
}

// forceStopMotors run on routine while stop motor not completed and event if get errors.
func (h *DFPHandler) forceStopMotors(ctx context.Context, chCancel chan bool) (res interface{}, err error) {
	isOk := false
	for !isOk {
		select {
		case <-chCancel:
			return nil, nil
		default:
			isOk = true
			err := h.relayDrum.Off(ctx)
			if err != nil {
				log.Errorf("Error appear when try to stop drum: %s", err.Error())
				isOk = false
			}

			err = h.relayPump.Off(ctx)
			if err != nil {
				log.Errorf("Error appear when try to stop pump: %s", err.Error())
				isOk = false
			}
		}

	}

	return nil, nil
}

// turnOnLight turn on led on buttons and on LCD
func (h *DFPHandler) turnOnLight(ctx context.Context) {

	// Buttons led
	for i, led := range h.ledButtons {
		err := led.TurnOn(ctx)
		if err != nil {
			log.Errorf("Error appear when turn on button led %d: %s", i, err.Error())
		}
	}
}

// turnOffLight turn off led on buttons and on LCD
func (h *DFPHandler) turnOffLight(ctx context.Context) {

	// Buttons led
	for i, led := range h.ledButtons {
		err := led.TurnOff(ctx)
		if err != nil {
			log.Errorf("Error appear when turn off button led %d: %s", i, err.Error())
		}
	}
}

func (h *DFPHandler) stopDFP(ctx context.Context, chCancel chan bool) (res interface{}, err error) {

	h.state.IsWashed = false
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// Handel cancel
	go func() {
		select {
		case <-chCancel:
			cancel()
		}
	}()

	err = h.relayDrum.Off(ctx)
	if err != nil {
		return nil, err
	}
	err = h.relayPump.Off(ctx)
	if err != nil {
		return nil, err
	}

	err = h.ledGreen.TurnOff(ctx)
	if err != nil {
		return nil, err
	}
	err = h.ledRed.TurnOn(ctx)
	if err != nil {
		return nil, err
	}

	// Close routine
	chCancel <- true

	return nil, nil

}

func (h *DFPHandler) init() error {

}
