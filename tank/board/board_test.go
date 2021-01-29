package tankboard

import (
	"context"
	"testing"
	"time"

	"github.com/disaster37/gobot-fat/helper"

	"github.com/disaster37/gobot-fat/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TankBoardTestSuite struct {
	suite.Suite
	board   *TankBoard
	adaptor *helper.MockPlateform
}

func TestTankBoardTestSuite(t *testing.T) {
	suite.Run(t, new(TankBoardTestSuite))
}

func (s *TankBoardTestSuite) SetupSuite() {
	s.board, s.adaptor = initTestBoard()
	if err := s.board.Start(context.Background()); err != nil {
		panic(err)
	}

	// wait initialized
	for !s.board.isInitialized {
		time.Sleep(1 * time.Second)
	}
}

// Put default state for i/o
func (s *TankBoardTestSuite) SetupTest() {

	// Return the right type for drivers
	s.adaptor.SetValueReadState("isRebooted", false)
	s.adaptor.SetValueReadState("distance", float64(0))

	// Config
	s.board.config = &models.TankConfig{
		Depth:        100,
		LiterPerCm:   1,
		SensorHeight: 0,
	}
}

func (s *TankBoardTestSuite) TestStartStopIsOnline() {
	board, _ := initTestBoard()

	err := board.Start(context.Background())
	assert.NoError(s.T(), err)
	assert.True(s.T(), board.IsOnline())

	err = board.Stop(context.Background())
	assert.NoError(s.T(), err)
	assert.False(s.T(), board.IsOnline())
}

func (s *TankBoardTestSuite) TestGetBoard() {
	assert.Equal(s.T(), "test", s.board.Board().Name)
	assert.True(s.T(), s.board.Board().IsOnline)
}

func (s *TankBoardTestSuite) TestName() {
	assert.Equal(s.T(), "test", s.board.Name())
}
