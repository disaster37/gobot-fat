package models

import (
	"github.com/jinzhu/gorm"
)

type Model interface {
	IsMoreRecentThan(data *gorm.Model) bool
}

// ModelGeneric is generic type
type ModelGeneric struct {
	gorm.Model

	// Version of configuration
	Version int64 `json:"version" gorm:"column:version;type:bigint" validate:"required"`
}

// IsMoreRecentThan return true if current object is more recent.
func (h *ModelGeneric) IsMoreRecentThan(data *gorm.Model) bool {
	return h.UpdatedAt.After(data.UpdatedAt)
}
