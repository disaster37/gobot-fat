package event

import (
	"context"

	"github.com/disaster37/gobot-fat/models"
)

// Repository represent the event's repository contract
type Repository interface {
	Fetch(ctx context.Context, cursor string, num int64) (res []*models.Event, nextCursor string, err error)
	GetByID(ctx context.Context, id string) (*models.Event, error)
	GetByLast(ctx context.Context, eventType string, kind string, sourceID string) (*models.Event, error)
	Update(ctx context.Context, object *models.Event) error
	Store(ctx context.Context, object *models.Event) error
	Delete(ctx context.Context, id string) error
}
