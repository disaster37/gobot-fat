package usecase

import (
	"context"
	"time"

	"github.com/disaster37/gobot-fat/models"
	"github.com/disaster37/gobot-fat/tfp"
	"github.com/disaster37/gobot-fat/tfpconfig"
	"github.com/disaster37/gobot-fat/tfpstate"
	"github.com/disaster37/gobot-fat/usecase"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

const (
	blisterUVC1  = "UVC1"
	blisterUVC2  = "UVC2"
	blisterOzone = "Ozone"
)

type tfpUsecase struct {
	tfp            tfp.Board
	config         usecase.UsecaseCRUD
	state          usecase.UsecaseCRUD
	contextTimeout time.Duration
}

// NewTFPUsecase will create new tfpUsecase object of tfp.Usecase interface
func NewTFPUsecase(handler tfp.Board, config usecase.UsecaseCRUD, state usecase.UsecaseCRUD, timeout time.Duration) tfp.Usecase {
	return &tfpUsecase{
		tfp:            handler,
		config:         config,
		contextTimeout: timeout,
		state:          state,
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
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*30)
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

	state := h.tfp.State()

	// Reflect waterfall auto on state for hass usecase
	state.IsWaterfallAuto = h.tfp.Config().IsWaterfallAuto

	return state, nil
}

// UVC1BlisterNew update the date when blister changed
func (h *tfpUsecase) UVC1BlisterNew(ctx context.Context) error {
	return h.blisterNew(ctx, blisterUVC1)
}

// UVC2BlisterNew update the date when blister changed
func (h *tfpUsecase) UVC2BlisterNew(ctx context.Context) error {
	return h.blisterNew(ctx, blisterUVC2)
}

// UVC2BlisterNew update the date when blister changed
func (h *tfpUsecase) OzoneBlisterNew(ctx context.Context) error {
	return h.blisterNew(ctx, blisterOzone)
}

// WaterfallAuto permit to enable or disable the waterfall auto
func (h *tfpUsecase) WaterfallAuto(c context.Context, status bool) error {
	ctx, cancel := context.WithTimeout(c, h.contextTimeout)
	defer cancel()

	if status {
		log.Debugf("Enable waterfall auto is required by API")
		return h.waterfallAuto(ctx, status)
	} else {
		log.Debugf("Disable waterfall auto required by API")
		return h.waterfallAuto(ctx, status)
	}

}

func (h *tfpUsecase) blisterNew(ctx context.Context, blisterName string) error {
	state := &models.TFPState{}
	if err := h.state.Get(ctx, tfpstate.ID, state); err != nil {
		return err
	}

	config := &models.TFPConfig{}
	if err := h.config.Get(ctx, tfpconfig.ID, config); err != nil {
		return err
	}

	switch blisterName {
	case blisterUVC1:
		config.UVC1BlisterTime = time.Now()
		state.UVC1BlisterNbHour = 0
		break
	case blisterUVC2:
		config.UVC2BlisterTime = time.Now()
		state.UVC2BlisterNbHour = 0
		break
	case blisterOzone:
		config.OzoneBlisterTime = time.Now()
		state.OzoneBlisterNbHour = 0
		break
	default:
		return errors.Errorf("Blister %s not found", blisterName)

	}

	if err := h.state.Update(ctx, state); err != nil {
		return err
	}

	if err := h.config.Update(ctx, config); err != nil {
		return err
	}

	return nil
}

// GetIO return the current IO of DFP
func (h *tfpUsecase) GetIO(ctx context.Context) (models.TFPIO, error) {
	return h.tfp.IO(), nil
}

func (h *tfpUsecase) waterfallAuto(ctx context.Context, state bool) error {

	config := &models.TFPConfig{}
	if err := h.config.Get(ctx, tfpconfig.ID, config); err != nil {
		return err
	}

	if config.IsWaterfallAuto != state {
		config.IsWaterfallAuto = state
		if err := h.config.Update(ctx, config); err != nil {
			return err
		}

	}

	return nil
}
