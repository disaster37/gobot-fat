package tfpboard

import (
	"time"

	"github.com/disaster37/gobot-fat/mock"
	"github.com/disaster37/gobot-fat/models"
	"github.com/disaster37/gobot-fat/usecase"
	"github.com/spf13/viper"
	"gobot.io/x/gobot"
)

func initTestBoard() (*TFPBoard, *mock.MockPlateform) {
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
	eventUsecaseMock := usecase.NewMockUsecasetBase()
	mockBoard := mock.NewMockPlateform()
	usecaseTFPMock := usecase.NewMockUsecasetBase()

	// Return the right type for drivers
	mockBoard.SetValueReadState("isRebooted", false)

	board := newTFP(mockBoard, configHandler, configTFP, stateTFP, eventUsecaseMock, usecaseTFPMock, eventer, 1*time.Millisecond)

	return board.(*TFPBoard), mockBoard
}
