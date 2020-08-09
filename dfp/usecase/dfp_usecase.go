package usecase

import (
	"context"
	"time"

	"github.com/disaster37/gobot-fat/dfp"
	"github.com/disaster37/gobot-fat/models"
	log "github.com/sirupsen/logrus"
)

type dfpUsecase struct {
	dfp            dfp.Board
	contextTimeout time.Duration
}

// NewDFPUsecase will create new dfpUsecase object of dfp.Usecase interface
func NewDFPUsecase(handler dfp.Board, timeout time.Duration) dfp.Usecase {
	return &dfpUsecase{
		dfp:            handler,
		contextTimeout: timeout,
	}
}

// Wash will force washing cycle if possible
func (h *dfpUsecase) Wash(c context.Context) error {
	log.Debugf("Washing is required")
	ctx, cancel := context.WithTimeout(c, h.contextTimeout)
	defer cancel()

	return h.dfp.ForceWashing(ctx)

}

// Stop will set stop mode
func (h *dfpUsecase) Stop(c context.Context) error {
	log.Debugf("Stop is required")
	ctx, cancel := context.WithTimeout(c, h.contextTimeout)
	defer cancel()

	return h.dfp.StopDFP(ctx)
}

// Auto witl set auto mode
func (h *dfpUsecase) Auto(c context.Context) error {
	log.Debugf("Auto is required")
	ctx, cancel := context.WithTimeout(c, h.contextTimeout)
	defer cancel()

	return h.dfp.Auto(ctx)
}

// ManualDrum will start / stop the drum motor
func (h *dfpUsecase) ManualDrum(c context.Context, status bool) error {
	ctx, cancel := context.WithTimeout(c, h.contextTimeout)
	defer cancel()

	if status {
		// Start drum
		log.Debugf("Start drum")
		return h.dfp.StartManualDrum(ctx)

	}
	// Stop drum
	log.Debugf("Stop drum")
	return h.dfp.StopManualDrum(ctx)

}

// ManualPump will start / stop the pump
func (h *dfpUsecase) ManualPump(c context.Context, status bool) error {
	ctx, cancel := context.WithTimeout(c, h.contextTimeout)
	defer cancel()

	if status {
		log.Debugf("Start pump")
		return h.dfp.StartManualPump(ctx)

	}
	log.Debugf("Stop pump")
	return h.dfp.StopManualPump(ctx)
}

// GetState return the current state of DFP
func (h *dfpUsecase) GetState(ctx context.Context) (models.DFPState, error) {
	return h.dfp.State(), nil
}
