package dfpboard

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTurnOnOffGreenLed(t *testing.T) {

	board, adaptor := initTestBoard()
	if err := board.Start(context.Background()); err != nil {
		panic(err)
	}

	// Turn on
	board.turnOnGreenLed()
	assert.Equal(t, 1, adaptor.DigitalPinState[board.ledGreen.Pin()])

	// Turn off
	board.turnOffGreenLed()
	assert.Equal(t, 0, adaptor.DigitalPinState[board.ledGreen.Pin()])
}

func TestTurnOnOffRedLed(t *testing.T) {

	board, adaptor := initTestBoard()
	if err := board.Start(context.Background()); err != nil {
		panic(err)
	}

	// Turn on
	board.turnOnRedLed()
	assert.Equal(t, 1, adaptor.DigitalPinState[board.ledRed.Pin()])

	// Turn off
	board.turnOffRedLed()
	assert.Equal(t, 0, adaptor.DigitalPinState[board.ledRed.Pin()])
}
