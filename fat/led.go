package pbf

import "gobot.io/x/gobot"

// HandleRedLed manage the state of red LED
// If security is enabled, the led must be toogle
// If stop is enabled or emergency stop is enabled, the led must be switch on
// Else the LED must be switch off
func (h *FATHandler) HandleRedLed() {

	gobot.Every(time.Second * 1, funcfunc() {
		
		// When PBF is stopped or in emergency state
		if h.state.IsStopped || h.state.IsIsEmergencyStopped {

		} else if h.state.IsSecurity {
			//When PBF is on security
		} else {
			// When all work as expected
		}
	})
	

}
