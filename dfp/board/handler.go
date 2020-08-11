package dfpboard

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
)

// handleReboot permit to check on background if board is rebooted
// If board is rebooted, it wil reset with current state
func handleReboot(handler *DFPHandler) func() {
	return func() {

		data, err := handler.board.ReadValue("isRebooted")
		if err != nil {
			log.Errorf("Error when read value isRebooted: %s", err.Error())
			handler.isOnline = false
			return
		}

		if data.(bool) {
			log.Info("Board %s has been rebooted, reset state", handler.Name())

			// Running / stopped
			if handler.state.IsRunning && !handler.state.IsEmergencyStopped {
				err := handler.StartDFP(context.Background())
				if err != nil {
					log.Errorf("Error when reset start mode: %s", err.Error())
				} else {
					log.Info("Successfully reset start mode")
				}

			} else {
				// Stop / Ermergency mode
				err := handler.StopDFP(context.Background())
				if err != nil {
					log.Errorf("Error when reset stop mode: %s", err.Error())
				} else {
					log.Info("Seccessfully reset stop mode")
				}
			}

			// Washing
			if handler.state.IsWashed {
				err := handler.ForceWashing(context.Background())
				if err != nil {
					log.Errorf("Error when reset wash: %s", err.Error())
				} else {
					log.Info("Successfully reset wash")
				}
			}

			// Acknolege reboot
			_, err := handler.board.CallFunction("acknoledgeRebooted", "")
			if err != nil {
				log.Errorf("Error when aknoledge reboot on board %s: %s", handler.Name(), err.Error())
			}

			handler.isOnline = true

		}
	}
}

func handleState(h *DFPHandler) {
	for h.isRunning {

		// Read all values
		err := h.buttonStart.Read()
		if err != nil {
			log.Errorf("Error when read button auto: %s", err.Error())
		}
		err = h.buttonForceDrum.Read()
		if err != nil {
			log.Errorf("Error when read button force drum: %s", err.Error())
		}
		err = h.buttonForcePump.Read()
		if err != nil {
			log.Errorf("Error when read button force pump: %s", err.Error())
		}
		err = h.buttonSet.Read()
		if err != nil {
			log.Errorf("Error when read button set: %s", err.Error())
		}
		err = h.buttonStop.Read()
		if err != nil {
			log.Errorf("Error when read button stop: %s", err.Error())
		}
		err = h.buttonWash.Read()
		if err != nil {
			log.Errorf("Error when read button wash; %s", err.Error())
		}
		for i, captor := range h.captorWaters {
			err := captor.Read()
			if err != nil {
				log.Errorf("Error when read water captor %d: %s", i, err.Error())
			}
		}
		for i, captor := range h.captorSecurities {
			err = captor.Read()
			if err != nil {
				log.Errorf("Error when read security captor %d: %s", i, err.Error())
			}
		}

		// Manage Security captor first
		isSecrity := false
		for _, captor := range h.captorSecurities {
			if captor.IsDown() {
				isSecrity = true
				break
			}
		}
		if isSecrity != h.state.IsSecurity {
			// Stop motors and update state
			h.state.IsSecurity = isSecrity
			if h.state.IsRunning && h.state.Security() {
				err := h.stopDFP()
				if err != nil {
					h.forceStopMotors()
					log.Errorf("Error when stop motor because of security state")
				}

				// send event
				h.sendEvent("star_security", "captor")

				err = h.stateUsecase.Update(context.Background(), h.state)
				if err != nil {
					log.Errorf("Error when save DFP state: %s", err.Error())
				}
			} else if h.state.IsRunning && !h.state.Security() {
				err := h.StartDFP(context.Background())
				if err != nil {
					log.Errorf("Error when start DFP after security ended: %s", err.Error())
				}

				// send event
				h.sendEvent("stop_security", "captor")

				err = h.stateUsecase.Update(context.Background(), h.state)
				if err != nil {
					log.Errorf("Error when save DFP state: %s", err.Error())
				}
			}
		}

		buttonPushed := false

		// Stop / Auto button pushed
		if h.buttonStop.IsPushed() {
			buttonPushed = true
			err = h.StopDFP(context.Background())
			if err != nil {
				log.Errorf("Error when stop DFP: %s", err.Error())
			}
		} else if h.buttonStart.IsPushed() {
			buttonPushed = true
			err = h.StartDFP(context.Background())
			if err != nil {
				log.Errorf("Error when start DFP: %s", err.Error())
			}
		}

		// Force drum button
		if h.buttonForceDrum.IsPushed() {
			buttonPushed = true
			err = h.StartManualDrum(context.Background())
			if err != nil {
				log.Errorf("Error when force drum motor: %s", err.Error())
			}
		} else if h.buttonForceDrum.IsReleazed() {
			err = h.StopManualDrum(context.Background())
			if err != nil {
				log.Errorf("Error when stop drum motor: %s", err.Error())
			}
		}

		// Force pump
		if h.buttonForcePump.IsPushed() {
			buttonPushed = true
			err := h.StartManualPump(context.Background())
			if err != nil {
				log.Errorf("Error when force pump: %s", err.Error())
			}
		} else if h.buttonForcePump.IsReleazed() {
			err := h.StopManualPump(context.Background())
			if err != nil {
				log.Errorf("Error when stop pump: %s", err.Error)
			}
		}

		// Force wash
		if h.buttonWash.IsPushed() {
			buttonPushed = true
			err := h.ForceWashing(context.Background())
			if err != nil {
				log.Errorf("Error when force washing: %s", err.Error)
			}
		}

		//Set button
		if h.buttonSet.IsPushed() {
			buttonPushed = true
		}

		// Manage button led and screen
		if buttonPushed {
			h.timerLED.Stop()
			h.timerLED.Reset(60 * time.Second)
			h.turnOffLED = false
			go h.handleLed()
		} else {
			h.turnOffLED = true
		}

		// Manage captor state
		for _, captor := range h.captorWaters {
			if captor.IsPushed() && time.Now().After(h.state.LastWashing.Add(time.Duration(h.config.WaitTimeBetweenWashing)*time.Second)) {
				err := h.ForceWashing(context.Background())
				if err != nil {
					log.Errorf("Error when run wash: %s", err.Error())
				}
				break
			}
		}

		time.Sleep(time.Millisecond * 1)

	}
}

func (h *DFPHandler) handleLed() {
	h.turnOnLight()
	<-h.timerLED.C
	if h.turnOffLED {
		h.turnOffLight()
	}

}

func handleConfig(handler *DFPHandler) func() {
	return func() {

		config, err := handler.configUsecase.Get(context.Background())
		if err != nil {
			log.Errorf("Error when update dfp config: %s", err.Error())
		}

		handler.config = config

	}
}
