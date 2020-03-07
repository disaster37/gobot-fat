package models

import (
	"time"
)

// DFPConfig contain config data for Drum Filter Pond
type DFPConfig struct {
	// ForceWashingDurationInMinutes is the maximum time in minutes to wait before force a washing since last washing
	ForceWashingDurationInMinutes int `json:"force_washing_duration_in_minutes" validate:"required"`

	// ForceWashingDurationWhenFrozenInMinutes is the maximum time in minutes to wait before force a washing since last washing when tempeture frozen
	ForceWashingDurationWhenFrozenInMinutes int `json:"force_washing_duration_when_frozen_in_minutes" validate:"required"`

	// TemperatureThresholdWhenFrozen is the tempeture in degrees to consider that is frozen
	TemperatureThresholdWhenFrozen int `json:"temperature_threshold_when_frozen" validate:"required"`

	// WaitTimeBetweenWashingInMinutes is the minimal time in minutes to wait before start a washing since last washing
	WaitTimeBetweenWashingInMinutes int `json:"wait_time_between_washing_in_minutes" validate:"required"`

	// LastWashingTime is the datetime of the last washing
	LastWashingTime time.Time `json:"last_washing_time" validate:"required"`

	// WashingDurationInSeconds is the time in seconds of washing cycle
	WashingDurationInSeconds int `json:"washing_duration_in_seconds" validate:"required"`

	// StartWashingPumpBeforeWashingInSeconds is the time in seconds witch we start washing pump before run washing cycle
	StartWashingPumpBeforeWashingInSeconds int `json:"start_washing_pump_before_washing_in_seconds" validate:"required"`

	Version int64 `json:"version"`
	updated time.Time `json:"updated"`
}
