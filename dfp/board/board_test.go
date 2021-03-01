package dfpboard

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

type DFPBoardTestSuite struct {
	suite.Suite
	board   *DFPBoard
	adaptor *mock.MockPlateform
}

func TestDFPBoardTestSuite(t *testing.T) {
	suite.Run(t, new(DFPBoardTestSuite))
}

func (s *DFPBoardTestSuite) SetupSuite() {
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
func (s *DFPBoardTestSuite) SetupTest() {

	// Button
	s.adaptor.SetDigitalPinState(s.board.buttonStart.Pin(), 1)
	s.adaptor.SetDigitalPinState(s.board.buttonStop.Pin(), 1)
	s.adaptor.SetDigitalPinState(s.board.buttonWash.Pin(), 1)
	s.adaptor.SetDigitalPinState(s.board.buttonEmergencyStop.Pin(), 1)
	s.adaptor.SetDigitalPinState(s.board.buttonForceDrum.Pin(), 1)
	s.adaptor.SetDigitalPinState(s.board.buttonForcePump.Pin(), 1)

	// Captor
	s.adaptor.SetDigitalPinState(s.board.captorSecurityUpper.Pin(), 1)
	s.adaptor.SetDigitalPinState(s.board.captorSecurityUnder.Pin(), 0)
	s.adaptor.SetDigitalPinState(s.board.captorWaterUpper.Pin(), 1)
	s.adaptor.SetDigitalPinState(s.board.captorWaterUnder.Pin(), 0)

	// Relay
	s.board.relayDrum.Off()
	s.board.relayPump.Off()

	// Led
	s.board.ledGreen.On()
	s.board.ledRed.Off()

	// State
	s.board.state = &models.DFPState{
		IsRunning: true,
	}

	// Config
	s.board.config = &models.DFPConfig{
		ForceWashingDuration:           180,
		ForceWashingDurationWhenFrozen: 60,
		StartWashingPumpBeforeWashing:  1,
		WaitTimeBetweenWashing:         1,
		WashingDuration:                1,
		WaitTimeBeforeUnsetSecurity:    1,
		TemperatureSensorPolling:       1,
	}

}

func (s *DFPBoardTestSuite) TestStartStopIsOnline() {

	// Normal start with all stopped on state and running
	assert.True(s.T(), s.board.IsOnline())
	assert.Equal(s.T(), 1, s.board.buttonForceDrum.DefaultState)
	assert.Equal(s.T(), 1, s.board.buttonForcePump.DefaultState)
	assert.Equal(s.T(), 1, s.board.buttonEmergencyStop.DefaultState)
	assert.Equal(s.T(), 1, s.board.buttonStart.DefaultState)
	assert.Equal(s.T(), 1, s.board.buttonStop.DefaultState)
	assert.Equal(s.T(), 1, s.board.buttonWash.DefaultState)

	assert.Equal(s.T(), 0, s.board.captorSecurityUnder.DefaultState)
	assert.Equal(s.T(), 1, s.board.captorSecurityUpper.DefaultState)
	assert.Equal(s.T(), 0, s.board.captorWaterUnder.DefaultState)
	assert.Equal(s.T(), 1, s.board.captorWaterUpper.DefaultState)

	assert.Equal(s.T(), 0, s.adaptor.GetDigitalPinState(s.board.relayDrum.Pin()))
	assert.Equal(s.T(), 0, s.adaptor.GetDigitalPinState(s.board.relayPump.Pin()))
	assert.Equal(s.T(), 1, s.adaptor.GetDigitalPinState(s.board.ledGreen.Pin()))
	assert.Equal(s.T(), 0, s.adaptor.GetDigitalPinState(s.board.ledRed.Pin()))

	// Normal start with all stopped on state and stopped
	board, adaptor := initTestBoard()
	board.state.IsRunning = false
	err := board.Start(context.Background())
	assert.NoError(s.T(), err)

	assert.True(s.T(), board.IsOnline())
	assert.Equal(s.T(), 1, board.buttonForceDrum.DefaultState)
	assert.Equal(s.T(), 1, board.buttonForcePump.DefaultState)
	assert.Equal(s.T(), 1, board.buttonEmergencyStop.DefaultState)
	assert.Equal(s.T(), 1, board.buttonStart.DefaultState)
	assert.Equal(s.T(), 1, board.buttonStop.DefaultState)
	assert.Equal(s.T(), 1, board.buttonWash.DefaultState)

	assert.Equal(s.T(), 0, board.captorSecurityUnder.DefaultState)
	assert.Equal(s.T(), 1, board.captorSecurityUpper.DefaultState)
	assert.Equal(s.T(), 0, board.captorWaterUnder.DefaultState)
	assert.Equal(s.T(), 1, board.captorWaterUpper.DefaultState)

	assert.Equal(s.T(), 0, adaptor.GetDigitalPinState(board.relayDrum.Pin()))
	assert.Equal(s.T(), 0, adaptor.GetDigitalPinState(board.relayPump.Pin()))
	assert.Equal(s.T(), 0, adaptor.GetDigitalPinState(board.ledGreen.Pin()))
	assert.Equal(s.T(), 1, adaptor.GetDigitalPinState(board.ledRed.Pin()))
	board.Stop(context.Background())

	// Start with wash and running)
	board, adaptor = initTestBoard()
	board.state.IsWashed = true
	status := mock.WaitEvent(board.Eventer, EventWash, 5*time.Second)
	err = board.Start(context.Background())
	assert.NoError(s.T(), err)
	assert.True(s.T(), <-status)

	// Stop
	// It emit event
	status = mock.WaitEvent(board.Eventer, EventBoardStop, 1*time.Second)
	err = board.Stop(context.Background())
	assert.NoError(s.T(), err)
	assert.True(s.T(), <-status)

}

func (s *DFPBoardTestSuite) TestGetBoard() {
	assert.Equal(s.T(), "test", s.board.Board().Name)
	assert.True(s.T(), s.board.Board().IsOnline)
}

func (s *DFPBoardTestSuite) TestName() {
	assert.Equal(s.T(), "test", s.board.Name())
}

func (s *DFPBoardTestSuite) TestState() {
	assert.True(s.T(), reflect.DeepEqual(models.DFPState{IsRunning: true}, s.board.State()))
}
