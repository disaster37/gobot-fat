package models

import "encoding/json"

type DFPIO struct {
	GreenLed            string `json:"green_led"`
	RedLed              string `json:"red_led"`
	DrumRelay           string `json:"drum_relay"`
	PumpRelay           string `json:"pump_relay"`
	StartButton         string `json:"start_button"`
	StopButton          string `json:"stop_button"`
	WashButton          string `json:"wash_button"`
	EmergencyButton     string `json:"emergency_button"`
	ForceDrumButton     string `json:"force_drum_button"`
	ForcePumpButton     string `json:"force_pump_button"`
	WaterTemperature    string `json:"water_temperature"`
	WaterCaptorUpper    string `json:"water_captor_upper"`
	WaterCaptorUnder    string `json:"water_captor_under"`
	SecurityCaptorUpper string `json:"security_captor_upper"`
	SecurityCaptorUnder string `json:"security_captor_under"`
}

func (h DFPIO) String() string {
	str, err := json.Marshal(h)
	if err != nil {
		panic(err)
	}
	return string(str)
}
