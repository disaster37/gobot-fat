package board

import "github.com/disaster37/gobot-fat/models"

// Board represent generic board
type Board interface {
	// IsOnline permit to know if board is online
	IsOnline() bool

	// Start start the main function
	Start() error

	// Stop interrupt the main function
	Stop() error

	// Name return the board name
	Name() string

	// Board return the board data
	Board() *models.Board
}
