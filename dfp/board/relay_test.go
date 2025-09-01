package dfpboard

import (
	"context"
	"time"

	"github.com/disaster37/gobot-fat/mock"
	"github.com/stretchr/testify/assert"
)

func (s *DFPBoardTestSuite) TestStartStopDFP() {

	// Start DFP when already started
	status := mock.WaitEvent(s.board.Eventer, EventStartDFP, 1*time.Second)
	err := s.board.StartDFP(context.Background())
	assert.NoError(s.T(), err)
	assert.False(s.T(), <-status)
	assert.True(s.T(), s.board.state.IsRunning)
	assert.Equal(s.T(), 1, s.adaptor.GetDigitalPinState(s.board.ledGreen.Pin()))
	assert.Equal(s.T(), 0, s.adaptor.GetDigitalPinState(s.board.ledRed.Pin()))

	// Start DFP when stopped
	s.board.state.IsRunning = false
	status = mock.WaitEvent(s.board.Eventer, EventStartDFP, 1*time.Second)
	err = s.board.StartDFP(context.Background())
	assert.NoError(s.T(), err)
	assert.True(s.T(), <-status)
	assert.True(s.T(), s.board.state.IsRunning)
	assert.Equal(s.T(), 1, s.adaptor.GetDigitalPinState(s.board.ledGreen.Pin()))
	assert.Equal(s.T(), 0, s.adaptor.GetDigitalPinState(s.board.ledRed.Pin()))

	// Stop DFP when running
	s.board.state.IsRunning = true
	status = mock.WaitEvent(s.board.Eventer, EventStopDFP, 1*time.Second)
	err = s.board.StopDFP(context.Background())
	assert.NoError(s.T(), err)
	assert.True(s.T(), <-status)
	assert.False(s.T(), s.board.state.IsRunning)
	assert.Equal(s.T(), 0, s.adaptor.GetDigitalPinState(s.board.ledGreen.Pin()))
	assert.Equal(s.T(), 0, s.adaptor.GetDigitalPinState(s.board.ledRed.Pin()))

	// Stop DFP when stopped
	s.board.state.IsRunning = false
	status = mock.WaitEvent(s.board.Eventer, EventStopDFP, 1*time.Second)
	err = s.board.StopDFP(context.Background())
	assert.NoError(s.T(), err)
	assert.False(s.T(), <-status)
	assert.False(s.T(), s.board.state.IsRunning)
	assert.Equal(s.T(), 0, s.adaptor.GetDigitalPinState(s.board.ledGreen.Pin()))
	assert.Equal(s.T(), 0, s.adaptor.GetDigitalPinState(s.board.ledRed.Pin()))

}

func (s *DFPBoardTestSuite) TestStartStopManualDrum() {

	// Start manual drum
	err := s.board.StartManualDrum(context.Background())
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 1, s.adaptor.GetDigitalPinState(s.board.relayDrum.Pin()))

	// Stop manual drum
	err = s.board.StopManualDrum(context.Background())
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 0, s.adaptor.GetDigitalPinState(s.board.relayDrum.Pin()))

	// Can start when emergency stop
	s.board.state.IsEmergencyStopped = true
	err = s.board.StartManualDrum(context.Background())
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 0, s.adaptor.GetDigitalPinState(s.board.relayDrum.Pin()))
}

func (s *DFPBoardTestSuite) TestStartStopManualPomp() {

	// Start manual pump
	err := s.board.StartManualPump(context.Background())
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 1, s.adaptor.GetDigitalPinState(s.board.relayPump.Pin()))

	// Stop manual pump
	err = s.board.StopManualPump(context.Background())
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 0, s.adaptor.GetDigitalPinState(s.board.relayPump.Pin()))

	// Can start when emergency stop
	s.board.state.IsEmergencyStopped = true
	err = s.board.StartManualPump(context.Background())
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 0, s.adaptor.GetDigitalPinState(s.board.relayPump.Pin()))
}

