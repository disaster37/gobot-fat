package dfpboard

import (
	"context"
	"testing"
	"time"

	"github.com/disaster37/gobot-fat/dfpconfig"
	"github.com/disaster37/gobot-fat/dfpstate"
	"github.com/disaster37/gobot-fat/helper"
	"github.com/disaster37/gobot-fat/models"
	"github.com/stretchr/testify/assert"
)

func TestWash(t *testing.T) {

	var status chan bool
	board, adaptor := initTestBoard()
	if err := board.Start(context.Background()); err != nil {
		panic(err)
	}

	// Control initial state
	assert.Equal(t, 0, adaptor.DigitalPinState[board.relayPump.Pin()])
	assert.Equal(t, 0, adaptor.DigitalPinState[board.relayDrum.Pin()])
	assert.Equal(t, 1, adaptor.DigitalPinState[board.ledGreen.Pin()])

	// Test that all wash routine run
	board.config.StartWashingPumpBeforeWashing = 2
	board.config.WashingDuration = 2
	status = helper.WaitEvent(board.Eventer, EventWash, 10*time.Second)
	board.wash()

	// Test in internal routine
	time.Sleep(500 * time.Millisecond)
	assert.Equal(t, 1, adaptor.DigitalPinState[board.relayPump.Pin()])
	assert.Equal(t, 0, adaptor.DigitalPinState[board.relayDrum.Pin()])

	time.Sleep(2 * time.Second)
	assert.Equal(t, 1, adaptor.DigitalPinState[board.relayPump.Pin()])
	assert.Equal(t, 1, adaptor.DigitalPinState[board.relayDrum.Pin()])

	assert.True(t, <-status)
	assert.Equal(t, 0, adaptor.DigitalPinState[board.relayPump.Pin()])
	assert.Equal(t, 0, adaptor.DigitalPinState[board.relayDrum.Pin()])
	assert.Equal(t, 1, adaptor.DigitalPinState[board.ledGreen.Pin()])
	board.Stop(context.Background())

	// When stop during process
	board, adaptor = initTestBoard()
	if err := board.Start(context.Background()); err != nil {
		panic(err)
	}
	board.config.StartWashingPumpBeforeWashing = 1
	board.config.WashingDuration = 10
	status = helper.WaitEvent(board.Eventer, EventWash, 5*time.Second)
	board.wash()

	time.Sleep(1 * time.Second)
	err := board.StopDFP(context.Background())
	assert.NoError(t, err)

	assert.False(t, <-status)
	assert.Equal(t, 0, adaptor.DigitalPinState[board.relayPump.Pin()])
	assert.Equal(t, 0, adaptor.DigitalPinState[board.relayDrum.Pin()])
	assert.Equal(t, 0, adaptor.DigitalPinState[board.ledGreen.Pin()])
	board.Stop(context.Background())
}

