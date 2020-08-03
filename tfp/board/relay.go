package tfpboard

import (
	"context"
	"errors"
	"time"

	"github.com/disaster37/gobot-fat/models"
	log "github.com/sirupsen/logrus"
)

var ErrRelayCanNotStart = errors.New("Relay can't start because of current state")

func (h *TFPHandler) canStartRelay() bool {
	if !h.state.IsEmergencyStopped && (!h.state.IsSecurity || h.state.IsDisableSecurity) {
		return true
	}
	return false
}

func (h *TFPHandler) sendEvent(eventType string, eventKind string) {
	event := &models.Event{
		SourceID:   h.state.Name,
		SourceName: h.state.Name,
		Timestamp:  time.Now(),
		EventType:  eventType,
		EventKind:  eventKind,
	}
	err := h.eventUsecase.Store(context.Background(), event)
	if err != nil {
		log.Errorf("Error when store new event: %s", err.Error())
	}
}

// StartPondPump permit to run pond pump
// The pump start only if no emergency and no security
func (h *TFPHandler) StartPondPump(ctx context.Context) error {
	if h.canStartRelay() {
		log.Debug("Start pond pump")
		err := h.relayPompPond.On()
		if err != nil {
			return err
		}

		h.sendEvent("start_pond_pump", "pump")

		// Save state only if state change
		if !h.state.PondPumpRunning {
			h.state.PondPumpRunning = true
			err = h.stateUsecase.Update(ctx, h.state)
			if err != nil {
				return err
			}
		}

		log.Info("Start pond pump successfully")
	} else {
		log.Info("Pond pump not started because of state not permit it")
		return ErrRelayCanNotStart
	}

	return nil
}

// StartUVC1 permit to run UVC1
// The UVC start only if no emergency and no security
func (h *TFPHandler) StartUVC1(ctx context.Context) error {

	if h.canStartRelay() && h.state.PondPumpRunning {
		log.Debug("Start UVC1")
		err := h.relayUVC1.On()
		if err != nil {
			return err
		}

		h.sendEvent("start_uvc1", "uvc")

		// Save state only if state change
		if !h.state.UVC1Running {
			h.state.UVC1Running = true
			err = h.stateUsecase.Update(ctx, h.state)
			if err != nil {
				return err
			}
		}

		log.Info("Start UVC1 successfully")
	} else {
		log.Info("UVC1 not started because of state not permit it")
		return ErrRelayCanNotStart
	}

	return nil
}

// StartUVC2 permit to run UVC2
// The UVC start only if no emergency and no security
func (h *TFPHandler) StartUVC2(ctx context.Context) error {
	if h.canStartRelay() && h.state.PondPumpRunning {
		log.Debug("Start UVC2")
		err := h.relayUVC2.On()
		if err != nil {
			return err
		}

		h.sendEvent("start_uvc2", "uvc")

		// Save state only if state change
		if !h.state.UVC2Running {
			h.state.UVC2Running = true
			err = h.stateUsecase.Update(ctx, h.state)
			if err != nil {
				return err
			}
		}

		log.Info("Start UVC2 successfully")
	} else {
		log.Info("UVC2 not started because of state not permit it")
		return ErrRelayCanNotStart
	}

	return nil
}

// StartPondPumpWithUVC permit to start pond pump with UVC
// The pump start only if no emergency and no security
func (h *TFPHandler) StartPondPumpWithUVC(ctx context.Context) error {
	if h.canStartRelay() {
		log.Debug("Start pond pump with UVC")

		err := h.StartPondPump(ctx)
		if err != nil {
			return err
		}
		err = h.StartUVC1(ctx)
		if err != nil {
			return err
		}
		err = h.StartUVC2(ctx)
		if err != nil {
			return err
		}

		log.Info("Start pond pump with UVCs successfully")
	} else {
		log.Info("Pond pump with UVC not started because of state not permit it")
		return ErrRelayCanNotStart
	}

	return nil
}

// StopUVC1 permit to stop UVC1
func (h *TFPHandler) StopUVC1(ctx context.Context) error {
	log.Debug("Stop UVC1")

	err := h.relayUVC1.Off()
	if err != nil {
		return err
	}

	h.sendEvent("stop_uvc1", "uvc")

	// Save state only if state change
	if h.state.UVC1Running {
		h.state.UVC1Running = false
		err = h.stateUsecase.Update(ctx, h.state)
		if err != nil {
			return err
		}
	}

	log.Info("Stop UVC1 successfully")

	return nil

}

// StopUVC2 permit to stop UVC2
// It will try while not stopped
func (h *TFPHandler) StopUVC2(ctx context.Context) error {
	log.Debug("Stop UVC2")

	err := h.relayUVC2.Off()
	if err != nil {
		return err
	}

	h.sendEvent("stop_uvc2", "uvc")

	// Save state only if state change
	if h.state.UVC2Running {
		h.state.UVC2Running = false
		err = h.stateUsecase.Update(ctx, h.state)
		if err != nil {
			return err
		}
	}

	log.Info("Stop UVC2 successfully")

	return nil
}

