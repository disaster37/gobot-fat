package dfp

import (
	"time"

	log "github.com/sirupsen/logrus"
)

// StartWashingPump permit to run washing pump
// The pump start only if no emergency and no security
func (h *DFPHandler) StartWashingPump() {
	if h.state.CanStartMotor() {
		log.Debug("Start whashing pump")
		err := h.relayWashingPump.On()
		if err != nil {
			log.Errorf("Error appear when try to start washing pump: %s", err)
		}
	} else {
		log.Debug("Washing pump not started because of state not permit it")
	}
}

// StopWashingPump permit to stop whashing pump
// It will try while not stopped
func (h *DFPHandler) StopWashingPump() {
	log.Debug("Stop whashing pump")

	isStopped := false
	for isStopped == false {
		err := h.relayWashingPump.Off()
		if err != nil {
			log.Errorf("Error when stop whashing pump: %s", err)
			time.Sleep(1 * time.Second)
		} else {
			isStopped = true
		}
	}

	log.Info("Stop whashing pump successfully")

}

// StartBarrelMotor permit to start barrel motor
// The motor start only if not emmergency and no security
func (h *DFPHandler) StartBarrelMotor() {
	if h.state.CanStartMotor() {
		log.Debug("Start barrel motor")
		err := h.relayBarrelMotor.On()
		if err != nil {
			log.Errorf("Error appear when try to start barrel motor: %s", err)
		}
	} else {
		log.Debug("Barrel motor not started because of state not permit it")
	}

}

// StopBarrelMotor permit to stop barrel motor
// It will try while is not stopped
func (h *DFPHandler) StopBarrelMotor() {
	log.Debug("Stop barrel motor")

	isStopped := false
	for isStopped == false {
		err := h.relayBarrelMotor.Off()
		if err != nil {
			log.Errorf("Error when stop barrel motor: %s", err)
			time.Sleep(1 * time.Second)
		} else {
			isStopped = true
		}
	}

	log.Info("Stop barrel motor successfully")
}

// HandleMotor manage the motor state
func (h *DFPHandler) HandleMotor() {

	//Handle washing
	h.eventer.On("stateChange", func(data interface{}) {
		event := data.(string)

		log.Debugf("Receive event %s", event)

		if h.state.ShouldWash() && h.state.CanWash() {
			log.Info("Start washing cycle")
			h.state.SetWashed()

			// Start washing pump and wait
			h.StartWashingPump()
			time.Sleep(time.Second * time.Duration(h.config.GetInt("dfp.washing.wait_time_washing_pump")))
			if !h.state.CanStartMotor() || !h.state.ShouldWash() {
				return
			}

			// Start barrel motor and wait
			h.StartBarrelMotor()
			time.Sleep(time.Second * time.Duration(h.config.GetInt("dfp.washing.duration")))
			if !h.state.CanStartMotor() || !h.state.ShouldWash() {
				return
			}

			h.StopWashingPump()
			h.StopBarrelMotor()

			if h.state.IsWashed() {
				h.state.UnsetShouldWash()
				h.state.UnsetWashed()
				h.state.UpdateLastWashing()
			}
		}
	})

	// Handle stop
	h.eventer.On("stateChange", func(data interface{}) {
		event := data.(string)

		log.Debugf("Receive event %s", event)

		// Stop motors
		if h.state.IsStopped() || h.state.IsEmergencyStopped() || h.state.IsSecurity() {
			h.StopMotors()
		}
	})
}

// StopMotors stop all motors
func (h *DFPHandler) StopMotors() {
	log.Info("Stop all motors")
	h.StopWashingPump()
	h.StopBarrelMotor()
}
