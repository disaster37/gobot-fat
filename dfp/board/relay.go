package dfpboard

import (
	"context"
	"time"

	"github.com/disaster37/gobot-fat/helper"
	log "github.com/sirupsen/logrus"
)

// StartDFP put dfp on auto
func (h *DFPBoard) StartDFP(ctx context.Context) (err error) {

	if !h.state.IsRunning {
		h.state.IsRunning = true

		if err = h.ledGreen.On(); err != nil {
			return err
		}

		if err := h.stateUsecase.Update(ctx, h.state); err != nil {
			return err
		}

		h.sendEvent(ctx, "board", "dfp_start")

		// Publish internal event
		h.Publish(EventStartDFP, nil)

		// Publish global event
		h.globalEventer.Publish(EventStartDFP, nil)
	}

	return
}

// StopDFP stop dfp and disable auto
func (h *DFPBoard) StopDFP(ctx context.Context) (err error) {

	h.forceStopRelais()

	if h.state.IsRunning {
		h.state.IsRunning = false

		if err = h.ledGreen.Off(); err != nil {
			return err
		}

		if err := h.stateUsecase.Update(ctx, h.state); err != nil {
			return err
		}

		h.sendEvent(ctx, "board", "dfp_stop")

		// Publish internal event
		h.Publish(EventStopDFP, nil)

		// Publish global event
		h.globalEventer.Publish(EventStopDFP, nil)
	}

	return
}

// SetEmergencyStop put DFP on emergency stop
// It send a global event to inform antoher board
func (h *DFPBoard) SetEmergencyStop(ctx context.Context) (err error) {

	// Stops all relais
	h.forceStopRelais()

	if !h.state.IsEmergencyStopped {
		h.state.IsEmergencyStopped = true

		if err = h.ledRed.On(); err != nil {
			return err
		}

		if err := h.stateUsecase.Update(ctx, h.state); err != nil {
			return err
		}

		h.sendEvent(ctx, "board", "dfp_set_emergency_stop")

		// Publish internal event
		h.Publish(EventSetEmergencyStop, nil)

		// Publish global event
		h.globalEventer.Publish(helper.SetEmergencyStop, nil)
	}

	return
}

// UnsetEmergencyStop remove the emergency stop
// It send global event to inform another board
func (h *DFPBoard) UnsetEmergencyStop(ctx context.Context) (err error) {

	if h.state.IsEmergencyStopped {
		h.state.IsEmergencyStopped = false

		// Turn off red led
		if !h.state.IsSecurity {
			if err = h.ledRed.Off(); err != nil {
				return err
			}
		}

		if err := h.stateUsecase.Update(ctx, h.state); err != nil {
			return err
		}

		h.sendEvent(ctx, "board", "dfp_unset_emergency_stop")

		// Publish internal event
		h.Publish(EventUnsetEmergencyStop, nil)

		// Publish global event
		h.globalEventer.Publish(helper.UnsetEmergencyStop, nil)
	}

	return
}

// SetSecurity put DFP on security
// It send a global event to inform antoher board
func (h *DFPBoard) SetSecurity(ctx context.Context) (err error) {

	if !h.state.IsSecurity {
		h.state.IsSecurity = true

		if err = h.ledRed.On(); err != nil {
			return err
		}

		if err := h.stateUsecase.Update(ctx, h.state); err != nil {
			return err
		}

		h.sendEvent(ctx, "board", "dfp_set_security")

		// Publish internal event
		h.Publish(EventSetSecurity, nil)

		// Publish global event
		h.globalEventer.Publish(EventSetSecurity, nil)
	}

	return
}

// UnsetSecurity remove the security
// It send global event to inform another board
func (h *DFPBoard) UnsetSecurity(ctx context.Context) (err error) {

	if h.state.IsSecurity {
		h.state.IsSecurity = false

		if !h.state.IsEmergencyStopped {
			// Turn off red led
			if err = h.ledRed.Off(); err != nil {
				return err
			}

		}

		if err := h.stateUsecase.Update(ctx, h.state); err != nil {
			return err
		}

		h.sendEvent(ctx, "board", "dfp_unset_security")

		// Publish internal event
		h.Publish(EventUnsetSecurity, nil)

		// Publish global event
		h.globalEventer.Publish(EventUnsetSecurity, nil)
	}

	return
}

// ForceWashing start a washing cycle
func (h *DFPBoard) ForceWashing(ctx context.Context) (err error) {
	if !h.state.IsWashed && !h.state.IsEmergencyStopped {
		log.Debug("Run force wash")
		h.wash()
	}

	return
}

// StartManualDrum force start drum motor
// Only if not already wash and is not on emergency stopped
func (h *DFPBoard) StartManualDrum(ctx context.Context) (err error) {

	if !h.state.IsWashed && !h.state.IsEmergencyStopped {

		log.Debug("Run force drum")

		if err = h.relayDrum.On(); err != nil {
			return err
		}

	}
	return
}

// StopManualDrum force stop drum motor
// Only if not current washing
func (h *DFPBoard) StopManualDrum(ctx context.Context) (err error) {

	if !h.state.IsWashed {
		log.Debug("Stop force drum")

		if err = h.relayDrum.Off(); err != nil {
			return err
		}
	}
	return
}

// StartManualPump force start pump
// Only if not already wash and is not on emergency stopped
func (h *DFPBoard) StartManualPump(ctx context.Context) (err error) {

	if !h.state.IsWashed && !h.state.IsEmergencyStopped {
		log.Debug("Run force pump")

		if err = h.relayPump.On(); err != nil {
			return err
		}
	}

	return
}

// StopManualPump force stop pump
// Only if not already wash
func (h *DFPBoard) StopManualPump(ctx context.Context) (err error) {

	// Stop force pump
	if !h.state.IsWashed {
		log.Debug("Stop force pump")

		if err = h.relayPump.Off(); err != nil {
			return err
		}
	}

	return
}

func (h *DFPBoard) forceStopRelais() {

	forceStopRelais := func() {
		isErr := true

		for isErr {
			isErr = false
			if err := h.relayDrum.Off(); err != nil {
				log.Errorf("Error when stop drump: %s", err.Error())
				isErr = true
			}

			if err := h.relayPump.Off(); err != nil {
				log.Errorf("Error when stop pump: %s", err.Error())
				isErr = true
			}

			if isErr {
				time.Sleep(100 * time.Millisecond)
			}
		}
	}

	// Lauch only routine in stop failed
	if err := h.relayDrum.Off(); err != nil {
		go forceStopRelais()
		return
	}
	if err := h.relayPump.Off(); err != nil {
		go forceStopRelais()
		return
	}
}
