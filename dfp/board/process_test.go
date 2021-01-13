package dfpboard

import (
	"context"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWash(t *testing.T) {

	sem := make(chan bool, 0)
	board, adaptor := initTestBoard()
	if err := board.Start(context.Background()); err != nil {
		panic(err)
	}

	// Control initial state
	assert.Equal(t, 0, adaptor.DigitalPinState[board.relayPump.Pin()])
	assert.Equal(t, 0, adaptor.DigitalPinState[board.relayDrum.Pin()])

	// Test that all wash routine run
	board.Once(NewWash, func(s interface{}) {
		sem <- true
	})
	board.config.StartWashingPumpBeforeWashing = 10
	board.config.WashingDuration = 10
	board.wash()

	// Test in internal routine
	time.Sleep(5 * time.Second)
	assert.Equal(t, 1, adaptor.DigitalPinState[board.relayPump.Pin()])
	assert.Equal(t, 0, adaptor.DigitalPinState[board.relayDrum.Pin()])

	time.Sleep(10 * time.Second)
	assert.Equal(t, 1, adaptor.DigitalPinState[board.relayPump.Pin()])
	assert.Equal(t, 1, adaptor.DigitalPinState[board.relayDrum.Pin()])

	select {
	case <-sem:
		assert.Equal(t, 0, adaptor.DigitalPinState[board.relayPump.Pin()])
		assert.Equal(t, 0, adaptor.DigitalPinState[board.relayDrum.Pin()])
		assert.Equal(t, 1, adaptor.DigitalPinState[board.ledGreen.Pin()])
	case <-time.After(10 * time.Second):
		t.Errorf("DFP wash not started")
	}

	// When stop during process
	board, adaptor = initTestBoard()
	if err := board.Start(context.Background()); err != nil {
		panic(err)
	}
	board.Once(NewWash, func(s interface{}) {
		sem <- true
	})
	board.config.StartWashingPumpBeforeWashing = 1
	board.config.WashingDuration = 10
	board.wash()

	time.Sleep(1 * time.Second)
	err := board.StopDFP(context.Background())
	assert.NoError(t, err)

	select {
	case <-sem:
		t.Errorf("DFP wash must be stopped")

	case <-time.After(5 * time.Second):
		assert.Equal(t, 0, adaptor.DigitalPinState[board.relayPump.Pin()])
		assert.Equal(t, 0, adaptor.DigitalPinState[board.relayDrum.Pin()])
		assert.Equal(t, 0, adaptor.DigitalPinState[board.ledGreen.Pin()])
	}
}

func TestWork(t *testing.T) {

	sem := make(chan bool, 0)
	board, adaptor := initTestBoard()
	if err := board.Start(context.Background()); err != nil {
		panic(err)
	}

	// Test button action on normal state
	adaptor.DigitalPinState[board.buttonForceDrum.Pin()] = 0
	board.Once(NewWash, func(s interface{}) {
		sem <- true
	})
	select {
	case <-sem:
		assert.Equal(t, 0, adaptor.DigitalPinState[board.relayPump.Pin()])
		assert.Equal(t, 0, adaptor.DigitalPinState[board.relayDrum.Pin()])
	case <-time.After(30 * time.Second):
		t.Errorf("DFP wash not started")
	}

}
