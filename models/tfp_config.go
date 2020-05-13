package models

import (
	"encoding/json"

	"github.com/jinzhu/gorm"
)

// TFPConfig contain config data for Technical Filter Pond
type TFPConfig struct {
	gorm.Model

	// UVC1BlisterMaxTime is the max usage in hour of UVC1 blister
	UVC1BlisterMaxTime int64 `json:"uvc1_blister_max_time" gorm:"column:uvc1_blister_max_time" validate:"required"`

	// UVC1BlisterMaxTime is the max usage in hour of UVC2 blister
	UVC2BlisterMaxTime int64 `json:"uvc2_blister_max_time" gorm:"column:uvc2_blister_max_time" validate:"required"`

	// OzoneBlisterMaxTime is the max usage in hour of ozonne blister
	OzoneBlisterMaxTime int64 `json:"ozone_blister_max_time" gorm:"column:ozone_blister_max_time" validate:"required"`

	// IsWaterfallAuto permit to start / stop waterfall pump automatically
	IsWaterfallAuto bool `json:"is_waterfall_auto" gorm:"column:is_waterfall_auto" validate:"required"`

	// StartTimeWaterfall is the hour of day when start waterfall pump
	StartTimeWaterfall string `json:"start_time_waterfall" gorm:"column:start_time_waterfall" validate:"required"`

	// StopTimeWaterfall is the hour of day when stop waterfall pump
	StopTimeWaterfall string `json:"stop_time_waterfall" gorm:"column:stop_time_waterfall" validate:"required"`

	//Mode is ozone, or UVC or none
	Mode string `json:"mode" gorm:"column:mode" validate:"required"`
}

func (h *TFPConfig) String() string {
	str, err := json.Marshal(h)
	if err != nil {
		panic(err)
	}
	return string(str)
}

func (h TFPConfig) TableName() string {
	return "tfpconfig"
}
