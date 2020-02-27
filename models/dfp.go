package models

import (
	"time"
)

// DFP  describe the current state of drum filter pond
type DFP struct {
	ID                 string      `json:"id"`
	Name               string      `json:"name"`
	IsWashed           bool        `json:"is_running"`
	IsAuto             bool        `json:"is_started"`
	IsStopped          bool        `json:"is_stopped"`
	IsSecurity         bool        `json:"is_security"`
	IsEmergencyStopped bool        `json:"is_emmergency"`
	LastWashing        time.Time   `json:"last_washing"`
	WashingHistory     []time.Time `json:"washing_history"`
}

// NewDFPState return new PBF handler
func NewDFPState() *DFP {
	return &DFP{}
}

// CanWash handle if wash can start or not
// Only if not in emergency stop, not stopped, not in security, not already in wash
func (h *DFP) CanWash() bool {
	if !h.IsEmergencyStopped && !h.IsStopped && !h.IsWashed && !h.IsSecurity {
		return true
	}
	return false
}

// CanSetSecurity handle if security can be set
// Only if not emergency stop, not stopped and not already on security
func (h *DFP) CanSetSecurity() bool {
	if !h.IsEmergencyStopped && !h.IsStopped && !h.IsSecurity {
		return true
	}
	return false
}

// CanUnsetSecurity handle if security can be unset
// Only if not emergency stop, not stopped and already on security
func (h *DFP) CanUnsetSecurity() bool {
	if !h.IsEmergencyStopped && !h.IsStopped && h.IsSecurity {
		return true
	}
	return false
}

// CanSetStop handle if stop can be set
// Only if not emergency stop and not already stopped
func (h *DFP) CanSetStop() bool {
	if !h.IsEmergencyStopped && !h.IsStopped {
		return true
	}
	return false
}

// CanUnsetStop handle if stop can be unset
// Only if not emergency stop and already stopped
func (h *DFP) CanUnsetStop() bool {
	if !h.IsEmergencyStopped && h.IsStopped {
		return true
	}
	return false
}

// CanSetEmergencyStop handle if emergency stop can be set
// Only if emergency stop is not already set
func (h *DFP) CanSetEmergencyStop() bool {
	if !h.IsEmergencyStopped {
		return true
	}
	return false
}

// CanUnsetEmergencyStop handle if emergency stop can be unset
// Only if emergency stop is already set
func (h *DFP) CanUnsetEmergencyStop() bool {
	if h.IsEmergencyStopped {
		return true
	}
	return false
}

// CanStartMotor handle if motor can be started
// Only if not emergency stop, not stopped and not security
func (h *DFP) CanStartMotor() bool {
	if !h.IsEmergencyStopped && !h.IsStopped && !h.IsSecurity {
		return true
	}
	return false
}

// LastWashDurationSecond return the number of second from now to last wash
func (h *DFP) LastWashDurationSecond() uint64 {
	return uint64(time.Now().Sub(h.LastWashing).Seconds())
}

// AverageDurationSecond compute average duration
func (h *DFP) AverageDurationSecond() uint64 {
	if len(h.WashingHistory) > 0 {
		totalDuration := uint64(0)
		for i := 0; i < len(h.WashingHistory); i++ {
			if (i + 1) < len(h.WashingHistory) {
				totalDuration = totalDuration + uint64(h.WashingHistory[i+1].Sub(h.WashingHistory[i]).Seconds())
			} else {
				totalDuration = totalDuration + uint64(h.WashingHistory[0].Sub(h.WashingHistory[i]).Seconds())
			}
		}

		return (totalDuration / uint64(len(h.WashingHistory)))
	}

	return 0

}
