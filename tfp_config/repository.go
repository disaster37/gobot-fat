package tfpconfig

import (
	"context"

	"github.com/disaster37/gobot-fat/models"
)

// Repository represent the config repository contract
type Repository interface {
	Get(ctx context.Context) (*models.TFPConfig, error)
	Update(ctx context.Context, config *models.TFPConfig) error
	Create(ctx context.Context, config *models.TFPConfig) error
}
