package models

import (
	"encoding/json"
	"time"
)

// DFPState  describe the current state of drum filter pond
type DFPState struct {
	ModelGeneric
	Name               string    `json:"name"`
	IsWashed           bool      `json:"is_washed" gorm:"column:is_washed" validate:"required"`
	IsRunning          bool      `json:"is_running" gorm:"column:is_running" validate:"required"`
	IsSecurity         bool      `json:"is_security" gorm:"column:is_security" validate:"required"`
	IsEmergencyStopped bool      `json:"is_emmergency_stopped" gorm:"column:is_emmergency_stopped" validate:"required"`
	IsDisableSecurity  bool      `json:"is_disable_security" gorm:"column:is_disable_security" validate:"required"`
	LastWashing        time.Time `json:"last_washing" gorm:"column:last_washing" validate:"required"`
	WaterTemperature   float64   `json:"water_tempareture" gorm:"column:last_washing" validate:"required"`
	AmbientTemperature float64   `json:"ambient_tempareture" gorm:"column:last_washing" validate:"required"`
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
