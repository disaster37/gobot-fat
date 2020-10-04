package board

import (
	"context"

	"github.com/disaster37/gobot-fat/models"
)

// Usecase is the board usecase interface
type Usecase interface {
	// GetBoards return the public board data
	GetBoards(ctx context.Context) ([]*models.Board, error)

	// AddBoard add board on list
	AddBoard(board Board)

	// Starts start each board
	Starts(ctx context.Context)

	// Stops stop each board
	Stops(ctx context.Context)
}
