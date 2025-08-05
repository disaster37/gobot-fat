package tfpboard

import (
	"context"
	"errors"

	"github.com/disaster37/gobot-fat/helper"
	log "github.com/sirupsen/logrus"
)

var ErrRelayCanNotStart = errors.New("relay can't start because of current state")

func (h *TFPBoard) canStartRelay() bool {
	if !h.state.IsEmergencyStopped && (!h.state.IsSecurity || h.state.IsDisableSecurity) {
		return true
	}
	return false
}

// StartPondPump permit to run pond pump
// The pump start only if no emergency and no security
func (h *TFPBoard) StartPondPump(ctx context.Context) error {
	if h.canStartRelay() {
		log.Debug("Start pond pump")
		err := h.relayPompPond.On()
		if err != nil {
			return err
		}

		// Send event
		helper.SendEvent(ctx, h.eventUsecase, h.name, helper.KindEventStart, "pond_pump")

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
func (h *TFPBoard) StartUVC1(ctx context.Context) error {

	if h.canStartRelay() && h.state.PondPumpRunning {
		log.Debug("Start UVC1")
		err := h.relayUVC1.On()
		if err != nil {
			return err
		}

		// Send event
		helper.SendEvent(ctx, h.eventUsecase, h.name, helper.KindEventStart, "uvc1")

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
func (h *TFPBoard) StartUVC2(ctx context.Context) error {
	if h.canStartRelay() && h.state.PondPumpRunning {
		log.Debug("Start UVC2")
		err := h.relayUVC2.On()
		if err != nil {
			return err
		}

		helper.SendEvent(ctx, h.eventUsecase, h.name, helper.KindEventStart, "uvc2")

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
func (h *TFPBoard) StartPondPumpWithUVC(ctx context.Context) error {
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
func (h *TFPBoard) StopUVC1(ctx context.Context) error {
	log.Debug("Stop UVC1")

	err := h.relayUVC1.Off()
	if err != nil {
		return err
	}

	helper.SendEvent(ctx, h.eventUsecase, h.name, helper.KindEventStop, "uvc1")

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
func (h *TFPBoard) StopUVC2(ctx context.Context) error {
	log.Debug("Stop UVC2")

	err := h.relayUVC2.Off()
	if err != nil {
		return err
	}

	helper.SendEvent(ctx, h.eventUsecase, h.name, helper.KindEventStop, "uvc2")

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
func (h *TFPBoard) StopPondPump(ctx context.Context) error {

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

	helper.SendEvent(ctx, h.eventUsecase, h.name, helper.KindEventStop, "pond_pump")

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
func (h *TFPBoard) StartWaterfallPump(ctx context.Context) error {
	if h.canStartRelay() {
		log.Debug("Start waterfall pump")
		err := h.relayPompWaterfall.On()
		if err != nil {
			return err
		}

		helper.SendEvent(ctx, h.eventUsecase, h.name, helper.KindEventStart, "waterfall_pump")

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
func (h *TFPBoard) StopWaterfallPump(ctx context.Context) error {
	log.Debug("Stop waterfall pump")

	err := h.relayPompWaterfall.Off()
	if err != nil {
		return err
	}

	helper.SendEvent(ctx, h.eventUsecase, h.name, helper.KindEventStop, "watefall_pump")

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
func (h *TFPBoard) StartPondBubble(ctx context.Context) error {
	if !h.state.IsEmergencyStopped {
		log.Debug("Start pond bubble")
		err := h.relayBubblePond.On()
		if err != nil {
			return err
		}

		helper.SendEvent(ctx, h.eventUsecase, h.name, helper.KindEventStart, "pond_bubble")

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
func (h *TFPBoard) StopPondBubble(ctx context.Context) error {
	log.Debug("Stop pond bubble")

	err := h.relayBubblePond.Off()
	if err != nil {
		return err
	}

	helper.SendEvent(ctx, h.eventUsecase, h.name, helper.KindEventStop, "pond_bubble")

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
func (h *TFPBoard) StartFilterBubble(ctx context.Context) error {
	if !h.state.IsEmergencyStopped {
		log.Debug("Start filter bubble")
		err := h.relayBubbleFilter.On()
		if err != nil {
			return err
		}

		helper.SendEvent(ctx, h.eventUsecase, h.name, helper.KindEventStart, "filter_bubble")

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
func (h *TFPBoard) StopFilterBubble(ctx context.Context) error {
	log.Debug("Stop filter bubble")

	err := h.relayBubbleFilter.Off()
	if err != nil {
		return err
	}

	helper.SendEvent(ctx, h.eventUsecase, h.name, helper.KindEventStop, "filter_bubble")

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

// StopRelais stop all relais
func (h *TFPBoard) StopRelais(ctx context.Context) error {
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
