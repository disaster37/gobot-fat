package tankconfig

import (
	"context"

	"github.com/disaster37/gobot-fat/models"
)

// Usecase represent the config's usecase
type Usecase interface {
	List(ctx context.Context) ([]*models.TankConfig, error)
	Get(ctx context.Context, name string) (*models.TankConfig, error)
	Update(ctx context.Context, config *models.TankConfig) error
	Create(ctx context.Context, config *models.TankConfig) error
	Init(ctx context.Context, config *models.TankConfig) error
}
