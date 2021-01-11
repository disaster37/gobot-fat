package tfpboard

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

func initTestBoard() (*TFPBoard, *helper.MockPlateform) {
	configHandler := viper.New()
	configHandler.Set("name", "test")
	configHandler.Set("pin.relay.pond_pomp", 1)
	configHandler.Set("pin.relay.uvc1", 2)
	configHandler.Set("pin.relay.uvc2", 3)
	configHandler.Set("pin.relay.pond_bubble", 4)
	configHandler.Set("pin.relay.filter_bubble", 5)
	configHandler.Set("pin.relay.waterfall_pomp", 6)
	configTFP := &models.TFPConfig{
		IsWaterfallAuto: false,
	}
	stateTFP := &models.TFPState{}
	eventer := gobot.NewEventer()
	eventUsecaseMock := eventusecase.NewMockEventBase()
	mockBoard := helper.NewMockPlateform()
	usecaseTFPMock := usecase.NewMockUsecasetBase()

	// Return the right type for drivers
	mockBoard.ValueReadState["isRebooted"] = false

	board := newTFP(mockBoard, configHandler, configTFP, stateTFP, eventUsecaseMock, usecaseTFPMock, eventer, 1*time.Millisecond)

	return board.(*TFPBoard), mockBoard
}

func TestStartStopIsOnline(t *testing.T) {
	board, adaptor := initTestBoard()

	// Normal start with all stopped on state
	err := board.Start(context.Background())
	assert.NoError(t, err)
	assert.True(t, board.IsOnline())
	assert.True(t, board.relayPompPond.Inverted)
	assert.False(t, board.relayPompWaterfall.Inverted)
	assert.True(t, board.relayUVC1.Inverted)
	assert.True(t, board.relayUVC2.Inverted)
	assert.True(t, board.relayBubbleFilter.Inverted)
	assert.True(t, board.relayBubblePond.Inverted)
	assert.Equal(t, 1, adaptor.DigitalPinState["1"])
	assert.Equal(t, 1, adaptor.DigitalPinState["2"])
	assert.Equal(t, 1, adaptor.DigitalPinState["3"])
	assert.Equal(t, 1, adaptor.DigitalPinState["4"])
	assert.Equal(t, 1, adaptor.DigitalPinState["5"])
	assert.Equal(t, 0, adaptor.DigitalPinState["6"])

	// Start with all started on state
	board, adaptor = initTestBoard()
	board.state.PondPumpRunning = true
	board.state.PondBubbleRunning = true
	board.state.FilterBubbleRunning = true
	board.state.WaterfallPumpRunning = true
	board.state.UVC1Running = true
	board.state.UVC2Running = true
	err = board.Start(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 0, adaptor.DigitalPinState["1"])
	assert.Equal(t, 0, adaptor.DigitalPinState["2"])
	assert.Equal(t, 0, adaptor.DigitalPinState["3"])
	assert.Equal(t, 0, adaptor.DigitalPinState["4"])
	assert.Equal(t, 0, adaptor.DigitalPinState["5"])
	assert.Equal(t, 1, adaptor.DigitalPinState["6"])

	// Stop
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
	assert.True(t, reflect.DeepEqual(models.TFPState{}, board.State()))
}
