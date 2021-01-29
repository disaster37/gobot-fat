package tankboard

import (
	"time"

	"github.com/disaster37/gobot-fat/helper"
	"github.com/disaster37/gobot-fat/models"
	"github.com/disaster37/gobot-fat/usecase"
	"github.com/spf13/viper"
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
	eventUsecaseMock := usecase.NewMockUsecasetBase()
	mockBoard := helper.NewMockPlateform()

	// Return the right type for drivers
	mockBoard.SetValueReadState("isRebooted", false)
	mockBoard.SetValueReadState("distance", float64(0))

	board := newTank(mockBoard, configHandler, configTank, eventUsecaseMock, eventer, 1*time.Millisecond)

	return board.(*TankBoard), mockBoard
}
