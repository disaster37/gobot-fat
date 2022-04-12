package models

import (
	"time"
)

// ModelGeneric is generic type
type ModelGeneric struct {
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time `sql:"index"`

	// Version of configuration
	Version int64 `json:"version" gorm:"column:version;type:bigint" validate:"required"`
}

type Model interface {
	SetVersion(version int64)
	GetVersion() int64
	SetUpdatedAt(date time.Time)
	GetID() uint
	SetID(id uint)
}

// SetVersion permit to set version
func (h *ModelGeneric) SetVersion(version int64) {
	h.Version = version
}

// GetVersion permit to get current version
func (h *ModelGeneric) GetVersion() int64 {
	return h.Version
}

// SetUpdatedAt permit to set updated date
func (h *ModelGeneric) SetUpdatedAt(date time.Time) {
	h.UpdatedAt = date
}
