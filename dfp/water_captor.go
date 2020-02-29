package dfp

import (
	log "github.com/sirupsen/logrus"
	"gobot.io/x/gobot/drivers/gpio"
)

// HandleSecurityWaterCaptor manage the security water captor
// Check if water level is ok
func (h *DFPHandler) HandleSecurityWaterCaptor() {

	// Top captor
	// Send event only if not on Emergency stop
	h.captorWaterSecurityTop.On(gpio.ButtonPush, func(data interface{}) {
		log.Debugf("Captor water security top pushed")
		h.state.SetSecurity()
	})

	h.captorWaterSecurityTop.On(gpio.ButtonRelease, func(data interface{}) {
		log.Debugf("Captor water security top released")
		h.state.UnsetSecurity()
	})

	// Under captor
	h.captorWaterSecurityUnder.On(gpio.ButtonPush, func(data interface{}) {
		log.Debugf("Captor water security under pushed")
		h.state.SetSecurity()
	})

	h.captorWaterSecurityUnder.On(gpio.ButtonRelease, func(data interface{}) {
		log.Debugf("Captor water security under released")
		h.state.UnsetSecurity()
	})
}

// HandleWaterCaptor manage the water captor
// Check if must washing
func (h *DFPHandler) HandleWaterCaptor() {

	// Top captor
	h.captorWaterTop.On(gpio.ButtonPush, func(data interface{}) {
		log.Debugf("Captor water top pushed")
		if h.state.IsAuto() && (h.state.LastWashDurationSecond() > h.config.GetUint64("dfp.washing.wait_time_between_wash")) {
			h.state.SetShouldWash()
		}
	})

	h.captorWaterTop.On(gpio.ButtonRelease, func(data interface{}) {
		log.Debugf("Captor water top released")
	})

	// Under captor
	h.captorWaterUnder.On(gpio.ButtonPush, func(data interface{}) {
		log.Debugf("Captor water under pushed")
		if h.state.IsAuto() && (h.state.LastWashDurationSecond() > h.config.GetUint64("dfp.washing.wait_time_between_wash")) {
			h.state.SetShouldWash()
		}
	})

	h.captorWaterUnder.On(gpio.ButtonRelease, func(data interface{}) {
		log.Debugf("Captor water under released")
	})
}
