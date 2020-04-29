package tfpgobot

import (
	"context"
	"time"

	"github.com/disaster37/gobot-fat/models"
	log "github.com/sirupsen/logrus"
)

// StartPondPump permit to run pond pump
// The pump start only if no emergency and no security
func (h *TFPHandler) StartPondPump() error {
	if h.stateRepository.CanStartRelay() {
		log.Debug("Start pond pump")
		err := h.relayPompPond.On()
		if err != nil {
			return err
		}

		event := &models.Event{
			SourceID:   h.stateRepository.State().ID,
			SourceName: h.stateRepository.State().Name,
			Timestamp:  time.Now(),
			EventType:  "start_pond_pump",
			EventKind:  "pump",
		}
		err = h.eventUsecase.Store(context.Background(), event)
		if err != nil {
			log.Errorf("Error when store new event: %s", err.Error())
		}

		log.Info("Start pond pump successfully")
	} else {
		log.Info("Pond pump not started because of state not permit it")
	}

	return nil
}

// StartUVC1 permit to run UVC1
// The UVC start only if no emergency and no security
func (h *TFPHandler) StartUVC1() error {
	if h.stateRepository.CanStartRelay() && h.stateRepository.State().PondPumpRunning {
		log.Debug("Start UVC1")
		err := h.relayUVC1.On()
		if err != nil {
			return err
		}

		event := &models.Event{
			SourceID:   h.stateRepository.State().ID,
			SourceName: h.stateRepository.State().Name,
			Timestamp:  time.Now(),
			EventType:  "start_uvc1",
			EventKind:  "uvc",
		}
		err = h.eventUsecase.Store(context.Background(), event)
		if err != nil {
			log.Errorf("Error when store new event: %s", err.Error())
		}

		log.Info("Start UVC1 successfully")
	} else {
		log.Info("UVC1 not started because of state not permit it")
	}

	return nil
}

// StartUVC2 permit to run UVC2
// The UVC start only if no emergency and no security
func (h *TFPHandler) StartUVC2() error {
	if h.stateRepository.CanStartRelay() && h.stateRepository.State().PondPumpRunning {
		log.Debug("Start UVC2")
		err := h.relayUVC2.On()
		if err != nil {
			return err
		}

		event := &models.Event{
			SourceID:   h.stateRepository.State().ID,
			SourceName: h.stateRepository.State().Name,
			Timestamp:  time.Now(),
			EventType:  "start_uvc2",
			EventKind:  "uvc",
		}
		err = h.eventUsecase.Store(context.Background(), event)
		if err != nil {
			log.Errorf("Error when store new event: %s", err.Error())
		}

		log.Info("Start UVC2 successfully")
	} else {
		log.Info("UVC2 not started because of state not permit it")
	}

	return nil
}

// StartPondPumpWithUVC permit to start pond pump with UVC
// The pump start only if no emergency and no security
func (h *TFPHandler) StartPondPumpWithUVC() error {
	if h.stateRepository.CanStartRelay() {
		log.Debug("Start pond pump with UVC")

		err := h.StartPondPump()
		if err != nil {
			return err
		}
		err = h.StartUVC1()
		if err != nil {
			return err
		}
		err = h.StartUVC2()
		if err != nil {
			return err
		}

		log.Info("Start pond pump with UVCs successfully")
	} else {
		log.Info("Pond pump with UVC not started because of state not permit it")
	}

	return nil
}

// StopUVC1 permit to stop UVC1
// It will try while not stopped
func (h *TFPHandler) StopUVC1() {
	log.Debug("Stop UVC1")

	isStopped := false
	for isStopped == false {
		err := h.relayUVC1.Off()
		if err != nil {
			log.Errorf("Error when stop UVC1: %s", err)
			time.Sleep(1 * time.Second)
		} else {
			isStopped = true
		}
	}

	event := &models.Event{
		SourceID:   h.stateRepository.State().ID,
		SourceName: h.stateRepository.State().Name,
		Timestamp:  time.Now(),
		EventType:  "stop_uvc1",
		EventKind:  "uvc",
	}
	err := h.eventUsecase.Store(context.Background(), event)
	if err != nil {
		log.Errorf("Error when store new event: %s", err.Error())
	}

	log.Info("Stop UVC1 successfully")

}

// StopUVC2 permit to stop UVC2
// It will try while not stopped
func (h *TFPHandler) StopUVC2() {
	log.Debug("Stop UVC2")

	isStopped := false
	for isStopped == false {
		err := h.relayUVC2.Off()
		if err != nil {
			log.Errorf("Error when stop UVC2: %s", err)
			time.Sleep(1 * time.Second)
		} else {
			isStopped = true
		}
	}

	event := &models.Event{
		SourceID:   h.stateRepository.State().ID,
		SourceName: h.stateRepository.State().Name,
		Timestamp:  time.Now(),
		EventType:  "stop_uvc2",
		EventKind:  "uvc",
	}
	err := h.eventUsecase.Store(context.Background(), event)
	if err != nil {
		log.Errorf("Error when store new event: %s", err.Error())
	}

	log.Info("Stop UVC2 successfully")
}

