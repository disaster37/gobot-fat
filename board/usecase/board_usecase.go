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
	boardsData := make([]*models.Board, 0, len(h.boards))

	for _, board := range h.boards {
		boardsData = append(boardsData, board.Board())
	}

	return boardsData, nil
}

func (h *boardUsecase) AddBoard(board board.Board) {
	h.boards = append(h.boards, board)
}

func (h *boardUsecase) Starts() {

	for _, board := range h.boards {
		go h.startBoard(board)
	}
}

func (h *boardUsecase) Stops() {
	for _, board := range h.boards {
		err := board.Stop()
		if err != nil {
			log.Errorf("Failed to stop successfully board %s: %s", board.Name, err.Error())
		}
	}
}

func (h *boardUsecase) startBoard(board board.Board) {

	log.Infof("Start board %s", board.Name())
	for {

		err := board.Start()
		if err != nil {
			log.Errorf("Failed to init board: %s", board.Name(), err.Error())
			time.Sleep(10 * time.Second)
		} else {
			break
		}
	}
}
