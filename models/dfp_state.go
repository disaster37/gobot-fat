package models

import (
	"encoding/json"
	"time"
)

// DFPState  describe the current state of drum filter pond
type DFPState struct {
	ModelGeneric

	ID                 uint      `json:"id" jsonapi:"primary,dfp-states" gorm:"primary_key"`
	Name               string    `json:"name" jsonapi:"attr,name" gorm:"column:name" validate:"required"`
	IsWashed           bool      `json:"is_washed" jsonapi:"attr,is_washed" gorm:"column:is_washed" validate:"required"`
	IsRunning          bool      `json:"is_running" jsonapi:"attr,is_running" gorm:"column:is_running" validate:"required"`
	IsSecurity         bool      `json:"is_security" jsonapi:"attr,is_security" gorm:"column:is_security" validate:"required"`
	IsEmergencyStopped bool      `json:"is_emmergency_stopped" jsonapi:"attr,is_emmergency_stopped" gorm:"column:is_emmergency_stopped" validate:"required"`
	IsDisableSecurity  bool      `json:"is_disable_security" jsonapi:"attr,is_disable_security" gorm:"column:is_disable_security" validate:"required"`
	IsForceDrum        bool      `json:"is_force_drum" jsonapi:"attr,is_force_drum" gorm:"column:is_force_drum" validate:"required"`
	IsForcePump        bool      `json:"is_force_pump" jsonapi:"attr,is_force_pump" gorm:"column:is_force_pump" validate:"required"`
	LastWashing        time.Time `json:"last_washing" jsonapi:"attr,last_washing,iso8601" gorm:"column:last_washing" validate:"required"`
	WaterTemperature   float64   `json:"water_tempareture" jsonapi:"attr,water_tempareture" gorm:"column:water_tempareture" validate:"required"`
	AmbientTemperature float64   `json:"ambient_tempareture" jsonapi:"attr,ambient_tempareture" gorm:"column:ambient_tempareture" validate:"required"`
}

func (h DFPState) TableName() string {
	return "dfpstate"
}

func (h *DFPState) String() string {
	data, err := json.Marshal(h)
	if err != nil {
		panic(err)
	}
	return string(data)
}

// ShouldWash is true if washing can start
func (h *DFPState) ShouldWash() bool {
	if h.IsWashed || !h.IsRunning || h.IsEmergencyStopped || (h.IsSecurity && !h.IsDisableSecurity) {
		return false
	}

	return true

}

// ShouldMotorStart is true if motor can start
func (h *DFPState) ShouldMotorStart() bool {
	if !h.IsRunning || h.IsEmergencyStopped || (h.IsSecurity && !h.IsDisableSecurity) {
		return false
	}

	return true
}

// Security return true id security and security not disabled
func (h *DFPState) Security() bool {
	if h.IsSecurity && !h.IsDisableSecurity {
		return true
	}

	return false
}

func (h *DFPState) SetID(id uint) {
	h.ID = id
}

func (h *DFPState) GetID() uint {
	return h.ID
}
