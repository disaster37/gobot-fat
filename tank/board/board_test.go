package tankboard

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/disaster37/gobot-arest/drivers/extra"
	"github.com/disaster37/gobot-fat/tankconfig"

	"github.com/disaster37/gobot-fat/event/usecase"
	"github.com/disaster37/gobot-fat/helper"
	"github.com/disaster37/gobot-fat/models"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"gobot.io/x/gobot"
)

func initTestBoard() (*TankBoard, *helper.MockPlateform) {
	configHandler := viper.New()
	configHandler.Set("name", "test")
	configTank := &models.TankConfig{
		Depth:        100,
		LiterPerCm:   1,
		SensorHeight: 0,
	}
	eventer := gobot.NewEventer()
	eventUsecaseMock := usecase.NewMockEventBase()
	mockBoard := helper.NewMockPlateform()

	// Return the right type for drivers
	mockBoard.ValueReadState["isRebooted"] = false
	mockBoard.ValueReadState["distance"] = float64(0)

	board := newTank(mockBoard, configHandler, configTank, eventUsecaseMock, eventer, 1*time.Millisecond)

	return board.(*TankBoard), mockBoard
}

func TestStartStopIsOnline(t *testing.T) {
	board, _ := initTestBoard()

	err := board.Start(context.Background())
	assert.NoError(t, err)
	assert.True(t, board.IsOnline())

	err = board.Stop(context.Background())
	assert.NoError(t, err)
	assert.False(t, board.IsOnline())
}

func TestName(t *testing.T) {
	board, _ := initTestBoard()
	assert.Equal(t, "test", board.Name())
}

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

	// Check update distance
	board.Once(NewDistance, func(s interface{}) {
		assert.Equal(t, int64(50), s)
		data, err := board.GetData(context.Background())
		assert.NoError(t, err)
		assert.Equal(t, 50, data.Level)
		sem <- true
	})
	adapter.ValueReadState["distance"] = float64(50)
	select {
	case <-sem:
	case <-time.After(waitDuration):
		t.Errorf("Read \"distance\" not updated")
	}

	// Check update local config on event
	newConfig := &models.TankConfig{
		Depth:        5,
		LiterPerCm:   5,
		SensorHeight: 5,
	}
	board.Once(NewConfig, func(s interface{}) {
		assert.Equal(t, newConfig, board.config)
		sem <- true
	})
	board.globalEventer.Publish(tankconfig.NewTankConfig, newConfig)
	select {
	case <-sem:
	case <-time.After(waitDuration):
		t.Errorf("Tank config not updated")
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
