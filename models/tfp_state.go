package models

import (
	"encoding/json"
)

// TFPState  describe the current state of drum filter pond
type TFPState struct {
	ID                   string `json:"id"`
	Name                 string `json:"name"`
	UVC1Running          bool   `json:"uvc1_running"`
	UVC2Running          bool   `json:"uvc2_running"`
	PondPumpRunning      bool   `json:"pond_pump_running"`
	WaterfallPumpRunning bool   `json:"waterfall_pump_running"`
	PondBubbleRunning    bool   `json:"pond_bubble_running"`
	FilterBubbleRunning  bool   `json:"filter_bubble_running"`
	IsSecurity           bool   `json:"is_security"`
	IsEmergencyStopped   bool   `json:"is_emmergency_stopped"`
	IsDisableSecurity    bool   `json:"is_disable_security"`
}

// NewTFPState return new TFP handler
func NewTFPState(ID string, name string, UVC1Running bool, UVC2Running bool, PondPumpRunning bool, WaterfallPumpRunning bool, PondBubbleRunning bool, FilterBubbleRunning bool, isEmergencyStopped bool, isDisableSecurity bool) *TFPState {
	return &TFPState{
		ID:                   ID,
		Name:                 name,
		UVC1Running:          UVC1Running,
		UVC2Running:          UVC2Running,
		PondPumpRunning:      PondPumpRunning,
		WaterfallPumpRunning: WaterfallPumpRunning,
		PondBubbleRunning:    PondBubbleRunning,
		FilterBubbleRunning:  FilterBubbleRunning,
		IsSecurity:           false,
		IsEmergencyStopped:   isEmergencyStopped,
		IsDisableSecurity:    isDisableSecurity,
	}
}

func (h *TFPState) String() string {
	data, err := json.Marshal(h)
	if err != nil {
		panic(err)
	}
	return string(data)
}
