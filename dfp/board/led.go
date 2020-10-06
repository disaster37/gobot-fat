package dfpboard

import log "github.com/sirupsen/logrus"

func (h *DFPBoard) turnOnGreenLed() {
	err := h.ledGreen.On()
	if err != nil {
		log.Errorf("Error when turn on GreenLed: %s", err.Error())
		return
	}

	if log.IsLevelEnabled(log.DebugLevel) {
		log.Debug("Turn on GreenLed successfully")
	}
}

func (h *DFPBoard) turnOffGreenLed() {
	err := h.ledGreen.Off()
	if err != nil {
		log.Errorf("Error when turn off GreenLed: %s", err.Error())
		return
	}

	if log.IsLevelEnabled(log.DebugLevel) {
		log.Debug("Turn off GreenLed successfully")
	}
}

func (h *DFPBoard) turnOnRedLed() {
	err := h.ledRed.On()
	if err != nil {
		log.Errorf("Error when turn on RedLed: %s", err.Error())
		return
	}

	if log.IsLevelEnabled(log.DebugLevel) {
		log.Debug("Turn on RedLed successfully")
	}
}

func (h *DFPBoard) turnOffRedLed() {
	err := h.ledRed.Off()
	if err != nil {
		log.Errorf("Error when turn of RedLed: %s", err.Error())
		return
	}

	if log.IsLevelEnabled(log.DebugLevel) {
		log.Debug("Turn off RedLed successfully")
	}
}