func TestWorkButton(t *testing.T) {

	var status chan bool
	board, adaptor := initTestBoard()
	if err := board.Start(context.Background()); err != nil {
		panic(err)
	}

	// Test button Force drum ON
	assert.Equal(t, 0, adaptor.DigitalPinState[board.relayDrum.Pin()])
	assert.Equal(t, 1, adaptor.DigitalPinState[board.buttonForceDrum.Pin()])
	status = helper.WaitEvent(board.Eventer, EventNewInput, 1*time.Second)
	adaptor.DigitalPinState[board.buttonForceDrum.Pin()] = 0
	assert.True(t, <-status)
	assert.Equal(t, 1, adaptor.DigitalPinState[board.relayDrum.Pin()])

	// Test button Force drum OFF
	assert.Equal(t, 1, adaptor.DigitalPinState[board.relayDrum.Pin()])
	assert.Equal(t, 0, adaptor.DigitalPinState[board.buttonForceDrum.Pin()])
	status = helper.WaitEvent(board.Eventer, EventNewInput, 1*time.Second)
	adaptor.DigitalPinState[board.buttonForceDrum.Pin()] = 1
	assert.True(t, <-status)
	assert.Equal(t, 0, adaptor.DigitalPinState[board.relayDrum.Pin()])

	// Test button Force pump ON
	assert.Equal(t, 0, adaptor.DigitalPinState[board.relayPump.Pin()])
	assert.Equal(t, 1, adaptor.DigitalPinState[board.buttonForcePump.Pin()])
	status = helper.WaitEvent(board.Eventer, EventNewInput, 1*time.Second)
	adaptor.DigitalPinState[board.buttonForcePump.Pin()] = 0
	assert.True(t, <-status)
	assert.Equal(t, 1, adaptor.DigitalPinState[board.relayPump.Pin()])

	// Test button Force pump OFF
	assert.Equal(t, 1, adaptor.DigitalPinState[board.relayPump.Pin()])
	assert.Equal(t, 0, adaptor.DigitalPinState[board.buttonForcePump.Pin()])
	status = helper.WaitEvent(board.Eventer, EventNewInput, 1*time.Second)
	adaptor.DigitalPinState[board.buttonForcePump.Pin()] = 1
	assert.True(t, <-status)
	assert.Equal(t, 0, adaptor.DigitalPinState[board.relayPump.Pin()])

	// Test button wash
	assert.Equal(t, 1, adaptor.DigitalPinState[board.buttonWash.Pin()])
	assert.False(t, board.state.IsWashed)
	assert.False(t, board.state.IsEmergencyStopped)
	status = helper.WaitEvent(board.Eventer, EventWash, 5*time.Second)
	adaptor.DigitalPinState[board.buttonWash.Pin()] = 0
	assert.True(t, <-status)

	// Test button stop
	assert.Equal(t, 1, adaptor.DigitalPinState[board.buttonStop.Pin()])
	assert.True(t, board.state.IsRunning)
	status = helper.WaitEvent(board.Eventer, EventNewInput, 1*time.Second)
	adaptor.DigitalPinState[board.buttonStop.Pin()] = 0
	assert.True(t, <-status)
	assert.False(t, board.state.IsRunning)
	adaptor.DigitalPinState[board.buttonStop.Pin()] = 1

	// Test button start
	assert.Equal(t, 1, adaptor.DigitalPinState[board.buttonStart.Pin()])
	assert.False(t, board.state.IsRunning)
	status = helper.WaitEvent(board.Eventer, EventNewInput, 1*time.Second)
	adaptor.DigitalPinState[board.buttonStart.Pin()] = 0
	assert.True(t, <-status)
	assert.True(t, board.state.IsRunning)
	adaptor.DigitalPinState[board.buttonStart.Pin()] = 1

	// Test button emergency stop ON
	assert.Equal(t, 1, adaptor.DigitalPinState[board.buttonEmergencyStop.Pin()])
	assert.False(t, board.state.IsEmergencyStopped)
	status = helper.WaitEvent(board.Eventer, EventNewInput, 1*time.Second)
	adaptor.DigitalPinState[board.buttonEmergencyStop.Pin()] = 0
	assert.True(t, <-status)
	assert.True(t, board.state.IsEmergencyStopped)

	// Test button emergency stop OFF
	assert.Equal(t, 0, adaptor.DigitalPinState[board.buttonEmergencyStop.Pin()])
	assert.True(t, board.state.IsEmergencyStopped)
	status = helper.WaitEvent(board.Eventer, EventNewInput, 1*time.Second)
	adaptor.DigitalPinState[board.buttonEmergencyStop.Pin()] = 1
	assert.True(t, <-status)
	assert.False(t, board.state.IsEmergencyStopped)

	board.Stop(context.Background())

}

