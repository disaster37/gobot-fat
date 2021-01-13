package dfpboard

import (
	log "github.com/sirupsen/logrus"
	"github.com/stianeikeland/go-rpio/v4"
	"gobot.io/x/gobot/drivers/gpio"
	"gobot.io/x/gobot/platforms/raspi"
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
	defer rpio.Close()

	for _, button := range listPins {

		// Need to translate pin
		translatedPin, err := translatePin(button.Pin(), "3")
		if err != nil {
			return err
		}
		pin := rpio.Pin(translatedPin)
		pin.Input()
		pin.PullUp()

		button.DefaultState = 1
	}

	log.Infof("RPIO initialized")
	return
}
