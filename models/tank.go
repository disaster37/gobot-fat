package models

// Tank store values measured per distance sensor
type Tank struct {
	// The current level of water in cm
	Level int

	// The current volume of water in liter
	Volume int

	// The ratio of water
	Percent float64
}
