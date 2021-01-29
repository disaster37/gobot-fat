package tfpboard

import (
	"context"
	"errors"

	"github.com/stretchr/testify/assert"
)

func (s *TFPBoardTestSuite) TestStartStopPondPump() {

	// Start pond pomp when already started
	s.adaptor.SetDigitalPinState(s.board.relayPompPond.Pin(), 0)
	err := s.board.StartPondPump(context.Background())
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 0, s.adaptor.GetDigitalPinState(s.board.relayPompPond.Pin()))

	// Start pond pomp when stopped
	s.adaptor.SetDigitalPinState(s.board.relayPompPond.Pin(), 1)
	err = s.board.StartPondPump(context.Background())
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 0, s.adaptor.GetDigitalPinState(s.board.relayPompPond.Pin()))

	// Start pond pomp when error
	s.adaptor.SetError(errors.New("test"))
	err = s.board.StartPondPump(context.Background())
	assert.Error(s.T(), err)
	s.adaptor.SetError(nil)

	// Stop pond pomp when already stopped
	s.adaptor.SetDigitalPinState(s.board.relayPompPond.Pin(), 1)
	err = s.board.StopPondPump(context.Background())
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 1, s.adaptor.GetDigitalPinState(s.board.relayPompPond.Pin()))

	// Stop pond pomp when started
	s.adaptor.SetDigitalPinState(s.board.relayPompPond.Pin(), 0)
	err = s.board.StopPondPump(context.Background())
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 1, s.adaptor.GetDigitalPinState(s.board.relayPompPond.Pin()))

	// Stop pond pomp when error
	s.adaptor.SetError(errors.New("test"))
	err = s.board.StopPondPump(context.Background())
	assert.Error(s.T(), err)
	s.adaptor.SetError(nil)

	// Stop UVC1 and UVC2 when stop pond pump
	s.board.state.UVC1Running = true
	s.board.state.UVC2Running = true
	err = s.board.StopPondPump(context.Background())
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 1, s.adaptor.GetDigitalPinState(s.board.relayUVC1.Pin()))
	assert.Equal(s.T(), 1, s.adaptor.GetDigitalPinState(s.board.relayUVC2.Pin()))

	// Can't start pond pomp when emergency stop
	s.board.state.IsEmergencyStopped = true
	s.adaptor.SetDigitalPinState(s.board.relayPompPond.Pin(), 1)
	err = s.board.StartPondPump(context.Background())
	assert.Error(s.T(), err)
	assert.Equal(s.T(), 1, s.adaptor.GetDigitalPinState(s.board.relayPompPond.Pin()))

	// Can't start pond pomp when security
	s.board.state.IsEmergencyStopped = false
	s.board.state.IsSecurity = true
	s.adaptor.SetDigitalPinState(s.board.relayPompPond.Pin(), 1)
	err = s.board.StartPondPump(context.Background())
	assert.Error(s.T(), err)
	assert.Equal(s.T(), 1, s.adaptor.GetDigitalPinState(s.board.relayPompPond.Pin()))

	// Can start  pond pomp when security and disable security
	s.board.state.IsEmergencyStopped = false
	s.board.state.IsSecurity = true
	s.board.state.IsDisableSecurity = true
	s.adaptor.SetDigitalPinState(s.board.relayPompPond.Pin(), 1)
	err = s.board.StartPondPump(context.Background())
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 0, s.adaptor.GetDigitalPinState(s.board.relayPompPond.Pin()))

}