// StopPondPump permit to stop pond pump
// It will try while not stopped
// It will stop all UVC
func (h *TFPHandler) StopPondPump(ctx context.Context) error {

	err := h.StopUVC1(ctx)
	if err != nil {
		return err
	}

	err = h.StopUVC2(ctx)
	if err != nil {
		return err
	}

	err = h.relayPompPond.Off()
	if err != nil {
		return err
	}

	h.sendEvent("stop_pond_pump", "pump")

	// Save state only if state change
	if h.state.PondPumpRunning {
		h.state.PondPumpRunning = false
		err = h.stateUsecase.Update(ctx, h.state)
		if err != nil {
			return err
		}
	}

	log.Info("Stop pond pump with UVCs successfully")

	return nil
}

// StartWaterfallPump permit to start waterfall pump
// The motor start only if not emmergency and no security
func (h *TFPHandler) StartWaterfallPump(ctx context.Context) error {
	if h.canStartRelay() {
		log.Debug("Start waterfall pump")
		err := h.relayPompWaterfall.On()
		if err != nil {
			return err
		}

		h.sendEvent("start_watterfall_pump", "pump")

		// Save state only if state change
		if !h.state.WaterfallPumpRunning {
			h.state.WaterfallPumpRunning = true
			err = h.stateUsecase.Update(ctx, h.state)
			if err != nil {
				return err
			}
		}

		log.Info("Start waterfall pump successfully")
	} else {
		log.Info("Waterfall pump not started because of state not permit it")
		return ErrRelayCanNotStart
	}

	return nil

}

// StopWaterfallPump permit to stop waterfall pump
// It will try while is not stopped
func (h *TFPHandler) StopWaterfallPump(ctx context.Context) error {
	log.Debug("Stop waterfall pump")

	err := h.relayPompWaterfall.Off()
	if err != nil {
		return err
	}

	h.sendEvent("stop_waterfall_pump", "pump")

	// Save state only if state change
	if h.state.WaterfallPumpRunning {
		h.state.WaterfallPumpRunning = false
		err = h.stateUsecase.Update(ctx, h.state)
		if err != nil {
			return err
		}
	}

	log.Info("Stop waterfall pump successfully")

	return nil
}

// StartPondBubble permit to start pond bubble
// The motor start only if not emmergency and no security
func (h *TFPHandler) StartPondBubble(ctx context.Context) error {
	if h.canStartRelay() {
		log.Debug("Start pond bubble")
		err := h.relayBubblePond.On()
		if err != nil {
			return err
		}

		h.sendEvent("start_pond_bubble", "bubble")

		// Save state only if state change
		if !h.state.PondBubbleRunning {
			h.state.PondBubbleRunning = true
			err = h.stateUsecase.Update(ctx, h.state)
			if err != nil {
				return err
			}
		}

		log.Info("Start pond bubble successfully")
	} else {
		log.Info("Pond bubble not started because of state not permit it")
		return ErrRelayCanNotStart
	}

	return nil

}

// StopPondBubble permit to stop pond bubble
// It will try while is not stopped
func (h *TFPHandler) StopPondBubble(ctx context.Context) error {
	log.Debug("Stop pond bubble")

	err := h.relayBubblePond.Off()
	if err != nil {
		return err
	}

	h.sendEvent("stop_pond_bubble", "bubble")

	// Save state only if state change
	if h.state.PondBubbleRunning {
		h.state.PondBubbleRunning = false
		err = h.stateUsecase.Update(ctx, h.state)
		if err != nil {
			return err
		}
	}

	log.Info("Stop pond bubble successfully")

	return nil
}

// StartFilterBubble permit to start filter bubble
// The motor start only if not emmergency and no security
func (h *TFPHandler) StartFilterBubble(ctx context.Context) error {
	if h.canStartRelay() {
		log.Debug("Start filter bubble")
		err := h.relayBubbleFilter.On()
		if err != nil {
			return err
		}

		h.sendEvent("start_filter_bubble", "bubble")

		// Save state only if state change
		if !h.state.FilterBubbleRunning {
			h.state.FilterBubbleRunning = true
			err = h.stateUsecase.Update(ctx, h.state)
			if err != nil {
				return err
			}
		}

		log.Info("Start filter bubble successfully")
	} else {
		log.Info("Filter bubble not started because of state not permit it")
		return ErrRelayCanNotStart
	}

	return nil

}

