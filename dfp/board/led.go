package dfpboard

import log "github.com/sirupsen/logrus"

func (h *DFPBoard) turnOnGreenLed() {
	if err := h.ledGreen.On(); err != nil {
		log.Errorf("Error when turn on GreenLed: %s", err.Error())
		return
	}

	log.Debug("Turn on GreenLed successfully")

}

func (h *DFPBoard) turnOffGreenLed() {
	if err := h.ledGreen.Off(); err != nil {
		log.Errorf("Error when turn off GreenLed: %s", err.Error())
		return
	}

	log.Debug("Turn off GreenLed successfully")

}

func (h *DFPBoard) turnOnRedLed() {
	if err := h.ledRed.On(); err != nil {
		log.Errorf("Error when turn on RedLed: %s", err.Error())
		return
	}

	log.Debug("Turn on RedLed successfully")

}

func (h *DFPBoard) turnOffRedLed() {
	if err := h.ledRed.Off(); err != nil {
		log.Errorf("Error when turn of RedLed: %s", err.Error())
		return
	}

	log.Debug("Turn off RedLed successfully")

}
