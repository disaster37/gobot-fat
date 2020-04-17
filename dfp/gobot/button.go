package dfpgobot

import (
	log "github.com/sirupsen/logrus"
	"gobot.io/x/gobot/drivers/gpio"
)

// HandleButtonWash manage wash button
func (h *DFPHandler) HandleButtonWash() {
	h.buttonWash.On(gpio.ButtonPush, func(data interface{}) {
		log.Debugf("Button wash is pushed")
		h.stateRepository.SetShouldWash()
	})

	h.buttonWash.On(gpio.ButtonRelease, func(data interface{}) {
		log.Debugf("Button wash is released")
	})
}

// HandleButtonAuto manage auto button
func (h *DFPHandler) HandleButtonAuto() {
	h.buttonAuto.On(gpio.ButtonPush, func(data interface{}) {
		log.Debugf("Button auto is pushed")
		if h.stateRepository.State().IsAuto {
			log.Infof("Unset auto mode")
			h.stateRepository.UnsetAuto()
		} else {
			log.Infof("Set auto mode")
			h.stateRepository.SetAuto()
		}
	})

	h.buttonAuto.On(gpio.ButtonRelease, func(data interface{}) {
		log.Debugf("Button auto is released")
	})
}

// HandleButtonStop manage stop button
func (h *DFPHandler) HandleButtonStop() {
	h.buttonStop.On(gpio.ButtonPush, func(data interface{}) {
		log.Debugf("Button stop is pushed")
		if h.stateRepository.State().IsStopped {
			log.Infof("Unset stop mode")
			h.stateRepository.UnsetStop()
		} else {
			log.Infof("Set stop mode")
			h.stateRepository.SetStop()
		}
	})

	h.buttonStop.On(gpio.ButtonRelease, func(data interface{}) {
		log.Debugf("Button stop is released")
	})
}

// HandleButtonEmergencyStop manage emergency stop button
func (h *DFPHandler) HandleButtonEmergencyStop() {
	h.buttonEmergencyStop.On(gpio.ButtonPush, func(data interface{}) {
		log.Debugf("Button emergency stop is pushed")
		if h.stateRepository.State().IsEmergencyStopped {
			log.Infof("Unset emergency stopped mode")
			h.stateRepository.UnsetEmergencyStop()
		} else {
			log.Infof("Set emergency stopped mode")
			h.stateRepository.SetEmergencyStop()
		}
	})

	h.buttonEmergencyStop.On(gpio.ButtonRelease, func(data interface{}) {
		log.Debugf("Button emergency stop is released")
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
