package models

import "encoding/json"

type DFPIO struct {
	GreenLed            bool `json:"green_led" jsonapi:"attr,green_led"`
	RedLed              bool `json:"red_led" jsonapi:"attr,red_led"`
	DrumRelay           bool `json:"drum_relay" jsonapi:"attr,drum_relay"`
	PumpRelay           bool `json:"pump_relay" jsonapi:"attr,pump_relay"`
	StartButton         bool `json:"start_button" jsonapi:"attr,start_button"`
	StopButton          bool `json:"stop_button" jsonapi:"attr,stop_button"`
	WashButton          bool `json:"wash_button" jsonapi:"attr,wash_button"`
	EmergencyButton     bool `json:"emergency_button" jsonapi:"attr,emergency_button"`
	ForceDrumButton     bool `json:"force_drum_button" jsonapi:"attr,force_drum_button"`
	ForcePumpButton     bool `json:"force_pump_button" jsonapi:"attr,force_pump_button"`
	WaterCaptorUpper    bool `json:"water_captor_upper" jsonapi:"attr,water_captor_upper"`
	WaterCaptorUnder    bool `json:"water_captor_under" jsonapi:"attr,water_captor_under"`
	SecurityCaptorUpper bool `json:"security_captor_upper" jsonapi:"attr,security_captor_upper"`
	SecurityCaptorUnder bool `json:"security_captor_under" jsonapi:"attr,security_captor_under"`
}

func (h DFPIO) String() string {
	str, err := json.Marshal(h)
	if err != nil {
		panic(err)
	}
	return string(str)
}
