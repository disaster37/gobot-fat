package models

import (
	"encoding/json"
	"time"

	"github.com/jinzhu/gorm"
)

// TFPState  describe the current state of drum filter pond
type TFPState struct {
	gorm.Model

	Name string `json:"name"`

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

	// UVC1BlisterTime is the date when replace UVC1 blister
	UVC1BlisterTime time.Time `json:"uvc1_blister_time" gorm:"column:uvc1_blister_time" validate:"required"`

	// UVC1BlisterNbHour is the blister usage in hour of UVC1
	UVC1BlisterNbHour int64 `json:"uvc1_blister_nb_hour" gorm:"column:uvc1_blister_nb_hour" validate:"required"`

	// UVC2BlisterTime is the date when replace UVC2 blister
	UVC2BlisterTime time.Time `json:"uvc2_blister_time" gorm:"column:uvc2_blister_time" validate:"required"`

	// UVC2BlisterNbHour is the blister usage in hour of UVC2
	UVC2BlisterNbHour int64 `json:"uvc2_blister_nb_hour" gorm:"column:uvc2_blister_nb_hour" validate:"required"`

	// OzoneBlisterTime is the date when replace Ozone blister
	OzoneBlisterTime time.Time `json:"ozone_blister_time" gorm:"column:ozone_blister_time" validate:"required"`

	// OzoneBlisterNbHour is the blister usage in hour of Ozone
	OzoneBlisterNbHour int64 `json:"ozone_blister_nb_hour" gorm:"column:ozone_blister_nb_hour" validate:"required"`

	// IsSecurity is true when security is fire
	IsSecurity bool `json:"is_security" gorm:"column:is_security" validate:"required"`

	// IsEmergencyStopped is stop when all must be stopped
	IsEmergencyStopped bool `json:"is_emmergency_stopped" gorm:"column:is_emmergency_stopped" validate:"required"`

	// IsDisableSecurity permit to not handle security state
	IsDisableSecurity bool `json:"is_disable_security" gorm:"column:is_disable_security" validate:"required"`

	// BacteriumTime is the time when introduce bacterium to power off UVC during 48h
	BacteriumTime time.Time `json:"bacterium_time" gorm:"column:bacterium_time" validate:"required"`

	AcknoledgeWaterfallAuto bool `json:"acknoledge_waterfall_auto" gorm:"column:acknoledge_waterfall_auto" validate:"required"`
}

func (h TFPState) TableName() string {
	return "tfpstate"
}

func (h *TFPState) String() string {
	data, err := json.Marshal(h)
	if err != nil {
		panic(err)
	}
	return string(data)
}
