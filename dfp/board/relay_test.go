package dfpboard

import (
	"context"
	"testing"
	"time"

	"github.com/disaster37/gobot-fat/helper"
	"github.com/stretchr/testify/assert"
)

func Test(t *testing.T) {

}

func TestStartStopDFP(t *testing.T) {
	var status chan bool

	// Start DFP when already started
	// Must not emit event
	board, adaptor := initTestBoard()
	if err := board.Start(context.Background()); err != nil {
		panic(err)
	}
	board.state.IsRunning = true
	status = helper.WaitEvent(board.Eventer, EventStartDFP, 1*time.Second)
	err := board.StartDFP(context.Background())
	assert.NoError(t, err)
	assert.False(t, <-status)
	assert.True(t, board.state.IsRunning)
	board.Stop(context.Background())

	// Start DFP when stopped
	// Must emit event
	board, adaptor = initTestBoard()
	if err := board.Start(context.Background()); err != nil {
		panic(err)
	}
	board.state.IsRunning = false
	status = helper.WaitEvent(board.Eventer, EventStartDFP, 1*time.Second)
	err = board.StartDFP(context.Background())
	assert.NoError(t, err)
	assert.True(t, <-status)
	assert.True(t, board.state.IsRunning)
	assert.Equal(t, 1, adaptor.DigitalPinState[board.ledGreen.Pin()])
	assert.Equal(t, 0, adaptor.DigitalPinState[board.ledRed.Pin()])
	board.Stop(context.Background())

	// Stop DFP when stopped
	// Must not emit event
	board, adaptor = initTestBoard()
	if err := board.Start(context.Background()); err != nil {
		panic(err)
	}
	board.state.IsRunning = false
	status = helper.WaitEvent(board.Eventer, EventStopDFP, 1*time.Second)
	err = board.StopDFP(context.Background())
	assert.NoError(t, err)
	assert.False(t, <-status)
	assert.False(t, board.state.IsRunning)
	board.Stop(context.Background())

	// Stop DFP when running
	// Must emit event
	board, adaptor = initTestBoard()
	if err := board.Start(context.Background()); err != nil {
		panic(err)
	}
	board.state.IsRunning = true
	status = helper.WaitEvent(board.Eventer, EventStopDFP, 1*time.Second)
	err = board.StopDFP(context.Background())
	assert.NoError(t, err)
	assert.True(t, <-status)
	assert.False(t, board.state.IsRunning)
	assert.Equal(t, 0, adaptor.DigitalPinState[board.ledGreen.Pin()])
	assert.Equal(t, 0, adaptor.DigitalPinState[board.ledRed.Pin()])
	board.Stop(context.Background())

}

func TestStartStopManualDrum(t *testing.T) {
	board, adaptor := initTestBoard()
	if err := board.Start(context.Background()); err != nil {
		panic(err)
	}

	// Start manual drum
	err := board.StartManualDrum(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 1, adaptor.DigitalPinState[board.relayDrum.Pin()])

	// Stop manual drum
	err = board.StopManualDrum(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 0, adaptor.DigitalPinState[board.relayDrum.Pin()])

	// Can start when emergency stop
	board.state.IsEmergencyStopped = true
	err = board.StartManualDrum(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 0, adaptor.DigitalPinState[board.relayDrum.Pin()])

	board.Stop(context.Background())
}

func TestStartStopManualPomp(t *testing.T) {

	board, adaptor := initTestBoard()
	if err := board.Start(context.Background()); err != nil {
		panic(err)
	}

	// Start manual pump
	err := board.StartManualPump(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 1, adaptor.DigitalPinState[board.relayPump.Pin()])

	// Stop manual pump
	err = board.StopManualPump(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 0, adaptor.DigitalPinState[board.relayPump.Pin()])

	// Can start when emergency stop
	board.state.IsEmergencyStopped = true
	err = board.StartManualPump(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 0, adaptor.DigitalPinState[board.relayPump.Pin()])

	board.Stop(context.Background())
}

