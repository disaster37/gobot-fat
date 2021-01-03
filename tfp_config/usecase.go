package tfpconfig

import (
	"context"

	"github.com/disaster37/gobot-fat/models"
	"github.com/disaster37/gobot-fat/usecase"
)

// Usecase represent the config's usecase
type Usecase interface {
	Get(ctx context.Context) (*models.TFPConfig, error)
	Update(ctx context.Context, config *models.TFPConfig) error
	Create(ctx context.Context, config *models.TFPConfig) error
	usecase.Usecase
}
