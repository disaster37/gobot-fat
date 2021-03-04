package models

import "encoding/json"

type DFPIO struct {
	GreenLed            bool `json:"green_led"`
	RedLed              bool `json:"red_led"`
	DrumRelay           bool `json:"drum_relay"`
	PumpRelay           bool `json:"pump_relay"`
	StartButton         bool `json:"start_button"`
	StopButton          bool `json:"stop_button"`
	WashButton          bool `json:"wash_button"`
	EmergencyButton     bool `json:"emergency_button"`
	ForceDrumButton     bool `json:"force_drum_button"`
	ForcePumpButton     bool `json:"force_pump_button"`
	WaterCaptorUpper    bool `json:"water_captor_upper"`
	WaterCaptorUnder    bool `json:"water_captor_under"`
	SecurityCaptorUpper bool `json:"security_captor_upper"`
	SecurityCaptorUnder bool `json:"security_captor_under"`
}

func (h DFPIO) String() string {
	str, err := json.Marshal(h)
	if err != nil {
		panic(err)
	}
	return string(str)
}
