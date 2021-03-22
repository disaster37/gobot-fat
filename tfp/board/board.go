package tfpboard

import (
	"context"
	"time"

	"github.com/disaster37/gobot-arest/drivers/extra"
	"github.com/disaster37/gobot-arest/plateforms/arest"
	"github.com/disaster37/gobot-fat/helper"
	"github.com/disaster37/gobot-fat/models"
	"github.com/disaster37/gobot-fat/tfp"
	"github.com/disaster37/gobot-fat/usecase"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
)

const (
	EventBoardStop            = "board-stop"
	EventBoardReboot          = "board-reboot"
	EventBoardOffline         = "board-offline"
	EventNewState             = "new-state"
	EventNewConfig            = "new-config"
	EventSetSecurity          = "set-security"
	EventUnsetSecurity        = "unset-security"
	EventSetDisableSecurity   = "set-disable-security"
	EventUnsetDisableSecurity = "unset-disable-security"
	EventSetEmergencyStop     = "set-emergency-stop"
	EventUnsetEmergencyStop   = "unset-emergency-stop"
)

// TFPAdaptor is TFP board interface
type TFPAdaptor interface {
	gobot.Adaptor
	gpio.DigitalReader
	gpio.DigitalWriter
	extra.ExtraReader
	Reconnect() error
}

// TFPBoard manage TFP board
type TFPBoard struct {
	gobot              *gobot.Robot
	name               string
	board              TFPAdaptor
	state              *models.TFPState
	config             *models.TFPConfig
	eventUsecase       usecase.UsecaseCRUD
	stateUsecase       usecase.UsecaseCRUD
	relayPompPond      *gpio.RelayDriver
	relayPompWaterfall *gpio.RelayDriver
	relayBubblePond    *gpio.RelayDriver
	relayBubbleFilter  *gpio.RelayDriver
	relayUVC1          *gpio.RelayDriver
	relayUVC2          *gpio.RelayDriver
	valueRebooted      *extra.ValueDriver
	functionRebooted   *extra.FunctionDriver
	configHandler      *viper.Viper
	isOnline           bool
	isInitialized      bool
	schedulingRoutines []*time.Ticker
	globalEventer      gobot.Eventer
	gobot.Eventer
}

// NewTFP create board to manage TFP
func NewTFP(configHandler *viper.Viper, config *models.TFPConfig, state *models.TFPState, eventUsecase usecase.UsecaseCRUD, tfpStateUsecase usecase.UsecaseCRUD, eventer gobot.Eventer) (tfpBoard tfp.Board) {

	//Create client
	c := arest.NewHTTPAdaptor(configHandler.GetString("url"))

	return newTFP(c, configHandler, config, state, eventUsecase, tfpStateUsecase, eventer, 1*time.Second)

}

func newTFP(board TFPAdaptor, configHandler *viper.Viper, config *models.TFPConfig, state *models.TFPState, eventUsecase usecase.UsecaseCRUD, tfpStateUsecase usecase.UsecaseCRUD, eventer gobot.Eventer, wait time.Duration) (tfpHandler tfp.Board) {

	// Create struct
	tfpBoard := &TFPBoard{
		board:              board,
		eventUsecase:       eventUsecase,
		stateUsecase:       tfpStateUsecase,
		configHandler:      configHandler,
		name:               configHandler.GetString("name"),
		config:             config,
		state:              state,
		isOnline:           false,
		isInitialized:      false,
		globalEventer:      eventer,
		relayPompPond:      gpio.NewRelayDriver(board, configHandler.GetString("pin.relay.pond_pomp")),
		relayPompWaterfall: gpio.NewRelayDriver(board, configHandler.GetString("pin.relay.waterfall_pomp")),
		relayUVC1:          gpio.NewRelayDriver(board, configHandler.GetString("pin.relay.uvc1")),
		relayUVC2:          gpio.NewRelayDriver(board, configHandler.GetString("pin.relay.uvc2")),
		relayBubbleFilter:  gpio.NewRelayDriver(board, configHandler.GetString("pin.relay.filter_bubble")),
		relayBubblePond:    gpio.NewRelayDriver(board, configHandler.GetString("pin.relay.pond_bubble")),
		valueRebooted:      extra.NewValueDriver(board, "isRebooted", wait),
		functionRebooted:   extra.NewFunctionDriver(board, "acknoledgeRebooted", ""),
		Eventer:            gobot.NewEventer(),
		schedulingRoutines: make([]*time.Ticker, 0),
	}

	tfpBoard.gobot = gobot.NewRobot(
		tfpBoard.Name(),
		[]gobot.Connection{tfpBoard.board},
		[]gobot.Device{
			tfpBoard.relayBubbleFilter,
			tfpBoard.relayBubblePond,
			tfpBoard.relayPompPond,
			tfpBoard.relayPompWaterfall,
			tfpBoard.relayUVC1,
			tfpBoard.relayUVC2,
			tfpBoard.valueRebooted,
			tfpBoard.functionRebooted,
		},
		tfpBoard.work,
	)

	tfpBoard.AddEvent(EventNewConfig)
	tfpBoard.AddEvent(EventNewState)
	tfpBoard.AddEvent(EventBoardReboot)
	tfpBoard.AddEvent(EventBoardOffline)
	tfpBoard.AddEvent(EventBoardStop)
	tfpBoard.AddEvent(EventNewConfig)
	tfpBoard.AddEvent(EventNewState)
	tfpBoard.AddEvent(EventSetDisableSecurity)
	tfpBoard.AddEvent(EventSetEmergencyStop)
	tfpBoard.AddEvent(EventSetSecurity)
	tfpBoard.AddEvent(EventUnsetDisableSecurity)
	tfpBoard.AddEvent(EventUnsetEmergencyStop)
	tfpBoard.AddEvent(EventUnsetSecurity)

	log.Infof("Board %s initialized successfully", tfpBoard.Name())

	return tfpBoard

}

