package repository

import (
	"context"
)

// ElasticsearchRepository is a generic repository
type ElasticsearchRepository interface {
	Get(ctx context.Context, id string, data interface{}) error
	List(ctx context.Context, listData interface{}) error
	Update(ctx context.Context, id string, data interface{}) error
	Create(ctx context.Context, id string, data interface{}) error
}

// SQLRepository is a generic repository
type SQLRepository interface {
	Get(ctx context.Context, id string, data interface{}) error
	List(ctx context.Context, listData interface{}) error
	Update(ctx context.Context, data interface{}) error
	Create(ctx context.Context, data interface{}) error
}
