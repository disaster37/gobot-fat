package tfpboard

import (
	"context"
	"errors"
	"time"

	"github.com/disaster37/gobot-arest/drivers/extra"
	"github.com/disaster37/gobot-fat/helper"
	"github.com/disaster37/gobot-fat/mock"
	"github.com/disaster37/gobot-fat/models"
	"github.com/disaster37/gobot-fat/tfpconfig"
	"github.com/disaster37/gobot-fat/tfpstate"
	"github.com/stretchr/testify/assert"
)

func (s *TFPBoardTestSuite) TestWork() {
	waitDuration := 100 * time.Millisecond

	// Check update local config on event
	newConfig := &models.TFPConfig{
		UVC1BlisterMaxTime: 1000,
	}
	status := mock.WaitEvent(s.board.Eventer, EventNewConfig, waitDuration)
	s.board.globalEventer.Publish(tfpconfig.NewTFPConfig, newConfig)
	assert.True(s.T(), <-status)

	// Check update local state on event
	newState := &models.TFPState{
		OzoneBlisterNbHour: 100,
		UVC1BlisterNbHour:  200,
		UVC2BlisterNbHour:  300,
		IsEmergencyStopped: true,
	}
	status = mock.WaitEvent(s.board.Eventer, EventNewState, waitDuration)
	s.board.globalEventer.Publish(tfpstate.NewTFPState, newState)
	assert.True(s.T(), <-status)
	assert.NotEqual(s.T(), newState, s.board.state)
	assert.Equal(s.T(), newState.OzoneBlisterNbHour, s.board.state.OzoneBlisterNbHour)
	assert.Equal(s.T(), newState.UVC1BlisterNbHour, s.board.state.UVC1BlisterNbHour)
	assert.Equal(s.T(), newState.UVC2BlisterNbHour, s.board.state.UVC2BlisterNbHour)

	// Check detect reboot
	isReconnectCalled := false
	s.adaptor.TestReconnect(func() error {
		isReconnectCalled = true
		return nil
	})
	status = mock.WaitEvent(s.board.Eventer, EventBoardReboot, waitDuration)
	s.adaptor.SetValueReadState("isRebooted", true)
	assert.True(s.T(), <-status)
	assert.True(s.T(), isReconnectCalled)

	// Check offline
	status = mock.WaitEvent(s.board.Eventer, EventBoardOffline, waitDuration)
	s.board.valueRebooted.Publish(extra.Error, errors.New("test"))
	assert.True(s.T(), <-status)
	assert.False(s.T(), s.board.IsOnline())

}

func (s *TFPBoardTestSuite) TestHandleBlisterTime() {

	// When all stopped en mode none
	s.board.config.Mode = "none"
	s.board.handleBlisterTime()
	assert.Equal(s.T(), int64(0), s.board.state.OzoneBlisterNbHour)
	assert.Equal(s.T(), int64(0), s.board.state.UVC1BlisterNbHour)
	assert.Equal(s.T(), int64(0), s.board.state.UVC2BlisterNbHour)

	// When all stopped en mode uvc
	s.board.config.Mode = "uvc"
	s.board.handleBlisterTime()
	assert.Equal(s.T(), int64(0), s.board.state.OzoneBlisterNbHour)
	assert.Equal(s.T(), int64(0), s.board.state.UVC1BlisterNbHour)
	assert.Equal(s.T(), int64(0), s.board.state.UVC2BlisterNbHour)

	// When all stopped en mode ozone
	s.board.config.Mode = "ozone"
	s.board.handleBlisterTime()
	assert.Equal(s.T(), int64(0), s.board.state.OzoneBlisterNbHour)
	assert.Equal(s.T(), int64(0), s.board.state.UVC1BlisterNbHour)
	assert.Equal(s.T(), int64(0), s.board.state.UVC2BlisterNbHour)

	// When all started and mode none
	s.board.config.Mode = "none"
	s.board.state.UVC1Running = true
	s.board.state.UVC2Running = true
	s.board.handleBlisterTime()
	assert.Equal(s.T(), int64(0), s.board.state.OzoneBlisterNbHour)
	assert.Equal(s.T(), int64(0), s.board.state.UVC1BlisterNbHour)
	assert.Equal(s.T(), int64(0), s.board.state.UVC2BlisterNbHour)

	// When all started and mode uvc
	s.board.config.Mode = "uvc"
	s.board.state.UVC1Running = true
	s.board.state.UVC2Running = true
	s.board.handleBlisterTime()
	assert.Equal(s.T(), int64(0), s.board.state.OzoneBlisterNbHour)
	assert.Equal(s.T(), int64(1), s.board.state.UVC1BlisterNbHour)
	assert.Equal(s.T(), int64(1), s.board.state.UVC2BlisterNbHour)

	// When all started and mode ozone
	s.board.config.Mode = "ozone"
	s.board.state = &models.TFPState{}
	s.board.state.UVC1Running = true
	s.board.state.UVC2Running = true
	s.board.handleBlisterTime()
	assert.Equal(s.T(), int64(1), s.board.state.OzoneBlisterNbHour)
	assert.Equal(s.T(), int64(1), s.board.state.UVC1BlisterNbHour)
	assert.Equal(s.T(), int64(0), s.board.state.UVC2BlisterNbHour)
}