func (s *TFPBoardTestSuite) TestStartStopUVC1() {

	s.board.state.PondPumpRunning = true

	// Start UVC1 when already started
	s.adaptor.SetDigitalPinState(s.board.relayUVC1.Pin(), 0)
	err := s.board.StartUVC1(context.Background())
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 0, s.adaptor.GetDigitalPinState(s.board.relayUVC1.Pin()))

	// Start UVC1 when stopped
	s.adaptor.SetDigitalPinState(s.board.relayUVC1.Pin(), 1)
	err = s.board.StartUVC1(context.Background())
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 0, s.adaptor.GetDigitalPinState(s.board.relayUVC1.Pin()))

	// Start UVC1 when error
	s.adaptor.SetError(errors.New("test"))
	err = s.board.StartUVC1(context.Background())
	assert.Error(s.T(), err)
	s.adaptor.SetError(nil)

	// Stop UVC1 when already stopped
	s.adaptor.SetDigitalPinState(s.board.relayUVC1.Pin(), 1)
	err = s.board.StopUVC1(context.Background())
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 1, s.adaptor.GetDigitalPinState(s.board.relayUVC1.Pin()))

	// Stop UVC1 when started
	s.adaptor.SetDigitalPinState(s.board.relayUVC1.Pin(), 0)
	err = s.board.StopUVC1(context.Background())
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 1, s.adaptor.GetDigitalPinState(s.board.relayUVC1.Pin()))

	// Stop UVC1 when error
	s.adaptor.SetError(errors.New("test"))
	err = s.board.StopUVC1(context.Background())
	assert.Error(s.T(), err)
	s.adaptor.SetError(nil)

	// Can't start UVC1 when pond pump is stopped
	s.board.state.PondPumpRunning = false
	s.adaptor.SetDigitalPinState(s.board.relayUVC1.Pin(), 1)
	err = s.board.StartUVC1(context.Background())
	assert.Error(s.T(), err)
	assert.Equal(s.T(), 1, s.adaptor.GetDigitalPinState(s.board.relayUVC1.Pin()))

	// Can stop UVC1 when pond pump is started
	s.board.state.PondPumpRunning = false
	s.adaptor.SetDigitalPinState(s.board.relayUVC1.Pin(), 0)
	err = s.board.StopUVC1(context.Background())
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 1, s.adaptor.GetDigitalPinState(s.board.relayUVC1.Pin()))

	// Can't start UVC1 when emergency stop
	s.board.state.PondPumpRunning = true
	s.board.state.IsEmergencyStopped = true
	s.adaptor.SetDigitalPinState(s.board.relayUVC1.Pin(), 1)
	err = s.board.StartUVC1(context.Background())
	assert.Error(s.T(), err)
	assert.Equal(s.T(), 1, s.adaptor.GetDigitalPinState(s.board.relayUVC1.Pin()))

	// Can't start UVC1 when security
	s.board.state.IsEmergencyStopped = false
	s.board.state.IsSecurity = true
	s.adaptor.SetDigitalPinState(s.board.relayUVC1.Pin(), 1)
	err = s.board.StartUVC1(context.Background())
	assert.Error(s.T(), err)
	assert.Equal(s.T(), 1, s.adaptor.GetDigitalPinState(s.board.relayUVC1.Pin()))

	// Can start  UVC1 when security and disable security
	s.board.state.IsEmergencyStopped = false
	s.board.state.IsSecurity = true
	s.board.state.IsDisableSecurity = true
	s.adaptor.SetDigitalPinState(s.board.relayUVC1.Pin(), 1)
	err = s.board.StartUVC1(context.Background())
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 0, s.adaptor.GetDigitalPinState(s.board.relayUVC1.Pin()))
}

func (s *TFPBoardTestSuite) TestStartStopUVC2() {
	s.board.state.PondPumpRunning = true

	// Start UVC2 when already started
	s.adaptor.SetDigitalPinState(s.board.relayUVC2.Pin(), 0)
	err := s.board.StartUVC2(context.Background())
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 0, s.adaptor.GetDigitalPinState(s.board.relayUVC2.Pin()))

	// Start UVC2 when stopped
	s.adaptor.SetDigitalPinState(s.board.relayUVC2.Pin(), 1)
	err = s.board.StartUVC2(context.Background())
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 0, s.adaptor.GetDigitalPinState(s.board.relayUVC2.Pin()))

	// Start UVC2 when error
	s.adaptor.SetError(errors.New("test"))
	err = s.board.StartUVC2(context.Background())
	assert.Error(s.T(), err)
	s.adaptor.SetError(nil)

	// Stop UVC2 when already stopped
	s.adaptor.SetDigitalPinState(s.board.relayUVC2.Pin(), 1)
	err = s.board.StopUVC2(context.Background())
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 1, s.adaptor.GetDigitalPinState(s.board.relayUVC2.Pin()))

	// Stop UVC2 when started
	s.adaptor.SetDigitalPinState(s.board.relayUVC2.Pin(), 0)
	err = s.board.StopUVC2(context.Background())
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 1, s.adaptor.GetDigitalPinState(s.board.relayUVC2.Pin()))

	// Stop UVC2 when error
	s.adaptor.SetError(errors.New("test"))
	err = s.board.StopUVC2(context.Background())
	assert.Error(s.T(), err)
	s.adaptor.SetError(nil)

	// Can't start UVC2 when pond pump is stopped
	s.board.state.PondPumpRunning = false
	s.adaptor.SetDigitalPinState(s.board.relayUVC2.Pin(), 1)
	err = s.board.StartUVC2(context.Background())
	assert.Error(s.T(), err)
	assert.Equal(s.T(), 1, s.adaptor.GetDigitalPinState(s.board.relayUVC2.Pin()))

	// Can stop UVC2 when pond pump is started
	s.board.state.PondPumpRunning = false
	s.adaptor.SetDigitalPinState(s.board.relayUVC2.Pin(), 0)
	err = s.board.StopUVC2(context.Background())
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 1, s.adaptor.GetDigitalPinState(s.board.relayUVC2.Pin()))

	// Can't start UVC2 when emergency stop
	s.board.state.PondPumpRunning = true
	s.board.state.IsEmergencyStopped = true
	s.adaptor.SetDigitalPinState(s.board.relayUVC2.Pin(), 1)
	err = s.board.StartUVC2(context.Background())
	assert.Error(s.T(), err)
	assert.Equal(s.T(), 1, s.adaptor.GetDigitalPinState(s.board.relayUVC2.Pin()))

	// Can't start UVC2 when security
	s.board.state.IsEmergencyStopped = false
	s.board.state.IsSecurity = true
	s.adaptor.SetDigitalPinState(s.board.relayUVC2.Pin(), 1)
	err = s.board.StartUVC2(context.Background())
	assert.Error(s.T(), err)
	assert.Equal(s.T(), 1, s.adaptor.GetDigitalPinState(s.board.relayUVC2.Pin()))

	// Can start  UVC2 when security and disable security
	s.board.state.IsEmergencyStopped = false
	s.board.state.IsSecurity = true
	s.board.state.IsDisableSecurity = true
	s.adaptor.SetDigitalPinState(s.board.relayUVC2.Pin(), 1)
	err = s.board.StartUVC2(context.Background())
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 0, s.adaptor.GetDigitalPinState(s.board.relayUVC2.Pin()))
}

