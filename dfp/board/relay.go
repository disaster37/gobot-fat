package dfpboard

import (
	"context"
	"time"

	log "github.com/sirupsen/logrus"
)

// StartDFP put dfp on auto
func (h *DFPBoard) StartDFP(ctx context.Context) (err error) {

	if !h.state.IsRunning {
		h.state.IsRunning = true
		if err = h.ledGreen.On(); err != nil {
			return err
		}

		if err = h.stateUsecase.Update(ctx, h.state); err != nil {
			return err
		}

		h.sendEvent(ctx, "board", "dfp_start")
	}

	return
}

// StopDFP stop dfp and disable auto
func (h *DFPBoard) StopDFP(ctx context.Context) (err error) {

	if h.state.IsRunning {
		h.state.IsRunning = false
		if err = h.ledGreen.Off(); err != nil {
			return err
		}

		if err := h.stateUsecase.Update(ctx, h.state); err != nil {
			return err
		}
		h.sendEvent(ctx, "board", "dfp_stop")
	}

	h.Publish(Stop, nil)

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

func (h *DFPBoard) startDrum() {
	if err := h.relayDrum.On(); err != nil {
		log.Errorf("Error when start drum: %s", err.Error())
		return
	}

	log.Debug("Start drum successfully")

}

func (h *DFPBoard) stopDrum() {
	if err := h.relayDrum.Off(); err != nil {
		log.Errorf("Error when stop drum: %s", err.Error())
		return
	}

	log.Debug("Stop drum successfully")

}

func (h *DFPBoard) startPump() {
	if err := h.relayPump.On(); err != nil {
		log.Errorf("Error when start pump: %s", err.Error())
		return
	}

	log.Debug("Start pump successfully")

}

func (h *DFPBoard) stopPump() {
	if err := h.relayPump.Off(); err != nil {
		log.Errorf("Error when stop pump: %s", err.Error())
		return
	}

	log.Debug("Stop pump successfully")

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
