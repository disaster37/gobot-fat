package tfpstate

import (
	"context"

	"github.com/disaster37/gobot-fat/models"
)

// Usecase represent the state's usecase
type Usecase interface {
	Get(ctx context.Context) (*models.TFPState, error)
	Update(ctx context.Context, state *models.TFPState) error
	Create(ctx context.Context, state *models.TFPState) error
	Init(ctx context.Context, state *models.TFPState) error
}
