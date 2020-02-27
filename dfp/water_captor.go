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
		log.Infof("Captor water security top pushed")
		if h.state.CanSetSecurity() {
			h.state.IsWashed = false
			h.Publish(SecurityEvent, data)
		}
		h.state.IsSecurity = true
	})

	h.captorWaterSecurityTop.On(gpio.ButtonRelease, func(data interface{}) {
		log.Infof("Captor water security top released")
		if h.state.CanUnsetSecurity() {
			h.Publish(UnSecurityEvent, data)
		}
		h.state.IsSecurity = false
	})

	// Under captor
	h.captorWaterSecurityUnder.On(gpio.ButtonPush, func(data interface{}) {
		log.Infof("Captor water security under pushed")
		if h.state.CanSetSecurity() {
			h.state.IsWashed = false
			h.Publish(SecurityEvent, data)
		}
		h.state.IsSecurity = true
	})

	h.captorWaterSecurityUnder.On(gpio.ButtonRelease, func(data interface{}) {
		log.Infof("Captor water security under released")
		if h.state.CanUnsetSecurity() {
			h.Publish(UnSecurityEvent, data)
		}
		h.state.IsSecurity = false
	})
}

// HandleWaterCaptor manage the water captor
// Check if must washing
func (h *DFPHandler) HandleWaterCaptor() {

	// Top captor
	h.captorWaterTop.On(gpio.ButtonPush, func(data interface{}) {
		log.Infof("Captor water top pushed")
		if h.state.CanWash() && h.state.IsAuto && (h.state.LastWashDurationSecond() > h.config.GetUint64("fat.washing.wait_time_between_wash")) {
			h.state.IsWashed = true
			h.Publish(WashingEvent, data)
		}
	})

	h.captorWaterTop.On(gpio.ButtonRelease, func(data interface{}) {
		log.Infof("Captor water top released")
	})

	// Under captor
	h.captorWaterUnder.On(gpio.ButtonPush, func(data interface{}) {
		log.Infof("Captor water under pushed")
		if h.state.CanWash() && h.state.IsAuto && (h.state.LastWashDurationSecond() > h.config.GetUint64("fat.washing.wait_time_between_wash")) {
			h.state.IsWashed = true
			h.Publish(WashingEvent, data)
		}
	})

	h.captorWaterUnder.On(gpio.ButtonRelease, func(data interface{}) {
		log.Infof("Captor water under released")
	})
}