func (s *TFPBoardTestSuite) TestStartPondPumpWithUVC() {

	// When already started
	err := s.board.StartPondPumpWithUVC(context.Background())
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 0, s.adaptor.GetDigitalPinState(s.board.relayPompPond.Pin()))
	assert.Equal(s.T(), 0, s.adaptor.GetDigitalPinState(s.board.relayUVC1.Pin()))
	assert.Equal(s.T(), 0, s.adaptor.GetDigitalPinState(s.board.relayUVC2.Pin()))

	// When all stopped
	s.adaptor.SetDigitalPinState(s.board.relayPompPond.Pin(), 1)
	s.adaptor.SetDigitalPinState(s.board.relayUVC1.Pin(), 1)
	s.adaptor.SetDigitalPinState(s.board.relayUVC2.Pin(), 1)
	err = s.board.StartPondPumpWithUVC(context.Background())
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 0, s.adaptor.GetDigitalPinState(s.board.relayPompPond.Pin()))
	assert.Equal(s.T(), 0, s.adaptor.GetDigitalPinState(s.board.relayUVC1.Pin()))
	assert.Equal(s.T(), 0, s.adaptor.GetDigitalPinState(s.board.relayUVC2.Pin()))
}

func (s *TFPBoardTestSuite) TestStartStopPondBubble() {

	// Start pond bubble when already started
	s.adaptor.SetDigitalPinState(s.board.relayBubblePond.Pin(), 0)
	err := s.board.StartPondBubble(context.Background())
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 0, s.adaptor.GetDigitalPinState(s.board.relayBubblePond.Pin()))

	// Start pond bubble when stopped
	s.adaptor.SetDigitalPinState(s.board.relayBubblePond.Pin(), 1)
	err = s.board.StartPondBubble(context.Background())
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 0, s.adaptor.GetDigitalPinState(s.board.relayBubblePond.Pin()))

	// Start pond bubble when error
	s.adaptor.SetError(errors.New("test"))
	err = s.board.StartPondBubble(context.Background())
	assert.Error(s.T(), err)
	s.adaptor.SetError(nil)

	// Stop pond bubble when already stopped
	s.adaptor.SetDigitalPinState(s.board.relayBubblePond.Pin(), 1)
	err = s.board.StopPondBubble(context.Background())
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 1, s.adaptor.GetDigitalPinState(s.board.relayBubblePond.Pin()))

	// Stop pond bubble when started
	s.adaptor.SetDigitalPinState(s.board.relayBubblePond.Pin(), 0)
	err = s.board.StopPondBubble(context.Background())
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 1, s.adaptor.GetDigitalPinState(s.board.relayBubblePond.Pin()))

	// Stop pond bubble when error
	s.adaptor.SetError(errors.New("test"))
	err = s.board.StopPondBubble(context.Background())
	assert.Error(s.T(), err)
	s.adaptor.SetError(nil)

	// Can't start pond bubble when emergency stop
	s.board.state.IsEmergencyStopped = true
	s.adaptor.SetDigitalPinState(s.board.relayBubblePond.Pin(), 1)
	err = s.board.StartPondBubble(context.Background())
	assert.Error(s.T(), err)
	assert.Equal(s.T(), 1, s.adaptor.GetDigitalPinState(s.board.relayBubblePond.Pin()))

	// Can't start pond bubble when security
	s.board.state.IsEmergencyStopped = false
	s.board.state.IsSecurity = true
	s.adaptor.SetDigitalPinState(s.board.relayBubblePond.Pin(), 1)
	err = s.board.StartPondBubble(context.Background())
	assert.Error(s.T(), err)
	assert.Equal(s.T(), 1, s.adaptor.GetDigitalPinState(s.board.relayBubblePond.Pin()))

	// Can start  pond pomp when security and disable security
	s.board.state.IsEmergencyStopped = false
	s.board.state.IsSecurity = true
	s.board.state.IsDisableSecurity = true
	s.adaptor.SetDigitalPinState(s.board.relayBubblePond.Pin(), 1)
	err = s.board.StartPondBubble(context.Background())
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 0, s.adaptor.GetDigitalPinState(s.board.relayBubblePond.Pin()))
}

