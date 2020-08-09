package tank

import (
	"context"

	"github.com/disaster37/gobot-fat/models"
)

// Usecase represent the tfp usecase
type Usecase interface {
	Tanks(ctx context.Context) (values map[string]*models.Tank, err error)
	Tank(ctx context.Context, name string) (value *models.Tank, err error)
}
