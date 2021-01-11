package tankboard

import (
	"context"
	"fmt"
	"time"

	"github.com/disaster37/gobot-arest/drivers/extra"
	"github.com/disaster37/gobot-arest/plateforms/arest"
	"github.com/disaster37/gobot-fat/event"
	"github.com/disaster37/gobot-fat/models"
	"github.com/disaster37/gobot-fat/tank"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
)

const (
	NewDistance = "new-distance"
	NewConfig   = "new-config"
	NewReboot   = "new-reboot"
	NewOffline  = "new-offline"
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
	eventUsecase     event.Usecase
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
func NewTank(configHandler *viper.Viper, config *models.TankConfig, eventUsecase event.Usecase, eventer gobot.Eventer) (tankHandler tank.Board) {

	//Create client
	c := arest.NewHTTPAdaptor(configHandler.GetString("url"))

	return newTank(c, configHandler, config, eventUsecase, eventer, 10*time.Second)

}

func newTank(board TankAdaptor, configHandler *viper.Viper, config *models.TankConfig, eventUsecase event.Usecase, eventer gobot.Eventer, wait time.Duration) (tankHandler tank.Board) {

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

	tankBoard.AddEvent(NewDistance)
	tankBoard.AddEvent(NewConfig)
	tankBoard.AddEvent(NewReboot)
	tankBoard.AddEvent(NewOffline)

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

	h.sendEvent(ctx, fmt.Sprintf("start_%s", h.name), "board", 0)

	return nil
}

// Stop stop the functions handle by board
func (h *TankBoard) Stop(ctx context.Context) (err error) {

	err = h.gobot.Stop()
	if err != nil {
		return err
	}

	h.isOnline = false
	h.isInitialized = false

	h.sendEvent(ctx, fmt.Sprintf("stop_%s", h.name), "board", 0)

	return nil

}
