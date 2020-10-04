package dfpboard

import (
	"context"
	"time"

	"github.com/disaster37/gobot-fat/board"
	log "github.com/sirupsen/logrus"
)

// handleReboot permit to check on background if board is rebooted
// If board is rebooted, it wil reset with current state
func (h *DFPHandler) handleReboot(ctx context.Context) {
	log.Debug("Handle reboot fired")

	data, err := h.board.ReadValue(ctx, "isRebooted")
	if err != nil {
		log.Errorf("Error when read value isRebooted: %s", err.Error())
		h.isOnline = false
		return
	}

	if data.(bool) {
		log.Infof("Board %s has been rebooted, reset state", h.Name())

		if h.isRunning {
			h.Stop(ctx)
			h.Start(ctx)
		}

		// Acknolege reboot
		_, err := h.board.CallFunction(ctx, "acknoledgeRebooted", "")
		if err != nil {
			log.Errorf("Error when aknoledge reboot on board %s: %s", h.Name(), err.Error())
		}

		h.isOnline = true
	}
}

func (h *DFPHandler) handleState(ctx context.Context) {
	if h.isRunning {

		// Read all values
		err := h.buttonStart.Read(ctx)
		if err != nil {
			log.Errorf("Error when read button auto: %s", err.Error())
		}
		err = h.buttonForceDrum.Read(ctx)
		if err != nil {
			log.Errorf("Error when read button force drum: %s", err.Error())
		}
		err = h.buttonForcePump.Read(ctx)
		if err != nil {
			log.Errorf("Error when read button force pump: %s", err.Error())
		}
		err = h.buttonSet.Read(ctx)
		if err != nil {
			log.Errorf("Error when read button set: %s", err.Error())
		}
		err = h.buttonStop.Read(ctx)
		if err != nil {
			log.Errorf("Error when read button stop: %s", err.Error())
		}
		err = h.buttonWash.Read(ctx)
		if err != nil {
			log.Errorf("Error when read button wash; %s", err.Error())
		}
		for i, captor := range h.captorWaters {
			err := captor.Read(ctx)
			if err != nil {
				log.Errorf("Error when read water captor %d: %s", i, err.Error())
			}
		}
		for i, captor := range h.captorSecurities {
			err = captor.Read(ctx)
			if err != nil {
				log.Errorf("Error when read security captor %d: %s", i, err.Error())
			}
		}

		// Manage Security captor first
		isSecrity := false
		for _, captor := range h.captorSecurities {
			if captor.IsDown() {
				log.Info("Security captor fired")
				isSecrity = true
				break
			}
		}
		if isSecrity != h.state.IsSecurity {
			// Stop motors and update state
			h.state.IsSecurity = isSecrity
			if h.state.IsRunning && h.state.Security() {
				log.Info("DFP set security")
				routine := board.NewRoutine(ctx, h.stopDFP)
				// handle the stop without block
				go func() {
					select {
					case err := <-routine.Error():
						routine = board.NewRoutine(ctx, h.forceStopMotors)
						select {
						case err := <-routine.Error():
							log.Errorf("Error when force stop motors: %s", err.Error())
						case <-routine.Result():
							break
						}
						log.Errorf("Error when stop motor because of security state: %s", err.Error())
					case <-routine.Result():
						log.Info("DFP stopped because of security state")
						break
					}
				}()

				// send event
				h.sendEvent(ctx, "star_security", "captor")

				err = h.stateUsecase.Update(ctx, h.state)
				if err != nil {
					log.Errorf("Error when save DFP state: %s", err.Error())
				}
			} else if h.state.IsRunning && !h.state.Security() {
				log.Info("DFP unset security")
				err := h.StartDFP(ctx)
				if err != nil {
					log.Errorf("Error when start DFP after security ended: %s", err.Error())
				}

				// send event
				h.sendEvent(ctx, "stop_security", "captor")

				err = h.stateUsecase.Update(ctx, h.state)
				if err != nil {
					log.Errorf("Error when save DFP state: %s", err.Error())
				}
			}
		}

		buttonPushed := false

		// Stop / Auto button pushed
		if h.buttonStop.IsPushed() {
			log.Info("Button stop pushed")
			buttonPushed = true
			err = h.StopDFP(ctx)
			if err != nil {
				log.Errorf("Error when stop DFP: %s", err.Error())
			}
		} else if h.buttonStart.IsPushed() {
			log.Info("Button start pushed")
			buttonPushed = true
			err = h.StartDFP(ctx)
			if err != nil {
				log.Errorf("Error when start DFP: %s", err.Error())
			}
		}

		// Force drum button
		if h.buttonForceDrum.IsPushed() {
			log.Info("Button force drum pushed")
			buttonPushed = true
			err = h.StartManualDrum(ctx)
			if err != nil {
				log.Errorf("Error when force drum motor: %s", err.Error())
			}
		} else if h.buttonForceDrum.IsReleazed() {
			log.Info("Button force drum releazed")
			err = h.StopManualDrum(ctx)
			if err != nil {
				log.Errorf("Error when stop drum motor: %s", err.Error())
			}
		}

		// Force pump
		if h.buttonForcePump.IsPushed() {
			log.Info("Button force pump pushed")
			buttonPushed = true
			err := h.StartManualPump(ctx)
			if err != nil {
				log.Errorf("Error when force pump: %s", err.Error())
			}
		} else if h.buttonForcePump.IsReleazed() {
			log.Info("Button force pump releazed")
			err := h.StopManualPump(ctx)
			if err != nil {
				log.Errorf("Error when stop pump: %s", err.Error())
			}
		}

		// Force wash
		if h.buttonWash.IsPushed() {
			log.Info("Button force wash pushed")
			buttonPushed = true
			if h.state.ShouldWash() {
				routine := board.NewRoutine(ctx, h.wash)
				// Handle result
				go func() {
					select {
					case err := <-routine.Error():
						log.Errorf("Error when force washing: %s", err.Error())
					case <-routine.Result():
						log.Info("Force wash successfully finished")
					}
				}()
			}
		}

		//Set button
		if h.buttonSet.IsPushed() {
			log.Debug("Button set pushed")
			buttonPushed = true
		} else if h.buttonSet.IsReleazed() {
			log.Debug("Button set releazed")
		}

		// Manage button led and screen
		if buttonPushed {
			h.timerLED.Stop()
			h.timerLED.Reset(60 * time.Second)
			h.turnOffLED = false
			go h.handleLed(ctx)
		} else {
			h.turnOffLED = true
		}

		// Manage captor state
		for _, captor := range h.captorWaters {
			if captor.IsPushed() && time.Now().After(h.state.LastWashing.Add(time.Duration(h.config.WaitTimeBetweenWashing)*time.Second)) {
				log.Info("Captor start auto washing")
				err := h.ForceWashing(ctx)
				if err != nil {
					log.Errorf("Error when run wash: %s", err.Error())
				}
				break
			}
		}
	}
}

func (h *DFPHandler) handleLed(ctx context.Context) {
	h.turnOnLight(ctx)
	<-h.timerLED.C
	if h.turnOffLED {
		h.turnOffLight(ctx)
	}
}

func (h *DFPHandler) handleConfig(ctx context.Context) {

	config, err := h.configUsecase.Get(ctx)
	if err != nil {
		log.Errorf("Error when update dfp config: %s", err.Error())
	}

	h.config = config
}
