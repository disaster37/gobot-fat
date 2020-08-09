package board

import (
	"context"

	"github.com/disaster37/gobot-fat/models"
)

type Usecase interface {
	GetBoards(ctx context.Context) ([]*models.Board, error)
	AddBoard(board Board)
	Starts()
	Stops()
}
