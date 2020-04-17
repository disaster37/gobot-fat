package repository

import (
	"time"

	"github.com/disaster37/gobot-fat/dfp"
	"github.com/disaster37/gobot-fat/models"
	"gobot.io/x/gobot"
)

type dfpRepository struct {
	state   *models.DFPState
	eventer gobot.Eventer
}

// NewDFPRepository instanciate DFPRepository interface
func NewDFPRepository(state *models.DFPState, eventer gobot.Eventer) dfp.Repository {
	return &dfpRepository{
		state:   state,
		eventer: eventer,
	}
}

// SetWashed set washed state to true
func (h *dfpRepository) SetWashed() (bool, error) {
	if !h.state.IsWashed {
		h.state.IsWashed = true
		h.eventer.Publish("stateChange", "isWashed")
		return true, nil
	}
	return false, nil
}

// SetShouldWash set should wash state
func (h *dfpRepository) SetShouldWash() (bool, error) {
	if !h.state.ShouldWash {
		h.state.ShouldWash = true
		h.eventer.Publish("stateChange", "shouldWash")
		return true, nil
	}
	return false, nil
}

// UnsetShouldWash unset should wash state
func (h *dfpRepository) UnsetShouldWash() (bool, error) {
	if h.state.ShouldWash {
		h.state.ShouldWash = false
		h.eventer.Publish("stateChange", "shouldNotWash")
		return true, nil
	}
	return false, nil
}

// UnsetWashed set washed state to false
func (h *dfpRepository) UnsetWashed() (bool, error) {
	if h.state.IsWashed {
		h.state.IsWashed = false
		h.eventer.Publish("stateChange", "isNotWashed")
		return true, nil
	}
	return false, nil
}

// SetSecurity set security state to true
func (h *dfpRepository) SetSecurity() (bool, error) {
	if !h.state.IsSecurity {
		h.state.IsSecurity = true
		h.state.IsWashed = false
		h.eventer.Publish("stateChange", "isSecurity")
		return true, nil
	}
	return false, nil
}

// UnsetSecurity set security state to false
func (h *dfpRepository) UnsetSecurity() (bool, error) {
	if h.state.IsSecurity {
		h.state.IsSecurity = false
		h.state.ShouldWash = false
		h.state.IsWashed = false
		h.eventer.Publish("stateChange", "isNotSecurity")
		return true, nil
	}
	return false, nil
}

// SetAuto set auto state to true
func (h *dfpRepository) SetAuto() (bool, error) {
	if !h.state.IsAuto {
		h.state.IsAuto = true
		h.eventer.Publish("stateChange", "isAuto")
		return true, nil
	}
	return false, nil
}

// UnsetAuto set auto state to false
func (h *dfpRepository) UnsetAuto() (bool, error) {
	if h.state.IsAuto {
		h.state.IsAuto = false
		h.eventer.Publish("stateChange", "isNotAuto")
		return true, nil
	}
	return false, nil
}

// SetStop set stop state to true
func (h *dfpRepository) SetStop() (bool, error) {
	if !h.state.IsStopped {
		h.state.IsWashed = false
		h.state.IsStopped = true
		h.eventer.Publish("stateChange", "isStop")
		return true, nil
	}
	return false, nil
}

// UnsetStop set stop state to false
func (h *dfpRepository) UnsetStop() (bool, error) {
	if h.state.IsStopped {
		h.state.IsStopped = false
		h.state.ShouldWash = false
		h.state.IsWashed = false
		h.eventer.Publish("stateChange", "isNotStop")
		return true, nil
	}
	return false, nil
}

// SetEmergencyStop set emergency state to true
func (h *dfpRepository) SetEmergencyStop() (bool, error) {
	if !h.state.IsEmergencyStopped {
		h.state.IsWashed = false
		h.state.IsEmergencyStopped = true
		h.eventer.Publish("stateChange", "isEmergencyStop")
		return true, nil
	}
	return false, nil
}

// UnsetEmergencyStop set emergency state to false
func (h *dfpRepository) UnsetEmergencyStop() (bool, error) {
	if h.state.IsEmergencyStopped {
		h.state.ShouldWash = false
		h.state.IsWashed = false
		h.state.IsEmergencyStopped = false
		h.eventer.Publish("stateChange", "isNotEmergencyStop")
		return true, nil
	}
	return false, nil
}

// SetDisableSecurity set disable security state to true
func (h *dfpRepository) SetDisableSecurity() (bool, error) {
	if !h.state.IsDisableSecurity {
		h.state.IsDisableSecurity = true
		h.eventer.Publish("stateChange", "isDisableSecurity")
		return true, nil
	}
	return false, nil
}

// UnsetDisableSecurity set disable security state to false
func (h *dfpRepository) UnsetDisableSecurity() (bool, error) {
	if h.state.IsDisableSecurity {
		h.state.IsDisableSecurity = false
		h.eventer.Publish("stateChange", "isNotDisableSecurity")
		return true, nil
	}
	return false, nil
}

// CanWash handle if wash can start or not
// Only if not in emergency stop, not stopped, not in security, not already in wash
func (h *dfpRepository) CanWash() bool {
	if !h.state.IsEmergencyStopped && !h.state.IsStopped && !h.state.IsWashed && (!h.state.IsSecurity || h.state.IsDisableSecurity) {
		return true
	}
	return false
}

// CanStartMotor handle if motor can be started
// Only if not emergency stop, not stopped and not security
func (h *dfpRepository) CanStartMotor() bool {
	if !h.state.IsEmergencyStopped && !h.state.IsStopped && (!h.state.IsSecurity || h.state.IsDisableSecurity) {
		return true
	}
	return false
}

// LastWashDurationSecond return the number of second from now to last wash
func (h *dfpRepository) LastWashDurationSecond() uint64 {
	return uint64(time.Now().Sub(h.state.LastWashing).Seconds())
}

// UpdateLastWashing update the last washing time with current time
func (h *dfpRepository) UpdateLastWashing() error {
	h.state.LastWashing = time.Now()

	return nil
}

func (h *dfpRepository) String() string {
	return h.state.String()
}

func (h *dfpRepository) State() *models.DFPState {
	return h.state
}
