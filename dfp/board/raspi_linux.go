package dfpboard

import (
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/stianeikeland/go-rpio/v4"
	"gobot.io/x/gobot/v2/drivers/gpio"
	"gobot.io/x/gobot/v2/platforms/adaptors"
	"gobot.io/x/gobot/v2/platforms/raspi"
)

type RaspiAdaptor struct {
	raspi.Adaptor
}

func NewRaspiAdaptor(configHandler *viper.Viper) *RaspiAdaptor {
	return &RaspiAdaptor{
		Adaptor: *raspi.NewAdaptor(
			adaptors.WithGpiosPullUp(
				configHandler.GetString("pin.button.emergency_stop"),
				configHandler.GetString("pin.button.start"),
				configHandler.GetString("pin.button.stop"),
				configHandler.GetString("pin.button.wash"),
				configHandler.GetString("pin.button.force_drum"),
				configHandler.GetString("pin.button.force_pump"),
				configHandler.GetString("pin.captor.security_upper"),
				configHandler.GetString("pin.captor.security_under"),
				configHandler.GetString("pin.captor.water_upper"),
				configHandler.GetString("pin.captor.water_under"),
			),
		),
	}
}

// SetInputPullup permit to set pins ad input pullup
func (h *RaspiAdaptor) SetInputPullup(listPins []*gpio.ButtonDriver) (err error) {

	if err := rpio.Open(); err != nil {
		log.Errorf("Error when open rpio: %s", err.Error())
		return err
	}
	defer func() { _ = rpio.Close() }()

	for _, button := range listPins {

		// Need to translate pin
		translatedPin, err := translatePin(button.Pin(), "3")
		if err != nil {
			return errors.Wrapf(err, "Error when configure pin %s as Input Pullup", button.Pin())
		}
		pin := rpio.Pin(translatedPin)
		pin.Input()
		pin.PullUp()
	}

	log.Infof("GPIO initialized")
	return
}
