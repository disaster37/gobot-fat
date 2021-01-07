package models

import (
	"encoding/json"
)

// DFPConfig contain config data for Drum Filter Pond
type DFPConfig struct {
	ModelGeneric

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
}

func (h *DFPConfig) String() string {
	str, err := json.Marshal(h)
	if err != nil {
		panic(err)
	}
	return string(str)
}

func (DFPConfig) TableName() string {
	return "dfpconfig"
}
