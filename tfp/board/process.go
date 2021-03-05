package tfpboard

import (
	"context"
	"time"

	"github.com/disaster37/gobot-fat/helper"
	"github.com/disaster37/gobot-fat/tfpstate"
	"gobot.io/x/gobot"

	"github.com/disaster37/gobot-arest/drivers/extra"
	"github.com/disaster37/gobot-fat/models"
	"github.com/disaster37/gobot-fat/tfpconfig"
	"github.com/labstack/gommon/log"
)

func (h *TFPBoard) work() {

	ctx := context.Background()

	// Handle config
	h.on(h.globalEventer, tfpconfig.NewTFPConfig, func(s interface{}) {
		tfpConfig := s.(*models.TFPConfig)
		log.Debugf("New config received for board %s, we update it", h.name)

		h.config = tfpConfig

		// Publish internal event
		h.Publish(EventNewConfig, tfpConfig)
	})

	// Handle state
	h.on(h.globalEventer, tfpstate.NewTFPState, func(s interface{}) {

		tfpState := s.(*models.TFPState)
		log.Debugf("New state received for board %s, we update it", h.name)
		h.state.UVC1BlisterNbHour = tfpState.UVC1BlisterNbHour
		h.state.UVC2BlisterNbHour = tfpState.UVC2BlisterNbHour
		h.state.OzoneBlisterNbHour = tfpState.OzoneBlisterNbHour

		// Publish internal event
		h.Publish(EventNewState, h.state)
	})

	// Handle board reboot
	h.on(h.valueRebooted, extra.NewValue, func(s interface{}) {
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

			// Send event
			helper.SendEvent(ctx, h.eventUsecase, h.name, helper.KindEventRebootBoard, h.name)

			h.isOnline = true

			// Publish internal event
			h.Publish(EventBoardReboot, nil)
		}
	})

	// Handle board error / offline
	h.on(h.valueRebooted, extra.Error, func(s interface{}) {
		h.isOnline = false

		err := s.(error)
		log.Errorf("Board %s is offline: %s", h.name, err.Error())

		// Send event
		helper.SendEvent(ctx, h.eventUsecase, h.name, helper.KindEventOfflineBoard, h.name)

		// Publish internal event
		h.Publish(EventBoardOffline, nil)

	})

	// Handle set emergency stop
	h.on(h.globalEventer, helper.SetEmergencyStop, func(s interface{}) {
		h.state.IsEmergencyStopped = true

		// Stop UVC1
		if err := h.relayUVC1.Off(); err != nil {
			log.Errorf("Error when stop UVC1: %s", err.Error())
		}

		// Stop UVC2
		if err := h.relayUVC2.Off(); err != nil {
			log.Errorf("Error when stop UVC2: %s", err.Error())
		}

		// Stop pond pump
		if err := h.relayPompPond.Off(); err != nil {
			log.Errorf("Error when stop pond pump: %s", err.Error())
		}

		// Stop pond bubble
		if err := h.relayBubblePond.Off(); err != nil {
			log.Errorf("Error when stop pond bubble: %s", err.Error())
		}

		// Stop filter bubble
		if err := h.relayBubbleFilter.Off(); err != nil {
			log.Errorf("Error when stop filter bubble: %s", err.Error())
		}

		// Stop waterfall pump
		if err := h.relayPompWaterfall.Off(); err != nil {
			log.Errorf("Error when stop waterfall pump: %s", err.Error())
		}

		// Send event
		helper.SendEvent(ctx, h.eventUsecase, h.name, helper.KindEventSetEmergencyStop, h.name)

		// Publish internal event
		h.Publish(EventSetEmergencyStop, nil)
	})

	// Handle unset emergency stop
	h.on(h.globalEventer, helper.UnsetEmergencyStop, func(s interface{}) {

		h.state.IsEmergencyStopped = false

		h.handleUnsetSecurityOrEmergencyStop()

		// Send event
		helper.SendEvent(ctx, h.eventUsecase, h.name, helper.KindEventUnsetEmergencyStop, h.name)

		// Publish internal event
		h.Publish(EventUnsetEmergencyStop, nil)
	})

	// Handle set secrutity
	h.on(h.globalEventer, helper.SetSecurity, func(s interface{}) {
		h.state.IsSecurity = true

		// Stop UVC1
		if err := h.relayUVC1.Off(); err != nil {
			log.Errorf("Error when stop UVC1: %s", err.Error())
		}

		// Stop UVC2
		if err := h.relayUVC2.Off(); err != nil {
			log.Errorf("Error when stop UVC2: %s", err.Error())
		}

		// Stop pond pump
		if err := h.relayPompPond.Off(); err != nil {
			log.Errorf("Error when stop pond pump: %s", err.Error())
		}

		// Stop waterfall pump
		if err := h.relayPompWaterfall.Off(); err != nil {
			log.Errorf("Error when stop waterfall pond pump: %s", err.Error())
		}

		// Send event
		helper.SendEvent(ctx, h.eventUsecase, h.name, helper.KindEventSetSecurity, h.name)

		// Publish internal event
		h.Publish(EventSetSecurity, nil)
	})

	// Handler unset security
	h.on(h.globalEventer, helper.UnsetSecurity, func(data interface{}) {

		h.state.IsSecurity = false

		h.handleUnsetSecurityOrEmergencyStop()

		// Send event
		helper.SendEvent(ctx, h.eventUsecase, h.name, helper.KindEventUnsetSecurity, h.name)

		// Publish internal event
		h.Publish(EventUnsetSecurity, nil)

	})

	// Handle blister time
	h.schedulingRoutines = append(h.schedulingRoutines, gobot.Every(1*time.Hour, h.handleBlisterTime))

	// Handle waterfall auto
	h.schedulingRoutines = append(h.schedulingRoutines, gobot.Every(1*time.Minute, h.handleWaterfallAuto))

	log.Debugf("TFP IO:\n %s", h.IO().String())
	log.Debugf("TFP state: %s", h.state.String())

	h.isInitialized = true
}

