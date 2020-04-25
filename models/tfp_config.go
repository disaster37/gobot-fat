package models

import (
	"encoding/json"
	"time"

	"github.com/jinzhu/gorm"
)

// TFPConfig contain config data for Technical Filter Pond
type TFPConfig struct {
	gorm.Model

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

	// UVC1BlisterTime is the time usage of UVC1 blister
	UVC1BlisterTime time.Time

	// UVC2BlisterTime is the time usage of UVC2 blister
	UVC2BlisterTime time.Time

	// UVC1BlisterMaxTime is the max usage in hour of UVC1 blister
	UVC1BlisterMaxTime int64

	// UVC1BlisterMaxTime is the max usage in hour of UVC2 blister
	UVC2BlisterMaxTime int64

	// BacteriumTime is the time when introduce bacterium to power off UVC during 48h
	BacteriumTime time.Time

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
