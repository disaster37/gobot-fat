package models

import (
	"encoding/json"
	"time"
)

// Event contain data event
type Event struct {
	ModelGeneric

	ID                      uint      `jsonapi:"primary,events" gorm:"primary_key"`
	SourceID                string    `json:"source_id" validate:"required"`
	SourceName              string    `json:"source_name" validate:"required"`
	Timestamp               time.Time `json:"timestamp" validate:"required"`
	EventType               string    `json:"type" validate:"required"`
	EventKind               string    `json:"kind" validate:"required"`
	Temperature             float64   `json:"temperature,omitempty"`
	Humidity                float64   `json:"humidity,omitempty"`
	Duration                int64     `json:"duration,omitempty"`
	DurationFromLastWashing int64     `json:"duration_from_last,omitempty"`
	Level                   int64     `json:"level,omitempty"`
}

func (h *Event) String() string {
	str, err := json.Marshal(h)
	if err != nil {
		panic(err)
	}
	return string(str)
}

func (h *Event) SetID(id uint) {
	h.ID = id
}

func (h *Event) GetID() uint {
	return h.ID
}