func TestSetUnsetSecurity(t *testing.T) {
	var status chan bool

	// Set security when no security
	board, adaptor := initTestBoard()
	if err := board.Start(context.Background()); err != nil {
		panic(err)
	}
	board.state.IsSecurity = false
	status = helper.WaitEvent(board.Eventer, EventSetSecurity, 1*time.Second)
	err := board.SetSecurity(context.Background())
	assert.NoError(t, err)
	assert.True(t, <-status)
	assert.True(t, board.state.IsSecurity)
	assert.Equal(t, 1, adaptor.DigitalPinState[board.ledRed.Pin()])
	board.Stop(context.Background())

	// Set security when security
	board, adaptor = initTestBoard()
	if err := board.Start(context.Background()); err != nil {
		panic(err)
	}
	board.state.IsSecurity = true
	status = helper.WaitEvent(board.Eventer, EventSetSecurity, 1*time.Second)
	err = board.SetSecurity(context.Background())
	assert.NoError(t, err)
	assert.False(t, <-status)
	assert.True(t, board.state.IsSecurity)
	board.Stop(context.Background())

	// Unset security when security
	board, adaptor = initTestBoard()
	if err := board.Start(context.Background()); err != nil {
		panic(err)
	}
	board.state.IsSecurity = true
	status = helper.WaitEvent(board.Eventer, EventUnsetSecurity, 1*time.Second)
	err = board.UnsetSecurity(context.Background())
	assert.NoError(t, err)
	assert.True(t, <-status)
	assert.False(t, board.state.IsSecurity)
	assert.Equal(t, 0, adaptor.DigitalPinState[board.ledRed.Pin()])
	board.Stop(context.Background())

	// Unset security when no security
	board, adaptor = initTestBoard()
	if err := board.Start(context.Background()); err != nil {
		panic(err)
	}
	board.state.IsSecurity = false
	status = helper.WaitEvent(board.Eventer, EventUnsetSecurity, 1*time.Second)
	err = board.UnsetSecurity(context.Background())
	assert.NoError(t, err)
	assert.False(t, <-status)
	assert.False(t, board.state.IsSecurity)
	board.Stop(context.Background())
}

func TestSetUnsetEmergencySTop(t *testing.T) {
	var status chan bool

	// Set emergency stop when no emergency stop
	board, adaptor := initTestBoard()
	if err := board.Start(context.Background()); err != nil {
		panic(err)
	}
	board.state.IsEmergencyStopped = false
	status = helper.WaitEvent(board.Eventer, EventSetEmergencyStop, 1*time.Second)
	err := board.SetEmergencyStop(context.Background())
	assert.NoError(t, err)
	assert.True(t, <-status)
	assert.True(t, board.state.IsEmergencyStopped)
	assert.Equal(t, 1, adaptor.DigitalPinState[board.ledRed.Pin()])
	board.Stop(context.Background())

	// Set emergency stop when emergency stop
	board, adaptor = initTestBoard()
	if err := board.Start(context.Background()); err != nil {
		panic(err)
	}
	board.state.IsEmergencyStopped = true
	status = helper.WaitEvent(board.Eventer, EventSetEmergencyStop, 1*time.Second)
	err = board.SetEmergencyStop(context.Background())
	assert.False(t, <-status)
	assert.True(t, board.state.IsEmergencyStopped)
	board.Stop(context.Background())

	// Unset emergency stop when emergency stop
	board, adaptor = initTestBoard()
	if err := board.Start(context.Background()); err != nil {
		panic(err)
	}
	board.state.IsEmergencyStopped = true
	status = helper.WaitEvent(board.Eventer, EventUnsetEmergencyStop, 1*time.Second)
	err = board.UnsetEmergencyStop(context.Background())
	assert.NoError(t, err)
	assert.True(t, <-status)
	assert.False(t, board.state.IsEmergencyStopped)
	assert.Equal(t, 0, adaptor.DigitalPinState[board.ledRed.Pin()])
	board.Stop(context.Background())

	// Unset emergency stop when no emergency stop
	board, adaptor = initTestBoard()
	if err := board.Start(context.Background()); err != nil {
		panic(err)
	}
	board.state.IsEmergencyStopped = false
	status = helper.WaitEvent(board.Eventer, EventUnsetEmergencyStop, 1*time.Second)
	err = board.UnsetEmergencyStop(context.Background())
	assert.NoError(t, err)
	assert.False(t, <-status)
	assert.False(t, board.state.IsEmergencyStopped)
	board.Stop(context.Background())
}

func TestForceWashing(t *testing.T) {
	var status chan bool

	// When normal use case
	board, _ := initTestBoard()
	if err := board.Start(context.Background()); err != nil {
		panic(err)
	}
	board.state.IsWashed = false
	status = helper.WaitEvent(board.Eventer, EventWash, 5*time.Second)
	err := board.ForceWashing(context.Background())
	assert.NoError(t, err)
	assert.True(t, <-status)
	board.Stop(context.Background())

	// When is already on wash cycle, skip
	board, _ = initTestBoard()
	if err = board.Start(context.Background()); err != nil {
		panic(err)
	}
	board.state.IsWashed = true
	status = helper.WaitEvent(board.Eventer, EventWash, 5*time.Second)
	err = board.ForceWashing(context.Background())
	assert.NoError(t, err)
	assert.False(t, <-status)
	board.Stop(context.Background())

	// When emergency stop
	board, _ = initTestBoard()
	if err = board.Start(context.Background()); err != nil {
		panic(err)
	}
	board.state.IsEmergencyStopped = true
	status = helper.WaitEvent(board.Eventer, EventWash, 5*time.Second)
	err = board.ForceWashing(context.Background())
	assert.NoError(t, err)
	assert.False(t, <-status)
	board.Stop(context.Background())
}
