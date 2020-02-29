package models

import (
	"encoding/json"
	"time"

	"gobot.io/x/gobot"
)

// DFPState  describe the current state of drum filter pond
type DFPState struct {
	id                 string    `json:"id"`
	name               string    `json:"name"`
	isWashed           bool      `json:"is_washed"`
	shouldWash         bool      `json:"should_wash"`
	isAuto             bool      `json:"is_started"`
	isStopped          bool      `json:"is_stopped"`
	isSecurity         bool      `json:"is_security"`
	isEmergencyStopped bool      `json:"is_emmergency_stopped"`
	lastWashing        time.Time `json:"last_washing"`
	eventer            gobot.Eventer
}

// NewDFPState return new PBF handler
func NewDFPState(ID string, name string, eventer gobot.Eventer) *DFPState {
	return &DFPState{
		id:                 ID,
		name:               name,
		isWashed:           false,
		shouldWash:         false,
		isAuto:             false,
		isStopped:          false,
		isSecurity:         false,
		isEmergencyStopped: false,
		eventer:            eventer,
	}
}

func (h *DFPState) String() string {
	temp := map[string]interface{}{
		"id":                   h.ID(),
		"name":                 h.Name(),
		"is_washed":            h.IsWashed(),
		"should_wash":          h.ShouldWash(),
		"is_auto":              h.IsAuto(),
		"is_stopped":           h.IsStopped(),
		"is_security":          h.IsSecurity(),
		"is_emergency_stopped": h.IsEmergencyStopped(),
	}
	data, err := json.Marshal(temp)
	if err != nil {
		panic(err)
	}
	return string(data)
}

// ID return the robot ID
func (h *DFPState) ID() string {
	return h.id
}

// Name return the robot name
func (h *DFPState) Name() string {
	return h.name
}

// IsWashed return the washing state
func (h *DFPState) IsWashed() bool {
	return h.isWashed
}

// ShouldWash return the should wash state
func (h *DFPState) ShouldWash() bool {
	return h.shouldWash
}

// IsAuto return the auto state
func (h *DFPState) IsAuto() bool {
	return h.isAuto
}

// IsStopped return the stop state
func (h *DFPState) IsStopped() bool {
	return h.isStopped
}

// IsEmergencyStopped return the emergency state
func (h *DFPState) IsEmergencyStopped() bool {
	return h.isEmergencyStopped
}

// IsSecurity return the security state
func (h *DFPState) IsSecurity() bool {
	return h.isSecurity
}

// LastWashing return the time of last washing
func (h *DFPState) LastWashing() time.Time {
	return h.lastWashing
}

// UpdateLastWashing update the last washing time to now
func (h *DFPState) UpdateLastWashing() {
	h.lastWashing = time.Now()
}

// SetWashed set washed state to true
func (h *DFPState) SetWashed() {
	if !h.IsWashed() {
		h.isWashed = true
		h.eventer.Publish("stateChange", "isWashed")
	}
}

// SetShouldWash set should wash state
func (h *DFPState) SetShouldWash() {
	if !h.ShouldWash() {
		h.shouldWash = true
		h.eventer.Publish("stateChange", "shouldWash")
	}
}

// UnsetShouldWash unset should wash state
func (h *DFPState) UnsetShouldWash() {
	if h.ShouldWash() {
		h.shouldWash = false
		h.eventer.Publish("stateChange", "shouldNotWash")
	}
}

// UnsetWashed set washed state to false
func (h *DFPState) UnsetWashed() {
	if h.IsWashed() {
		h.isWashed = false
		h.eventer.Publish("stateChange", "isNotWashed")
	}
}

// SetSecurity set security state to true
func (h *DFPState) SetSecurity() {
	if !h.IsSecurity() {
		h.isSecurity = true
		h.isWashed = false
		h.eventer.Publish("stateChange", "isSecurity")
	}
}

// UnsetSecurity set security state to false
func (h *DFPState) UnsetSecurity() {
	if h.IsSecurity() {
		h.isSecurity = false
		h.shouldWash = false
		h.isWashed = false
		h.eventer.Publish("stateChange", "isNotSecurity")
	}
}

// SetAuto set auto state to true
func (h *DFPState) SetAuto() {
	if !h.IsAuto() {
		h.isAuto = true
		h.eventer.Publish("stateChange", "isAuto")
	}
}

// UnsetAuto set auto state to false
func (h *DFPState) UnsetAuto() {
	if h.IsAuto() {
		h.isAuto = false
		h.eventer.Publish("stateChange", "isNotAuto")
	}
}

// SetStop set stop state to true
func (h *DFPState) SetStop() {
	if !h.IsStopped() {
		h.isWashed = false
		h.isStopped = true
		h.eventer.Publish("stateChange", "isStop")
	}
}

// UnsetStop set stop state to false
func (h *DFPState) UnsetStop() {
	if h.IsStopped() {
		h.isStopped = false
		h.shouldWash = false
		h.isWashed = false
		h.eventer.Publish("stateChange", "isNotStop")
	}
}

// SetEmergencyStop set emergency state to true
func (h *DFPState) SetEmergencyStop() {
	if !h.IsEmergencyStopped() {
		h.isWashed = false
		h.isEmergencyStopped = true
		h.eventer.Publish("stateChange", "isEmergencyStop")
	}
}

// UnsetEmergencyStop set emergency state to false
func (h *DFPState) UnsetEmergencyStop() {
	if h.IsEmergencyStopped() {
		h.shouldWash = false
		h.isWashed = false
		h.isEmergencyStopped = false
		h.eventer.Publish("stateChange", "isNotEmergencyStop")
	}
}

// CanWash handle if wash can start or not
// Only if not in emergency stop, not stopped, not in security, not already in wash
func (h *DFPState) CanWash() bool {
	if !h.IsEmergencyStopped() && !h.IsStopped() && !h.IsWashed() && !h.IsSecurity() {
		return true
	}
	return false
}

// CanStartMotor handle if motor can be started
// Only if not emergency stop, not stopped and not security
func (h *DFPState) CanStartMotor() bool {
	if !h.IsEmergencyStopped() && !h.IsStopped() && !h.IsSecurity() {
		return true
	}
	return false
}

// LastWashDurationSecond return the number of second from now to last wash
func (h *DFPState) LastWashDurationSecond() uint64 {
	return uint64(time.Now().Sub(h.LastWashing()).Seconds())
}