// Name permit to get the board name
func (h *TFPBoard) Name() string {
	return h.name
}

// Board get board info as object
func (h *TFPBoard) Board() *models.Board {
	return &models.Board{
		Name:     h.name,
		IsOnline: h.isOnline,
	}
}

// IsOnline permit to know is board is online
func (h *TFPBoard) IsOnline() bool {
	return h.isOnline
}

// Start run the main function
func (h *TFPBoard) Start(ctx context.Context) (err error) {

	// Start connection on board
	err = h.board.Connect()
	if err != nil {
		return err
	}

	// Relay relayPompPond is Normaly Close
	h.relayPompPond.Inverted = true
	if h.state.PondPumpRunning {
		err = h.relayPompPond.On()
	} else {
		err = h.relayPompPond.Off()
	}
	if err != nil {
		return err
	}

	// Relay relayUVC1 is Normaly Close
	h.relayUVC1.Inverted = true
	if h.state.UVC1Running {
		err = h.relayUVC1.On()
	} else {
		err = h.relayUVC1.Off()
	}
	if err != nil {
		return err
	}

	// Relay relayUVC2 is Normaly Close
	h.relayUVC2.Inverted = true
	if h.state.UVC2Running {
		err = h.relayUVC2.On()
	} else {
		err = h.relayUVC2.Off()
	}
	if err != nil {
		return err
	}

	// Relay relayBubblePond  is Normaly Close
	h.relayBubblePond.Inverted = true
	if h.state.PondBubbleRunning {
		err = h.relayBubblePond.On()
	} else {
		err = h.relayBubblePond.Off()
	}
	if err != nil {
		return err
	}

	// Relay relayBubbleFilter is Normaly Close
	h.relayBubbleFilter.Inverted = true
	if h.state.FilterBubbleRunning {
		err = h.relayBubbleFilter.On()
	} else {
		err = h.relayBubbleFilter.Off()
	}
	if err != nil {
		return err
	}

	// Relay relayPompWaterfall is Normaly Open
	if h.state.WaterfallPumpRunning {
		err = h.relayPompWaterfall.On()
	} else {
		err = h.relayPompWaterfall.Off()
	}
	if err != nil {
		return err
	}

	err = h.gobot.Start(false)
	if err != nil {
		return err
	}
	h.isOnline = true

	// Send event
	helper.SendEvent(ctx, h.eventUsecase, h.name, helper.KindEventStartBoard, h.name)

	return nil
}

// Stop stop the functions handle by board
func (h *TFPBoard) Stop(ctx context.Context) (err error) {

	// Internal event
	h.Publish(EventBoardStop, nil)

	// Stop scheduling routines
	for _, ticker := range h.schedulingRoutines {
		ticker.Stop()
	}
	h.schedulingRoutines = make([]*time.Ticker, 0)

	err = h.gobot.Stop()
	if err != nil {
		return err
	}

	h.isOnline = false
	h.isInitialized = false

	// Send event
	helper.SendEvent(ctx, h.eventUsecase, h.name, helper.KindEventStopBoard, h.name)

	return nil

}

// State return the current state
func (h *TFPBoard) State() models.TFPState {
	return *h.state
}

// State return the current state
func (h *TFPBoard) Config() models.TFPConfig {
	return *h.config
}

// IO return current IO state
func (h *TFPBoard) IO() models.TFPIO {
	io := models.TFPIO{}

	// Relais state
	io.PondPumpRelay = h.relayPompPond.State()
	io.WaterfallPumpRelay = h.relayPompWaterfall.State()
	io.UVC1Relay = h.relayUVC1.State()
	io.UVC2Relay = h.relayUVC2.State()
	io.FilterBubble = h.relayBubbleFilter.State()
	io.PondBubble = h.relayBubblePond.State()

	return io
}
