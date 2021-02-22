package dfpboard

import (
	"context"
	"time"

	"github.com/disaster37/gobot-fat/dfpconfig"
	"github.com/disaster37/gobot-fat/dfpstate"
	"github.com/disaster37/gobot-fat/mock"
	"github.com/disaster37/gobot-fat/models"
	"github.com/stretchr/testify/assert"
)

func (s *DFPBoardTestSuite) TestButtonStart() {

	// When DFP  stopped
	s.board.state.IsRunning = false
	status := mock.WaitEvent(s.board.Eventer, EventStartDFP, 1*time.Second)
	s.adaptor.SetDigitalPinState(s.board.buttonStart.Pin(), 0)
	assert.True(s.T(), <-status)
	assert.True(s.T(), s.board.state.IsRunning)
	s.adaptor.SetDigitalPinState(s.board.buttonStart.Pin(), 1)

	// When DFP alreay running
	status = mock.WaitEvent(s.board.Eventer, EventStartDFP, 1*time.Second)
	s.adaptor.SetDigitalPinState(s.board.buttonStart.Pin(), 0)
	assert.False(s.T(), <-status)
	assert.True(s.T(), s.board.state.IsRunning)

}

func (s *DFPBoardTestSuite) TestButtonStop() {

	// When DFP running
	status := mock.WaitEvent(s.board.Eventer, EventStopDFP, 1*time.Second)
	s.adaptor.SetDigitalPinState(s.board.buttonStop.Pin(), 0)
	assert.True(s.T(), <-status)
	assert.False(s.T(), s.board.state.IsRunning)
	s.adaptor.SetDigitalPinState(s.board.buttonStop.Pin(), 1)

	// When DFP already stopped
	status = mock.WaitEvent(s.board.Eventer, EventStopDFP, 1*time.Second)
	s.adaptor.SetDigitalPinState(s.board.buttonStop.Pin(), 0)
	assert.False(s.T(), <-status)
	assert.False(s.T(), s.board.state.IsRunning)
}

func (s *DFPBoardTestSuite) TestButtonForceDrum() {
	// Test button Force drum ON
	status := mock.WaitEvent(s.board.Eventer, EventNewInput, 1*time.Second)
	s.adaptor.SetDigitalPinState(s.board.buttonForceDrum.Pin(), 0)
	assert.True(s.T(), <-status)
	assert.Equal(s.T(), 1, s.adaptor.GetDigitalPinState(s.board.relayDrum.Pin()))

	// Test button Force drum OFF
	status = mock.WaitEvent(s.board.Eventer, EventNewInput, 1*time.Second)
	s.adaptor.SetDigitalPinState(s.board.buttonForceDrum.Pin(), 1)
	assert.True(s.T(), <-status)
	assert.Equal(s.T(), 0, s.adaptor.GetDigitalPinState(s.board.relayDrum.Pin()))
}

func (s *DFPBoardTestSuite) TestButtonForcePump() {
	// Test button Force pump ON
	status := mock.WaitEvent(s.board.Eventer, EventNewInput, 1*time.Second)
	s.adaptor.SetDigitalPinState(s.board.buttonForcePump.Pin(), 0)
	assert.True(s.T(), <-status)
	assert.Equal(s.T(), 1, s.adaptor.GetDigitalPinState(s.board.relayPump.Pin()))

	// Test button Force pump OFF
	status = mock.WaitEvent(s.board.Eventer, EventNewInput, 1*time.Second)
	s.adaptor.SetDigitalPinState(s.board.buttonForcePump.Pin(), 1)
	assert.True(s.T(), <-status)
	assert.Equal(s.T(), 0, s.adaptor.GetDigitalPinState(s.board.relayPump.Pin()))
}

func (s *DFPBoardTestSuite) TestButtonForceWash() {

	status := mock.WaitEvent(s.board.Eventer, EventWash, 5*time.Second)
	s.adaptor.SetDigitalPinState(s.board.buttonWash.Pin(), 0)
	assert.True(s.T(), <-status)
}