// StopFilterBubble permit to stop filter bubble
// It will try while is not stopped
func (h *TFPHandler) StopFilterBubble(ctx context.Context) error {
	log.Debug("Stop filter bubble")

	err := h.relayBubbleFilter.Off()
	if err != nil {
		return err
	}

	h.sendEvent("start_filter_bubble", "bubble")

	// Save state only if state change
	if h.state.FilterBubbleRunning {
		h.state.FilterBubbleRunning = false
		err = h.stateUsecase.Update(ctx, h.state)
		if err != nil {
			return err
		}
	}

	log.Info("Stop filter bubble successfully")

	return err
}

/*
// HandleRelay manage the relay state
func (h *TFPHandler) HandleRelay() {

	// Handle ermergency stop
	h.eventer.On("stateChange", func(data interface{}) {
		event := data.(string)

		log.Debugf("Receive event %s", event)

		// Stop relais
		if h.state.IsEmergencyStopped {
			h.StopRelais(context.Background())
		}
	})

	// Handle security
	h.eventer.On("stateChange", func(data interface{}) {
		event := data.(string)

		log.Debugf("Receive event %s", event)

		// Stop pump
		if h.state.IsSecurity && !h.state.IsDisableSecurity {
			err := h.StopPondPump(context.Background())
			if err != nil {
				log.Errorf("Failed to stop Pond pump: %s", err.Error())
			}

			err = h.StopWaterfallPump(context.Background())
			if err != nil {
				log.Errorf("Failed to stop waterfall pump: %s", err.Error())
			}
		}
	})

	// Handle relay
	h.eventer.On("stateChange", func(data interface{}) {
		event := data.(string)

		log.Debugf("Receive event %s", event)

		if event == "initTFP" || event == "reconnectTFP" {
			// Manage pond pump
			if h.state.PondPumpRunning {
				err := h.StartPondPump(context.Background())
				if err != nil {
					log.Errorf("Failed to start pond pump: %s", err.Error())
				}
			} else {
				err := h.StopPondPump(context.Background())
				if err != nil {
					log.Errorf("Failed to stop pond pump: %s", err.Error())
				}
			}

			// Manage UVC1
			if h.state.UVC1Running {
				err := h.StartUVC1(context.Background())
				if err != nil {
					log.Errorf("Failed to start UVC1: %s", err.Error())
				}
			} else {
				err := h.StopUVC1(context.Background())
				if err != nil {
					log.Errorf("Failed to sop UVC1: %s", err.Error())
				}
			}

			// Manage UVC2
			if h.state.UVC2Running {
				err := h.StartUVC2(context.Background())
				if err != nil {
					log.Errorf("Failed to start UVC2: %s", err.Error())
				}
			} else {
				err := h.StopUVC2(context.Background())
				if err != nil {
					log.Errorf("Failed to stop UVC2: %s", err.Error())
				}
			}

			// Manage waterfall pump
			if h.state.WaterfallPumpRunning {
				err := h.StartWaterfallPump(context.Background())
				if err != nil {
					log.Errorf("Failed to start Waterfall pump: %s", err.Error())
				}
			} else {
				err := h.StopWaterfallPump(context.Background())
				if err != nil {
					log.Errorf("Failed to stop waterfall pump: %s", err.Error())
				}
			}

			// Manage pond bubble
			if h.state.PondBubbleRunning {
				err := h.StartPondBubble(context.Background())
				if err != nil {
					log.Errorf("Failed to start pond bubble: %s", err.Error())
				}
			} else {
				err := h.StopPondBubble(context.Background())
				if err != nil {
					log.Errorf("Failed to stop Pond bubble: %s", err.Error())
				}
			}

			// Manage filter bubble
			if h.state.FilterBubbleRunning {
				err := h.StartFilterBubble(context.Background())
				if err != nil {
					log.Errorf("Failed to start filter bubble: %s", err.Error())
				}
			} else {
				err := h.StopFilterBubble(context.Background())
				if err != nil {
					log.Errorf("Failed to stop filter bubble")
				}
			}
		}
	})
}
*/

// StopRelais stop all relais
func (h *TFPHandler) StopRelais(ctx context.Context) error {
	log.Info("Stop all relais")
	err := h.StopPondPump(ctx)
	if err != nil {
		log.Errorf("Error when stop pond pump: %s", err.Error())
		return err
	}
	err = h.StopWaterfallPump(ctx)
	if err != nil {
		log.Errorf("Error when stop waterfall pump: %s", err.Error())
		return err
	}
	err = h.StopPondBubble(ctx)
	if err != nil {
		log.Errorf("Error when stop pond bubble: %s", err.Error())
		return err
	}
	err = h.StopFilterBubble(ctx)
	if err != nil {
		log.Errorf("Error when stop filter bubble: %s", err.Error())
		return err
	}
	err = h.StopUVC1(ctx)
	if err != nil {
		log.Errorf("Error when stop UVC1: %s", err.Error())
		return err
	}
	err = h.StopUVC2(ctx)
	if err != nil {
		log.Errorf("Error when stop UVC2: %s", err.Error())
		return err
	}

	return nil
}
