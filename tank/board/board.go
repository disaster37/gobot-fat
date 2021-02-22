package tankboard

import (
	"context"
	"time"

	"github.com/disaster37/gobot-arest/drivers/extra"
	"github.com/disaster37/gobot-arest/plateforms/arest"
	"github.com/disaster37/gobot-fat/helper"
	"github.com/disaster37/gobot-fat/models"
	"github.com/disaster37/gobot-fat/tank"
	"github.com/disaster37/gobot-fat/usecase"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
)

const (
	EventNewDistance  = "new-distance"
	EventNewConfig    = "new-config"
	EventBoardReboot  = "board-reboot"
	EventBoardOffline = "board-offline"
	EventBoardStop    = "board-stop"
)

type TankAdaptor interface {
	gobot.Adaptor
	gpio.DigitalReader
	extra.ExtraReader
	Reconnect() error
}

// TankBoard manage all i/o on Tank
type TankBoard struct {
	gobot            *gobot.Robot
	board            TankAdaptor
	eventUsecase     usecase.UsecaseCRUD
	configHandler    *viper.Viper
	config           *models.TankConfig
	data             *models.Tank
	name             string
	isOnline         bool
	isInitialized    bool
	valueRebooted    *extra.ValueDriver
	valueDistance    *extra.ValueDriver
	functionRebooted *extra.FunctionDriver
	globalEventer    gobot.Eventer
	gobot.Eventer
}

// NewTank create handler to manage Tank
func NewTank(configHandler *viper.Viper, config *models.TankConfig, eventUsecase usecase.UsecaseCRUD, eventer gobot.Eventer) (tankHandler tank.Board) {

	//Create client
	c := arest.NewHTTPAdaptor(configHandler.GetString("url"))

	return newTank(c, configHandler, config, eventUsecase, eventer, 10*time.Second)

}

func newTank(board TankAdaptor, configHandler *viper.Viper, config *models.TankConfig, eventUsecase usecase.UsecaseCRUD, eventer gobot.Eventer, wait time.Duration) (tankHandler tank.Board) {

	// Create struct
	tankBoard := &TankBoard{
		board:            board,
		eventUsecase:     eventUsecase,
		configHandler:    configHandler,
		name:             configHandler.GetString("name"),
		config:           config,
		data:             &models.Tank{},
		isOnline:         false,
		isInitialized:    false,
		globalEventer:    eventer,
		valueRebooted:    extra.NewValueDriver(board, "isRebooted", wait),
		valueDistance:    extra.NewValueDriver(board, "distance", wait),
		functionRebooted: extra.NewFunctionDriver(board, "acknoledgeRebooted", ""),
		Eventer:          gobot.NewEventer(),
	}

	tankBoard.gobot = gobot.NewRobot(
		tankBoard.Name(),
		[]gobot.Connection{tankBoard.board},
		[]gobot.Device{
			tankBoard.valueRebooted,
			tankBoard.valueDistance,
			tankBoard.functionRebooted,
		},
		tankBoard.work,
	)

	tankBoard.AddEvent(EventNewDistance)
	tankBoard.AddEvent(EventNewConfig)
	tankBoard.AddEvent(EventBoardReboot)
	tankBoard.AddEvent(EventBoardOffline)
	tankBoard.AddEvent(EventBoardStop)

	log.Infof("Board %s initialized successfully", tankBoard.Name())

	return tankBoard

}

// Name permit to get the board name
func (h *TankBoard) Name() string {
	return h.name
}

// Board get board info as object
func (h *TankBoard) Board() *models.Board {
	return &models.Board{
		Name:     h.name,
		IsOnline: h.isOnline,
	}
}

// GetData permit to read current level on tank
func (h *TankBoard) GetData(ctx context.Context) (data *models.Tank, err error) {
	return h.data, nil
}

// IsOnline permit to know is board is online
func (h *TankBoard) IsOnline() bool {
	return h.isOnline
}

// Start run the main function
func (h *TankBoard) Start(ctx context.Context) (err error) {

	// Start connection on board
	err = h.board.Connect()
	if err != nil {
		return err
	}

	err = h.gobot.Start(false)
	if err != nil {
		return err
	}
	h.isOnline = true

	helper.SendEvent(ctx, h.eventUsecase, h.name, helper.KindEventStartBoard, h.name)

	return nil
}

// Stop stop the functions handle by board
func (h *TankBoard) Stop(ctx context.Context) (err error) {

	// Internal event
	h.Publish(EventBoardStop, nil)

	err = h.gobot.Stop()
	if err != nil {
		return err
	}

	h.isOnline = false
	h.isInitialized = false

	helper.SendEvent(ctx, h.eventUsecase, h.name, helper.KindEventStopBoard, h.name)

	return nil

}
