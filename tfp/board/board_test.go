package tfpboard

import (
	"context"
	"reflect"
	"testing"
	"time"

	"github.com/disaster37/gobot-fat/mock"
	"github.com/disaster37/gobot-fat/models"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
)

type TFPBoardTestSuite struct {
	suite.Suite
	board   *TFPBoard
	adaptor *mock.MockPlateform
}

func TestDFPBoardTestSuite(t *testing.T) {
	suite.Run(t, new(TFPBoardTestSuite))
}

func (s *TFPBoardTestSuite) SetupSuite() {
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
func (s *TFPBoardTestSuite) SetupTest() {

	// Relays
	s.adaptor.SetDigitalPinState(s.board.relayBubbleFilter.Pin(), 1)
	s.adaptor.SetDigitalPinState(s.board.relayPompPond.Pin(), 1)
	s.adaptor.SetDigitalPinState(s.board.relayPompWaterfall.Pin(), 0)
	s.adaptor.SetDigitalPinState(s.board.relayBubblePond.Pin(), 1)
	s.adaptor.SetDigitalPinState(s.board.relayUVC1.Pin(), 1)
	s.adaptor.SetDigitalPinState(s.board.relayUVC2.Pin(), 1)

	// State
	s.board.state = &models.TFPState{}

	// Config
	s.board.config = &models.TFPConfig{
		IsWaterfallAuto: false,
	}

	// Return the right type for drivers
	s.adaptor.SetValueReadState("isRebooted", false)
}

func (s *TFPBoardTestSuite) TestStartStopIsOnline() {
	board, adaptor := initTestBoard()

	// Normal start with all stopped on state
	err := board.Start(context.Background())
	assert.NoError(s.T(), err)
	assert.True(s.T(), board.IsOnline())
	assert.True(s.T(), board.relayPompPond.IsInverted())
	assert.False(s.T(), board.relayPompWaterfall.IsInverted())
	assert.True(s.T(), board.relayUVC1.IsInverted())
	assert.True(s.T(), board.relayUVC2.IsInverted())
	assert.True(s.T(), board.relayBubbleFilter.IsInverted())
	assert.True(s.T(), board.relayBubblePond.IsInverted())
	assert.Equal(s.T(), 1, adaptor.GetDigitalPinState(board.relayPompPond.Pin()))
	assert.Equal(s.T(), 1, adaptor.GetDigitalPinState(board.relayUVC1.Pin()))
	assert.Equal(s.T(), 1, adaptor.GetDigitalPinState(board.relayUVC2.Pin()))
	assert.Equal(s.T(), 1, adaptor.GetDigitalPinState(board.relayBubblePond.Pin()))
	assert.Equal(s.T(), 1, adaptor.GetDigitalPinState(board.relayBubbleFilter.Pin()))
	assert.Equal(s.T(), 0, adaptor.GetDigitalPinState(board.relayPompWaterfall.Pin()))
	err = board.Stop(context.Background())
	assert.NoError(s.T(), err)

	// Start with all started on state
	board, adaptor = initTestBoard()
	board.state.PondPumpRunning = true
	board.state.PondBubbleRunning = true
	board.state.FilterBubbleRunning = true
	board.state.WaterfallPumpRunning = true
	board.state.UVC1Running = true
	board.state.UVC2Running = true
	err = board.Start(context.Background())
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 0, adaptor.GetDigitalPinState(board.relayPompPond.Pin()))
	assert.Equal(s.T(), 0, adaptor.GetDigitalPinState(board.relayUVC1.Pin()))
	assert.Equal(s.T(), 0, adaptor.GetDigitalPinState(board.relayUVC2.Pin()))
	assert.Equal(s.T(), 0, adaptor.GetDigitalPinState(board.relayBubblePond.Pin()))
	assert.Equal(s.T(), 0, adaptor.GetDigitalPinState(board.relayBubbleFilter.Pin()))
	assert.Equal(s.T(), 1, adaptor.GetDigitalPinState(board.relayPompWaterfall.Pin()))

	// Stop
	err = board.Stop(context.Background())
	assert.NoError(s.T(), err)
	assert.False(s.T(), board.IsOnline())
}

func (s *TFPBoardTestSuite) TestGetBoard() {
	assert.Equal(s.T(), "test", s.board.Board().Name)
	assert.True(s.T(), s.board.Board().IsOnline)
}

func (s *TFPBoardTestSuite) TestName() {
	assert.Equal(s.T(), "test", s.board.Name())
}

func (s *TFPBoardTestSuite) TestState() {
	assert.True(s.T(), reflect.DeepEqual(models.TFPState{}, s.board.State()))
}

func (s *TFPBoardTestSuite) TestConfig() {
	assert.True(s.T(), reflect.DeepEqual(models.TFPConfig{IsWaterfallAuto: false}, s.board.Config()))
}
