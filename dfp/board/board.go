package dfpboard

import (
	"context"
	"sync"
	"time"

	"github.com/disaster37/gobot-arest/plateforms/arest"
	"github.com/disaster37/gobot-arest/plateforms/arest/client"
	dfpconfig "github.com/disaster37/gobot-fat/dfp_config"
	dfpstate "github.com/disaster37/gobot-fat/dfp_state"
	"github.com/disaster37/gobot-fat/event"
	"github.com/disaster37/gobot-fat/models"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
)

// DFPBoard is the DFP board
type DFPBoard struct {
	state               *models.DFPState
	board               *arest.Adaptor
	gobot               *gobot.Robot
	configUsecase       dfpconfig.Usecase
	eventUsecase        event.Usecase
	stateUsecase        dfpstate.Usecase
	configHandler       *viper.Viper
	isOnline            bool
	relayDrum           *gpio.RelayDriver
	relayPump           *gpio.RelayDriver
	ledGreen            *gpio.LedDriver
	ledRed              *gpio.LedDriver
	ledButtons          []*gpio.LedDriver
	buttonStart         *gpio.ButtonDriver
	buttonStop          *gpio.ButtonDriver
	buttonForceDrum     *gpio.ButtonDriver
	buttonForcePump     *gpio.ButtonDriver
	buttonWash          *gpio.ButtonDriver
	buttonSet           *gpio.ButtonDriver
	captorWaterUpper    *gpio.ButtonDriver
	captorWaterUnder    *gpio.ButtonDriver
	captorSecurityUpper *gpio.ButtonDriver
	captorSecurityUnder *gpio.ButtonDriver
	config              *models.DFPConfig
	isRunning           bool
	mutexState          sync.Mutex
	gobot.Eventer
}

//NewDFPBoard return the DFP board with all IO created but not started
func NewDFPBoard(configHandler *viper.Viper, configUsecase dfpconfig.Usecase, eventUsecase event.Usecase, stateUsecase dfpstate.Usecase, state *models.DFPState) (dfpBoard *DFPBoard) {

	dfpBoard = &DFPBoard{
		configHandler: configHandler,
		configUsecase: configUsecase,
		eventUsecase:  eventUsecase,
		stateUsecase:  stateUsecase,
		state:         state,
		board:         arest.NewSerialAdaptor(configHandler.GetString("port"), configHandler.GetString("name"), configHandler.GetDuration("timeout")*time.Second),
		Eventer:       gobot.NewEventer(),
	}

	// Create relay
	dfpBoard.relayDrum = gpio.NewRelayDriver(dfpBoard.board, configHandler.GetString("pin.relay.drum"))
	dfpBoard.relayPump = gpio.NewRelayDriver(dfpBoard.board, configHandler.GetString("pin.relay.pump"))

	// Create LED
	dfpBoard.ledGreen = gpio.NewLedDriver(dfpBoard.board, configHandler.GetString("pin.led.green"))
	dfpBoard.ledRed = gpio.NewLedDriver(dfpBoard.board, configHandler.GetString("pin.led.red"))

	// Create button
	dfpBoard.buttonSet = gpio.NewButtonDriver(dfpBoard.board, configHandler.GetString("pin.button.set"))
	dfpBoard.buttonSet.DefaultState = 1
	dfpBoard.buttonStart = gpio.NewButtonDriver(dfpBoard.board, configHandler.GetString("pin.button.start"))
	dfpBoard.buttonStart.DefaultState = 1
	dfpBoard.buttonStop = gpio.NewButtonDriver(dfpBoard.board, configHandler.GetString("pin.button.stop"))
	dfpBoard.buttonStop.DefaultState = 1
	dfpBoard.buttonWash = gpio.NewButtonDriver(dfpBoard.board, configHandler.GetString("pin.button.wash"))
	dfpBoard.buttonWash.DefaultState = 1
	dfpBoard.buttonForceDrum = gpio.NewButtonDriver(dfpBoard.board, configHandler.GetString("pin.button.force_drump"))
	dfpBoard.buttonForceDrum.DefaultState = 1
	dfpBoard.buttonForcePump = gpio.NewButtonDriver(dfpBoard.board, configHandler.GetString("pin.button.force_pump"))
	dfpBoard.buttonForcePump.DefaultState = 1

	// Create water captors
	dfpBoard.captorSecurityUpper = gpio.NewButtonDriver(dfpBoard.board, configHandler.GetString("pin.captor.security_upper"))
	dfpBoard.captorSecurityUnder = gpio.NewButtonDriver(dfpBoard.board, configHandler.GetString("pin.captor.security_under"))
	dfpBoard.captorSecurityUnder.DefaultState = 1
	dfpBoard.captorWaterUpper = gpio.NewButtonDriver(dfpBoard.board, configHandler.GetString("pin.captor.water_upper"))
	dfpBoard.captorWaterUnder = gpio.NewButtonDriver(dfpBoard.board, configHandler.GetString("pin.captor.water_under"))
	dfpBoard.captorWaterUnder.DefaultState = 1

	dfpBoard.gobot = gobot.NewRobot(
		configHandler.GetString("name"),
		[]gobot.Connection{dfpBoard.board},
		[]gobot.Device{
			dfpBoard.relayDrum,
			dfpBoard.relayPump,
			dfpBoard.ledGreen,
			dfpBoard.ledRed,
			dfpBoard.buttonSet,
			dfpBoard.buttonStart,
			dfpBoard.buttonStop,
			dfpBoard.buttonWash,
			dfpBoard.buttonForceDrum,
			dfpBoard.buttonForcePump,
			dfpBoard.captorSecurityUpper,
			dfpBoard.captorSecurityUnder,
			dfpBoard.captorWaterUpper,
			dfpBoard.captorWaterUnder,
		},
		dfpBoard.work,
	)

	dfpBoard.AddEvent("state")

	return dfpBoard

}