// StopPondPump permit to stop pond pump
// It will try while not stopped
// It will stop all UVC
func (h *TFPHandler) StopPondPump() {

	h.StopUVC1()
	h.StopUVC2()

	isStopped := false
	for isStopped == false {
		err := h.relayPompPond.Off()
		if err != nil {
			log.Errorf("Error when stop pond pump: %s", err)
			time.Sleep(1 * time.Second)
		} else {
			isStopped = true
		}
	}

	event := &models.Event{
		SourceID:   h.stateRepository.State().ID,
		SourceName: h.stateRepository.State().Name,
		Timestamp:  time.Now(),
		EventType:  "stop_pond_pump",
		EventKind:  "pump",
	}
	err := h.eventUsecase.Store(context.Background(), event)
	if err != nil {
		log.Errorf("Error when store new event: %s", err.Error())
	}

	log.Info("Stop pond pump successfully")
}

// StartWaterfallPump permit to start waterfall pump
// The motor start only if not emmergency and no security
func (h *TFPHandler) StartWaterfallPump() error {
	if h.stateRepository.CanStartRelay() {
		log.Debug("Start waterfall pump")
		err := h.relayPompWaterfall.On()
		if err != nil {
			return err
		}

		event := &models.Event{
			SourceID:   h.stateRepository.State().ID,
			SourceName: h.stateRepository.State().Name,
			Timestamp:  time.Now(),
			EventType:  "start_watterfall_pump",
			EventKind:  "pump",
		}
		err = h.eventUsecase.Store(context.Background(), event)
		if err != nil {
			log.Errorf("Error when store new event: %s", err.Error())
		}

		log.Info("Start waterfall pump successfully")
	} else {
		log.Info("Waterfall pump not started because of state not permit it")
	}

	return nil

}

// StopWaterfallPump permit to stop waterfall pump
// It will try while is not stopped
func (h *TFPHandler) StopWaterfallPump() {
	log.Debug("Stop waterfall pump")

	isStopped := false
	for isStopped == false {
		err := h.relayPompWaterfall.Off()
		if err != nil {
			log.Errorf("Error when stop waterfall pump: %s", err)
			time.Sleep(1 * time.Second)
		} else {
			isStopped = true
		}
	}

	event := &models.Event{
		SourceID:   h.stateRepository.State().ID,
		SourceName: h.stateRepository.State().Name,
		Timestamp:  time.Now(),
		EventType:  "stop_waterfall_pump",
		EventKind:  "pump",
	}
	err := h.eventUsecase.Store(context.Background(), event)
	if err != nil {
		log.Errorf("Error when store new event: %s", err.Error())
	}

	log.Info("Stop waterfall pump successfully")
}

// StartPondBubble permit to start pond bubble
// The motor start only if not emmergency and no security
func (h *TFPHandler) StartPondBubble() error {
	if h.stateRepository.CanStartRelay() {
		log.Debug("Start pond bubble")
		err := h.relayBubblePond.On()
		if err != nil {
			return err
		}

		event := &models.Event{
			SourceID:   h.stateRepository.State().ID,
			SourceName: h.stateRepository.State().Name,
			Timestamp:  time.Now(),
			EventType:  "start_pond_bubble",
			EventKind:  "bubble",
		}
		err = h.eventUsecase.Store(context.Background(), event)
		if err != nil {
			log.Errorf("Error when store new event: %s", err.Error())
		}

		log.Info("Start pond bubble successfully")
	} else {
		log.Info("Pond bubble not started because of state not permit it")
	}

	return nil

}

// StopPondBubble permit to stop pond bubble
// It will try while is not stopped
func (h *TFPHandler) StopPondBubble() {
	log.Debug("Stop pond bubble")

	isStopped := false
	for isStopped == false {
		err := h.relayBubblePond.Off()
		if err != nil {
			log.Errorf("Error when stop pond bubble: %s", err)
			time.Sleep(1 * time.Second)
		} else {
			isStopped = true
		}
	}

	event := &models.Event{
		SourceID:   h.stateRepository.State().ID,
		SourceName: h.stateRepository.State().Name,
		Timestamp:  time.Now(),
		EventType:  "stop_pund_bubble",
		EventKind:  "bubble",
	}
	err := h.eventUsecase.Store(context.Background(), event)
	if err != nil {
		log.Errorf("Error when store new event: %s", err.Error())
	}

	log.Info("Stop pond bubble successfully")
}

