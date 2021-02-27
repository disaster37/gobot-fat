package models

import "encoding/json"

type DFPIO struct {
	GreenLed         string `json:"green_led"`
	RedLed           string `json:"red_led"`
	DrumRelay        string `json:"drum_relay"`
	PumpRelay        string `json:"pump_relay"`
	StartButton      string `json:"start_button"`
	StopButton       string `json:"stop_button"`
	WashButton       string `json:"wash_button"`
	EmergencyButton  string `json:"emergency_button"`
	ForceDrumButton  string `json:"force_drum_button"`
	ForcePumpButton  string `json:"force_pump_button"`
	WaterTemperature string `json:"water_temperature"`
}

func (h *DFPIO) String() string {
	str, err := json.Marshal(h)
	if err != nil {
		panic(err)
	}
	return string(str)
}
