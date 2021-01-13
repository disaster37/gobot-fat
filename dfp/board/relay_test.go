package dfpboard

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestStartStopDFP(t *testing.T) {

	board, adaptor := initTestBoard()
	if err := board.Start(context.Background()); err != nil {
		panic(err)
	}

	// Start DFP
	err := board.StartDFP(context.Background())
	assert.NoError(t, err)
	assert.True(t, board.state.IsRunning)
	assert.Equal(t, 1, adaptor.DigitalPinState[board.ledGreen.Pin()])
	assert.Equal(t, 0, adaptor.DigitalPinState[board.ledRed.Pin()])

	// Stop DFP
	err = board.StopDFP(context.Background())
	assert.NoError(t, err)
	assert.False(t, board.state.IsRunning)
	assert.Equal(t, 0, adaptor.DigitalPinState[board.ledGreen.Pin()])
	assert.Equal(t, 0, adaptor.DigitalPinState[board.ledRed.Pin()])
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
}

func TestForceWashing(t *testing.T) {
	sem := make(chan bool, 0)
	board, _ := initTestBoard()
	if err := board.Start(context.Background()); err != nil {
		panic(err)
	}

	// When normal use case
	board.Once(NewWash, func(s interface{}) {
		sem <- true
	})
	err := board.ForceWashing(context.Background())
	assert.NoError(t, err)
	select {
	case <-sem:
	case <-time.After(10 * time.Second):
		t.Errorf("DFP wash not started")
	}

	// When is already on wash cycle, skip
	board, _ = initTestBoard()
	if err = board.Start(context.Background()); err != nil {
		panic(err)
	}
	board.state.IsWashed = true
	board.Once(NewWash, func(s interface{}) {
		t.Errorf("DFP wash must not started")
		sem <- true
	})
	err = board.ForceWashing(context.Background())
	assert.NoError(t, err)
	select {
	case <-sem:

	case <-time.After(10 * time.Second):
	}

	// When emergency stop
	board, _ = initTestBoard()
	if err = board.Start(context.Background()); err != nil {
		panic(err)
	}
	board.state.IsEmergencyStopped = true
	board.Once(NewWash, func(s interface{}) {
		t.Errorf("DFP wash must not started")
		sem <- true
	})
	err = board.ForceWashing(context.Background())
	assert.NoError(t, err)
	select {
	case <-sem:

	case <-time.After(10 * time.Second):
	}

}