// StartFilterBubble permit to start filter bubble
// The motor start only if not emmergency and no security
func (h *TFPHandler) StartFilterBubble() error {
	if h.stateRepository.CanStartRelay() {
		log.Debug("Start filter bubble")
		err := h.relayBubbleFilter.On()
		if err != nil {
			return err
		}

		event := &models.Event{
			SourceID:   h.stateRepository.State().ID,
			SourceName: h.stateRepository.State().Name,
			Timestamp:  time.Now(),
			EventType:  "start_filter_bubble",
			EventKind:  "bubble",
		}
		err = h.eventUsecase.Store(context.Background(), event)
		if err != nil {
			log.Errorf("Error when store new event: %s", err.Error())
		}

		log.Info("Start filter bubble successfully")
	} else {
		log.Info("Filter bubble not started because of state not permit it")
	}

	return nil

}

// StopFilterBubble permit to stop filter bubble
// It will try while is not stopped
func (h *TFPHandler) StopFilterBubble() {
	log.Debug("Stop filter bubble")

	isStopped := false
	for isStopped == false {
		err := h.relayBubbleFilter.Off()
		if err != nil {
			log.Errorf("Error when stop filter bubble: %s", err)
			time.Sleep(1 * time.Second)
		} else {
			isStopped = true
		}
	}

	event := &models.Event{
		SourceID:   h.stateRepository.State().ID,
		SourceName: h.stateRepository.State().Name,
		Timestamp:  time.Now(),
		EventType:  "stop_filter_bubble",
		EventKind:  "bubble",
	}
	err := h.eventUsecase.Store(context.Background(), event)
	if err != nil {
		log.Errorf("Error when store new event: %s", err.Error())
	}

	log.Info("Stop filter bubble successfully")
}

// HandleRelay manage the relay state
func (h *TFPHandler) HandleRelay() {

	// Handle ermergency stop
	h.eventer.On("stateChange", func(data interface{}) {
		event := data.(string)

		log.Debugf("Receive event %s", event)

		// Stop relais
		if h.stateRepository.State().IsEmergencyStopped {
			h.StopRelais()
		}
	})

	// Handle security
	h.eventer.On("stateChange", func(data interface{}) {
		event := data.(string)

		log.Debugf("Receive event %s", event)

		// Stop pump
		if h.stateRepository.State().IsSecurity && !h.stateRepository.State().IsDisableSecurity {
			h.StopPondPump()
			h.StopWaterfallPump()
		}
	})

	// Handle relay
	h.eventer.On("stateChange", func(data interface{}) {
		event := data.(string)

		log.Debugf("Receive event %s", event)

		if event == "initTFP" {
			// Manage pond pump
			if h.stateRepository.State().PondPumpRunning {
				err := h.StartPondPump()
				if err != nil {
					log.Errorf("Failed to start pond pump: %s", err.Error())
				}
			} else {
				h.StopPondPump()
			}

			// Manage UVC1
			if h.stateRepository.State().UVC1Running {
				err := h.StartUVC1()
				if err != nil {
					log.Errorf("Failed to start UVC1: %s", err.Error())
				}
			} else {
				h.StopUVC1()
			}

			// Manage UVC2
			if h.stateRepository.State().UVC2Running {
				err := h.StartUVC2()
				if err != nil {
					log.Errorf("Failed to start UVC2: %s", err.Error())
				}
			} else {
				h.StopUVC2()
			}

			// Manage waterfall pump
			if h.stateRepository.State().WaterfallPumpRunning {
				err := h.StartWaterfallPump()
				if err != nil {
					log.Errorf("Failed to start Waterfall pump: %s", err.Error())
				}
			} else {
				h.StopWaterfallPump()
			}

			// Manage pond bubble
			if h.stateRepository.State().PondBubbleRunning {
				err := h.StartPondBubble()
				if err != nil {
					log.Errorf("Failed to start pond bubble: %s", err.Error())
				}
			} else {
				h.StopPondBubble()
			}

			// Manage filter bubble
			if h.stateRepository.State().FilterBubbleRunning {
				err := h.StartFilterBubble()
				if err != nil {
					log.Errorf("Failed to start filter bubble: %s", err.Error())
				}
			} else {
				h.StopFilterBubble()
			}
		}
	})
}

// StopRelais stop all relais
func (h *TFPHandler) StopRelais() {
	log.Info("Stop all relais")
	h.StopPondPump()
	h.StopWaterfallPump()
	h.StopPondBubble()
	h.StopFilterBubble()
	h.StopUVC1()
	h.StopUVC2()
}
