package dfpboard

import (
	"context"

	log "github.com/sirupsen/logrus"
)

// StartDFP put dfp on auto
func (h *DFPBoard) StartDFP(ctx context.Context) (err error) {

	if !h.state.IsRunning {
		h.state.IsRunning = true
		err = h.ledGreen.On()
		if err != nil {
			return
		}
		h.Publish("state", h.state)
		h.sendEvent(ctx, "board", "dfp_start")
	}

	return
}

// StopDFP stop dfp and disable auto
func (h *DFPBoard) StopDFP(ctx context.Context) (err error) {

	if h.state.IsRunning {
		h.state.IsRunning = false
		err = h.ledGreen.Off()
		if err != nil {
			return
		}
		h.Publish("state", h.state)
		h.sendEvent(ctx, "board", "dfp_stop")
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
		if log.IsLevelEnabled(log.DebugLevel) {
			log.Debug("Run force drum")
		}

		err = h.relayDrum.On()
		if err != nil {
			return
		}

	}
	return
}

// StopManualDrum force stop drum motor
// Only if not current washing
func (h *DFPBoard) StopManualDrum(ctx context.Context) (err error) {

	if !h.state.IsWashed {
		if log.IsLevelEnabled(log.DebugLevel) {
			log.Debug("Stop force drum")
		}

		err = h.relayDrum.Off()
		if err != nil {
			return
		}
	}
	return
}

// StartManualPump force start pump
// Only if not already wash and is not on emergency stopped
func (h *DFPBoard) StartManualPump(ctx context.Context) (err error) {

	if !h.state.IsWashed && !h.state.IsEmergencyStopped {
		if log.IsLevelEnabled(log.DebugLevel) {
			log.Debug("Run force pump")
		}

		err = h.relayPump.On()
		if err != nil {
			return
		}
	}

	return
}

// StopManualPump force stop pump
// Only if not already wash
func (h *DFPBoard) StopManualPump(ctx context.Context) (err error) {

	// Stop force pump
	if !h.state.IsWashed {
		if log.IsLevelEnabled(log.DebugLevel) {
			log.Debug("Stop force pump")
		}

		err = h.relayPump.Off()
		if err != nil {
			return
		}
	}

	return
}

func (h *DFPBoard) startDrum() {
	err := h.relayDrum.On()
	if err != nil {
		log.Errorf("Error when start drum: %s", err.Error())
		return
	}

	if log.IsLevelEnabled(log.DebugLevel) {
		log.Debug("Start drum successfully")
	}
}

func (h *DFPBoard) stopDrum() {
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
		isErr := true

		for isErr {
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
