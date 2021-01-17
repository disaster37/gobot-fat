package dfpboard

import (
	"context"
	"reflect"
	"testing"
	"time"

	eventusecase "github.com/disaster37/gobot-fat/event/usecase"
	"github.com/disaster37/gobot-fat/helper"
	"github.com/disaster37/gobot-fat/models"
	"github.com/disaster37/gobot-fat/usecase"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
	"gobot.io/x/gobot"
)

func initTestBoard() (*DFPBoard, *helper.MockPlateform) {
	configHandler := viper.New()
	configHandler.Set("name", "test")
	configHandler.Set("button_polling", 10)
	configHandler.Set("pin.relay.drum", 3)
	configHandler.Set("pin.relay.pomp", 5)
	configHandler.Set("pin.led.green", 7)
	configHandler.Set("pin.led.red", 8)
	configHandler.Set("pin.button.emergency_stop", 10)
	configHandler.Set("pin.button.start", 11)
	configHandler.Set("pin.button.stop", 12)
	configHandler.Set("pin.button.wash", 13)
	configHandler.Set("pin.button.force_drum", 15)
	configHandler.Set("pin.button.force_pump", 16)
	configHandler.Set("pin.captor.security_upper", 18)
	configHandler.Set("pin.captor.security_under", 19)
	configHandler.Set("pin.captor.water_upper", 21)
	configHandler.Set("pin.captor.water_under", 22)
	dfpConfig := &models.DFPConfig{
		ForceWashingDuration:           180,
		ForceWashingDurationWhenFrozen: 60,
		StartWashingPumpBeforeWashing:  1,
		WaitTimeBetweenWashing:         1,
		WashingDuration:                1,
	}
	dfpState := &models.DFPState{
		IsRunning: true,
	}
	eventer := gobot.NewEventer()
	eventUsecaseMock := eventusecase.NewMockEventBase()
	mockBoard := helper.NewMockPlateform()
	usecaseDFPMock := usecase.NewMockUsecasetBase()

	mockBoard.SetInvertInitialPinState(configHandler.GetString("pin.captor.security_upper"))
	mockBoard.SetInvertInitialPinState(configHandler.GetString("pin.captor.water_upper"))

	board := newDFP(mockBoard, configHandler, dfpConfig, dfpState, eventUsecaseMock, usecaseDFPMock, eventer)

	return board.(*DFPBoard), mockBoard
}

func TestStartStopIsOnline(t *testing.T) {
	sem := make(chan bool, 0)
	//waitDuration := 100 * time.Millisecond

	// Normal start with all stopped on state and running
	board, adaptor := initTestBoard()
	err := board.Start(context.Background())
	assert.NoError(t, err)

	assert.True(t, board.IsOnline())
	assert.Equal(t, 1, board.buttonForceDrum.DefaultState)
	assert.Equal(t, 1, board.buttonForcePump.DefaultState)
	assert.Equal(t, 1, board.buttonEmergencyStop.DefaultState)
	assert.Equal(t, 1, board.buttonStart.DefaultState)
	assert.Equal(t, 1, board.buttonStop.DefaultState)
	assert.Equal(t, 1, board.buttonWash.DefaultState)

	assert.Equal(t, 1, board.captorSecurityUnder.DefaultState)
	assert.Equal(t, 0, board.captorSecurityUpper.DefaultState)
	assert.Equal(t, 1, board.captorWaterUnder.DefaultState)
	assert.Equal(t, 0, board.captorWaterUpper.DefaultState)

	assert.Equal(t, 0, adaptor.DigitalPinState[board.relayDrum.Pin()])
	assert.Equal(t, 0, adaptor.DigitalPinState[board.relayPump.Pin()])
	assert.Equal(t, 1, adaptor.DigitalPinState[board.ledGreen.Pin()])
	assert.Equal(t, 0, adaptor.DigitalPinState[board.ledRed.Pin()])
	board.Stop(context.Background())

	// Normal start with all stopped on state and stopped
	board, adaptor = initTestBoard()
	board.state.IsRunning = false
	err = board.Start(context.Background())
	assert.NoError(t, err)

	assert.True(t, board.IsOnline())
	assert.Equal(t, 1, board.buttonForceDrum.DefaultState)
	assert.Equal(t, 1, board.buttonForcePump.DefaultState)
	assert.Equal(t, 1, board.buttonEmergencyStop.DefaultState)
	assert.Equal(t, 1, board.buttonStart.DefaultState)
	assert.Equal(t, 1, board.buttonStop.DefaultState)
	assert.Equal(t, 1, board.buttonWash.DefaultState)

	assert.Equal(t, 1, board.captorSecurityUnder.DefaultState)
	assert.Equal(t, 0, board.captorSecurityUpper.DefaultState)
	assert.Equal(t, 1, board.captorWaterUnder.DefaultState)
	assert.Equal(t, 0, board.captorWaterUpper.DefaultState)

	assert.Equal(t, 0, adaptor.DigitalPinState[board.relayDrum.Pin()])
	assert.Equal(t, 0, adaptor.DigitalPinState[board.relayPump.Pin()])
	assert.Equal(t, 0, adaptor.DigitalPinState[board.ledGreen.Pin()])
	assert.Equal(t, 1, adaptor.DigitalPinState[board.ledRed.Pin()])
	board.Stop(context.Background())

	// Start with wash and running
	board, adaptor = initTestBoard()
	board.state.IsWashed = true

	board.Once(NewWash, func(s interface{}) {
		sem <- true
	})
	err = board.Start(context.Background())
	assert.NoError(t, err)
	select {
	case <-sem:
	case <-time.After(10 * time.Second):
		t.Errorf("DFP wash not started")
	}
	board.Stop(context.Background())

	// Stop
	board, adaptor = initTestBoard()
	err = board.Start(context.Background())
	assert.NoError(t, err)
	err = board.Stop(context.Background())
	assert.NoError(t, err)
	assert.False(t, board.IsOnline())
}

func TestGetBoard(t *testing.T) {
	board, _ := initTestBoard()
	assert.Equal(t, "test", board.Board().Name)
	assert.False(t, board.Board().IsOnline)
}

func TestName(t *testing.T) {
	board, _ := initTestBoard()
	assert.Equal(t, "test", board.Name())
}

func TestState(t *testing.T) {
	board, _ := initTestBoard()
	assert.True(t, reflect.DeepEqual(models.DFPState{IsRunning: true}, board.State()))
}
