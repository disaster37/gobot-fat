package dfpboard

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestWash(t *testing.T) {

	sem := make(chan bool, 1)
	board, adaptor := initTestBoard()
	if err := board.Start(context.Background()); err != nil {
		panic(err)
	}

	// Control initial state
	assert.Equal(t, 0, adaptor.DigitalPinState[board.relayPump.Pin()])
	assert.Equal(t, 0, adaptor.DigitalPinState[board.relayDrum.Pin()])
	assert.Equal(t, 1, adaptor.DigitalPinState[board.ledGreen.Pin()])

	// Test that all wash routine run
	board.Once(NewWash, func(s interface{}) {
		sem <- true
	})
	board.config.StartWashingPumpBeforeWashing = 2
	board.config.WashingDuration = 2
	board.wash()

	// Test in internal routine
	time.Sleep(500 * time.Millisecond)
	assert.Equal(t, 1, adaptor.DigitalPinState[board.relayPump.Pin()])
	assert.Equal(t, 0, adaptor.DigitalPinState[board.relayDrum.Pin()])

	time.Sleep(2 * time.Second)
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

	board.Stop(context.Background())
}

func TestWorkButton(t *testing.T) {

	sem := make(chan bool, 1)
	board, adaptor := initTestBoard()
	if err := board.Start(context.Background()); err != nil {
		panic(err)
	}
	board.config.StartWashingPumpBeforeWashing = 1
	board.config.WashingDuration = 1
	board.config.WaitTimeBetweenWashing = 2

	// wait routine launch
	time.Sleep(5 * time.Second)
	if !board.isInitialized {
		panic(errors.New("Board not initialized"))
	}

	// Test button Force drum ON
	assert.Equal(t, 0, adaptor.DigitalPinState[board.relayDrum.Pin()])
	assert.Equal(t, 1, adaptor.DigitalPinState[board.buttonForceDrum.Pin()])
	board.Once(NewInput, func(s interface{}) {
		sem <- true
	})
	adaptor.DigitalPinState[board.buttonForceDrum.Pin()] = 0
	select {
	case <-sem:
		assert.Equal(t, 1, adaptor.DigitalPinState[board.relayDrum.Pin()])
	case <-time.After(5 * time.Second):
		t.Errorf("DFP force drum not started")
	}

	// Test button Force drum OFF
	assert.Equal(t, 1, adaptor.DigitalPinState[board.relayDrum.Pin()])
	assert.Equal(t, 0, adaptor.DigitalPinState[board.buttonForceDrum.Pin()])
	board.Once(NewInput, func(s interface{}) {
		sem <- true
	})
	adaptor.DigitalPinState[board.buttonForceDrum.Pin()] = 1
	select {
	case <-sem:
		assert.Equal(t, 0, adaptor.DigitalPinState[board.relayDrum.Pin()])
	case <-time.After(5 * time.Second):
		t.Errorf("DFP stop force drum not started")
	}

	// Test button Force pump ON
	assert.Equal(t, 0, adaptor.DigitalPinState[board.relayPump.Pin()])
	assert.Equal(t, 1, adaptor.DigitalPinState[board.buttonForcePump.Pin()])
	board.Once(NewInput, func(s interface{}) {
		sem <- true
	})
	adaptor.DigitalPinState[board.buttonForcePump.Pin()] = 0
	select {
	case <-sem:
		assert.Equal(t, 1, adaptor.DigitalPinState[board.relayPump.Pin()])
	case <-time.After(5 * time.Second):
		t.Errorf("DFP force pump not started")
	}

	// Test button Force pump OFF
	assert.Equal(t, 1, adaptor.DigitalPinState[board.relayPump.Pin()])
	assert.Equal(t, 0, adaptor.DigitalPinState[board.buttonForcePump.Pin()])
	board.Once(NewInput, func(s interface{}) {
		sem <- true
	})
	adaptor.DigitalPinState[board.buttonForcePump.Pin()] = 1
	select {
	case <-sem:
		assert.Equal(t, 0, adaptor.DigitalPinState[board.relayPump.Pin()])
	case <-time.After(5 * time.Second):
		t.Errorf("DFP force stop pump not started")
	}

	// Test button wash
	assert.Equal(t, 1, adaptor.DigitalPinState[board.buttonWash.Pin()])
	assert.False(t, board.state.IsWashed)
	assert.False(t, board.state.IsEmergencyStopped)
	board.Once(NewWash, func(s interface{}) {
		sem <- true
	})
	adaptor.DigitalPinState[board.buttonWash.Pin()] = 0
	select {
	case <-sem:
	case <-time.After(10 * time.Second):
		t.Errorf("DFP force wash not started")
	}

	// Test button stop
	assert.Equal(t, 1, adaptor.DigitalPinState[board.buttonStop.Pin()])
	assert.True(t, board.state.IsRunning)
	board.Once(NewInput, func(s interface{}) {
		sem <- true
	})
	adaptor.DigitalPinState[board.buttonStop.Pin()] = 0
	select {
	case <-sem:
		assert.False(t, board.state.IsRunning)
	case <-time.After(5 * time.Second):
		t.Errorf("DFP stop not started")
	}
	adaptor.DigitalPinState[board.buttonStop.Pin()] = 1

	// Test button start
	assert.Equal(t, 1, adaptor.DigitalPinState[board.buttonStart.Pin()])
	assert.False(t, board.state.IsRunning)
	board.Once(NewInput, func(s interface{}) {
		sem <- true
	})
	adaptor.DigitalPinState[board.buttonStart.Pin()] = 0
	select {
	case <-sem:
		assert.True(t, board.state.IsRunning)
	case <-time.After(5 * time.Second):
		t.Errorf("DFP start not started")
	}
	adaptor.DigitalPinState[board.buttonStart.Pin()] = 1

	// Test button emergency stop ON
	assert.Equal(t, 1, adaptor.DigitalPinState[board.buttonEmergencyStop.Pin()])
	assert.False(t, board.state.IsEmergencyStopped)
	board.Once(NewInput, func(s interface{}) {
		sem <- true
	})
	adaptor.DigitalPinState[board.buttonEmergencyStop.Pin()] = 0
	select {
	case <-sem:
		assert.True(t, board.state.IsEmergencyStopped)
	case <-time.After(5 * time.Second):
		t.Errorf("DFP emergency stop ON not started")
	}

	// Test button emergency stop OFF
	assert.Equal(t, 0, adaptor.DigitalPinState[board.buttonEmergencyStop.Pin()])
	assert.True(t, board.state.IsEmergencyStopped)
	board.Once(NewInput, func(s interface{}) {
		sem <- true
	})
	adaptor.DigitalPinState[board.buttonEmergencyStop.Pin()] = 1
	select {
	case <-sem:
		assert.False(t, board.state.IsEmergencyStopped)
	case <-time.After(5 * time.Second):
		t.Errorf("DFP emergency stop OFF not started")
	}

	board.Stop(context.Background())

}