func (s *TFPBoardTestSuite) TestStartStopFilterBubble() {

	// Start filter bubble when already started
	s.adaptor.SetDigitalPinState(s.board.relayBubbleFilter.Pin(), 0)
	err := s.board.StartFilterBubble(context.Background())
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 0, s.adaptor.GetDigitalPinState(s.board.relayBubbleFilter.Pin()))

	// Start filter bubble when stopped
	s.adaptor.SetDigitalPinState(s.board.relayBubbleFilter.Pin(), 1)
	err = s.board.StartFilterBubble(context.Background())
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 0, s.adaptor.GetDigitalPinState(s.board.relayBubbleFilter.Pin()))

	// Start filter bubble when error
	s.adaptor.SetError(errors.New("test"))
	err = s.board.StartFilterBubble(context.Background())
	assert.Error(s.T(), err)
	s.adaptor.SetError(nil)

	// Stop filter bubble when already stopped
	s.adaptor.SetDigitalPinState(s.board.relayBubbleFilter.Pin(), 1)
	err = s.board.StopFilterBubble(context.Background())
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 1, s.adaptor.GetDigitalPinState(s.board.relayBubbleFilter.Pin()))

	// Stop filter bubble when started
	s.adaptor.SetDigitalPinState(s.board.relayBubbleFilter.Pin(), 0)
	err = s.board.StopFilterBubble(context.Background())
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 1, s.adaptor.GetDigitalPinState(s.board.relayBubbleFilter.Pin()))

	// Stop filter bubble when error
	s.adaptor.SetError(errors.New("test"))
	err = s.board.StopFilterBubble(context.Background())
	assert.Error(s.T(), err)
	s.adaptor.SetError(nil)

	// Can't start filter bubble when emergency stop
	s.board.state.IsEmergencyStopped = true
	s.adaptor.SetDigitalPinState(s.board.relayBubbleFilter.Pin(), 1)
	err = s.board.StartFilterBubble(context.Background())
	assert.Error(s.T(), err)
	assert.Equal(s.T(), 1, s.adaptor.GetDigitalPinState(s.board.relayBubbleFilter.Pin()))

	// Can't start filter bubble when security
	s.board.state.IsEmergencyStopped = false
	s.board.state.IsSecurity = true
	s.adaptor.SetDigitalPinState(s.board.relayBubbleFilter.Pin(), 1)
	err = s.board.StartFilterBubble(context.Background())
	assert.Error(s.T(), err)
	assert.Equal(s.T(), 1, s.adaptor.GetDigitalPinState(s.board.relayBubbleFilter.Pin()))

	// Can start  filter bubble when security and disable security
	s.board.state.IsEmergencyStopped = false
	s.board.state.IsSecurity = true
	s.board.state.IsDisableSecurity = true
	s.adaptor.SetDigitalPinState(s.board.relayBubbleFilter.Pin(), 1)
	err = s.board.StartFilterBubble(context.Background())
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 0, s.adaptor.GetDigitalPinState(s.board.relayBubbleFilter.Pin()))
}

