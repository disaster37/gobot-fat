package dfp

import (
	log "github.com/sirupsen/logrus"
	"gobot.io/x/gobot/drivers/gpio"
)

// HandleButtonWash manage wash button
func (h *DFPHandler) HandleButtonWash() {
	h.buttonWash.On(gpio.ButtonPush, func(data interface{}) {
		log.Infof("Button wash is pushed")
		if h.state.CanWash() {
			h.state.IsWashed = true
			h.Publish(WashingEvent, data)
		}
	})

	h.buttonWash.On(gpio.ButtonRelease, func(data interface{}) {
		log.Infof("Button wash is released")
	})
}

// HandleButtonAuto manage auto button
func (h *DFPHandler) HandleButtonAuto() {
	h.buttonAuto.On(gpio.ButtonPush, func(data interface{}) {
		log.Infof("Button auto is pushed")
		h.state.IsAuto = true
		h.Publish(AutoEvent, data)
	})

	h.buttonAuto.On(gpio.ButtonRelease, func(data interface{}) {
		log.Infof("Button auto is released")
		h.state.IsAuto = false
		h.Publish(UnAutoEvent, data)
	})
}

// HandleButtonStop manage stop button
func (h *DFPHandler) HandleButtonStop() {
	h.buttonStop.On(gpio.ButtonPush, func(data interface{}) {
		log.Infof("Button stop is pushed")
		if h.state.CanSetStop() {
			h.state.IsWashed = false
			h.Publish(StopEvent, data)
		}
		h.state.IsStopped = true
	})

	h.buttonStop.On(gpio.ButtonRelease, func(data interface{}) {
		log.Infof("Button stop is released")
		if h.state.CanUnsetStop() {
			h.Publish(UnStopEvent, data)
		}
		h.state.IsStopped = false
	})
}

// HandleButtonEmergencyStop manage emergency stop button
func (h *DFPHandler) HandleButtonEmergencyStop() {
	h.buttonEmergencyStop.On(gpio.ButtonPush, func(data interface{}) {
		log.Infof("Button emergency stop is pushed")
		if h.state.CanSetEmergencyStop() {
			h.state.IsWashed = false
			h.Publish(EmergencyStopEvent, data)
		}
		h.state.IsEmergencyStopped = true
	})

	h.buttonEmergencyStop.On(gpio.ButtonRelease, func(data interface{}) {
		log.Infof("Button emergency stop is released")
		if h.state.CanUnsetEmergencyStop() {
			h.Publish(UnEmergencyStopEvent, data)
		}
		h.state.IsEmergencyStopped = false
	})
}

// HandleButtonForceMotor manage the button that permit to force to start motors
func (h *DFPHandler) HandleButtonForceMotor() {

	// Force washing pump
	h.buttonForceWashingPump.On(gpio.ButtonPush, func(data interface{}) {
		log.Infof("Button force washing pump is pushed")
		if h.state.CanStartMotor() {
			h.StartWashingPump()
		}
	})

	h.buttonForceWashingPump.On(gpio.ButtonRelease, func(data interface{}) {
		log.Infof("Button force washing pump is released")
		h.StopWashingPump()
	})

	// Force barrel motor
	h.buttonForceBarrelMotor.On(gpio.ButtonPush, func(data interface{}) {
		log.Infof("Button force barrel motor is pushed")
		if h.state.CanStartMotor() {
			h.StartBarrelMotor()
		}
	})

	h.buttonForceBarrelMotor.On(gpio.ButtonRelease, func(data interface{}) {
		log.Infof("Button force barrel motor is released")
		h.StopBarrelMotor()
	})
}
