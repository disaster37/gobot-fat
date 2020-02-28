package dfp

import (
	log "github.com/sirupsen/logrus"
	"gobot.io/x/gobot/drivers/gpio"
)

// HandleButtonWash manage wash button
func (h *DFPHandler) HandleButtonWash() {
	h.buttonWash.On(gpio.ButtonPush, func(data interface{}) {
		log.Debugf("Button wash is pushed")
		h.state.SetShouldWash()
	})

	h.buttonWash.On(gpio.ButtonRelease, func(data interface{}) {
		log.Debugf("Button wash is released")
	})
}

// HandleButtonAuto manage auto button
func (h *DFPHandler) HandleButtonAuto() {
	h.buttonAuto.On(gpio.ButtonPush, func(data interface{}) {
		log.Debugf("Button auto is pushed")
		h.state.SetAuto()
	})

	h.buttonAuto.On(gpio.ButtonRelease, func(data interface{}) {
		log.Debugf("Button auto is released")
		h.state.UnsetAuto()
	})
}

// HandleButtonStop manage stop button
func (h *DFPHandler) HandleButtonStop() {
	h.buttonStop.On(gpio.ButtonPush, func(data interface{}) {
		log.Debugf("Button stop is pushed")
		h.state.SetStop()
	})

	h.buttonStop.On(gpio.ButtonRelease, func(data interface{}) {
		log.Debugf("Button stop is released")
		h.state.UnsetStop()
	})
}

// HandleButtonEmergencyStop manage emergency stop button
func (h *DFPHandler) HandleButtonEmergencyStop() {
	h.buttonEmergencyStop.On(gpio.ButtonPush, func(data interface{}) {
		log.Debugf("Button emergency stop is pushed")
		h.state.SetEmergencyStop()
	})

	h.buttonEmergencyStop.On(gpio.ButtonRelease, func(data interface{}) {
		log.Debugf("Button emergency stop is released")
		h.state.UnsetEmergencyStop()
	})
}

// HandleButtonForceMotor manage the button that permit to force to start motors
func (h *DFPHandler) HandleButtonForceMotor() {

	// Force washing pump
	h.buttonForceWashingPump.On(gpio.ButtonPush, func(data interface{}) {
		log.Debugf("Button force washing pump is pushed")
		h.StartWashingPump()
	})

	h.buttonForceWashingPump.On(gpio.ButtonRelease, func(data interface{}) {
		log.Debugf("Button force washing pump is released")
		h.StopWashingPump()
	})

	// Force barrel motor
	h.buttonForceBarrelMotor.On(gpio.ButtonPush, func(data interface{}) {
		log.Debugf("Button force barrel motor is pushed")
		h.StartBarrelMotor()
	})

	h.buttonForceBarrelMotor.On(gpio.ButtonRelease, func(data interface{}) {
		log.Debugf("Button force barrel motor is released")
		h.StopBarrelMotor()
	})
}