func (h *DFPBoard) work() {

	/****************
	 * Init state
	 */

	// If run normally
	if h.state.IsRunning && !h.state.IsSecurity && !h.state.IsEmergencyStopped {
		h.turnOnGreenLed()
		h.turnOffRedLed()
	} else {
		// It stopped or security
		h.forceStopRelais()
		h.turnOffGreenLed()
		h.turnOnRedLed()
	}

	// If on current wash
	if h.state.IsWashed {
		h.wash()
	}

	/*******
	 * Process on button events
	 */

	// When button start
	h.buttonStart.On(gpio.ButtonPush, func(s interface{}) {
		if log.IsLevelEnabled(log.DebugLevel) {
			log.Debug("Button start pushed")
		}

		if !h.state.IsRunning {
			h.state.IsRunning = true
			h.updateState()
			h.sendEvent("board", "dfp_start")
		}
	})

	// When button stop

	// When button wash

	// When button force drum

	// When button force pump

	// When button set

}

// Start will init some item, like INPUT_PULLUP button, then start gobot
func (h *DFPBoard) Start(ctx context.Context) (err error) {

	// Start connection on board and set INPUT_PULLUP on some pins
	err = h.board.Connect()
	if err != nil {
		return err
	}

	listPins := []int{
		h.configHandler.GetInt("pin.button.set"),
		h.configHandler.GetInt("pin.button.start"),
		h.configHandler.GetInt("pin.button.stop"),
		h.configHandler.GetInt("pin.button.wash"),
		h.configHandler.GetInt("pin.button.force_drum"),
		h.configHandler.GetInt("pin.button.force_pump"),
		h.configHandler.GetInt("pin.captor.water_upper"),
		h.configHandler.GetInt("pin.captor.water_under"),
		h.configHandler.GetInt("pin.captor.security_upper"),
		h.configHandler.GetInt("pin.captor.security_under"),
	}

	for _, pin := range listPins {
		err = h.board.Board.SetPinMode(ctx, pin, client.ModeInputPullup)
		if err != nil {
			return err
		}
	}

	return h.gobot.Start(false)

}
