package dfpboard

import (
	"github.com/disaster37/gobot-fat/dfpstate"
	"github.com/disaster37/gobot-fat/mock"

	"github.com/disaster37/gobot-fat/dfpconfig"

	"github.com/disaster37/gobot-fat/models"
	"github.com/disaster37/gobot-fat/usecase"
	"github.com/spf13/viper"
	"gobot.io/x/gobot/v2"
)

func initTestBoard() (*DFPBoard, *mock.MockPlateform) {
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
		WaitTimeBeforeUnsetSecurity:    1,
		TemperatureSensorPolling:       1,
	}
	dfpState := &models.DFPState{
		IsRunning: true,
	}
	eventer := gobot.NewEventer()
	eventUsecaseMock := usecase.NewMockUsecasetBase()
	mockBoard := mock.NewMockPlateform()
	usecaseDFPMock := usecase.NewMockUsecasetBase()
	mockMail := mock.NewMockMail()

	mockBoard.SetInvertInitialPinState(configHandler.GetString("pin.captor.security_upper"))
	mockBoard.SetInvertInitialPinState(configHandler.GetString("pin.captor.water_upper"))

	eventer.AddEvent(dfpconfig.NewDFPConfig)
	eventer.AddEvent(dfpstate.NewDFPState)

	board := newDFP(mockBoard, configHandler, dfpConfig, dfpState, eventUsecaseMock, usecaseDFPMock, eventer, mockMail)

	return board.(*DFPBoard), mockBoard
}
