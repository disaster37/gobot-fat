package dfp

import (
	"time"

	"github.com/disaster37/gobot-fat/models"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/firmata"
)

// DFPHandler manage all i/o on FAT
type DFPHandler struct {
	state                    *models.DFPState
	arduino                  *firmata.Adaptor
	config                   *viper.Viper
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
func NewDFP(configHandler *viper.Viper) (dfpHandler *DFPHandler, err error) {
	arduino := firmata.NewAdaptor(configHandler.GetString("dfp.port"))
	eventer := gobot.NewEventer()

	// Initialise i/o
	dfpHandler = &DFPHandler{
		state:                    models.NewDFPState(configHandler.GetString("dfp.id"), configHandler.GetString("dfp.name"), eventer),
		arduino:                  arduino,
		config:                   configHandler,
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
	err = dfpHandler.captorWaterSecurityTop.SetInputPullup()
	if err != nil {
		return
	}
	err = dfpHandler.captorWaterSecurityUnder.SetInputPullup()
	if err != nil {
		return
	}
	err = dfpHandler.captorWaterTop.SetInputPullup()
	if err != nil {
		return
	}
	err = dfpHandler.captorWaterUnder.SetInputPullup()
	if err != nil {
		return
	}
	err = dfpHandler.buttonWash.SetInputPullup()
	if err != nil {
		return
	}
	err = dfpHandler.buttonAuto.SetInputPullup()
	if err != nil {
		return
	}
	err = dfpHandler.buttonStop.SetInputPullup()
	if err != nil {
		return
	}
	err = dfpHandler.buttonEmergencyStop.SetInputPullup()
	if err != nil {
		return
	}
	err = dfpHandler.buttonForceWashingPump.SetInputPullup()
	if err != nil {
		return
	}
	err = dfpHandler.buttonForceBarrelMotor.SetInputPullup()
	if err != nil {
		return
	}

	// Manage default state for button and Captor that work like button
	dfpHandler.captorWaterTop.DefaultState = 1
	//dfpHandler.captorWaterSecurityTop.DefaultState = 1

	// Set event
	dfpHandler.eventer.AddEvent("stateChange")

	// Initialize robot
	dfpHandler.robot = gobot.NewRobot(
		dfpHandler.state.Name,
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

	log.Infof("Robot %s initialized successfully", dfpHandler.state.Name())

	return

}

// State return the current state
func (h *DFPHandler) State() *models.DFPState {
	return h.state
}

// Start permit to run robot
func (h *DFPHandler) Start() {
	h.robot.Start()
}

func (h *DFPHandler) work() {

	// Debug
	h.eventer.On("stateChange", func(data interface{}) {
		log.Debugf("state: %s", h.state)
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

	// Halt button than can keep on position
	h.buttonAuto.Halt()
	h.buttonStop.Halt()
	h.buttonEmergencyStop.Halt()

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

	// Start button than can keep position
	h.buttonAuto.Start()
	h.buttonStop.Start()
	h.buttonEmergencyStop.Start()

	log.Infof("Robot %s started successfully", h.state.Name)
}
