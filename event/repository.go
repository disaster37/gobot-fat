package event

import (
	"context"

	"github.com/disaster37/gobot-fat/models"
)

// Repository represent the event's repository contract
type Repository interface {
	Fetch(ctx context.Context, from int, size int) (res []*models.Event, nextFrom int, err error)
	GetByID(ctx context.Context, id string) (*models.Event, error)
	Search(ctx context.Context, query map[string]interface{}, minimalScoring float64) ([]*models.Event, error)
	Update(ctx context.Context, object *models.Event) error
	Store(ctx context.Context, object *models.Event) error
	Delete(ctx context.Context, id string) error
}
