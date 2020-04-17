package dfpgobot

import (
	"time"

	log "github.com/sirupsen/logrus"
)

// HandleRedLed manage the state of red LED
// If security is enabled, the led must be toggle
// If stop is enabled or emergency stop is enabled, the led must be switch on
// Else the LED must be switch off
func (h *DFPHandler) HandleRedLed() {

	// Compute LED state
	h.eventer.On("stateChange", func(data interface{}) {

		// Turn on LED when emergency or stop mode
		if h.stateRepository.State().IsStopped || h.stateRepository.State().IsEmergencyStopped {
			err := h.ledRed.On()
			if err != nil {
				log.Errorf("Error appear when turn on red led: %s", err)
			}
		} else if h.stateRepository.State().IsSecurity {
			// Blink led when security mode
			for h.stateRepository.State().IsSecurity {
				err := h.ledRed.Toggle()
				if err != nil {
					log.Errorf("Error appear when toggle red led: %s", err)
				}
				time.Sleep(1 * time.Second)
			}
		} else {
			// Turn of LED
			err := h.ledRed.Off()
			if err != nil {
				log.Errorf("Error appear when turn off red led: %s", err)
			}
		}
	})
}

// HandleGreenLed manage the state of green led
// If auto is enabled, so the green led is witch on
// If washing is enabled, so the green led is toogle
// Else, green led if switch off
func (h *DFPHandler) HandleGreenLed() {

	// Compute LED state
	h.eventer.On("stateChange", func(data interface{}) {

		// Blink LED when washing
		if h.stateRepository.State().IsWashed {
			for h.stateRepository.State().IsWashed {
				err := h.ledGreen.Toggle()
				if err != nil {
					log.Errorf("Error appear when toggle green led: %s", err)
				}
				time.Sleep(1 * time.Second)
			}
		} else if h.stateRepository.State().IsAuto {
			// Turn on LED when auto
			err := h.ledGreen.On()
			if err != nil {
				log.Errorf("Error appear when turn on green led: %s", err)
			}
		} else {
			// Turn off LED when auto is disabled
			err := h.ledGreen.Off()
			if err != nil {
				log.Errorf("Error appear when turn off green led: %s", err)
			}
		}
	})
}
