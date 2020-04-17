package models

import (
	"encoding/json"
	"time"
)

// DFPState  describe the current state of drum filter pond
type DFPState struct {
	ID                 string    `json:"id"`
	Name               string    `json:"name"`
	IsWashed           bool      `json:"is_washed"`
	ShouldWash         bool      `json:"should_wash"`
	IsAuto             bool      `json:"is_started"`
	IsStopped          bool      `json:"is_stopped"`
	IsSecurity         bool      `json:"is_security"`
	IsEmergencyStopped bool      `json:"is_emmergency_stopped"`
	IsDisableSecurity  bool      `json:"is_disable_security"`
	LastWashing        time.Time `json:"last_washing"`
}

// NewDFPState return new PBF handler
func NewDFPState(ID string, name string, isAuto bool, isStopped bool, isEmergencyStopped bool, isDisableSecurity bool) *DFPState {
	return &DFPState{
		ID:                 ID,
		Name:               name,
		IsWashed:           false,
		ShouldWash:         false,
		IsAuto:             isAuto,
		IsStopped:          isStopped,
		IsSecurity:         false,
		IsEmergencyStopped: isEmergencyStopped,
		IsDisableSecurity:  isDisableSecurity,
	}
}

func (h *DFPState) String() string {
	data, err := json.Marshal(h)
	if err != nil {
		panic(err)
	}
	return string(data)
}
