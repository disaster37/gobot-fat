package models

// Tank store values measured per distance sensor
type Tank struct {

	// The tank name
	ID string `jsonapi:"primary,tanks"`

	// The current level of water in cm
	Level int `json:"level" jsonapi:"attr,level"`

	// The current volume of water in liter
	Volume int `json:"volume" jsonapi:"attr,volume"`

	// The ratio of water
	Percent float64 `json:"percent" jsonapi:"attr,percent"`

	// The current distance in cm
	Distance int `json:"distance" jsonapi:"attr,distance"`
}
