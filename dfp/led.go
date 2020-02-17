package dfp

import (
	"time"

	log "github.com/sirupsen/logrus"
)

// HandleRedLed manage the state of red LED
// If security is enabled, the led must be toggle
// If stop is enabled or emergency stop is enabled, the led must be switch on
// Else the LED must be switch off
func (h *DFPHandler) HandleRedLed() {

	// Set Mode security
	h.On(SecurityEvent, func(data interface{}) {
		err := h.ledRed.On()
		if err != nil {
			log.Errorf("Error appear when switch on red led: %s", err)
		}
		time.Sleep(1 * time.Second)

		err = h.ledRed.Off()
		if err != nil {
			log.Errorf("Error appear when switch off red led: %s", err)
		}
		time.Sleep(1 * time.Second)
	})

	// Unset mode security
	// Only if not on emergency state and not in stopped mode
	h.On(UnSecurityEvent, func(data interface{}) {
		err := h.ledRed.Off()
		if err != nil {
			log.Errorf("Error appear when remove toogle on red led: %s", err)
		}
	})

	// Set stop mode
	// Only if not on emergency state and not on security mode
	h.On(StopEvent, func(data interface{}) {
		err := h.ledRed.On()
		if err != nil {
			log.Errorf("Error appear when switch on red led: %s", err)
		}
	})

	// Unset stop mode
	// Only if not on emergency
	h.On(UnStopEvent, func(data interface{}) {
		if !h.state.IsSecurity {
			err := h.ledRed.Off()
			if err != nil {
				log.Errorf("Error appear when switch off red led: %s", err)
			}
		} else {
			err := h.ledRed.Toggle()
			if err != nil {
				log.Errorf("Error appear when toogle red led: %s", err)
			}
		}

	})

	// Set emergency stop mode
	h.On(EmergencyStopEvent, func(data interface{}) {
		err := h.ledRed.On()
		if err != nil {
			log.Errorf("Error appear when switch on red led: %s", err)
		}
	})

	// Unset emergency stop mode
	// Switch off led only if not on security or not in stopped mode
	h.On(UnEmergencyStopEvent, func(data interface{}) {

		if !h.state.IsStopped && !h.state.IsSecurity {
			err := h.ledRed.Off()
			if err != nil {
				log.Errorf("Error appear when switch off red led: %s", err)
			}
		} else if !h.state.IsStopped && h.state.IsSecurity {
			err := h.ledRed.Toggle()
			if err != nil {
				log.Errorf("Error appear when toogle red led: %s", err)
			}
		}
	})

}

// HandleGreenLed manage the state of green led
// If auto is enabled, so the green led is witch on
// If washing is enabled, so the green led is toogle
// Else, green led if switch off
func (h *DFPHandler) HandleGreenLed() {

	// Mode auto is set
	h.On(AutoEvent, func(data interface{}) {
		if !h.state.IsWashed {
			err := h.ledGreen.On()
			if err != nil {
				log.Errorf("Error appear when switch on green led: %s", err)
			}
		}
	})

	// Mode auto is unset
	h.On(UnAutoEvent, func(data interface{}) {
		if !h.state.IsWashed {
			err := h.ledGreen.Off()
			if err != nil {
				log.Errorf("Error appear when switch off green led: %s", err)
			}
		}
	})

	finishedGreenLedSwicthOnSwitchOff := false

	// Washing mode is set
	h.On(WashingEvent, func(data interface{}) {
		log.Info("Wash event fired")
		for h.state.IsWashed == true {
			finishedGreenLedSwicthOnSwitchOff = false
			err := h.ledGreen.Toggle()
			if err != nil {
				log.Errorf("Error appear when toogle green led: %s", err)
			}
			time.Sleep(1 * time.Second)
			finishedGreenLedSwicthOnSwitchOff = true
		}

	})

	// Washing mode is unset
	h.On(UnWashingEvent, func(data interface{}) {
		log.Info("Unwash event fired")

		// Wait switch on and switch off cycle is finished
		for finishedGreenLedSwicthOnSwitchOff != true {
			time.Sleep(1 * time.Millisecond)
		}
		if h.state.IsAuto {
			err := h.ledGreen.On()
			if err != nil {
				log.Errorf("Error appear when switch on green led")
			}
		} else {
			err := h.ledGreen.Off()
			if err != nil {
				log.Errorf("Error appear when switch of green led")
			}
		}
	})
}
