package dfpboard

import (
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"github.com/stianeikeland/go-rpio/v4"
	"gobot.io/x/gobot/v2/drivers/gpio"
	"gobot.io/x/gobot/v2/platforms/raspi"
)

type RaspiAdaptor struct {
	raspi.Adaptor
}

func NewRaspiAdaptor() *RaspiAdaptor {
	return &RaspiAdaptor{
		Adaptor: *raspi.NewAdaptor(),
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
			return errors.Wrapf(err, "Error when configure pin %d as Input Pullup", button.Pin())
		}
		pin := rpio.Pin(translatedPin)
		pin.Input()
		pin.PullUp()
	}

	log.Infof("GPIO initialized")
	return
}
