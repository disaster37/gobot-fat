package dfpconfig

import (
	"context"

	"github.com/disaster37/gobot-fat/models"
)

// Repository represent the config repository contract
type Repository interface {
	Get(ctx context.Context) (*models.DFPConfig, error)
	Update(ctx context.Context, config *models.DFPConfig) error
	Create(ctx context.Context, config *models.DFPConfig) error
}
