package pbf

import (
	"time"

	log "github.com/sirupsen/logrus"
	"gobot.io/x/gobot"
)

// StartWashingPump permit to run washing pump
// The pump start only if no emergency and no security
func (h *FATHandler) StartWashingPump() (err error) {
	if !h.state.IsEmergencyStopped && !h.state.IsSecurity {
		log.Debug("Start whashing pump")
		return h.relayWashingPump.On()
	}

	log.Debug("Washing pump not started because of state not permit it")

	return

}

// StopWashingPump permit to stop whashing pump
// It will try while not stopped
func (h *FATHandler) StopWashingPump() {
	log.Debug("Stop whashing pump")

	gobot.After(0, func() {
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

	})

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
// It will try while is not stopped
func (h *FATHandler) StopBarrelMotor() {
	log.Debug("Stop barrel motor")

	gobot.After(0, func() {
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

	})
}
