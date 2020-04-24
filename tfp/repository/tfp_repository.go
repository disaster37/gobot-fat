package repository

import (
	"time"

	"github.com/disaster37/gobot-fat/models"
	"github.com/disaster37/gobot-fat/tfp"
	"gobot.io/x/gobot"
)

type tfpRepository struct {
	state   *models.TFPState
	eventer gobot.Eventer
}

// NewTFPRepository instanciate TFPRepository interface
func NewTFPRepository(state *models.TFPState, eventer gobot.Eventer) tfp.Repository {
	return &tfpRepository{
		state:   state,
		eventer: eventer,
	}
}

func (h *ftpRepository) StartPondPump() (bool, error) {
	if !h.state.PondPumpRunning {
		h.state.PondPumpRunning = true
		h.eventer.Publish("stateChange", "pondPumpRunning")
		return true, nil
	}
	return false, nil
}

func (h *ftpRepository) StopPondPump() (bool, error) {
	if h.state.PondPumpRunning {
		h.state.PondPumpRunning = false
		h.eventer.Publish("stateChange", "pondPumpStopped")
		return true, nil
	}
	return false, nil
}

func (h *ftpRepository) StartWaterfallPump() (bool, error) {
	if !h.state.WaterfallPumpRunning {
		h.state.WaterfallPumpRunning = true
		h.eventer.Publish("stateChange", "waterfallPumpRunning")
		return true, nil
	}
	return false, nil
}

func (h *ftpRepository) StopWaterfallPump() (bool, error) {
	if h.state.WaterfallPumpRunning {
		h.state.WaterfallPumpRunning = false
		h.eventer.Publish("stateChange", "waterfallPumpStopped")
		return true, nil
	}
	return false, nil
}

func (h *ftpRepository) StartUVC1() (bool, error) {
	if !h.state.UVC1Running {
		h.state.UVC1Running = true
		h.eventer.Publish("stateChange", "uvc1Running")
		return true, nil
	}
	return false, nil
}

func (h *ftpRepository) StopUVC1() (bool, error) {
	if h.state.UVC1Running {
		h.state.UVC1Running = false
		h.eventer.Publish("stateChange", "uvc1Stopped")
		return true, nil
	}
	return false, nil
}

func (h *ftpRepository) StartUVC2() (bool, error) {
	if !h.state.UVC2Running {
		h.state.UVC2Running = true
		h.eventer.Publish("stateChange", "uvc2Running")
		return true, nil
	}
	return false, nil
}

func (h *ftpRepository) StopUVC2() (bool, error) {
	if h.state.UVC2Running {
		h.state.UVC2Running = false
		h.eventer.Publish("stateChange", "uvc2Stopped")
		return true, nil
	}
	return false, nil
}

func (h *ftpRepository) StartPondBubble() (bool, error) {
	if !h.state.PondBubbleRunning {
		h.state.PondBubbleRunning = true
		h.eventer.Publish("stateChange", "pondBubbleRunning")
		return true, nil
	}
	return false, nil
}

func (h *ftpRepository) StopPondBubble() (bool, error) {
	if h.state.PondBubbleRunning {
		h.state.PondBubbleRunning = false
		h.eventer.Publish("stateChange", "pondBubbleStopped")
		return true, nil
	}
	return false, nil
}

func (h *ftpRepository) StartFilterBubble() (bool, error) {
	if !h.state.FilterBubbleRunning {
		h.state.FilterBubbleRunning = true
		h.eventer.Publish("stateChange", "filterBubbleRunning")
		return true, nil
	}
	return false, nil
}

func (h *ftpRepository) StopFilterBubble() (bool, error) {
	if h.state.FilterBubbleRunning {
		h.state.FilterBubbleRunning = false
		h.eventer.Publish("stateChange", "filterBubbleStopped")
		return true, nil
	}
	return false, nil
}

// CanStartRelay handle if relay can be started
// Only if not emergency stop and not security
func (h *tfpRepository) CanStartRelay() bool {
	if !h.state.IsEmergencyStopped && (!h.state.IsSecurity || h.state.IsDisableSecurity) {
		return true
	}
	return false
}

func (h *tfpRepository) String() string {
	return h.state.String()
}

func (h *tfpRepository) State() *models.TFPState {
	return h.state
}
