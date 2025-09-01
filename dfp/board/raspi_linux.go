package dfpboard

import (
	"github.com/spf13/viper"
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
