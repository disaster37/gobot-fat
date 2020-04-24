package models

import (
	"encoding/json"
	"time"

	"github.com/jinzhu/gorm"
)

// TFPConfig contain config data for Technical Filter Pond
type TFPConfig struct {
	gorm.Model

	// EmergencyStopped is true if DFP must be stopped
	EmergencyStopped bool `json:"emergency_stopped" gorm:"column:emergency_stopped" validate:"required"`

	// UVC1Running is true if UVC1 running
	UVC1Running bool `json:"uvc1_running" gorm:"column:uvc1_running" validate:"required"`

	// UVC2Running is true if UVC2 running
	UVC2Running bool `json:"uvc2_running" gorm:"column:uvc2_running" validate:"required"`

	// PondPumpRunning is true if pond pump running
	PondPumpRunning bool `json:"pond_pump_running" gorm:"column:pond_pump_running" validate:"required"`

	// WaterfallPumpRunning is true if Waterfall pump running
	WaterfallPumpRunning bool `json:"waterfall_pump_running" gorm:"column:waterfall_pump_running" validate:"required"`

	// PondBubbleRunning is true if pond bubble running
	PondBubbleRunning bool `json:"pond_bubble_running" gorm:"column:pond_bubble_running" validate:"required"`

	// FilterBubbleRunning is true if filter bubble running
	FilterBubbleRunning bool `json:"filter_bubble_running" gorm:"column:filter_bubble_running" validate:"required"`

	// SecurityDisabled is true if security is disabled
	SecurityDisabled bool `json:"security_disabled" gorm:"column:security_disabled" validate:"required"`

	// Version of configuration
	Version int64 `json:"version" gorm:"column:version;type:bigint" validate:"required"`
}

func (h *TFPConfig) String() string {
	str, err := json.Marshal(h)
	if err != nil {
		panic(err)
	}
	return string(str)
}
