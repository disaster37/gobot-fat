package dfp

import (
	"time"

	log "github.com/sirupsen/logrus"
)

// StartWashingPump permit to run washing pump
// The pump start only if no emergency and no security
func (h *DFPHandler) StartWashingPump() (err error) {
	if h.state.CanStartMotor() {
		log.Debug("Start whashing pump")
		return h.relayWashingPump.On()
	}

	log.Debug("Washing pump not started because of state not permit it")

	return

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
func (h *DFPHandler) StartBarrelMotor() (err error) {
	if h.state.CanStartMotor() {
		log.Debug("Start barrel motor")
		return h.relayBarrelMotor.On()
	}
	log.Debug("Barrel motor not started because of state not permit it")

	return
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

// HandleWash permit to run one washing cycle
func (h *DFPHandler) HandleWash() {
	h.On(WashingEvent, func(data interface{}) {
		err := h.StartWashingPump()
		if err != nil {
			log.Errorf("Faild to start washing pump: %s", err)
		}
		time.Sleep(time.Second * time.Duration(h.config.GetInt("fat.washing.wait_time_washing_pump")))

		err = h.StartBarrelMotor()
		if err != nil {
			log.Errorf("Faild to start barrel motor: %s", err)
		}
		time.Sleep(time.Second * time.Duration(h.config.GetInt("fat.washing.duration")))

		h.StopWashingPump()
		h.StopBarrelMotor()

		if h.state.IsWashed {
			h.state.IsWashed = false
			h.Publish(UnWashingEvent, data)
		}

	})
}

// HandleStopMotor manage the stop motor events
func (h *DFPHandler) HandleStopMotor() {

	// Stop event
	h.On(StopEvent, func(data interface{}) {
		h.StopWashingPump()
		h.StopBarrelMotor()
	})

	// Security event
	h.On(SecurityEvent, func(data interface{}) {
		h.StopWashingPump()
		h.StopBarrelMotor()
	})

	// Emergency stop
	h.On(EmergencyStopEvent, func(data interface{}) {
		h.StopWashingPump()
		h.StopBarrelMotor()
	})
}
