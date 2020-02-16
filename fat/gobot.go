package pbf

import (
	"context"
	"time"

	"github.com/disaster37/gobot-fat/models"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"gobot.io/x/gobot"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/firmata"
)

// FATHandler manage all i/o on FAT
type FATHandler struct {
	state                    *models.FAT
	arduino                  *firmata.Adaptor
	robot                    *gobot.Robot
	captorWaterTop           *gpio.ButtonDriver
	captorWaterUnder         *gpio.ButtonDriver
	captorWaterSecurityTop   *gpio.ButtonDriver
	captorWaterSecurityUnder *gpio.ButtonDriver
	relayBarrelMotor         *gpio.RelayDriver
	relayWashingPump         *gpio.RelayDriver
	ledRed                   *gpio.LedDriver
	ledGreen                 *gpio.LedDriver
}

// NewFAT create handler to manage FAT
func NewFAT(adaptor string, configHandler *viper.Viper, fatState *models.FAT) (fatHandler *FATHandler, err error) {
	arduino := firmata.NewAdaptor(adaptor)

	// Initialise i/o
	fatHandler = &FATHandler{
		state:                    fatState,
		arduino:                  arduino,
		captorWaterSecurityTop:   gpio.NewDirectPinDriver(arduino, configHandler.GetString("fat.pin.captor.water_security_top")),
		captorWaterSecurityUnder: gpio.NewDirectPinDriver(arduino, configHandler.GetString("fat.pin.captor.water_security_under")),
		captorWaterTop:           gpio.NewButtonDriver(arduino, configHandler.GetString("fat.pin.captor.water_top")),
		captorWaterUnder:         gpio.NewButtonDriver(arduino, configHandler.GetString("fat.pin.captor.water_under")),
		relayBarrelMotor:         gpio.NewRelayDriver(arduino, configHandler.GetString("fat.pin.relay.barrel_motor")),
		relayWashingPump:         gpio.NewRelayDriver(arduino, configHandler.GetString("fat.pin.relay.washing_pump")),
		ledGreen:                 gpio.NewLedDriver(arduino, configHandler.GetString("fat.pin.led.green")),
		ledRed:                   gpio.NewLedDriver(arduino, configHandler.GetString("fat.pin.led.red")),
	}

	// Set INPUT_PULLUP on some captor
	err = fatHandler.captorWaterSecurityTop.SetInputPullup()
	if err != nil {
		return err
	}
	err = fatHandler.captorWaterSecurityUnder.SetInputPullup()
	if err != nil {
		return err
	}
	err = fatHandler.captorWaterTop.SetInputPullup()
	if err != nil {
		return err
	}
	err = fatHandler.captorWaterUnder.SetInputPullup()
	if err != nil {
		return err
	}

	// Manage default state for button and Captor that work like button
	fatHandler.captorWaterUnder.DefaultState = 0
	fatHandler.captorWaterTop.DefaultState = 1

	// Manage default relay state
	err = fatHandler.relayBarrelMotor.Off()
	if err != nil {
		return err
	}
	err = fatHandler.relayWashingPump.Off()
	if err != nil {
		return err
	}

	// Manage default led state
	err = fatHandler.ledGreen.Off()
	if err != nil {
		return err
	}
	err = fatHandler.ledRed.Off()
	if err != nil {
		return err
	}

	// Initialize robot
	fatHandler.robot = gobot.NewRobot(
		fatHandler.state.Name,
		[]gobot.Connection{fatHandler.arduino},
		[]gobot.Device{
			fatHandler.captorWaterSecurityTop,
			fatHandler.captorWaterSecurityUnder,
			fatHandler.captorWaterTop,
			fatHandler.captorWaterUnder,
			fatHandler.relayBarrelMotor,
			fatHandler.relayWashingPump,
		},
		fatHandler.work,
	)

	return

}

// Start permit to run robot
func (h *FATHandler) Start() {
	h.robot.Start()
}

func (h *FATHandler) work() {

	var err error

	log.Info("Fat handler is running.")

	//Read some value to init state
	h.state.IsEmergencyStopped = false
	h.state.IsSecurity = false
	h.state.IsStopped = false
	h.state.IsWashed = false
	h.state.IsStarted = true
	log.Infof("FAT is initialized: mode auto")

	// Manage captor for whashing
	h.captorWaterTop.On(gpio.ButtonPush, func(data interface{}) {
		err = h.wash()
		if err != nil {
			log.Errorf("Error during whashing: %s", err)
		}
	})
	h.captorWaterUnder.On(gpio.ButtonPush, func(data interface{}) {
		log.Debugf("CaptorWaterUnder pushed")

		err = h.wash()
		if err != nil {
			log.Errorf("Error during whashing: %s", err)
		}
	})

}

func (h *FATHandler) wash() (err error) {

	// First check if can run
	if h.state.IsStarted && !h.state.IsWashed && !h.state.IsStopped && !h.state.IsSecurity && !h.state.IsEmergencyStopped {

		h.state.IsWashed = true

		// Run whashing pump and wait 5s
		err = h.relayWashingPump.On()
		if err != nil {
			log.Errorf("Error when enable relay for whashing pump: %s", err)
			h.stopWashing()
			return err
		}
		time.Sleep(5 * time.Second)

		// Run barrel motor and wait 10s
		err = h.relayBarrelMotor.On()
		if err != nil {
			log.Errorf("Error when enable relay for barrel motor: %s", err)
			h.stopWashing()
			return err
		}
		time.Sleep(10 * time.Second)

		// Stop pump and barrel
		h.stopWashing()

	}

	return nil

}

func (h *FATHandler) stopWashing() {
	err := h.relayWashingPump.Off()
	if err != nil {
		log.Errorf("Error when stop relay for whashing pump: %s", err)
	}

	err = h.relayBarrelMotor.Off()
	if err != nil {
		log.Errorf("Error when stop relay for barrelMotor: %s", err)
	}

	h.state.IsWashed = false

}
