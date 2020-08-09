package dfpstate

import (
	"context"

	"github.com/disaster37/gobot-fat/models"
)

// Repository represent the state repository contract
type Repository interface {
	Get(ctx context.Context) (*models.DFPState, error)
	Update(ctx context.Context, config *models.DFPState) error
	Create(ctx context.Context, config *models.DFPState) error
}
