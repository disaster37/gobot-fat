package dfpgobot

import (
	"context"
	"time"

	"github.com/disaster37/gobot-fat/models"
	log "github.com/sirupsen/logrus"
)

// StartWashingPump permit to run washing pump
// The pump start only if no emergency and no security
func (h *DFPHandler) StartWashingPump() {
	if h.stateRepository.CanStartMotor() {
		log.Debug("Start whashing pump")
		err := h.relayWashingPump.On()
		if err != nil {
			log.Errorf("Error appear when try to start washing pump: %s", err)
		}

		event := &models.Event{
			SourceID:   h.stateRepository.State().ID,
			SourceName: h.stateRepository.State().Name,
			Timestamp:  time.Now(),
			EventType:  "start_washing_pump",
			EventKind:  "motor",
		}
		err = h.eventUsecase.Store(context.Background(), event)
		if err != nil {
			log.Errorf("Error when store new event: %s", err.Error())
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

	event := &models.Event{
		SourceID:   h.stateRepository.State().ID,
		SourceName: h.stateRepository.State().Name,
		Timestamp:  time.Now(),
		EventType:  "stop_washing_pump",
		EventKind:  "motor",
	}
	err := h.eventUsecase.Store(context.Background(), event)
	if err != nil {
		log.Errorf("Error when store new event: %s", err.Error())
	}

	log.Info("Stop whashing pump successfully")

}

// StartBarrelMotor permit to start barrel motor
// The motor start only if not emmergency and no security
func (h *DFPHandler) StartBarrelMotor() {
	if h.stateRepository.CanStartMotor() {
		log.Debug("Start barrel motor")
		err := h.relayBarrelMotor.On()
		if err != nil {
			log.Errorf("Error appear when try to start barrel motor: %s", err)
		}

		event := &models.Event{
			SourceID:   h.stateRepository.State().ID,
			SourceName: h.stateRepository.State().Name,
			Timestamp:  time.Now(),
			EventType:  "start_washing_drum",
			EventKind:  "motor",
		}
		err = h.eventUsecase.Store(context.Background(), event)
		if err != nil {
			log.Errorf("Error when store new event: %s", err.Error())
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

	event := &models.Event{
		SourceID:   h.stateRepository.State().ID,
		SourceName: h.stateRepository.State().Name,
		Timestamp:  time.Now(),
		EventType:  "stop_washing_drum",
		EventKind:  "motor",
	}
	err := h.eventUsecase.Store(context.Background(), event)
	if err != nil {
		log.Errorf("Error when store new event: %s", err.Error())
	}

	log.Info("Stop barrel motor successfully")
}

// HandleMotor manage the motor state
func (h *DFPHandler) HandleMotor() {

	//Handle washing
	h.eventer.On("stateChange", func(data interface{}) {
		event := data.(string)

		log.Debugf("Receive event %s", event)

		if h.stateRepository.State().ShouldWash && h.stateRepository.CanWash() {

			// Get current config
			config, err := h.configUsecase.Get(context.Background())
			if err != nil {
				log.Errorf("Error when get current config: %s", err.Error())
				return
			}

			log.Info("Start washing cycle")
			h.stateRepository.SetWashed()

			// Start washing pump and wait
			h.StartWashingPump()
			time.Sleep(time.Second * time.Duration(config.StartWashingPumpBeforeWashing))
			if !h.stateRepository.CanStartMotor() || !h.stateRepository.State().ShouldWash {
				return
			}

			// Start barrel motor and wait
			h.StartBarrelMotor()
			time.Sleep(time.Second * time.Duration(config.WashingDuration))
			if !h.stateRepository.CanStartMotor() || !h.stateRepository.State().ShouldWash {
				return
			}

			h.StopWashingPump()
			h.StopBarrelMotor()

			event := &models.Event{
				SourceID:                h.stateRepository.State().ID,
				SourceName:              h.stateRepository.State().Name,
				Timestamp:               time.Now(),
				EventType:               "wash",
				EventKind:               "motor",
				Duration:                int64(config.WashingDuration + config.StartWashingPumpBeforeWashing),
				DurationFromLastWashing: int64(h.stateRepository.LastWashDurationSecond()),
			}
			err = h.eventUsecase.Store(context.Background(), event)
			if err != nil {
				log.Errorf("Error when store new event: %s", err.Error())
			}

			h.stateRepository.UnsetShouldWash()
			h.stateRepository.UnsetWashed()
			h.stateRepository.UpdateLastWashing()
		}
	})

	// Handle stop
	h.eventer.On("stateChange", func(data interface{}) {
		event := data.(string)

		log.Debugf("Receive event %s", event)

		// Stop motors
		if h.stateRepository.State().IsStopped || h.stateRepository.State().IsEmergencyStopped || (h.stateRepository.State().IsSecurity && !h.stateRepository.State().IsDisableSecurity) {
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
