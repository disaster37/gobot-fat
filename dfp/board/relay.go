package dfpboard

import log "github.com/sirupsen/logrus"

func (h *DFPBoard) startDrump() {
	err := h.relayDrum.On()
	if err != nil {
		log.Errorf("Error when start drum: %s", err.Error())
		return
	}

	if log.IsLevelEnabled(log.DebugLevel) {
		log.Debug("Start drum successfully")
	}
}

func (h *DFPBoard) stopDrump() {
	err := h.relayDrum.Off()
	if err != nil {
		log.Errorf("Error when stop drum: %s", err.Error())
		return
	}

	if log.IsLevelEnabled(log.DebugLevel) {
		log.Debug("Stop drum successfully")
	}
}

func (h *DFPBoard) startPump() {
	err := h.relayPump.On()
	if err != nil {
		log.Errorf("Error when start pump: %s", err.Error())
		return
	}

	if log.IsLevelEnabled(log.DebugLevel) {
		log.Debug("Start pump successfully")
	}
}

func (h *DFPBoard) stopPump() {
	err := h.relayPump.Off()
	if err != nil {
		log.Errorf("Error when stop pump: %s", err.Error())
		return
	}

	if log.IsLevelEnabled(log.DebugLevel) {
		log.Debug("Stop pump successfully")
	}
}

func (h *DFPBoard) forceStopRelais() {
	go func() {
		isErr := false

		for !isErr {
			isErr = false
			err := h.relayDrum.Off()
			if err != nil {
				log.Errorf("Error when stop drump: %s", err.Error())
				isErr = true
			}

			err = h.relayPump.Off()
			if err != nil {
				log.Errorf("Error when stop pump: %s", err.Error())
				isErr = true
			}
		}
	}()
}
