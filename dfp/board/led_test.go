package dfpboard

import (
	"time"

	"github.com/stretchr/testify/assert"
)

func (s *DFPBoardTestSuite) TestTurnOnOffBlinkGreenLed() {

	// Turn on
	s.board.turnOnGreenLed()
	assert.Equal(s.T(), 1, s.adaptor.GetDigitalPinState(s.board.ledGreen.Pin()))

	// Turn off
	s.board.turnOffGreenLed()
	assert.Equal(s.T(), 0, s.adaptor.GetDigitalPinState(s.board.ledGreen.Pin()))

	// Blink
	lc := s.board.blinkGreenLed()
	time.Sleep(200 * time.Millisecond)
	currentState := s.board.ledGreen.State()
	time.Sleep(1 * time.Second)
	assert.NotEqual(s.T(), currentState, s.board.ledGreen.State())
	lc.Stop()
	lc.Wait()
	assert.False(s.T(), s.board.ledGreen.State())

}

func (s *DFPBoardTestSuite) TestTurnOnOffRedLed() {

	// Turn on
	s.board.turnOnRedLed()
	assert.Equal(s.T(), 1, s.adaptor.GetDigitalPinState(s.board.ledRed.Pin()))

	// Turn off
	s.board.turnOffRedLed()
	assert.Equal(s.T(), 0, s.adaptor.GetDigitalPinState(s.board.ledRed.Pin()))

	// Blink
	lc := s.board.blinkRedLed()
	time.Sleep(200 * time.Millisecond)
	currentState := s.board.ledRed.State()
	time.Sleep(1 * time.Second)
	assert.NotEqual(s.T(), currentState, s.board.ledRed.State())
	lc.Stop()
	lc.Wait()
	assert.False(s.T(), s.board.ledRed.State())
}