func (s *DFPBoardTestSuite) TestSetUnsetSecurity() {

	// Set security when no security
	s.board.state.IsSecurity = false
	status := mock.WaitEvent(s.board.Eventer, EventSetSecurity, 1*time.Second)
	err := s.board.SetSecurity(context.Background())
	assert.NoError(s.T(), err)
	assert.True(s.T(), <-status)
	assert.True(s.T(), s.board.state.IsSecurity)
	assert.Equal(s.T(), 1, s.adaptor.GetDigitalPinState(s.board.ledRed.Pin()))

	// Set security when security
	s.board.state.IsSecurity = true
	status = mock.WaitEvent(s.board.Eventer, EventSetSecurity, 1*time.Second)
	err = s.board.SetSecurity(context.Background())
	assert.NoError(s.T(), err)
	assert.False(s.T(), <-status)
	assert.True(s.T(), s.board.state.IsSecurity)

	// Unset security when security
	s.board.state.IsSecurity = true
	status = mock.WaitEvent(s.board.Eventer, EventUnsetSecurity, 1*time.Second)
	err = s.board.UnsetSecurity(context.Background())
	assert.NoError(s.T(), err)
	assert.True(s.T(), <-status)
	assert.False(s.T(), s.board.state.IsSecurity)
	assert.Equal(s.T(), 0, s.adaptor.GetDigitalPinState(s.board.ledRed.Pin()))

	// Unset security when no security
	s.board.state.IsSecurity = false
	status = mock.WaitEvent(s.board.Eventer, EventUnsetSecurity, 1*time.Second)
	err = s.board.UnsetSecurity(context.Background())
	assert.NoError(s.T(), err)
	assert.False(s.T(), <-status)
	assert.False(s.T(), s.board.state.IsSecurity)
}

func (s *DFPBoardTestSuite) TestSetUnsetEmergencyStop() {

	// Set emergency stop when no emergency stop
	s.board.state.IsEmergencyStopped = false
	status := mock.WaitEvent(s.board.Eventer, EventSetEmergencyStop, 1*time.Second)
	err := s.board.SetEmergencyStop(context.Background())
	assert.NoError(s.T(), err)
	assert.True(s.T(), <-status)
	assert.True(s.T(), s.board.state.IsEmergencyStopped)
	assert.Equal(s.T(), 1, s.adaptor.GetDigitalPinState(s.board.ledRed.Pin()))

	// Set emergency stop when emergency stop
	s.board.state.IsEmergencyStopped = true
	status = mock.WaitEvent(s.board.Eventer, EventSetEmergencyStop, 1*time.Second)
	err = s.board.SetEmergencyStop(context.Background())
	assert.NoError(s.T(), err)
	assert.False(s.T(), <-status)
	assert.True(s.T(), s.board.state.IsEmergencyStopped)

	// Unset emergency stop when emergency stop
	s.board.state.IsEmergencyStopped = true
	status = mock.WaitEvent(s.board.Eventer, EventUnsetEmergencyStop, 1*time.Second)
	err = s.board.UnsetEmergencyStop(context.Background())
	assert.NoError(s.T(), err)
	assert.True(s.T(), <-status)
	assert.False(s.T(), s.board.state.IsEmergencyStopped)
	assert.Equal(s.T(), 0, s.adaptor.GetDigitalPinState(s.board.ledRed.Pin()))

	// Unset emergency stop when no emergency stop
	s.board.state.IsEmergencyStopped = false
	status = mock.WaitEvent(s.board.Eventer, EventUnsetEmergencyStop, 1*time.Second)
	err = s.board.UnsetEmergencyStop(context.Background())
	assert.NoError(s.T(), err)
	assert.False(s.T(), <-status)
	assert.False(s.T(), s.board.state.IsEmergencyStopped)
}

func (s *DFPBoardTestSuite) TestForceWashing() {

	// When normal use case
	s.board.state.IsWashed = false
	status := mock.WaitEvent(s.board.Eventer, EventWash, 5*time.Second)
	err := s.board.ForceWashing(context.Background())
	assert.NoError(s.T(), err)
	assert.True(s.T(), <-status)

	// When is already on wash cycle, skip
	s.board.state.IsWashed = true
	status = mock.WaitEvent(s.board.Eventer, EventWash, 5*time.Second)
	err = s.board.ForceWashing(context.Background())
	assert.NoError(s.T(), err)
	// Temp delete test
	//assert.False(s.T(), <-status)

	// When emergency stop
	s.board.state.IsEmergencyStopped = true
	status = mock.WaitEvent(s.board.Eventer, EventWash, 5*time.Second)
	err = s.board.ForceWashing(context.Background())
	assert.NoError(s.T(), err)
	assert.False(s.T(), <-status)
}