func (s *TFPBoardTestSuite) TestHandleWaterfallAuto() {
	startTime := time.Now().Add(-1 * time.Hour)
	endTime := time.Now().Add(1 * time.Hour)
	s.board.config.StartTimeWaterfall = startTime.Format("15:04")
	s.board.config.StopTimeWaterfall = endTime.Format("15:04")

	// When waterfall auto is off
	s.board.config.IsWaterfallAuto = false
	s.board.handleWaterfallAuto()
	assert.Equal(s.T(), 0, s.adaptor.GetDigitalPinState(s.board.relayPompWaterfall.Pin()))

	// When waterfall auto is ON and need to start waterfall
	s.board.config.IsWaterfallAuto = true
	s.board.handleWaterfallAuto()
	assert.Equal(s.T(), 1, s.adaptor.GetDigitalPinState(s.board.relayPompWaterfall.Pin()))

	// When waterfall auto is On and need to stop waterfall
	s.board.config.IsWaterfallAuto = true
	s.board.config.StopTimeWaterfall = s.board.config.StartTimeWaterfall
	s.board.state.AcknoledgeWaterfallAuto = true
	s.board.handleWaterfallAuto()
	assert.Equal(s.T(), 0, s.adaptor.GetDigitalPinState(s.board.relayPompWaterfall.Pin()))
}

