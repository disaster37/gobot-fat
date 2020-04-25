package dfpgobot

import (
	"time"

	"github.com/disaster37/gobot-fat/dfp"
	"github.com/disaster37/gobot-fat/dfp_config"
	"github.com/disaster37/gobot-fat/event"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/firmata"
)

// DFPHandler manage all i/o on FAT
type DFPHandler struct {
	stateRepository          dfp.Repository
	arduino                  *firmata.Adaptor
	configUsecase            dfpconfig.Usecase
	eventUsecase             event.Usecase
	robot                    *gobot.Robot
	captorWaterTop           *gpio.ButtonDriver
	captorWaterUnder         *gpio.ButtonDriver
	captorWaterSecurityTop   *gpio.ButtonDriver
	captorWaterSecurityUnder *gpio.ButtonDriver
	relayBarrelMotor         *gpio.RelayDriver
	relayWashingPump         *gpio.RelayDriver
	ledRed                   *gpio.LedDriver
	ledGreen                 *gpio.LedDriver
	buttonAuto               *gpio.ButtonDriver
	buttonStop               *gpio.ButtonDriver
	buttonEmergencyStop      *gpio.ButtonDriver
	buttonWash               *gpio.ButtonDriver
	buttonForceWashingPump   *gpio.ButtonDriver
	buttonForceBarrelMotor   *gpio.ButtonDriver
	eventer                  gobot.Eventer
}

// NewDFP create handler to manage FAT
func NewDFP(configHandler *viper.Viper, configUsecase dfpconfig.Usecase, eventUsecase event.Usecase, stateRepository dfp.Repository, eventer gobot.Eventer) (dfp.Gobot, error) {
	arduino := firmata.NewAdaptor(configHandler.GetString("dfp.port"))

	// Initialise i/o
	dfpHandler := &DFPHandler{
		stateRepository:          stateRepository,
		arduino:                  arduino,
		configUsecase:            configUsecase,
		eventUsecase:             eventUsecase,
		eventer:                  eventer,
		captorWaterSecurityTop:   gpio.NewButtonDriver(arduino, configHandler.GetString("dfp.pin.captor.water_security_top")),
		captorWaterSecurityUnder: gpio.NewButtonDriver(arduino, configHandler.GetString("dfp.pin.captor.water_security_under")),
		captorWaterTop:           gpio.NewButtonDriver(arduino, configHandler.GetString("dfp.pin.captor.water_top")),
		captorWaterUnder:         gpio.NewButtonDriver(arduino, configHandler.GetString("dfp.pin.captor.water_under")),
		relayBarrelMotor:         gpio.NewRelayDriver(arduino, configHandler.GetString("dfp.pin.relay.barrel_motor")),
		relayWashingPump:         gpio.NewRelayDriver(arduino, configHandler.GetString("dfp.pin.relay.washing_pump")),
		ledGreen:                 gpio.NewLedDriver(arduino, configHandler.GetString("dfp.pin.led.green")),
		ledRed:                   gpio.NewLedDriver(arduino, configHandler.GetString("dfp.pin.led.red")),
		buttonAuto:               gpio.NewButtonDriver(arduino, configHandler.GetString("dfp.pin.button.auto")),
		buttonStop:               gpio.NewButtonDriver(arduino, configHandler.GetString("dfp.pin.button.stop")),
		buttonEmergencyStop:      gpio.NewButtonDriver(arduino, configHandler.GetString("dfp.pin.button.emergency_stop")),
		buttonWash:               gpio.NewButtonDriver(arduino, configHandler.GetString("dfp.pin.button.wash")),
		buttonForceWashingPump:   gpio.NewButtonDriver(arduino, configHandler.GetString("dfp.pin.button.force_washing_pump")),
		buttonForceBarrelMotor:   gpio.NewButtonDriver(arduino, configHandler.GetString("dfp.pin.button.force_barrel_motor")),
	}

	// Set INPUT_PULLUP on some captor
	err := dfpHandler.captorWaterSecurityTop.SetInputPullup()
	if err != nil {
		return nil, err
	}
	err = dfpHandler.captorWaterSecurityUnder.SetInputPullup()
	if err != nil {
		return nil, err
	}
	err = dfpHandler.captorWaterTop.SetInputPullup()
	if err != nil {
		return nil, err
	}
	err = dfpHandler.captorWaterUnder.SetInputPullup()
	if err != nil {
		return nil, err
	}
	err = dfpHandler.buttonWash.SetInputPullup()
	if err != nil {
		return nil, err
	}
	err = dfpHandler.buttonAuto.SetInputPullup()
	if err != nil {
		return nil, err
	}
	err = dfpHandler.buttonStop.SetInputPullup()
	if err != nil {
		return nil, err
	}
	err = dfpHandler.buttonEmergencyStop.SetInputPullup()
	if err != nil {
		return nil, err
	}
	err = dfpHandler.buttonForceWashingPump.SetInputPullup()
	if err != nil {
		return nil, err
	}
	err = dfpHandler.buttonForceBarrelMotor.SetInputPullup()
	if err != nil {
		return nil, err
	}

	// Manage default state for button and Captor that work like button
	dfpHandler.captorWaterTop.DefaultState = 1
	dfpHandler.captorWaterSecurityTop.DefaultState = 1

	// Set event
	dfpHandler.eventer.AddEvent("stateChange")

	// Initialize robot
	dfpHandler.robot = gobot.NewRobot(
		configHandler.GetString("dfp.name"),
		[]gobot.Connection{dfpHandler.arduino},
		[]gobot.Device{
			dfpHandler.buttonForceWashingPump,
			dfpHandler.buttonForceBarrelMotor,
			dfpHandler.buttonWash,
			dfpHandler.buttonAuto,
			dfpHandler.buttonStop,
			dfpHandler.buttonEmergencyStop,
			dfpHandler.captorWaterSecurityTop,
			dfpHandler.captorWaterSecurityUnder,
			dfpHandler.captorWaterTop,
			dfpHandler.captorWaterUnder,
			dfpHandler.relayBarrelMotor,
			dfpHandler.relayWashingPump,
			dfpHandler.ledGreen,
			dfpHandler.ledRed,
		},
		dfpHandler.work,
	)

	log.Infof("Robot %s initialized successfully", configHandler.GetString("dfp.name"))

	return dfpHandler, nil

}

