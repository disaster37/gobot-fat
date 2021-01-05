package repository

import (
	"context"

	"github.com/disaster37/gobot-fat/models"
)

// Repository is a generic repository
type Repository interface {
	Get(ctx context.Context, id uint, data models.Model) error
	List(ctx context.Context, listData interface{}) error
	Update(ctx context.Context, data models.Model) error
	Create(ctx context.Context, data models.Model) error
}
