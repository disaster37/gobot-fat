package models

import (
	"time"
)

// Event contain data event
type Event struct {
	SourceID   string    `json:"source_id" validate:"required"`
	SourceName string    `json:"source_name" validate:"required"`
	Timestamp  time.Time `json:"timestamp" validate:"required"`
	EventType  string    `json:"type" validate:"required"`
	Kind       string    `json:"kind" validate:"required"`
	Duration   uint64    `json:"duration,omitempty"`
	Data       int64     `json:"data,omitempty"`
}
