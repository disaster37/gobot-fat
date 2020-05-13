package tfpgobot

import (
	"context"

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
			h.state.IsEmergencyStopped = true
			isUpdate = true
		case "isNotEmergencyStop":
			h.state.IsEmergencyStopped = false
			isUpdate = true
		case "isSecurity":
			h.state.IsSecurity = true
			isUpdate = true
		case "isNotSecurity":
			h.state.IsSecurity = false
			isUpdate = true
		case "isDisableSecurity":
			h.state.IsDisableSecurity = true
			isUpdate = true
		case "isNotDisableSecurity":
			h.state.IsDisableSecurity = false
			isUpdate = true
		}

		// Publish event to handle new state
		if isUpdate {
			h.eventer.Publish("stateChange", "externalEventTFP")

			err := h.stateUsecase.Update(context.Background(), h.state)
			if err != nil {
				log.Errorf("Error when save TFP state: %s", err.Error)
			}
		}
	})
}
