package dfpgobot

import (
	"context"
	"time"

	"github.com/disaster37/gobot-fat/models"
	log "github.com/sirupsen/logrus"
	"gobot.io/x/gobot/drivers/gpio"
)

// HandleSecurityWaterCaptor manage the security water captor
// Check if water level is ok
func (h *DFPHandler) HandleSecurityWaterCaptor() {

	// Top captor
	// Send event only if not on Emergency stop
	h.captorWaterSecurityTop.On(gpio.ButtonPush, func(data interface{}) {
		log.Debugf("Captor water security top pushed")
		h.stateRepository.SetSecurity()

		event := &models.Event{
			SourceID:   h.stateRepository.State().ID,
			SourceName: h.stateRepository.State().Name,
			Timestamp:  time.Now(),
			EventType:  "security_top_on",
			EventKind:  "captor",
		}
		err := h.eventUsecase.Store(context.Background(), event)
		if err != nil {
			log.Errorf("Error when store new event: %s", err.Error())
		}
	})

	h.captorWaterSecurityTop.On(gpio.ButtonRelease, func(data interface{}) {
		log.Debugf("Captor water security top released")
		h.stateRepository.UnsetSecurity()

		event := &models.Event{
			SourceID:   h.stateRepository.State().ID,
			SourceName: h.stateRepository.State().Name,
			Timestamp:  time.Now(),
			EventType:  "security_top_off",
			EventKind:  "captor",
		}
		err := h.eventUsecase.Store(context.Background(), event)
		if err != nil {
			log.Errorf("Error when store new event: %s", err.Error())
		}
	})

	// Under captor
	h.captorWaterSecurityUnder.On(gpio.ButtonPush, func(data interface{}) {
		log.Debugf("Captor water security under pushed")
		h.stateRepository.SetSecurity()

		event := &models.Event{
			SourceID:   h.stateRepository.State().ID,
			SourceName: h.stateRepository.State().Name,
			Timestamp:  time.Now(),
			EventType:  "security_under_on",
			EventKind:  "captor",
		}
		err := h.eventUsecase.Store(context.Background(), event)
		if err != nil {
			log.Errorf("Error when store new event: %s", err.Error())
		}
	})

	h.captorWaterSecurityUnder.On(gpio.ButtonRelease, func(data interface{}) {
		log.Debugf("Captor water security under released")
		h.stateRepository.UnsetSecurity()

		event := &models.Event{
			SourceID:   h.stateRepository.State().ID,
			SourceName: h.stateRepository.State().Name,
			Timestamp:  time.Now(),
			EventType:  "security_under_off",
			EventKind:  "captor",
		}
		err := h.eventUsecase.Store(context.Background(), event)
		if err != nil {
			log.Errorf("Error when store new event: %s", err.Error())
		}
	})
}

// HandleWaterCaptor manage the water captor
// Check if must washing
func (h *DFPHandler) HandleWaterCaptor() {

	// Top captor
	h.captorWaterTop.On(gpio.ButtonPush, func(data interface{}) {
		log.Debugf("Captor water top pushed")

		// Get current config
		config, err := h.configUsecase.Get(context.Background())
		if err != nil {
			log.Errorf("Error when get current config: %s", err.Error())
			return
		}

		if h.stateRepository.State().IsAuto && (h.stateRepository.LastWashDurationSecond() > uint64(config.WaitTimeBetweenWashing)) {
			h.stateRepository.SetShouldWash()
		}

		event := &models.Event{
			SourceID:   h.stateRepository.State().ID,
			SourceName: h.stateRepository.State().Name,
			Timestamp:  time.Now(),
			EventType:  "water_top_on",
			EventKind:  "captor",
		}
		err = h.eventUsecase.Store(context.Background(), event)
		if err != nil {
			log.Errorf("Error when store new event: %s", err.Error())
		}
	})

	h.captorWaterTop.On(gpio.ButtonRelease, func(data interface{}) {
		log.Debugf("Captor water top released")

		event := &models.Event{
			SourceID:   h.stateRepository.State().ID,
			SourceName: h.stateRepository.State().Name,
			Timestamp:  time.Now(),
			EventType:  "water_top_off",
			EventKind:  "captor",
		}
		err := h.eventUsecase.Store(context.Background(), event)
		if err != nil {
			log.Errorf("Error when store new event: %s", err.Error())
		}
	})

	// Under captor
	h.captorWaterUnder.On(gpio.ButtonPush, func(data interface{}) {
		log.Debugf("Captor water under pushed")

		// Get current config
		config, err := h.configUsecase.Get(context.Background())
		if err != nil {
			log.Errorf("Error when get current config: %s", err.Error())
			return
		}

		if h.stateRepository.State().IsAuto && (h.stateRepository.LastWashDurationSecond() > uint64(config.WaitTimeBetweenWashing)) {
			h.stateRepository.SetShouldWash()
		}

		event := &models.Event{
			SourceID:   h.stateRepository.State().ID,
			SourceName: h.stateRepository.State().Name,
			Timestamp:  time.Now(),
			EventType:  "water_under_on",
			EventKind:  "captor",
		}
		err = h.eventUsecase.Store(context.Background(), event)
		if err != nil {
			log.Errorf("Error when store new event: %s", err.Error())
		}
	})

	h.captorWaterUnder.On(gpio.ButtonRelease, func(data interface{}) {
		log.Debugf("Captor water under released")

		event := &models.Event{
			SourceID:   h.stateRepository.State().ID,
			SourceName: h.stateRepository.State().Name,
			Timestamp:  time.Now(),
			EventType:  "water_under_off",
			EventKind:  "captor",
		}
		err := h.eventUsecase.Store(context.Background(), event)
		if err != nil {
			log.Errorf("Error when store new event: %s", err.Error())
		}
	})
}
