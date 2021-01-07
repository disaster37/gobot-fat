package repository

import (
	"context"

	"github.com/pkg/errors"
)

// ErrRecordNotFoundError is error when record not found on repository
var ErrRecordNotFoundError error = errors.New("Record not found")

// Repository is a generic repository
type Repository interface {
	Get(ctx context.Context, id uint, data interface{}) error
	List(ctx context.Context, listData interface{}) error
	Update(ctx context.Context, data interface{}) error
	Create(ctx context.Context, data interface{}) error
}

// IsRecordNotFoundError return true if current error is because of record not found on repository
func IsRecordNotFoundError(err error) bool {
	if err == ErrRecordNotFoundError {
		return true
	}
	return false
}
