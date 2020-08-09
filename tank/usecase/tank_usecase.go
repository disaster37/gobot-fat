package usecase

import (
	"context"
	"time"

	"github.com/disaster37/gobot-fat/models"
	"github.com/disaster37/gobot-fat/tank"
	log "github.com/sirupsen/logrus"
)

type tankUsecase struct {
	tanks          []tank.Board
	contextTimeout time.Duration
}

// NewTankUsecase will create new tankUsecase object of tank.Usecase interface
func NewTankUsecase(handlers []tank.Board, timeout time.Duration) tank.Usecase {
	return &tankUsecase{
		tanks:          handlers,
		contextTimeout: timeout,
	}
}

// Tanks return the current watter level on tank
func (h *tankUsecase) Tanks(ctx context.Context) (values map[string]*models.Tank, err error) {

	values = make(map[string]*models.Tank)

	for _, tank := range h.tanks {
		data, err := tank.GetData(ctx)
		if err != nil {
			return values, err
		}

		values[tank.Name()] = data
	}

	return values, err
}

// Tanks return the current watter level on tank
func (h *tankUsecase) Tank(ctx context.Context, name string) (value *models.Tank, err error) {

	log.Debugf("Name: %s", name)

	for _, tank := range h.tanks {
		if tank.Name() == name {
			return tank.GetData(ctx)
		}
	}

	return nil, nil
}
