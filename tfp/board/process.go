package tfpboard

import (
	"context"
	"fmt"
	"time"

	"github.com/disaster37/gobot-fat/tfpstate"

	"github.com/disaster37/gobot-arest/drivers/extra"
	"github.com/disaster37/gobot-fat/board"
	"github.com/disaster37/gobot-fat/models"
	"github.com/disaster37/gobot-fat/tfpconfig"
	"github.com/labstack/gommon/log"
)

func (h *TFPBoard) work() {

	ctx := context.Background()

	// Handle config
	h.globalEventer.On(tfpconfig.NewTFPConfig, func(s interface{}) {
		tfpConfig := s.(*models.TFPConfig)
		log.Debugf("New config received for board %s, we update it", h.name)

		h.config = tfpConfig

		// Publish internal event
		h.Publish(NewConfig, tfpConfig)
	})

	// Handle state
	h.globalEventer.On(tfpstate.NewTFPState, func(s interface{}) {

		tfpState := s.(*models.TFPState)
		log.Debugf("New state received for board %s, we update it", h.name)
		h.state.UVC1BlisterNbHour = tfpState.UVC1BlisterNbHour
		h.state.UVC2BlisterNbHour = tfpState.UVC2BlisterNbHour
		h.state.OzoneBlisterNbHour = tfpState.OzoneBlisterNbHour

		// Publish internal event
		h.Publish(NewState, h.state)
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
			h.sendEvent(ctx, fmt.Sprintf("reboot_%s", h.name), "board")

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
		h.sendEvent(ctx, fmt.Sprintf("offline_%s", h.name), "board")

		// Publish internal event
		h.Publish(NewOffline, nil)

	})

	// Handle blister time
	board.NewHandler(ctx, 1*time.Hour, h.chStop, h.handleBlisterTime)

	// Handle watrefall auto
	board.NewHandler(ctx, 1*time.Minute, h.chStop, h.handleWaterfallAuto)

	h.isInitialized = true
}

// handleBlisterTime permit to increment the number of hour of each blister enabled
func (h *TFPBoard) handleBlisterTime(ctx context.Context) {

	isUpdated := false

	switch h.config.Mode {
	case "ozone":
		log.Debug("Ozone mode detected")
		if h.state.UVC1Running {
			h.state.UVC1BlisterNbHour++
			isUpdated = true
		}
		if h.state.UVC2Running {
			h.state.OzoneBlisterNbHour++
			isUpdated = true
		}
	case "uvc":
		log.Debug("UVC mode detected")
		if h.state.UVC1Running {
			h.state.UVC1BlisterNbHour++
			isUpdated = true
		}
		if h.state.UVC2Running {
			h.state.UVC2BlisterNbHour++
			isUpdated = true
		}
	case "none":
		log.Debug("None mode detected")
		return
	default:
		log.Warn("Can't detect mode")
		return
	}

	if isUpdated {
		err := h.stateUsecase.Update(ctx, h.state)
		if err != nil {
			log.Errorf("Error when save blister time: %s", err.Error())
		}
	}
}

// handleWaterfall auto permit to start and stop waterfall automatically
func (h *TFPBoard) handleWaterfallAuto(ctx context.Context) {

	if h.config.IsWaterfallAuto {
		startDate, err := time.Parse("15:04", h.config.StartTimeWaterfall)
		if err != nil {
			log.Errorf("Error when parse StartTimeWaterfall: %s", err.Error())
			return
		}
		endDate, err := time.Parse("15:04", h.config.StopTimeWaterfall)
		if err != nil {
			log.Errorf("Error when parse StopTimeWaterfall: %s", err.Error())
			return
		}
		currentDate, err := time.Parse("15:04", time.Now().Format("15:04"))
		if err != nil {
			log.Errorf("Error when parse currentdata: %s", err.Error())
			return
		}

		isUpdated := false

		if startDate.Before(currentDate) && endDate.After(currentDate) {
			if h.state.AcknoledgeWaterfallAuto != true {
				log.Debug("Waterfall must be running")
				err := h.StartWaterfallPump(ctx)
				if err != nil {
					log.Errorf("Error when try to start automatically waterfall pomp: %s", err.Error())
					return
				}
				h.state.AcknoledgeWaterfallAuto = true
				isUpdated = true
			}

		} else {
			if h.state.AcknoledgeWaterfallAuto {
				log.Debug("Waterfall must be stopped")
				err := h.StopWaterfallPump(ctx)
				if err != nil {
					log.Errorf("Error when try to stop automatically waterfall pomp: %s", err.Error())
					return
				}
				h.state.AcknoledgeWaterfallAuto = false
				isUpdated = true
			}
		}

		if isUpdated {
			err := h.stateUsecase.Update(ctx, h.state)
			if err != nil {
				log.Errorf("Error when try to update tfp state after manage auto waterfall mode: %s", err.Error())
				return
			}
		}

	} else {
		log.Debug("Waterfall is on manual mode")
		return
	}

}

func (h *TFPBoard) sendEvent(ctx context.Context, eventType string, eventKind string) {
	event := &models.Event{
		SourceID:   h.name,
		SourceName: h.name,
		Timestamp:  time.Now(),
		EventType:  eventType,
		EventKind:  eventKind,
	}
	err := h.eventUsecase.Store(ctx, event)
	if err != nil {
		log.Errorf("Error when store new event: %s", err.Error())
	}
}