// handleBlisterTime permit to increment the number of hour of each blister enabled
func (h *TFPBoard) handleBlisterTime() {
	ctx := context.Background()
	isUpdated := false

	// If we can start relay, All UVC are already stopped
	if !h.canStartRelay() {
		return
	}

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
func (h *TFPBoard) handleWaterfallAuto() {
	ctx := context.Background()

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
				// Force state if security or emergency stop, to start after that
				h.state.WaterfallPumpRunning = true
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

// Use on instead gobot.Eventer.On because of it not close routine at board is stopped.
// So, if you start / stop / start board, you have so many routine
func (h *TFPBoard) on(driver gobot.Eventer, event string, f func(data interface{})) {

	halt := make(chan bool)

	// Detect stop board
	go func() {
		out := h.Subscribe()

		for {
			select {
			case evt := <-out:
				if evt.Name == EventBoardStop {
					halt <- true
					h.Unsubscribe(out)
					return
				}
			}
		}
	}()

	// Handle on event
	go func() {
		out := driver.Subscribe()
		for {
			select {
			case <-halt:
				driver.Unsubscribe(out)
				return
			case evt := <-out:
				if evt.Name == event {
					f(evt.Data)
				}
			}
		}

	}()
}

func (h *TFPBoard) handleUnsetSecurityOrEmergencyStop() {

	ctx := context.Background()

	// We can start bubbles
	if !h.state.IsEmergencyStopped {

		// Filter bubble
		if h.state.FilterBubbleRunning {
			if err := h.StartFilterBubble(ctx); err != nil {
				log.Errorf("When start filter bubble: %s", err.Error())
			}
		}

		// Pond bubble
		if h.state.PondBubbleRunning {
			if err := h.StartPondBubble(ctx); err != nil {
				log.Errorf("When start pond bubble: %s", err.Error())
			}
		}

	}

	// Start other relais
	if h.canStartRelay() {

		// Pond pump
		if h.state.PondPumpRunning {
			if err := h.StartPondPump(ctx); err != nil {
				log.Errorf("When start pond pump: %s", err.Error())
			}
		}

		// UVC1
		if h.state.UVC1Running {
			if err := h.StartUVC1(ctx); err != nil {
				log.Errorf("When start UVC1: %s", err.Error())
			}
		}

		// UVC2
		if h.state.UVC2Running {
			if err := h.StartUVC2(ctx); err != nil {
				log.Errorf("When start UVC2: %s", err.Error())
			}
		}

		// Waterfall pump
		if h.state.WaterfallPumpRunning {
			if err := h.StartWaterfallPump(ctx); err != nil {
				log.Errorf("When start waterfall pump: %s", err.Error())
			}
		}
	}
}
