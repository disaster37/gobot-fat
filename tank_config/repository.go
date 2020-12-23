package tankconfig

import (
	"context"

	"github.com/disaster37/gobot-fat/models"
)

// Repository represent the config repository contract
type Repository interface {
	List(ctx context.Context) ([]*models.TankConfig, error)
	Get(ctx context.Context, name string) (*models.TankConfig, error)
	Update(ctx context.Context, config *models.TankConfig) error
	Create(ctx context.Context, config *models.TankConfig) error
}
