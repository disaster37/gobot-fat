package tankboard

import (
	"context"
	"fmt"
	"time"

	"github.com/disaster37/gobot-arest/drivers/extra"
	"github.com/disaster37/gobot-fat/models"
	"github.com/disaster37/gobot-fat/tankconfig"
	log "github.com/sirupsen/logrus"
)

func (h *TankBoard) work() {

	ctx := context.Background()

	// Handle config
	h.globalEventer.On(tankconfig.NewTankConfig, func(s interface{}) {
		tankConfig := s.(*models.TankConfig)
		if tankConfig.ID == h.config.ID {
			log.Debugf("New config received for board %s, we update it", h.name)
			h.config = tankConfig

			// Publish internal event
			h.Publish(NewConfig, tankConfig)
		}
	})

	// Handle board reboot
	h.valueRebooted.On(extra.NewValue, func(s interface{}) {
		log.Debug("New value fired for isRebooted")

		isRebooted := s.(bool)
		if isRebooted {
			// Board rebooted
			log.Infof("Detect board %s is rebooted", h.name)

			// Force reconnect to init pin and set output as expected
			if err := h.board.Reconnect(); err != nil {
				log.Errorf("Error when reconnect on board %s: %s", h.name, err.Error())
			}

			// Nothink todo, juste acknoledge and send event
			if err := h.functionRebooted.Call(); err != nil {
				log.Errorf("Error when acknoledge reboot on board %s: %s", h.name, err.Error())
			}

			// Send rebooted event
			h.sendEvent(ctx, fmt.Sprintf("reboot_%s", h.name), "board", 0)

			h.isOnline = true

			// Publish internal event
			h.Publish(NewReboot, nil)
		}
	})

	// Handle board error / offline
	h.valueRebooted.On(extra.Error, func(s interface{}) {
		h.isOnline = false

		err := s.(error)
		log.Errorf("Board %s is offline: %s", h.name, err.Error())

		// Send offline event
		h.sendEvent(ctx, fmt.Sprintf("offline_%s", h.name), "board", 0)

		// Publish internal event
		h.Publish(NewOffline, nil)

	})

	// Handle read distance
	h.valueDistance.On(extra.NewValue, func(s interface{}) {
		log.Debugf("Distance change: %d", s)

		// Update distance
		distance := int64(s.(float64))
		log.Debugf("Distance on board %s: %d", h.name, distance)
		h.data.Level = int(h.config.Depth - (distance - h.config.SensorHeight))
		h.data.Volume = h.data.Level * int(h.config.LiterPerCm)
		h.data.Percent = float64(h.data.Level) / float64(h.config.Depth) * 100

		// Send event
		h.sendEvent(ctx, "read_distance", "sensor", h.data.Level)

		// Publish internal event
		h.Publish(NewDistance, distance)
	})

	// Handle error when read distance
	h.valueDistance.On(extra.Error, func(s interface{}) {
		err := s.(error)
		log.Errorf("Error when read value distance on board %s: %s", h.name, err.Error())
	})

	h.isInitialized = true
}

func (h *TankBoard) sendEvent(ctx context.Context, eventType string, eventKind string, distance int) {

	var event *models.Event
	if eventType == "read_distance" {
		event = &models.Event{
			SourceID:   fmt.Sprintf("%d", h.config.ID),
			SourceName: h.name,
			Timestamp:  time.Now(),
			EventType:  eventType,
			EventKind:  eventKind,
			Distance:   int64(distance),
		}
	} else {
		event = &models.Event{
			SourceID:   fmt.Sprintf("%d", h.config.ID),
			SourceName: h.name,
			Timestamp:  time.Now(),
			EventType:  eventType,
			EventKind:  eventKind,
		}
	}

	err := h.eventUsecase.Store(ctx, event)
	if err != nil {
		log.Errorf("Error when store new event: %s", err.Error())
	}
}
