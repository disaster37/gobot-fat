package models

import (
	"github.com/jinzhu/gorm"
)

// Model is generic type
type ModelGeneric struct {
	gorm.Model

	// Version of configuration
	Version int64 `json:"version" gorm:"column:version;type:bigint" validate:"required"`
}

type Model interface {
	SetVersion(version int64)
	GetVersion() int64
	GetModel() *ModelGeneric
}

// SetVersion permit to set version
func (h *ModelGeneric) SetVersion(version int64) {
	h.Version = version
}

// GetVersion permit to get current version
func (h *ModelGeneric) GetVersion() int64 {
	return h.Version
}

// GetModel return current model
func (h *ModelGeneric) GetModel() *ModelGeneric {
	return h
}
