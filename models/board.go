package models

// Board represent generic board
type Board struct {

	// Name is the board name
	Name string `json:"name" jsonapi:"primary,boards"`

	// IsOnline is true if board is online
	IsOnline bool `json:"is_online" jsonapi:"attr,is_online"`
}
