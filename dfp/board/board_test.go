package dfpboard

import (
	"context"
	"testing"
	"time"

	"github.com/disaster37/gobot-fat/mock"
	"github.com/disaster37/gobot-fat/models"
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
	s.adaptor.SetDigitalPinState(s.board.captorSecurityUpper.Pin(), 0)
	s.adaptor.SetDigitalPinState(s.board.captorSecurityUnder.Pin(), 1)
	s.adaptor.SetDigitalPinState(s.board.captorWaterUpper.Pin(), 0)
	s.adaptor.SetDigitalPinState(s.board.captorWaterUnder.Pin(), 1)

	// Relay
	err := s.board.relayDrum.Off()
	if err != nil {
		s.T().Fatal(err)
	}
	err = s.board.relayPump.Off()
	if err != nil {
		s.T().Fatal(err)
	}

	// Led
	err = s.board.ledGreen.On()
	if err != nil {
		s.T().Fatal(err)
	}
	err = s.board.ledRed.Off()
	if err != nil {
		s.T().Fatal(err)
	}

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

/*

func (s *DFPBoardTestSuite) TestStartStopIsOnline() {

	// Normal start with all stopped on state and running
	assert.True(s.T(), s.board.IsOnline())

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

	assert.Equal(s.T(), 0, adaptor.GetDigitalPinState(board.relayDrum.Pin()))
	assert.Equal(s.T(), 0, adaptor.GetDigitalPinState(board.relayPump.Pin()))
	assert.Equal(s.T(), 0, adaptor.GetDigitalPinState(board.ledGreen.Pin()))
	assert.Equal(s.T(), 0, adaptor.GetDigitalPinState(board.ledRed.Pin()))
	err = board.Stop(context.Background())
	if err != nil {
		s.T().Fatal(err)
	}

	// Start with wash and running)
	board, _ = initTestBoard()
	board.state.IsWashed = true
	status := mock.WaitEvent(board.Eventer, EventWash, 5*time.Second)
	err = board.Start(context.Background())
	assert.NoError(s.T(), err)
	// Temp delete test
	//assert.True(s.T(), <-status)

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

*/
