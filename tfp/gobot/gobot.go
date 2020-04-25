package tfpgobot

import (
	"time"

	"github.com/disaster37/gobot-fat/event"
	"github.com/disaster37/gobot-fat/tfp"
	"github.com/disaster37/gobot-fat/tfp_config"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/firmata"
)

// DFPHandler manage all i/o on FAT
type TFPHandler struct {
	stateRepository    tfp.Repository
	arduino            *firmata.TCPAdaptor
	configUsecase      tfpconfig.Usecase
	eventUsecase       event.Usecase
	robot              *gobot.Robot
	relayPompPond      *gpio.RelayDriver
	relayPompWaterfall *gpio.RelayDriver
	relayBubblePond    *gpio.RelayDriver
	relayBubbleFilter  *gpio.RelayDriver
	relayUVC1          *gpio.RelayDriver
	relayUVC2          *gpio.RelayDriver
	eventer            gobot.Eventer
}

// NewTFP create handler to manage FAT
func NewTFP(configHandler *viper.Viper, configUsecase tfpconfig.Usecase, eventUsecase event.Usecase, stateRepository tfp.Repository, eventer gobot.Eventer) (tfp.Gobot, error) {
	arduino := firmata.NewTCPAdaptor(configHandler.GetString("tfp.address"))

	// Initialise i/o
	tfpHandler := &TFPHandler{
		stateRepository:    stateRepository,
		arduino:            arduino,
		configUsecase:      configUsecase,
		eventUsecase:       eventUsecase,
		eventer:            eventer,
		relayPompPond:      gpio.NewRelayDriver(arduino, configHandler.GetString("tfp.pin.relay.pond_pomp")),
		relayPompWaterfall: gpio.NewRelayDriver(arduino, configHandler.GetString("tfp.pin.relay.waterfall_pomp")),
		relayBubblePond:    gpio.NewRelayDriver(arduino, configHandler.GetString("tfp.pin.relay.pond_bubble")),
		relayBubbleFilter:  gpio.NewRelayDriver(arduino, configHandler.GetString("tfp.pin.relay.filter_bubble")),
		relayUVC1:          gpio.NewRelayDriver(arduino, configHandler.GetString("tfp.pin.relay.uvc1")),
		relayUVC2:          gpio.NewRelayDriver(arduino, configHandler.GetString("tfp.pin.relay.uvc2")),
	}

	// Set event
	tfpHandler.eventer.AddEvent("stateChange")

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

// Start permit to run robot
func (h *TFPHandler) Start() {
	go h.start()
}

func (h *TFPHandler) start() {
	err := h.robot.Start()
	for err != nil {
		log.Errorf("Error when start Robot %s: %s", h.stateRepository.State().Name, err.Error())
		time.Sleep(10 * time.Second)
		err = h.robot.Start()
	}

	log.Infof("Robot %s started successfully", h.stateRepository.State().Name)
}

// Stop permit to stop robot
func (h *TFPHandler) Stop() error {
	return h.robot.Stop()
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