func (s *DFPBoardTestSuite) TestButtonEmergencyStop() {

	// Set emergency stop
	status := mock.WaitEvent(s.board.Eventer, EventSetEmergencyStop, 1*time.Second)
	s.adaptor.SetDigitalPinState(s.board.buttonEmergencyStop.Pin(), 0)
	assert.True(s.T(), <-status)
	assert.True(s.T(), s.board.state.IsEmergencyStopped)
	assert.Equal(s.T(), 1, s.adaptor.GetDigitalPinState(s.board.ledRed.Pin()))

	// Unset emergency stop
	status = mock.WaitEvent(s.board.Eventer, EventUnsetEmergencyStop, 1*time.Second)
	s.adaptor.SetDigitalPinState(s.board.buttonEmergencyStop.Pin(), 1)
	assert.True(s.T(), <-status)
	assert.False(s.T(), s.board.state.IsEmergencyStopped)
	assert.Equal(s.T(), 0, s.adaptor.GetDigitalPinState(s.board.ledRed.Pin()))
}

func (s *DFPBoardTestSuite) TestSecurityCaptor() {
	// Test secruity upper captor ON
	status := mock.WaitEvent(s.board.Eventer, EventSetSecurity, 1*time.Second)
	s.adaptor.SetDigitalPinState(s.board.captorSecurityUpper.Pin(), 1)
	assert.True(s.T(), <-status)
	assert.True(s.T(), s.board.state.IsSecurity)
	assert.Equal(s.T(), 1, s.adaptor.GetDigitalPinState(s.board.ledRed.Pin()))

	// Test secruity upper captor OFF
	time.Sleep(1 * time.Second)
	status = mock.WaitEvent(s.board.Eventer, EventUnsetSecurity, 1*time.Second)
	s.adaptor.SetDigitalPinState(s.board.captorSecurityUpper.Pin(), 0)
	assert.True(s.T(), <-status)
	assert.False(s.T(), s.board.state.IsSecurity)
	assert.Equal(s.T(), 0, s.adaptor.GetDigitalPinState(s.board.ledRed.Pin()))

	// Test secruity under captor ON
	status = mock.WaitEvent(s.board.Eventer, EventSetSecurity, 1*time.Second)
	s.adaptor.SetDigitalPinState(s.board.captorSecurityUnder.Pin(), 0)
	assert.True(s.T(), <-status)
	assert.True(s.T(), s.board.state.IsSecurity)
	assert.Equal(s.T(), 1, s.adaptor.GetDigitalPinState(s.board.ledRed.Pin()))

	// Test secruity under captor OFF
	time.Sleep(1 * time.Second)
	status = mock.WaitEvent(s.board.Eventer, EventUnsetSecurity, 1*time.Second)
	s.adaptor.SetDigitalPinState(s.board.captorSecurityUnder.Pin(), 1)
	assert.True(s.T(), <-status)
	assert.False(s.T(), s.board.state.IsSecurity)
	assert.Equal(s.T(), 0, s.adaptor.GetDigitalPinState(s.board.ledRed.Pin()))
}

func (s *DFPBoardTestSuite) TestWaterCaptor() {
	// Test water upper captor ON
	status := mock.WaitEvent(s.board.Eventer, EventWash, 5*time.Second)
	s.adaptor.SetDigitalPinState(s.board.captorWaterUpper.Pin(), 1)
	assert.True(s.T(), <-status)
	s.adaptor.SetDigitalPinState(s.board.captorWaterUpper.Pin(), 0)

	// Test water under captor ON
	status = mock.WaitEvent(s.board.Eventer, EventWash, 5*time.Second)
	s.adaptor.SetDigitalPinState(s.board.captorWaterUnder.Pin(), 0)
	assert.True(s.T(), <-status)

	// Don't run more wash after some time
	// First, run wash to update timer
	// Then control other wash not run
	s.board.config.WaitTimeBetweenWashing = 60
	status = mock.WaitEvent(s.board.Eventer, EventWash, 5*time.Second)
	s.adaptor.SetDigitalPinState(s.board.captorWaterUpper.Pin(), 1)
	assert.True(s.T(), <-status)
	s.adaptor.SetDigitalPinState(s.board.captorWaterUpper.Pin(), 0)
	status = mock.WaitEvent(s.board.Eventer, EventWash, 5*time.Second)
	s.adaptor.SetDigitalPinState(s.board.captorWaterUpper.Pin(), 1)
	assert.False(s.T(), <-status)

	// Sorry for that, it's just to be sure is not broke other tests.
	time.Sleep(10 * time.Second)

}

