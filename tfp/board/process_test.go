package tfpboard

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/disaster37/gobot-arest/drivers/extra"
	"github.com/disaster37/gobot-fat/models"
	"github.com/disaster37/gobot-fat/tfpconfig"
	"github.com/disaster37/gobot-fat/tfpstate"
	"github.com/stretchr/testify/assert"
)

func TestWork(t *testing.T) {
	board, adapter := initTestBoard()
	sem := make(chan bool, 0)
	waitDuration := 100 * time.Millisecond

	// Start board and wait it's initialized
	board.Start(context.Background())
	isLoop := true
	go func() {
		select {
		case <-sem:
		case <-time.After(10 * time.Second):
			isLoop = false
			t.Errorf("Board not initialzed after 10s")
		}
	}()
	for !board.isInitialized && isLoop {
		time.Sleep(1 * time.Millisecond)
	}
	sem <- true

	// Check update local config on event
	newConfig := &models.TFPConfig{
		UVC1BlisterMaxTime: 1000,
	}
	board.Once(NewConfig, func(s interface{}) {
		assert.Equal(t, newConfig, board.config)
		sem <- true
	})
	board.globalEventer.Publish(tfpconfig.NewTFPConfig, newConfig)
	select {
	case <-sem:
	case <-time.After(waitDuration):
		t.Errorf("TFP config not updated")
	}

	// Check update local state on event
	newState := &models.TFPState{
		OzoneBlisterNbHour: 100,
		UVC1BlisterNbHour:  200,
		UVC2BlisterNbHour:  300,
		IsEmergencyStopped: true,
	}
	board.Once(NewState, func(s interface{}) {
		assert.NotEqual(t, newState, board.state)
		assert.Equal(t, newState.OzoneBlisterNbHour, board.state.OzoneBlisterNbHour)
		assert.Equal(t, newState.UVC1BlisterNbHour, board.state.UVC1BlisterNbHour)
		assert.Equal(t, newState.UVC2BlisterNbHour, board.state.UVC2BlisterNbHour)
		sem <- true
	})
	board.globalEventer.Publish(tfpstate.NewTFPState, newState)
	select {
	case <-sem:
	case <-time.After(waitDuration):
		t.Errorf("TFP state not updated")
	}

	// Check detect reboot
	isReconnectCalled := false
	board.Once(NewReboot, func(s interface{}) {
		assert.True(t, isReconnectCalled)
		sem <- true
	})
	adapter.TestReconnect(func() error {
		isReconnectCalled = true
		return nil
	})
	adapter.ValueReadState["isRebooted"] = true
	select {
	case <-sem:
	case <-time.After(waitDuration):
		t.Errorf("Reboot not detected")
	}

	// Check offline
	board.Once(NewOffline, func(s interface{}) {
		assert.False(t, board.IsOnline())
		sem <- true
	})
	board.valueRebooted.Publish(extra.Error, errors.New("test"))
	select {
	case <-sem:
	case <-time.After(waitDuration):
		t.Errorf("Offline not detected")
	}
}

func TestHandleBlisterTime(t *testing.T) {
	board, _ := initTestBoard()

	// When all stopped en mode none
	board.config.Mode = "none"
	board.handleBlisterTime(context.Background())
	assert.Equal(t, int64(0), board.state.OzoneBlisterNbHour)
	assert.Equal(t, int64(0), board.state.UVC1BlisterNbHour)
	assert.Equal(t, int64(0), board.state.UVC2BlisterNbHour)

	// When all stopped en mode uvc
	board.config.Mode = "uvc"
	board.handleBlisterTime(context.Background())
	assert.Equal(t, int64(0), board.state.OzoneBlisterNbHour)
	assert.Equal(t, int64(0), board.state.UVC1BlisterNbHour)
	assert.Equal(t, int64(0), board.state.UVC2BlisterNbHour)

	// When all stopped en mode ozone
	board.config.Mode = "ozone"
	board.handleBlisterTime(context.Background())
	assert.Equal(t, int64(0), board.state.OzoneBlisterNbHour)
	assert.Equal(t, int64(0), board.state.UVC1BlisterNbHour)
	assert.Equal(t, int64(0), board.state.UVC2BlisterNbHour)

	// When all started and mode none
	board.config.Mode = "none"
	board.state.UVC1Running = true
	board.state.UVC2Running = true
	board.handleBlisterTime(context.Background())
	assert.Equal(t, int64(0), board.state.OzoneBlisterNbHour)
	assert.Equal(t, int64(0), board.state.UVC1BlisterNbHour)
	assert.Equal(t, int64(0), board.state.UVC2BlisterNbHour)

	// When all started and mode uvc
	board.config.Mode = "uvc"
	board.state.UVC1Running = true
	board.state.UVC2Running = true
	board.handleBlisterTime(context.Background())
	assert.Equal(t, int64(0), board.state.OzoneBlisterNbHour)
	assert.Equal(t, int64(1), board.state.UVC1BlisterNbHour)
	assert.Equal(t, int64(1), board.state.UVC2BlisterNbHour)

	// When all started and mode ozone
	board.config.Mode = "ozone"
	board.state = &models.TFPState{}
	board.state.UVC1Running = true
	board.state.UVC2Running = true
	board.handleBlisterTime(context.Background())
	assert.Equal(t, int64(1), board.state.OzoneBlisterNbHour)
	assert.Equal(t, int64(1), board.state.UVC1BlisterNbHour)
	assert.Equal(t, int64(0), board.state.UVC2BlisterNbHour)
}

func TestHandleWaterfallAuto(t *testing.T) {
	board, adaptor := initTestBoard()
	startTime := time.Now().Add(-1 * time.Hour)
	endTime := time.Now().Add(1 * time.Hour)
	board.config.StartTimeWaterfall = startTime.Format("15:04")
	board.config.StopTimeWaterfall = endTime.Format("15:04")

	// When waterfall auto is off
	board.config.IsWaterfallAuto = false
	board.handleWaterfallAuto(context.Background())
	assert.Equal(t, 0, adaptor.DigitalPinState["6"])

	// When waterfall auto is ON and need to start waterfall
	board.config.IsWaterfallAuto = true
	board.handleWaterfallAuto(context.Background())
	assert.Equal(t, 1, adaptor.DigitalPinState["6"])

	// When waterfall auto is On and need to stop waterfall
	board.config.IsWaterfallAuto = true
	board.config.StopTimeWaterfall = board.config.StartTimeWaterfall
	board.state.AcknoledgeWaterfallAuto = true
	board.handleWaterfallAuto(context.Background())
	assert.Equal(t, 0, adaptor.DigitalPinState["6"])
}
