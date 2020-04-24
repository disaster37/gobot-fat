package usecase

import (
	"context"

	"github.com/disaster37/gobot-fat/models"
	"github.com/disaster37/gobot-fat/tfp"
	"github.com/disaster37/gobot-fat/tfp_config"
	log "github.com/sirupsen/logrus"
)

type tfpUsecase struct {
	tfp    tfp.Gobot
	state  tfp.Repository
	config tfpconfig.Usecase
}

// NewTFPUsecase will create new tfpUsecase object of tfp.Usecase interface
func NewTFPUsecase(handler tfp.Gobot, repo tfp.Repository, config tfpconfig.Usecase) tfp.Usecase {
	return &tfpUsecase{
		tfp:    handler,
		state:  repo,
		config: config,
	}
}

func (h *tfpUsecase) PondPump(ctx context.Context, status bool) error {
	var isUpdate bool
	var err error
	if status {
		// Start ond pump
		log.Debugf("Start pond pump is required by API")
		isUpdate, err = h.state.StartPondPump()
		if err != nil {
			return err
		}
		err = h.tfp.StartPondPump()
	} else {
		// Unset stop
		log.Debugf("Stop pond pump is required by API")
		isUpdate, err = h.state.StopPondPump()
		h.tfp.StopPondPump()
	}

	if err != nil {
		return err
	}
	if isUpdate {
		config, err := h.config.Get(ctx)
		if err != nil {
			return err
		}
		config.PondPumpRunning = status
		err = h.config.Update(ctx, config)
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *tfpUsecase) WaterfallPump(ctx context.Context, status bool) error {
	var isUpdate bool
	var err error
	if status {
		// Start ond pump
		log.Debugf("Start waterfall pump is required by API")
		isUpdate, err = h.state.StartWaterfallPump()
		if err != nil {
			return err
		}
		err = h.tfp.StartWaterfallPump()
	} else {
		// Unset stop
		log.Debugf("Stop waterfall pump is required by API")
		isUpdate, err = h.state.StopWaterfallPump()
		h.tfp.StopWaterfallPump()
	}

	if err != nil {
		return err
	}
	if isUpdate {
		config, err := h.config.Get(ctx)
		if err != nil {
			return err
		}
		config.WaterfallPumpRunning = status
		err = h.config.Update(ctx, config)
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *tfpUsecase) UVC1(ctx context.Context, status bool) error {
	var isUpdate bool
	var err error
	if status {
		// Start ond pump
		log.Debugf("Start UVC1 is required by API")
		isUpdate, err = h.state.StartUVC1()
		if err != nil {
			return err
		}
		err = h.tfp.StartUVC1()
	} else {
		// Unset stop
		log.Debugf("Stop UVC1 is required by API")
		isUpdate, err = h.state.StopUVC1()
		h.tfp.StopUVC1()
	}

	if err != nil {
		return err
	}
	if isUpdate {
		config, err := h.config.Get(ctx)
		if err != nil {
			return err
		}
		config.UVC1Running = status
		err = h.config.Update(ctx, config)
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *tfpUsecase) UVC2(ctx context.Context, status bool) error {
	var isUpdate bool
	var err error
	if status {
		// Start ond pump
		log.Debugf("Start UVC2 is required by API")
		isUpdate, err = h.state.StartUVC2()
		if err != nil {
			return err
		}
		err = h.tfp.StartUVC2()
	} else {
		// Unset stop
		log.Debugf("Stop UVC2 is required by API")
		isUpdate, err = h.state.StopUVC2()
		h.tfp.StopUVC2()
	}

	if err != nil {
		return err
	}
	if isUpdate {
		config, err := h.config.Get(ctx)
		if err != nil {
			return err
		}
		config.UVC2Running = status
		err = h.config.Update(ctx, config)
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *tfpUsecase) PondBubble(ctx context.Context, status bool) error {
	var isUpdate bool
	var err error
	if status {
		// Start ond pump
		log.Debugf("Start pond bubble is required by API")
		isUpdate, err = h.state.StartPondBubble()
		if err != nil {
			return err
		}
		err = h.tfp.StartPondBubble()
	} else {
		// Unset stop
		log.Debugf("Stop pond bubble is required by API")
		isUpdate, err = h.state.StopPondBubble()
		h.tfp.StopPondBubble()
	}

	if err != nil {
		return err
	}
	if isUpdate {
		config, err := h.config.Get(ctx)
		if err != nil {
			return err
		}
		config.PondBubbleRunning = status
		err = h.config.Update(ctx, config)
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *tfpUsecase) FilterBubble(ctx context.Context, status bool) error {
	var isUpdate bool
	var err error
	if status {
		log.Debugf("Start filter bubble is required by API")
		isUpdate, err = h.state.StartFilterBubble()
		if err != nil {
			return err
		}
		err = h.tfp.StartFilterBubble()
	} else {
		log.Debugf("Stop filter bubble is required by API")
		isUpdate, err = h.state.StopFilterBubble()
		h.tfp.StopFilterBubble()
	}

	if err != nil {
		return err
	}
	if isUpdate {
		config, err := h.config.Get(ctx)
		if err != nil {
			return err
		}
		config.FilterBubbleRunning = status
		err = h.config.Update(ctx, config)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetState return the current state of TFP
func (h *tfpUsecase) GetState(ctx context.Context) (*models.TFPState, error) {
	return h.state.State(), nil
}

// StartRobot start the rebot that manage the TFP
func (h *tfpUsecase) StartRobot(ctx context.Context) error {
	h.tfp.Start()

	return nil
}

// StopRobot stop the robot that manage the TFP
func (h *tfpUsecase) StopRobot(ctx context.Context) error {
	return h.tfp.Stop()
}