func (s *TFPBoardTestSuite) TestStartStopWaterfallPomp() {

	// Start waterfall pomp when already started
	s.adaptor.SetDigitalPinState(s.board.relayPompWaterfall.Pin(), 1)
	err := s.board.StartWaterfallPump(context.Background())
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 1, s.adaptor.GetDigitalPinState(s.board.relayPompWaterfall.Pin()))

	// Start waterfall pomp when stopped
	s.adaptor.SetDigitalPinState(s.board.relayPompWaterfall.Pin(), 0)
	err = s.board.StartWaterfallPump(context.Background())
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 1, s.adaptor.GetDigitalPinState(s.board.relayPompWaterfall.Pin()))

	// Start waterfall pomp when error
	s.adaptor.SetError(errors.New("test"))
	err = s.board.StartWaterfallPump(context.Background())
	assert.Error(s.T(), err)
	s.adaptor.SetError(nil)

	// Stop waterfall pomp when already stopped
	s.adaptor.SetDigitalPinState(s.board.relayPompWaterfall.Pin(), 0)
	err = s.board.StopWaterfallPump(context.Background())
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 0, s.adaptor.GetDigitalPinState(s.board.relayPompWaterfall.Pin()))

	// Stop waterfall pomp when started
	s.adaptor.SetDigitalPinState(s.board.relayPompWaterfall.Pin(), 1)
	err = s.board.StopWaterfallPump(context.Background())
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 0, s.adaptor.GetDigitalPinState(s.board.relayPompWaterfall.Pin()))

	// Stop waterfall pomp when error
	s.adaptor.SetError(errors.New("test"))
	err = s.board.StopWaterfallPump(context.Background())
	assert.Error(s.T(), err)
	s.adaptor.SetError(nil)

	// Can't start waterfall pomp when emergency stop
	s.board.state.IsEmergencyStopped = true
	s.adaptor.SetDigitalPinState(s.board.relayPompWaterfall.Pin(), 0)
	err = s.board.StartWaterfallPump(context.Background())
	assert.Error(s.T(), err)
	assert.Equal(s.T(), 0, s.adaptor.GetDigitalPinState(s.board.relayPompWaterfall.Pin()))

	// Can't start waterfall pomp when security
	s.board.state.IsEmergencyStopped = false
	s.board.state.IsSecurity = true
	s.adaptor.SetDigitalPinState(s.board.relayPompWaterfall.Pin(), 0)
	err = s.board.StartWaterfallPump(context.Background())
	assert.Error(s.T(), err)
	assert.Equal(s.T(), 0, s.adaptor.GetDigitalPinState(s.board.relayPompWaterfall.Pin()))

	// Can start waterfall pomp when security and disable security
	s.board.state.IsEmergencyStopped = false
	s.board.state.IsSecurity = true
	s.board.state.IsDisableSecurity = true
	s.adaptor.SetDigitalPinState(s.board.relayPompWaterfall.Pin(), 0)
	err = s.board.StartWaterfallPump(context.Background())
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 1, s.adaptor.GetDigitalPinState(s.board.relayPompWaterfall.Pin()))
}
func (s *TFPBoardTestSuite) TestStopRelais() {

	s.adaptor.SetDigitalPinState(s.board.relayPompPond.Pin(), 0)
	s.adaptor.SetDigitalPinState(s.board.relayUVC1.Pin(), 0)
	s.adaptor.SetDigitalPinState(s.board.relayUVC2.Pin(), 0)
	s.adaptor.SetDigitalPinState(s.board.relayBubblePond.Pin(), 0)
	s.adaptor.SetDigitalPinState(s.board.relayBubbleFilter.Pin(), 0)
	s.adaptor.SetDigitalPinState(s.board.relayPompWaterfall.Pin(), 1)

	err := s.board.StopRelais(context.Background())
	assert.NoError(s.T(), err)
	assert.Equal(s.T(), 1, s.adaptor.GetDigitalPinState(s.board.relayPompPond.Pin()))
	assert.Equal(s.T(), 1, s.adaptor.GetDigitalPinState(s.board.relayUVC1.Pin()))
	assert.Equal(s.T(), 1, s.adaptor.GetDigitalPinState(s.board.relayUVC2.Pin()))
	assert.Equal(s.T(), 1, s.adaptor.GetDigitalPinState(s.board.relayBubblePond.Pin()))
	assert.Equal(s.T(), 1, s.adaptor.GetDigitalPinState(s.board.relayBubbleFilter.Pin()))
	assert.Equal(s.T(), 0, s.adaptor.GetDigitalPinState(s.board.relayPompWaterfall.Pin()))

}
