package usecase

import (
	"context"
	"time"

	"github.com/disaster37/gobot-fat/tank"
)

type tankUsecase struct {
	tank           tank.Board
	contextTimeout time.Duration
}

// NewTankUsecase will create new tankUsecase object of tank.Usecase interface
func NewTankUsecase(handler tank.Board, timeout time.Duration) tank.Usecase {
	return &tankUsecase{
		tank:           handler,
		contextTimeout: timeout,
	}
}

// Level return the current watter level on tank
func (h *tankUsecase) Level(ctx context.Context) (level int, err error) {
	return h.tank.Level(ctx)
}

// Volume return the current watter volume in liter on tank
func (h *tankUsecase) Volume(ctx context.Context) (volume int, err error) {

	level, err := h.tank.Level(ctx)
	if err != nil {
		return 0, err
	}

	// 50 liter per cm
	return 50 * level, nil
}