// Start permit to run robot
func (h *DFPHandler) Start() {
	go h.start()
}

func (h *DFPHandler) start() {
	err := h.robot.Start()
	for err != nil {
		log.Errorf("Error when start Robot %s: %s", h.stateRepository.State().Name, err.Error())
		time.Sleep(10 * time.Second)
		err = h.robot.Start()
	}

	log.Infof("Robot %s started successfully", h.stateRepository.State().Name)
}

// Stop permit to stop robot
func (h *DFPHandler) Stop() error {
	return h.robot.Stop()
}

func (h *DFPHandler) work() {

	// Debug
	h.eventer.On("stateChange", func(data interface{}) {
		log.Debugf("state: %s", h.stateRepository.String())
	})

	// Stop motors
	h.StopWashingPump()
	h.StopBarrelMotor()

	// Manage default led state
	h.ledGreen.Off()
	h.ledRed.Off()

	// Led handler
	h.HandleRedLed()
	h.HandleGreenLed()

	time.Sleep(1 * time.Second)

	// Button handler
	h.HandleButtonEmergencyStop()
	h.HandleButtonStop()
	h.HandleButtonAuto()
	h.HandleButtonWash()
	h.HandleButtonForceMotor()

	// Captor handler
	h.HandleSecurityWaterCaptor()
	h.HandleWaterCaptor()

	// Motor handler
	h.HandleMotor()

	// Fire event to init saved state
	h.eventer.Publish("stateChange", "initDFP")

	// Publish external event
	if h.stateRepository.State().IsEmergencyStopped {
		h.stateRepository.SetEmergencyStop()
	}
	if h.stateRepository.State().IsDisableSecurity {
		h.stateRepository.SetDisableSecurity()
	}

	log.Infof("Robot %s started successfully", h.stateRepository.State().Name)
}
