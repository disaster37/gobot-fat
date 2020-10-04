package usecase

import (
	"context"
	"time"

	"github.com/disaster37/gobot-fat/board"
	"github.com/disaster37/gobot-fat/models"
	log "github.com/sirupsen/logrus"
)

type boardUsecase struct {
	boards []board.Board
}

// NewBoardUsecase implement board usecase
func NewBoardUsecase() board.Usecase {

	return &boardUsecase{
		boards: make([]board.Board, 0, 1),
	}

}

// GetBoards return the public data for all boards
func (h *boardUsecase) GetBoards(ctx context.Context) ([]*models.Board, error) {

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	default:
		boardsData := make([]*models.Board, 0, len(h.boards))

		for _, board := range h.boards {
			boardsData = append(boardsData, board.Board())
		}

		return boardsData, nil
	}
}

// AddBoard add board on list
func (h *boardUsecase) AddBoard(board board.Board) {
	h.boards = append(h.boards, board)
}

// Starts start each board on background
// If board failed to start, it try again while context not canceled
func (h *boardUsecase) Starts(ctx context.Context) {
	select {
	case <-ctx.Done():
		log.Infof("Context canceled: %s", ctx.Err())
		return
	default:
		for _, board := range h.boards {
			go h.startBoard(ctx, board)
		}
		return
	}
}

// Stops stop each board
func (h *boardUsecase) Stops(ctx context.Context) {
	select {
	case <-ctx.Done():
		log.Infof("Context canceled: %s", ctx.Err())
		return
	default:
		for _, board := range h.boards {
			err := board.Stop(ctx)
			if err != nil {
				log.Errorf("Failed to stop successfully board %s: %s", board.Name(), err.Error())
			}
		}
		return
	}
}

// startBoard start board and try while context not canceled
func (h *boardUsecase) startBoard(ctx context.Context, board board.Board) {
	for {
		select {
		case <-ctx.Done():
			log.Infof("Context canceled: %s", ctx.Err())
			return
		default:
			log.Infof("Start board %s", board.Name())

			err := board.Start(ctx)
			if err != nil {
				log.Errorf("Failed to init board %s: %s", board.Name(), err.Error())
				time.Sleep(10 * time.Second)
			} else {
				return
			}
		}
	}
}
