package models

import "encoding/json"

type TFPIO struct {
	PondPumpRelay      bool `json:"pond_pump_relay" jsonapi:"attr,pond_pump_relay"`
	WaterfallPumpRelay bool `json:"waterfall_pump_relay" jsonapi:"attr,waterfall_pump_relay"`
	UVC1Relay          bool `json:"uvc1_relay" jsonapi:"attr,uvc1_relay"`
	UVC2Relay          bool `json:"uvc2_relay" jsonapi:"attr,uvc2_relay"`
	PondBubble         bool `json:"pond_bubble" jsonapi:"attr,pond_bubble"`
	FilterBubble       bool `json:"filter_bubble" jsonapi:"attr,filter_bubble"`
}

func (h TFPIO) String() string {
	str, err := json.Marshal(h)
	if err != nil {
		panic(err)
	}
	return string(str)
}
