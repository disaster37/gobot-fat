package pbf

import (
	"time"

	log "github.com/sirupsen/logrus"
	"gobot.io/x/gobot"
)

// HandleSecurityWaterCaptor manage the security water captor
// Check if water level is ok
func (h *FATHandler) HandleSecurityWaterCaptor() (err error) {

	gobot.Every(1*time.Millisecond, func() {

		// Handle top security water captor
		val, err := h.fatHandler.captorWaterSecurityTop.DigitalRead()
		if err != nil {
			log.Errorf("Error when read top security water captor: %s", err)
			doSecurity()

		}
		if val == 0 {
			log.Infof("Top security water captor is fired", err)
			doSecurity()
		} else if h.state.IsSecurity == true {
			log.Infof("Top security water captor is unfired", err)
			h.state.IsSecurity = false
		}

		// Handle under security water captor
		val, err := h.fatHandler.captorWaterSecurityUnder.DigitalRead()
		if err != nil {
			log.Errorf("Error when read under security water captor: %s", err)
			doSecurity()
		}
		if val == 1 {
			log.Infof("Under security water captor is fired", err)
			doSecurity()
		} else if h.state.IsSecurity == true {
			log.Infof("Under security water captor is unfired", err)
			h.state.IsSecurity = false
		}
	})

}