func (s *TFPBoardTestSuite) TestHandleSetUnsetEmergencyStop() {

	waitDuration := 100 * time.Millisecond

	// Set emergency stop
	err := s.board.StartPondPumpWithUVC(context.Background())
	if err != nil {
		s.T().Fatal(err)
	}
	err = s.board.StartFilterBubble(context.Background())
	if err != nil {
		s.T().Fatal(err)
	}
	err = s.board.StartPondBubble(context.Background())
	if err != nil {
		s.T().Fatal(err)
	}
	err = s.board.StartWaterfallPump(context.Background())
	if err != nil {
		s.T().Fatal(err)
	}
	status := mock.WaitEvent(s.board, EventSetEmergencyStop, waitDuration)
	s.board.globalEventer.Publish(helper.SetEmergencyStop, nil)
	assert.True(s.T(), <-status)
	assert.Equal(s.T(), 1, s.adaptor.GetDigitalPinState(s.board.relayPompPond.Pin()))
	assert.Equal(s.T(), 1, s.adaptor.GetDigitalPinState(s.board.relayUVC1.Pin()))
	assert.Equal(s.T(), 1, s.adaptor.GetDigitalPinState(s.board.relayUVC2.Pin()))
	assert.Equal(s.T(), 1, s.adaptor.GetDigitalPinState(s.board.relayBubbleFilter.Pin()))
	assert.Equal(s.T(), 1, s.adaptor.GetDigitalPinState(s.board.relayBubblePond.Pin()))
	assert.Equal(s.T(), 0, s.adaptor.GetDigitalPinState(s.board.relayPompWaterfall.Pin()))

	// Unset emergency stop
	status = mock.WaitEvent(s.board, EventUnsetEmergencyStop, waitDuration)
	s.board.globalEventer.Publish(helper.UnsetEmergencyStop, nil)
	assert.True(s.T(), <-status)
	assert.Equal(s.T(), 0, s.adaptor.GetDigitalPinState(s.board.relayPompPond.Pin()))
	assert.Equal(s.T(), 0, s.adaptor.GetDigitalPinState(s.board.relayUVC1.Pin()))
	assert.Equal(s.T(), 0, s.adaptor.GetDigitalPinState(s.board.relayUVC2.Pin()))
	assert.Equal(s.T(), 0, s.adaptor.GetDigitalPinState(s.board.relayBubbleFilter.Pin()))
	assert.Equal(s.T(), 0, s.adaptor.GetDigitalPinState(s.board.relayBubblePond.Pin()))
	assert.Equal(s.T(), 1, s.adaptor.GetDigitalPinState(s.board.relayPompWaterfall.Pin()))

}

func (s *TFPBoardTestSuite) TestHandleSetUnsetSecurity() {

	waitDuration := 100 * time.Millisecond

	// Set security
	err := s.board.StartPondPumpWithUVC(context.Background())
	if err != nil {
		s.T().Fatal(err)
	}
	err = s.board.StartFilterBubble(context.Background())
	if err != nil {
		s.T().Fatal(err)
	}
	err = s.board.StartPondBubble(context.Background())
	if err != nil {
		s.T().Fatal(err)
	}
	err = s.board.StartWaterfallPump(context.Background())
	if err != nil {
		s.T().Fatal(err)
	}
	status := mock.WaitEvent(s.board, EventSetSecurity, waitDuration)
	s.board.globalEventer.Publish(helper.SetSecurity, nil)
	assert.True(s.T(), <-status)
	assert.Equal(s.T(), 1, s.adaptor.GetDigitalPinState(s.board.relayPompPond.Pin()))
	assert.Equal(s.T(), 1, s.adaptor.GetDigitalPinState(s.board.relayUVC1.Pin()))
	assert.Equal(s.T(), 1, s.adaptor.GetDigitalPinState(s.board.relayUVC2.Pin()))
	assert.Equal(s.T(), 0, s.adaptor.GetDigitalPinState(s.board.relayBubbleFilter.Pin()))
	assert.Equal(s.T(), 0, s.adaptor.GetDigitalPinState(s.board.relayBubblePond.Pin()))
	assert.Equal(s.T(), 0, s.adaptor.GetDigitalPinState(s.board.relayPompWaterfall.Pin()))

	// Unset security
	status = mock.WaitEvent(s.board, EventUnsetSecurity, waitDuration)
	s.board.globalEventer.Publish(helper.UnsetSecurity, nil)
	assert.True(s.T(), <-status)
	assert.Equal(s.T(), 0, s.adaptor.GetDigitalPinState(s.board.relayPompPond.Pin()))
	assert.Equal(s.T(), 0, s.adaptor.GetDigitalPinState(s.board.relayUVC1.Pin()))
	assert.Equal(s.T(), 0, s.adaptor.GetDigitalPinState(s.board.relayUVC2.Pin()))
	assert.Equal(s.T(), 0, s.adaptor.GetDigitalPinState(s.board.relayBubbleFilter.Pin()))
	assert.Equal(s.T(), 0, s.adaptor.GetDigitalPinState(s.board.relayBubblePond.Pin()))
	assert.Equal(s.T(), 1, s.adaptor.GetDigitalPinState(s.board.relayPompWaterfall.Pin()))

}