func (s *DFPBoardTestSuite) TestWorkUpdateConfig() {

	// Send update config event
	newConfig := &models.DFPConfig{
		TemperatureThresholdWhenFrozen: 10,
		ForceWashingDuration:           120,
		ForceWashingDurationWhenFrozen: 180,
	}
	status := mock.WaitEvent(s.board.Eventer, EventNewConfig, 1*time.Second)
	s.board.globalEventer.Publish(dfpconfig.NewDFPConfig, newConfig)
	assert.True(s.T(), <-status)
	assert.Equal(s.T(), newConfig, s.board.config)
}

func (s *DFPBoardTestSuite) TestWorkUpdateState() {

	// Send update config event
	newState := &models.DFPState{
		IsDisableSecurity: true,
	}
	status := mock.WaitEvent(s.board.Eventer, EventNewState, 1*time.Second)
	s.board.globalEventer.Publish(dfpstate.NewDFPState, newState)
	assert.True(s.T(), <-status)
	assert.Equal(s.T(), newState.IsDisableSecurity, s.board.state.IsDisableSecurity)
}

func (s *DFPBoardTestSuite) TestWash() {

	s.board.config.StartWashingPumpBeforeWashing = 1
	s.board.config.WashingDuration = 2

	// Wash
	status := mock.WaitEvent(s.board.Eventer, EventWash, 5*time.Second)
	s.board.wash()
	time.Sleep(500 * time.Millisecond)
	assert.Equal(s.T(), 1, s.adaptor.GetDigitalPinState(s.board.relayPump.Pin()))
	assert.Equal(s.T(), 0, s.adaptor.GetDigitalPinState(s.board.relayDrum.Pin()))
	time.Sleep(1 * time.Second)
	assert.Equal(s.T(), 1, s.adaptor.GetDigitalPinState(s.board.relayPump.Pin()))
	assert.Equal(s.T(), 1, s.adaptor.GetDigitalPinState(s.board.relayDrum.Pin()))
	assert.True(s.T(), <-status)
	assert.Equal(s.T(), 0, s.adaptor.GetDigitalPinState(s.board.relayPump.Pin()))
	assert.Equal(s.T(), 0, s.adaptor.GetDigitalPinState(s.board.relayDrum.Pin()))
	assert.Equal(s.T(), 1, s.adaptor.GetDigitalPinState(s.board.ledGreen.Pin()))

	// When stop during process
	status = mock.WaitEvent(s.board.Eventer, EventWash, 5*time.Second)
	s.board.wash()
	time.Sleep(1 * time.Second)
	err := s.board.StopDFP(context.Background())
	assert.NoError(s.T(), err)
	assert.False(s.T(), <-status)
	assert.Equal(s.T(), 0, s.adaptor.GetDigitalPinState(s.board.relayPump.Pin()))
	assert.Equal(s.T(), 0, s.adaptor.GetDigitalPinState(s.board.relayDrum.Pin()))
	assert.Equal(s.T(), 0, s.adaptor.GetDigitalPinState(s.board.ledGreen.Pin()))
	s.board.state.IsRunning = true

	// When emergency stop during process
	status = mock.WaitEvent(s.board.Eventer, EventWash, 5*time.Second)
	s.board.wash()
	time.Sleep(1 * time.Second)
	err = s.board.SetEmergencyStop(context.Background())
	assert.NoError(s.T(), err)
	assert.False(s.T(), <-status)
	assert.Equal(s.T(), 0, s.adaptor.GetDigitalPinState(s.board.relayPump.Pin()))
	assert.Equal(s.T(), 0, s.adaptor.GetDigitalPinState(s.board.relayDrum.Pin()))
	assert.Equal(s.T(), 1, s.adaptor.GetDigitalPinState(s.board.ledGreen.Pin()))
	s.board.state.IsEmergencyStopped = false

	// When security during process
	status = mock.WaitEvent(s.board.Eventer, EventWash, 5*time.Second)
	s.board.wash()
	time.Sleep(1 * time.Second)
	err = s.board.SetSecurity(context.Background())
	assert.NoError(s.T(), err)
	assert.False(s.T(), <-status)
	assert.Equal(s.T(), 0, s.adaptor.GetDigitalPinState(s.board.relayPump.Pin()))
	assert.Equal(s.T(), 0, s.adaptor.GetDigitalPinState(s.board.relayDrum.Pin()))
	assert.Equal(s.T(), 1, s.adaptor.GetDigitalPinState(s.board.ledGreen.Pin()))

}
