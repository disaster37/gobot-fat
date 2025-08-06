package tankboard

import (
	"time"

	"github.com/disaster37/gobot-fat/mock"
	"github.com/disaster37/gobot-fat/models"
	"github.com/disaster37/gobot-fat/usecase"
	"github.com/spf13/viper"
	"gobot.io/x/gobot/v2"
)

func initTestBoard() (*TankBoard, *mock.MockPlateform) {
	configHandler := viper.New()
	configHandler.Set("name", "test")
	configTank := &models.TankConfig{
		Depth:        100,
		LiterPerCm:   1,
		SensorHeight: 0,
	}
	eventer := gobot.NewEventer()
	eventUsecaseMock := usecase.NewMockUsecasetBase()
	mockBoard := mock.NewMockPlateform()

	// Return the right type for drivers
	mockBoard.SetValueReadState("isRebooted", false)
	mockBoard.SetValueReadState("distance", float64(0))

	board := newTank(mockBoard, configHandler, configTank, eventUsecaseMock, eventer, 1*time.Millisecond)

	return board.(*TankBoard), mockBoard
}
