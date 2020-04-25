package usecase

import (
	"context"
	"time"

	"github.com/disaster37/gobot-fat/models"
	"github.com/disaster37/gobot-fat/tfp"
	tfpconfig "github.com/disaster37/gobot-fat/tfp_config"
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
	var err error
	if status {
		// Start ond pump
		log.Debugf("Start pond pump is required by API")
		err = h.state.StartPondPump()
		if err != nil {
			return err
		}
		err = h.tfp.StartPondPump()
	} else {
		// Unset stop
		log.Debugf("Stop pond pump is required by API")
		err = h.state.StopPondPump()
		h.tfp.StopPondPump()
	}

	if err != nil {
		return err
	}

	return nil
}

func (h *tfpUsecase) WaterfallPump(ctx context.Context, status bool) error {
	var err error
	if status {
		// Start ond pump
		log.Debugf("Start waterfall pump is required by API")
		err = h.state.StartWaterfallPump()
		if err != nil {
			return err
		}
		err = h.tfp.StartWaterfallPump()
	} else {
		// Unset stop
		log.Debugf("Stop waterfall pump is required by API")
		err = h.state.StopWaterfallPump()
		h.tfp.StopWaterfallPump()
	}

	if err != nil {
		return err
	}

	return nil
}

func (h *tfpUsecase) UVC1(ctx context.Context, status bool) error {
	var err error
	if status {
		// Start ond pump
		log.Debugf("Start UVC1 is required by API")
		err = h.state.StartUVC1()
		if err != nil {
			return err
		}
		err = h.tfp.StartUVC1()
	} else {
		// Unset stop
		log.Debugf("Stop UVC1 is required by API")
		err = h.state.StopUVC1()
		h.tfp.StopUVC1()
	}

	if err != nil {
		return err
	}

	return nil
}

func (h *tfpUsecase) UVC2(ctx context.Context, status bool) error {
	var err error
	if status {
		// Start ond pump
		log.Debugf("Start UVC2 is required by API")
		err = h.state.StartUVC2()
		if err != nil {
			return err
		}
		err = h.tfp.StartUVC2()
	} else {
		// Unset stop
		log.Debugf("Stop UVC2 is required by API")
		err = h.state.StopUVC2()
		h.tfp.StopUVC2()
	}

	if err != nil {
		return err
	}

	return nil
}

func (h *tfpUsecase) PondBubble(ctx context.Context, status bool) error {
	var err error
	if status {
		// Start ond pump
		log.Debugf("Start pond bubble is required by API")
		err = h.state.StartPondBubble()
		if err != nil {
			return err
		}
		err = h.tfp.StartPondBubble()
	} else {
		// Unset stop
		log.Debugf("Stop pond bubble is required by API")
		err = h.state.StopPondBubble()
		h.tfp.StopPondBubble()
	}

	if err != nil {
		return err
	}

	return nil
}

func (h *tfpUsecase) FilterBubble(ctx context.Context, status bool) error {
	var err error
	if status {
		log.Debugf("Start filter bubble is required by API")
		err = h.state.StartFilterBubble()
		if err != nil {
			return err
		}
		err = h.tfp.StartFilterBubble()
	} else {
		log.Debugf("Stop filter bubble is required by API")
		err = h.state.StopFilterBubble()
		h.tfp.StopFilterBubble()
	}

	if err != nil {
		return err
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

// UVC1BlisterStatus return true if UVC1 blister not greather than the max time of use
func (h *tfpUsecase) UVC1BlisterStatus(ctx context.Context) (bool, error) {
	// Get actual config
	config, err := h.config.Get(ctx)
	if err != nil {
		return false, err
	}

	if int64(time.Since(config.UVC1BlisterTime).Hours()) < config.UVC1BlisterMaxTime {
		return true, nil
	}

	return false, nil
}

// UVC2BlisterStatus return true if UVC2 blister not greather than the max time of use
func (h *tfpUsecase) UVC2BlisterStatus(ctx context.Context) (bool, error) {
	// Get actual config
	config, err := h.config.Get(ctx)
	if err != nil {
		return false, err
	}

	if int64(time.Since(config.UVC2BlisterTime).Hours()) < config.UVC2BlisterMaxTime {
		return true, nil
	}

	return false, nil
}

// UVC1BlisterNew update the date when blister changed
func (h *tfpUsecase) UVC1BlisterNew(ctx context.Context) error {
	config, err := h.config.Get(ctx)
	if err != nil {
		return err
	}

	config.UVC1BlisterTime = time.Now()

	err = h.config.Update(ctx, config)
	if err != nil {
		return err
	}

	return nil
}

// UVC2BlisterNew update the date when blister changed
func (h *tfpUsecase) UVC2BlisterNew(ctx context.Context) error {
	config, err := h.config.Get(ctx)
	if err != nil {
		return err
	}

	config.UVC2BlisterTime = time.Now()

	err = h.config.Update(ctx, config)
	if err != nil {
		return err
	}

	return nil
}
