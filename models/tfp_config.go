package models

import (
	"encoding/json"
	"time"
)

// TFPConfig contain config data for Technical Filter Pond
type TFPConfig struct {
	ModelGeneric

	ID uint `jsonapi:"primary,tfp-configs" gorm:"primary_key"`

	// Enable is set to true if board is enabled
	Enable bool `json:"enable" jsonapi:"attr,enable" gorm:"column:enable" validate:"required"`

	// UVC1BlisterMaxTime is the max usage in hour of UVC1 blister
	UVC1BlisterMaxTime int64 `json:"uvc1_blister_max_time" jsonapi:"attr,uvc1_blister_max_time" gorm:"column:uvc1_blister_max_time" validate:"required"`

	// UVC1BlisterMaxTime is the max usage in hour of UVC2 blister
	UVC2BlisterMaxTime int64 `json:"uvc2_blister_max_time" jsonapi:"attr,uvc2_blister_max_time" gorm:"column:uvc2_blister_max_time" validate:"required"`

	// OzoneBlisterMaxTime is the max usage in hour of ozonne blister
	OzoneBlisterMaxTime int64 `json:"ozone_blister_max_time" jsonapi:"attr,ozone_blister_max_time" gorm:"column:ozone_blister_max_time" validate:"required"`

	// IsWaterfallAuto permit to start / stop waterfall pump automatically
	IsWaterfallAuto bool `json:"is_waterfall_auto" jsonapi:"attr,is_waterfall_auto" gorm:"column:is_waterfall_auto" validate:"required"`

	// StartTimeWaterfall is the hour of day when start waterfall pump
	StartTimeWaterfall string `json:"start_time_waterfall" jsonapi:"attr,start_time_waterfall" gorm:"column:start_time_waterfall" validate:"required"`

	// StopTimeWaterfall is the hour of day when stop waterfall pump
	StopTimeWaterfall string `json:"stop_time_waterfall" jsonapi:"attr,stop_time_waterfall" gorm:"column:stop_time_waterfall" validate:"required"`

	//Mode is ozone, or UVC or none
	Mode string `json:"mode" gorm:"column:mode" jsonapi:"attr,mode" validate:"required"`

	// UVC1BlisterTime is the date when replace UVC1 blister
	UVC1BlisterTime time.Time `json:"uvc1_blister_time" jsonapi:"attr,uvc1_blister_time,iso8601" gorm:"column:uvc1_blister_time" validate:"required"`

	// UVC2BlisterTime is the date when replace UVC2 blister
	UVC2BlisterTime time.Time `json:"uvc2_blister_time" jsonapi:"attr,uvc2_blister_time,iso8601" gorm:"column:uvc2_blister_time" validate:"required"`

	// OzoneBlisterTime is the date when replace Ozone blister
	OzoneBlisterTime time.Time `json:"ozone_blister_time" jsonapi:"attr,ozone_blister_time,iso8601" gorm:"column:ozone_blister_time" validate:"required"`
}

// String print the current object as json
func (h *TFPConfig) String() string {
	str, err := json.Marshal(h)
	if err != nil {
		panic(err)
	}
	return string(str)
}

// TableName permit to return the current table name
func (h TFPConfig) TableName() string {
	return "tfpconfig"
}

func (h *TFPConfig) SetID(id uint) {
	h.ID = id
}

func (h *TFPConfig) GetID() uint {
	return h.ID
}
