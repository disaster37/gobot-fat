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
)

// TankBoard manage all i/o on Tank
type TankBoard struct {
	gobot            *gobot.Robot
	board            *arest.Adaptor
	eventUsecase     event.Usecase
	configHandler    *viper.Viper
	config           *models.TankConfig
	data             *models.Tank
	name             string
	isOnline         bool
	valueRebooted    *extra.ValueDriver
	valueDistance    *extra.ValueDriver
	functionRebooted *extra.FunctionDriver
	gobot.Eventer
}

// NewTank create handler to manage Tank
func NewTank(configHandler *viper.Viper, config *models.TankConfig, eventUsecase event.Usecase, eventer gobot.Eventer) (tankHandler tank.Board) {

	//Create client
	c := arest.NewHTTPAdaptor(configHandler.GetString("url"))

	// Create struct
	tankBoard := &TankBoard{
		board:            c,
		eventUsecase:     eventUsecase,
		configHandler:    configHandler,
		name:             configHandler.GetString("name"),
		config:           config,
		data:             &models.Tank{},
		isOnline:         false,
		valueRebooted:    extra.NewValueDriver(c, "isRebooted", 10*time.Second),
		valueDistance:    extra.NewValueDriver(c, "distance", 10*time.Second),
		functionRebooted: extra.NewFunctionDriver(c, "acknoledgeRebooted", ""),
	}
	tankBoard.Eventer = eventer

	// Create drivers

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

	h.sendEvent(ctx, fmt.Sprintf("start_%s", h.name), "board", 0)

	return nil
}

// Stop stop the functions handle by board
func (h *TankBoard) Stop(ctx context.Context) (err error) {

	err = h.gobot.Stop()
	if err != nil {
		return err
	}

	h.sendEvent(ctx, fmt.Sprintf("stop_%s", h.name), "board", 0)

	return nil

}
