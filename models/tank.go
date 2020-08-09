package models

// Tank store values measured per distance sensor
type Tank struct {
	// The current level of water in cm
	Level int `json:"level"`

	// The current volume of water in liter
	Volume int `json:"volume"`

	// The ratio of water
	Percent float64 `json:"percent"`
}
