package usecase

import (
	"context"
	"time"

	"github.com/disaster37/gobot-fat/models"
	"github.com/disaster37/gobot-fat/tfp"
	tfpconfig "github.com/disaster37/gobot-fat/tfp_config"
	tfpstate "github.com/disaster37/gobot-fat/tfp_state"
	log "github.com/sirupsen/logrus"
)

type tfpUsecase struct {
	tfp            tfp.Gobot
	config         tfpconfig.Usecase
	state          tfpstate.Usecase
	contextTimeout time.Duration
}

// NewTFPUsecase will create new tfpUsecase object of tfp.Usecase interface
func NewTFPUsecase(handler tfp.Gobot, config tfpconfig.Usecase) tfp.Usecase {
	return &tfpUsecase{
		tfp:    handler,
		config: config,
	}
}

func (h *tfpUsecase) PondPump(c context.Context, status bool) error {

	ctx, cancel := context.WithTimeout(c, h.contextTimeout)
	defer cancel()

	if status {
		// Start ond pump
		log.Debugf("Start pond pump is required by API")
		err := h.tfp.StartPondPump(ctx)
		if err != nil {
			return err
		}
	} else {
		// Unset stop
		log.Debugf("Stop pond pump is required by API")
		err := h.tfp.StopPondPump(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *tfpUsecase) WaterfallPump(c context.Context, status bool) error {

	ctx, cancel := context.WithTimeout(c, h.contextTimeout)
	defer cancel()

	if status {
		// Start ond pump
		log.Debugf("Start waterfall pump is required by API")
		err := h.tfp.StartWaterfallPump(ctx)
		if err != nil {
			return err
		}
	} else {
		// Unset stop
		log.Debugf("Stop waterfall pump is required by API")
		err := h.tfp.StopWaterfallPump(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *tfpUsecase) UVC1(c context.Context, status bool) error {
	ctx, cancel := context.WithTimeout(c, h.contextTimeout)
	defer cancel()

	if status {
		// Start ond pump
		log.Debugf("Start UVC1 is required by API")
		err := h.tfp.StartUVC1(ctx)
		if err != nil {
			return err
		}
	} else {
		// Unset stop
		log.Debugf("Stop UVC1 is required by API")
		err := h.tfp.StopUVC1(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *tfpUsecase) UVC2(c context.Context, status bool) error {
	ctx, cancel := context.WithTimeout(c, h.contextTimeout)
	defer cancel()

	if status {
		// Start ond pump
		log.Debugf("Start UVC2 is required by API")
		err := h.tfp.StartUVC2(ctx)
		if err != nil {
			return err
		}
	} else {
		// Unset stop
		log.Debugf("Stop UVC2 is required by API")
		err := h.tfp.StopUVC2(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *tfpUsecase) PondBubble(c context.Context, status bool) error {
	ctx, cancel := context.WithTimeout(c, h.contextTimeout)
	defer cancel()

	if status {
		// Start ond pump
		log.Debugf("Start pond bubble is required by API")
		err := h.tfp.StartPondBubble(ctx)
		if err != nil {
			return err
		}
	} else {
		// Unset stop
		log.Debugf("Stop pond bubble is required by API")
		err := h.tfp.StopPondBubble(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

func (h *tfpUsecase) FilterBubble(c context.Context, status bool) error {
	ctx, cancel := context.WithTimeout(c, h.contextTimeout)
	defer cancel()

	if status {
		log.Debugf("Start filter bubble is required by API")
		err := h.tfp.StartFilterBubble(ctx)
		if err != nil {
			return err
		}
	} else {
		log.Debugf("Stop filter bubble is required by API")
		err := h.tfp.StopFilterBubble(ctx)
		if err != nil {
			return err
		}
	}

	return nil
}

// GetState return the current state of TFP
func (h *tfpUsecase) GetState(c context.Context) (models.TFPState, error) {
	return h.tfp.State(), nil
}

// StartRobot start the rebot that manage the TFP
func (h *tfpUsecase) StartRobot(c context.Context) error {
	h.tfp.Start()

	return nil
}

// StopRobot stop the robot that manage the TFP
func (h *tfpUsecase) StopRobot(c context.Context) error {
	return h.tfp.Stop()
}

// UVC1BlisterNew update the date when blister changed
func (h *tfpUsecase) UVC1BlisterNew(c context.Context) error {
	state, err := h.state.Get(c)
	if err != nil {
		return err
	}

	state.UVC1BlisterTime = time.Now()

	err = h.state.Update(c, state)
	if err != nil {
		return err
	}

	return nil
}

// UVC2BlisterNew update the date when blister changed
func (h *tfpUsecase) UVC2BlisterNew(ctx context.Context) error {
	state, err := h.state.Get(ctx)
	if err != nil {
		return err
	}

	state.UVC2BlisterTime = time.Now()

	err = h.state.Update(ctx, state)
	if err != nil {
		return err
	}

	return nil
}

// UVC2BlisterNew update the date when blister changed
func (h *tfpUsecase) OzoneBlisterNew(ctx context.Context) error {
	state, err := h.state.Get(ctx)
	if err != nil {
		return err
	}

	state.OzoneBlisterTime = time.Now()

	err = h.state.Update(ctx, state)
	if err != nil {
		return err
	}

	return nil
}