func TestWorkCaptor(t *testing.T) {
	sem := make(chan bool, 1)
	board, adaptor := initTestBoard()
	if err := board.Start(context.Background()); err != nil {
		panic(err)
	}
	board.config.StartWashingPumpBeforeWashing = 1
	board.config.WashingDuration = 1
	board.config.WaitTimeBetweenWashing = 2

	// wait routine launch
	time.Sleep(5 * time.Second)
	if !board.isInitialized {
		panic(errors.New("Board not initialized"))
	}

	// Test secruity upper captor ON
	assert.Equal(t, 0, adaptor.DigitalPinState[board.captorSecurityUpper.Pin()])
	assert.False(t, board.state.IsSecurity)
	board.Once(NewInput, func(s interface{}) {
		sem <- true
	})
	adaptor.DigitalPinState[board.captorSecurityUpper.Pin()] = 1
	select {
	case <-sem:
		assert.True(t, board.state.IsSecurity)
	case <-time.After(10 * time.Second):
		t.Errorf("DFP security upper ON not started")
	}

	// Test secruity upper captor OFF
	assert.Equal(t, 1, adaptor.DigitalPinState[board.captorSecurityUpper.Pin()])
	assert.True(t, board.state.IsSecurity)
	board.Once(NewInput, func(s interface{}) {
		sem <- true
	})
	adaptor.DigitalPinState[board.captorSecurityUpper.Pin()] = 0
	select {
	case <-sem:
		assert.False(t, board.state.IsSecurity)
	case <-time.After(5 * time.Second):
		t.Errorf("DFP security upper OFF not started")
	}

	// Test secruity under captor ON
	assert.Equal(t, 1, adaptor.DigitalPinState[board.captorSecurityUnder.Pin()])
	assert.False(t, board.state.IsSecurity)
	assert.False(t, board.captorSecurityUnder.Active)
	board.Once(NewInput, func(s interface{}) {
		sem <- true
	})
	adaptor.DigitalPinState[board.captorSecurityUnder.Pin()] = 0
	select {
	case <-sem:
		assert.True(t, board.state.IsSecurity)
	case <-time.After(10 * time.Second):
		t.Errorf("DFP security under ON not started")
	}

	// Test secruity under captor OFF
	assert.Equal(t, 0, adaptor.DigitalPinState[board.captorSecurityUnder.Pin()])
	assert.True(t, board.state.IsSecurity)
	board.Once(NewInput, func(s interface{}) {
		sem <- true
	})
	adaptor.DigitalPinState[board.captorSecurityUnder.Pin()] = 1
	select {
	case <-sem:
		assert.False(t, board.state.IsSecurity)
	case <-time.After(5 * time.Second):
		t.Errorf("DFP security upper OFF not started")
	}

	// Test water upper captor ON
	assert.Equal(t, 0, adaptor.DigitalPinState[board.captorWaterUpper.Pin()])
	board.Once(NewInput, func(s interface{}) {
		sem <- true
	})
	adaptor.DigitalPinState[board.captorWaterUpper.Pin()] = 1
	select {
	case <-sem:
	case <-time.After(5 * time.Second):
		t.Errorf("DFP water upper ON not started")
	}
	adaptor.DigitalPinState[board.captorWaterUpper.Pin()] = 0

	// Test water under captor ON
	assert.Equal(t, 1, adaptor.DigitalPinState[board.captorWaterUnder.Pin()])
	board.Once(NewInput, func(s interface{}) {
		sem <- true
	})
	adaptor.DigitalPinState[board.captorWaterUnder.Pin()] = 0
	select {
	case <-sem:
	case <-time.After(5 * time.Second):
		t.Errorf("DFP water under ON not started")
	}
	adaptor.DigitalPinState[board.captorWaterUnder.Pin()] = 1

	board.Stop(context.Background())

}