func TestWorkCaptor(t *testing.T) {
	var status chan bool
	board, adaptor := initTestBoard()
	if err := board.Start(context.Background()); err != nil {
		panic(err)
	}

	// Test secruity upper captor ON
	assert.Equal(t, 0, adaptor.DigitalPinState[board.captorSecurityUpper.Pin()])
	assert.False(t, board.state.IsSecurity)
	status = helper.WaitEvent(board.Eventer, EventNewInput, 1*time.Second)
	adaptor.DigitalPinState[board.captorSecurityUpper.Pin()] = 1
	assert.True(t, <-status)
	assert.True(t, board.state.IsSecurity)

	// Test secruity upper captor OFF
	assert.Equal(t, 1, adaptor.DigitalPinState[board.captorSecurityUpper.Pin()])
	assert.True(t, board.state.IsSecurity)
	status = helper.WaitEvent(board.Eventer, EventNewInput, 1*time.Second)
	adaptor.DigitalPinState[board.captorSecurityUpper.Pin()] = 0
	assert.True(t, <-status)
	assert.False(t, board.state.IsSecurity)

	// Test secruity under captor ON
	assert.Equal(t, 1, adaptor.DigitalPinState[board.captorSecurityUnder.Pin()])
	assert.False(t, board.state.IsSecurity)
	assert.False(t, board.captorSecurityUnder.Active)
	status = helper.WaitEvent(board.Eventer, EventNewInput, 1*time.Second)
	adaptor.DigitalPinState[board.captorSecurityUnder.Pin()] = 0
	assert.True(t, <-status)
	assert.True(t, board.state.IsSecurity)

	// Test secruity under captor OFF
	assert.Equal(t, 0, adaptor.DigitalPinState[board.captorSecurityUnder.Pin()])
	assert.True(t, board.state.IsSecurity)
	status = helper.WaitEvent(board.Eventer, EventNewInput, 1*time.Second)
	adaptor.DigitalPinState[board.captorSecurityUnder.Pin()] = 1
	assert.True(t, <-status)
	assert.False(t, board.state.IsSecurity)

	// Test water upper captor ON
	assert.Equal(t, 0, adaptor.DigitalPinState[board.captorWaterUpper.Pin()])
	status = helper.WaitEvent(board.Eventer, EventNewInput, 1*time.Second)
	adaptor.DigitalPinState[board.captorWaterUpper.Pin()] = 1
	assert.True(t, <-status)
	adaptor.DigitalPinState[board.captorWaterUpper.Pin()] = 0

	// Test water under captor ON
	assert.Equal(t, 1, adaptor.DigitalPinState[board.captorWaterUnder.Pin()])
	status = helper.WaitEvent(board.Eventer, EventNewInput, 1*time.Second)
	adaptor.DigitalPinState[board.captorWaterUnder.Pin()] = 0
	assert.True(t, <-status)
	adaptor.DigitalPinState[board.captorWaterUnder.Pin()] = 1

	board.Stop(context.Background())

}

func TestWorkUpdateConfig(t *testing.T) {
	var status chan bool
	board, _ := initTestBoard()
	if err := board.Start(context.Background()); err != nil {
		panic(err)
	}

	// Send update config event
	newConfig := &models.DFPConfig{
		TemperatureThresholdWhenFrozen: 10,
		ForceWashingDuration:           120,
		ForceWashingDurationWhenFrozen: 180,
	}
	status = helper.WaitEvent(board.Eventer, EventNewConfig, 10*time.Second)
	board.globalEventer.Publish(dfpconfig.NewDFPConfig, newConfig)
	assert.True(t, <-status)
	assert.Equal(t, newConfig, board.config)
	board.Stop(context.Background())
}

func TestWorkUpdateState(t *testing.T) {
	var status chan bool
	board, _ := initTestBoard()
	if err := board.Start(context.Background()); err != nil {
		panic(err)
	}

	// Send update config event
	newState := &models.DFPState{
		IsDisableSecurity: true,
	}
	status = helper.WaitEvent(board.Eventer, EventNewState, 10*time.Second)
	board.globalEventer.Publish(dfpstate.NewDFPState, newState)
	assert.True(t, <-status)
	assert.Equal(t, newState.IsDisableSecurity, board.state.IsDisableSecurity)
	board.Stop(context.Background())
}
