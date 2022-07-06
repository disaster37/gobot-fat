package models

import (
	"encoding/json"
	"time"
)

// TFPState  describe the current state of drum filter pond
type TFPState struct {
	ModelGeneric

	ID uint `jsonapi:"primary,tfp-states" gorm:"primary_key"`

	Name string `json:"name" jsonapi:"attr,name" gorm:"column:name"`

	// UVC1Running is true if UVC1 running
	UVC1Running bool `json:"uvc1_running" jsonapi:"attr,uvc1_running" gorm:"column:uvc1_running" validate:"required"`

	// UVC2Running is true if UVC2 running
	UVC2Running bool `json:"uvc2_running" jsonapi:"attr,uvc2_running" gorm:"column:uvc2_running" validate:"required"`

	// PondPumpRunning is true if pond pump running
	PondPumpRunning bool `json:"pond_pump_running" jsonapi:"attr,pond_pump_running" gorm:"column:pond_pump_running" validate:"required"`

	// WaterfallPumpRunning is true if Waterfall pump running
	WaterfallPumpRunning bool `json:"waterfall_pump_running" jsonapi:"attr,waterfall_pump_running" gorm:"column:waterfall_pump_running" validate:"required"`

	// PondBubbleRunning is true if pond bubble running
	PondBubbleRunning bool `json:"pond_bubble_running" jsonapi:"attr,pond_bubble_running" gorm:"column:pond_bubble_running" validate:"required"`

	// FilterBubbleRunning is true if filter bubble running
	FilterBubbleRunning bool `json:"filter_bubble_running" jsonapi:"attr,filter_bubble_running" gorm:"column:filter_bubble_running" validate:"required"`

	// UVC1BlisterNbHour is the blister usage in hour of UVC1
	UVC1BlisterNbHour int64 `json:"uvc1_blister_nb_hour" jsonapi:"attr,uvc1_blister_nb_hour" gorm:"column:uvc1_blister_nb_hour" validate:"required"`

	// UVC2BlisterNbHour is the blister usage in hour of UVC2
	UVC2BlisterNbHour int64 `json:"uvc2_blister_nb_hour" jsonapi:"attr,uvc2_blister_nb_hour" gorm:"column:uvc2_blister_nb_hour" validate:"required"`

	// OzoneBlisterNbHour is the blister usage in hour of Ozone
	OzoneBlisterNbHour int64 `json:"ozone_blister_nb_hour" jsonapi:"attr,ozone_blister_nb_hour" gorm:"column:ozone_blister_nb_hour" validate:"required"`

	// IsSecurity is true when security is fire
	IsSecurity bool `json:"is_security" jsonapi:"attr,is_security" gorm:"column:is_security" validate:"required"`

	// IsEmergencyStopped is stop when all must be stopped
	IsEmergencyStopped bool `json:"is_emmergency_stopped" jsonapi:"attr,is_emmergency_stopped" gorm:"column:is_emmergency_stopped" validate:"required"`

	// IsDisableSecurity permit to not handle security state
	IsDisableSecurity bool `json:"is_disable_security" jsonapi:"attr,is_disable_security" gorm:"column:is_disable_security" validate:"required"`

	// BacteriumTime is the time when introduce bacterium to power off UVC during 48h
	BacteriumTime time.Time `json:"bacterium_time" jsonapi:"attr,bacterium_time,iso8601" gorm:"column:bacterium_time" validate:"required"`

	AcknoledgeWaterfallAuto bool `json:"acknoledge_waterfall_auto" jsonapi:"attr,acknoledge_waterfall_auto" gorm:"column:acknoledge_waterfall_auto" validate:"required"`

	// IsWaterfallAuto is managed by tfpConfig
	// It's here only to reflect state
	IsWaterfallAuto bool `json:"is_waterfall_auto" jsonapi:"attr,is_waterfall_auto" gorm:"column:is_waterfall_auto" validate:"required"`
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

func (h *TFPState) SetID(id uint) {
	h.ID = id
}

func (h *TFPState) GetID() uint {
	return h.ID
}
