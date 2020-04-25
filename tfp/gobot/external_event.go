package tfpgobot

import (
	log "github.com/sirupsen/logrus"
)

// HandleExternalEvent permit to handle external event like security or emergency stop
func (h *TFPHandler) HandleExternalEvent() {

	// Handle ermergency stop
	h.eventer.On("stateChange", func(data interface{}) {
		event := data.(string)

		log.Debugf("Receive event %s", event)

		isUpdate := false

		switch event {
		case "isEmergencyStop":
			h.stateRepository.State().IsEmergencyStopped = true
			isUpdate = true
		case "isNotEmergencyStop":
			h.stateRepository.State().IsEmergencyStopped = false
			isUpdate = true
		case "isSecurity":
			h.stateRepository.State().IsSecurity = true
			isUpdate = true
		case "isNotSecurity":
			h.stateRepository.State().IsSecurity = false
			isUpdate = true
		case "isDisableSecurity":
			h.stateRepository.State().IsDisableSecurity = true
			isUpdate = true
		case "isNotDisableSecurity":
			h.stateRepository.State().IsDisableSecurity = false
			isUpdate = true
		}

		// Publish event to handle new state
		if isUpdate {
			h.eventer.Publish("stateChange", "externalEventTFP")
		}
	})
}
