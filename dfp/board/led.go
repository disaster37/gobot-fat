package dfpboard

import (
	"sync"
	"time"

	log "github.com/sirupsen/logrus"
)

// LedControl permit to control led when start blink
type LedControl struct {
	chStop chan bool
	wait   *sync.WaitGroup
}

// Stop permit to stop blink
func (l *LedControl) Stop() {
	l.chStop <- true
}

// Wait permit to wait blink is stopped
func (l *LedControl) Wait() {
	l.wait.Wait()
}

// newLedControl return new ledControl
func newLedControl() *LedControl {
	return &LedControl{
		chStop: make(chan bool, 0),
		wait:   &sync.WaitGroup{},
	}
}

// turnOnGreenLed turn on green led
func (h *DFPBoard) turnOnGreenLed() {
	if err := h.ledGreen.On(); err != nil {
		log.Errorf("Error when turn on GreenLed: %s", err.Error())
		return
	}

	log.Debug("Turn on GreenLed successfully")

}

// turnOffGreenLed turn off green led
func (h *DFPBoard) turnOffGreenLed() {
	if err := h.ledGreen.Off(); err != nil {
		log.Errorf("Error when turn off GreenLed: %s", err.Error())
		return
	}

	log.Debug("Turn off GreenLed successfully")

}

// blinkGreenLed blink green led
// When stop, it put initial led state
func (h *DFPBoard) blinkGreenLed() *LedControl {
	lc := newLedControl()
	lc.wait.Add(1)

	currentState := h.ledGreen.State()

	go func() {
		for {
			select {
			case <-lc.chStop:
				if currentState {
					h.turnOnGreenLed()
				} else {
					h.turnOffGreenLed()
				}
				lc.wait.Done()
				return
			default:
				if err := h.ledGreen.Toggle(); err != nil {
					log.Errorf("Error when toggle on GreenLed: %s", err.Error())
				}
				time.Sleep(1 * time.Second)
			}
		}
	}()

	return lc

}

// turnOnRedLed tunr on red led
func (h *DFPBoard) turnOnRedLed() {
	if err := h.ledRed.On(); err != nil {
		log.Errorf("Error when turn on RedLed: %s", err.Error())
		return
	}

	log.Debug("Turn on RedLed successfully")

}

// turnOffRedLed turn off red led
func (h *DFPBoard) turnOffRedLed() {
	if err := h.ledRed.Off(); err != nil {
		log.Errorf("Error when turn of RedLed: %s", err.Error())
		return
	}

	log.Debug("Turn off RedLed successfully")

}

// blinkRedLed blink red led
// When stop, it put initial led state
func (h *DFPBoard) blinkRedLed() *LedControl {
	lc := newLedControl()
	lc.wait.Add(1)

	currentState := h.ledRed.State()

	go func() {
		for {
			select {
			case <-lc.chStop:
				if currentState {
					h.turnOnRedLed()
				} else {
					h.turnOffRedLed()
				}
				lc.wait.Done()
				return
			default:
				if err := h.ledRed.Toggle(); err != nil {
					log.Errorf("Error when toggle on RedLed: %s", err.Error())
				}
				time.Sleep(1 * time.Second)
			}
		}
	}()

	return lc

}
