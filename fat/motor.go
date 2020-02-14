package fat

import (
	log "github.com/sirupsen/logrus"
)

// StartWashingPump permit to run washing pump
// The pump start only if no emergency and no security
func (h *FATHandler) StartWashingPump() (err error) {
	if !h.state.IsEmergencyStopped && !h.state.IsSecurity {
		log.Debug("Start whashing pump")
		return h.relayWashingPump.On()
	}

	log.Debug("Washing pump not started because of state not permit it")

}

// StopWashingPump permit to stop whashing pump
func (h *FATHandler) StopWashingPump() (err error) {
	log.Debug("Stop whashing pump")
	return h.relayWashingPump.Off()
}

// StartBarrelMotor permit to start barrel motor
// The motor start only if not emmergency and no security
func (h *FATHandler) StartBarrelMotor() (err error) {
	if !h.state.IsEmergencyStopped && !h.state.IsSecurity {
		log.Debug("Start barrel motor")
		return h.relayBarrelMotor.On()
	}
	log.Debug("Barrel motor not started because of state not permit it")

	return
}

// StopBarrelMotor permit to stop barrel motor
func (h *FATHandler) StopBarrelMotor() (err error) {
	log.Debug("Stop barrel motor")
	return h.relayBarrelMotor.Off()
}

// SwitchOnGreenLed display green led
func (h *FATHandler) SwitchOnGreenLed() (err error) {

	return

}

// SwitchOffGreenLed not display green led
func (h *FATHandler) SwitchOffGreenLed() (err error) {

	return

}

// SwitchOnRedLed display red led
func (h *FATHandler) SwitchOnRedLed() (err error) {

	return

}

// SwitchOffRedLed not display red led
func (h *FATHandler) SwitchOffRedLed() (err error) {
	return
}
