package models

import "encoding/json"

type TFPIO struct {
	PondPumpRelay      bool `json:"pond_pump_relay"`
	WaterfallPumpRelay bool `json:"waterfall_pump_relay"`
	UVC1Relay          bool `json:"uvc1_relay"`
	UVC2Relay          bool `json:"uvc2_relay"`
	PondBubble         bool `json:"pond_bubble"`
	FilterBubble       bool `json:"filter_bubble"`
}

func (h TFPIO) String() string {
	str, err := json.Marshal(h)
	if err != nil {
		panic(err)
	}
	return string(str)
}
