package repository

import (
	"context"

	"github.com/disaster37/gobot-fat/models"
	"github.com/disaster37/gobot-fat/tfp"
	"github.com/disaster37/gobot-fat/tfp_config"
	"gobot.io/x/gobot"
)

type tfpRepository struct {
	state   *models.TFPState
	eventer gobot.Eventer
	config  tfpconfig.Usecase
}

// NewTFPRepository instanciate TFPRepository interface
func NewTFPRepository(state *models.TFPState, eventer gobot.Eventer, config tfpconfig.Usecase) tfp.Repository {
	return &tfpRepository{
		state:   state,
		eventer: eventer,
		config:  config,
	}
}

func (h *tfpRepository) StartPondPump() error {
	if !h.state.PondPumpRunning {
		h.state.PondPumpRunning = true
		h.eventer.Publish("stateChange", "pondPumpRunning")

		// Update config
		err := h.updateConfig()

		return err
	}
	return nil
}

func (h *tfpRepository) StopPondPump() error {
	if h.state.PondPumpRunning {
		h.state.PondPumpRunning = false
		h.eventer.Publish("stateChange", "pondPumpStopped")

		// Update config
		err := h.updateConfig()

		return err
	}
	return nil
}

func (h *tfpRepository) StartWaterfallPump() error {
	if !h.state.WaterfallPumpRunning {
		h.state.WaterfallPumpRunning = true
		h.eventer.Publish("stateChange", "waterfallPumpRunning")

		// Update config
		err := h.updateConfig()

		return err
	}
	return nil
}

func (h *tfpRepository) StopWaterfallPump() error {
	if h.state.WaterfallPumpRunning {
		h.state.WaterfallPumpRunning = false
		h.eventer.Publish("stateChange", "waterfallPumpStopped")

		// Update config
		err := h.updateConfig()

		return err
	}
	return nil
}

func (h *tfpRepository) StartUVC1() error {
	if !h.state.UVC1Running {
		h.state.UVC1Running = true
		h.eventer.Publish("stateChange", "uvc1Running")

		// Update config
		err := h.updateConfig()

		return err
	}
	return nil
}

func (h *tfpRepository) StopUVC1() error {
	if h.state.UVC1Running {
		h.state.UVC1Running = false
		h.eventer.Publish("stateChange", "uvc1Stopped")

		// Update config
		err := h.updateConfig()

		return err
	}
	return nil
}

func (h *tfpRepository) StartUVC2() error {
	if !h.state.UVC2Running {
		h.state.UVC2Running = true
		h.eventer.Publish("stateChange", "uvc2Running")

		// Update config
		err := h.updateConfig()

		return err
	}
	return nil
}

func (h *tfpRepository) StopUVC2() error {
	if h.state.UVC2Running {
		h.state.UVC2Running = false
		h.eventer.Publish("stateChange", "uvc2Stopped")

		// Update config
		err := h.updateConfig()

		return err
	}
	return nil
}

func (h *tfpRepository) StartPondBubble() error {
	if !h.state.PondBubbleRunning {
		h.state.PondBubbleRunning = true
		h.eventer.Publish("stateChange", "pondBubbleRunning")

		// Update config
		err := h.updateConfig()

		return err
	}
	return nil
}

func (h *tfpRepository) StopPondBubble() error {
	if h.state.PondBubbleRunning {
		h.state.PondBubbleRunning = false
		h.eventer.Publish("stateChange", "pondBubbleStopped")

		// Update config
		err := h.updateConfig()

		return err
	}
	return nil
}

func (h *tfpRepository) StartFilterBubble() error {
	if !h.state.FilterBubbleRunning {
		h.state.FilterBubbleRunning = true
		h.eventer.Publish("stateChange", "filterBubbleRunning")

		// Update config
		err := h.updateConfig()

		return err
	}
	return nil
}

func (h *tfpRepository) StopFilterBubble() error {
	if h.state.FilterBubbleRunning {
		h.state.FilterBubbleRunning = false
		h.eventer.Publish("stateChange", "filterBubbleStopped")

		// Update config
		err := h.updateConfig()

		return err
	}
	return nil
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

func (h *tfpRepository) updateConfig() error {
	ctx := context.Background()
	config, err := h.config.Get(ctx)
	if err != nil {
		return err
	}

	config.PondPumpRunning = h.State().PondPumpRunning
	config.UVC1Running = h.State().UVC1Running
	config.UVC2Running = h.State().UVC2Running
	config.PondBubbleRunning = h.State().PondBubbleRunning
	config.FilterBubbleRunning = h.State().FilterBubbleRunning
	config.WaterfallPumpRunning = h.State().WaterfallPumpRunning

	err = h.config.Update(ctx, config)
	if err != nil {
		return err
	}

	return nil
}
