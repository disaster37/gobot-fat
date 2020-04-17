package models

import (
	"encoding/json"
	"time"

	"github.com/jinzhu/gorm"
)

// DFPConfig contain config data for Drum Filter Pond
type DFPConfig struct {
	gorm.Model

	// ForceWashingDuration is the maximum time in minutes to wait before force a washing since last washing
	ForceWashingDuration int `json:"force_washing_duration" validate:"required" gorm:"column:force_washing_duration;type:bigint" validate:"required"`

	// ForceWashingDurationWhenFrozen is the maximum time in minutes to wait before force a washing since last washing when tempeture frozen
	ForceWashingDurationWhenFrozen int `json:"force_washing_duration_when_frozen" gorm:"column:force_washing_duration_when_frozen;type:bigint" validate:"required"`

	// TemperatureThresholdWhenFrozen is the tempeture in degrees to consider that is frozen
	TemperatureThresholdWhenFrozen int `json:"temperature_threshold_when_frozen" gorm:"column:temperature_threshold_when_frozen;type:bigint" validate:"required"`

	// WaitTimeBetweenWashing is the minimal time in second to wait before start a washing since last washing
	WaitTimeBetweenWashing int `json:"wait_time_between_washing" gorm:"column:wait_time_between_washing;type:bigint" validate:"required"`

	// WashingDuration is the time in seconds of washing cycle
	WashingDuration int `json:"washing_duration" validate:"required" gorm:"column:washing_duration;type:bigint" validate:"required"`

	// StartWashingPumpBeforeWashing is the time in seconds witch we start washing pump before run washing cycle
	StartWashingPumpBeforeWashing int `json:"start_washing_pump_before_washing" gorm:"column:start_washing_pump_before_washing;type:bigint" validate:"required"`

	// Stopped is true if DFP must be stopped
	Stopped bool `json:"stopped" gorm:"column:stopped" validate:"required"`

	// Emergency stopped is true if DFP must be stopped
	EmergencyStopped bool `json:"emergency_stopped" gorm:"column:emergency_stopped" validate:"required"`

	// Auto is true if DFP is in auto mode
	Auto bool `json:"auto" gorm:"column:auto" validate:"required"`

	// SecurityDisabled is true if security is disabled
	SecurityDisabled bool `json:"security_disabled" gorm:"column:security_disabled" validate:"required"`

	// LastWashing is the date of last washing time
	LastWashing time.Time `json:"last_washing" gorm:"column:lastwashing;type:datetime"`

	// Version of configuration
	Version int64 `json:"version" gorm:"column:version;type:bigint" validate:"required"`
}

func (h *DFPConfig) String() string {
	str, err := json.Marshal(h)
	if err != nil {
		panic(err)
	}
	return string(str)
}
