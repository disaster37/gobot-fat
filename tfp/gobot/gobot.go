package tfpgobot

import (
	"time"

	"github.com/disaster37/gobot-fat/event"
	"github.com/disaster37/gobot-fat/tfp"
	tfpconfig "github.com/disaster37/gobot-fat/tfp_config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/firmata"
)

// DFPHandler manage all i/o on FAT
type TFPHandler struct {
	stateRepository    tfp.Repository
	arduino            gobot.Adaptor
	configUsecase      tfpconfig.Usecase
	eventUsecase       event.Usecase
	robot              *gobot.Robot
	relayPompPond      *gpio.RelayDriver
	relayPompWaterfall *gpio.RelayDriver
	relayBubblePond    *gpio.RelayDriver
	relayBubbleFilter  *gpio.RelayDriver
	relayUVC1          *gpio.RelayDriver
	relayUVC2          *gpio.RelayDriver
	configHandler      *viper.Viper
	eventer            gobot.Eventer
}

// NewTFP create handler to manage FAT
func NewTFP(configHandler *viper.Viper, configUsecase tfpconfig.Usecase, eventUsecase event.Usecase, stateRepository tfp.Repository, eventer gobot.Eventer) (tfp.Gobot, error) {

	// Initialise i/o
	tfpHandler := &TFPHandler{
		stateRepository: stateRepository,
		configUsecase:   configUsecase,
		eventUsecase:    eventUsecase,
		configHandler:   configHandler,
		eventer:         eventer,
	}

	tfpHandler.init()

	// Initialize robot
	tfpHandler.robot = gobot.NewRobot(
		configHandler.GetString("tfp.name"),
		[]gobot.Connection{tfpHandler.arduino},
		[]gobot.Device{
			tfpHandler.relayPompPond,
			tfpHandler.relayPompWaterfall,
			tfpHandler.relayBubblePond,
			tfpHandler.relayBubbleFilter,
			tfpHandler.relayUVC1,
			tfpHandler.relayUVC2,
		},
		tfpHandler.work,
	)

	log.Infof("Robot %s initialized successfully", configHandler.GetString("tfp.name"))

	return tfpHandler, nil

}

func (h *TFPHandler) init() {
	arduino := firmata.NewTCPAdaptor(h.configHandler.GetString("tfp.address"))
	h.arduino = arduino
	h.relayPompPond = gpio.NewRelayDriver(arduino, h.configHandler.GetString("tfp.pin.relay.pond_pomp"))
	h.relayPompWaterfall = gpio.NewRelayDriver(arduino, h.configHandler.GetString("tfp.pin.relay.waterfall_pomp"))
	h.relayBubblePond = gpio.NewRelayDriver(arduino, h.configHandler.GetString("tfp.pin.relay.pond_bubble"))
	h.relayBubbleFilter = gpio.NewRelayDriver(arduino, h.configHandler.GetString("tfp.pin.relay.filter_bubble"))
	h.relayUVC1 = gpio.NewRelayDriver(arduino, h.configHandler.GetString("tfp.pin.relay.uvc1"))
	h.relayUVC2 = gpio.NewRelayDriver(arduino, h.configHandler.GetString("tfp.pin.relay.uvc2"))
}

// Start permit to run robot
func (h *TFPHandler) Start() error {
	go h.start()
	return nil
}

func (h *TFPHandler) start() {
	h.robot.Start(false)
}

// Stop permit to stop robot
func (h *TFPHandler) Stop() error {
	return h.robot.Stop()
}

func (h *TFPHandler) Reconnect() error {
	h.arduino.Finalize()
	h.init()
	err := h.arduino.Connect()
	if err != nil {
		return err
	}
	for _, device := range *h.robot.Devices() {
		device.Start()
	}

	return nil
}

func (h *TFPHandler) work() {

	// Debug
	h.eventer.On("stateChange", func(data interface{}) {
		log.Debugf("state: %s", h.stateRepository.String())
	})

	time.Sleep(1 * time.Second)

	// Relais handler
	h.HandleRelay()

	// External event handler
	h.HandleExternalEvent()

	// Fire event to init saved state
	h.eventer.Publish("stateChange", "initTFP")

	log.Infof("Robot %s started successfully", h.stateRepository.State().Name)
}
