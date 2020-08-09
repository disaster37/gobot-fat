package dfpstate

import (
	"context"

	"github.com/disaster37/gobot-fat/models"
)

// Usecase represent the state's usecase
type Usecase interface {
	Get(ctx context.Context) (*models.DFPState, error)
	Update(ctx context.Context, state *models.DFPState) error
	Create(ctx context.Context, state *models.DFPState) error
	Init(ctx context.Context, state *models.DFPState) error
}
