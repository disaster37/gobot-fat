package repository

import (
	"context"
	"time"

	"github.com/disaster37/gobot-fat/dfp"
	"github.com/disaster37/gobot-fat/dfp_config"
	"github.com/disaster37/gobot-fat/models"
	"gobot.io/x/gobot"
)

type dfpRepository struct {
	state   *models.DFPState
	eventer gobot.Eventer
	config  dfpconfig.Usecase
}

// NewDFPRepository instanciate DFPRepository interface
func NewDFPRepository(state *models.DFPState, eventer gobot.Eventer, config dfpconfig.Usecase) dfp.Repository {
	return &dfpRepository{
		state:   state,
		eventer: eventer,
		config:  config,
	}
}

// SetWashed set washed state to true
func (h *dfpRepository) SetWashed() error {
	if !h.state.IsWashed {
		h.state.IsWashed = true
		h.eventer.Publish("stateChange", "isWashed")
		return nil
	}
	return nil
}

// SetShouldWash set should wash state
func (h *dfpRepository) SetShouldWash() error {
	if !h.state.ShouldWash {
		h.state.ShouldWash = true
		h.eventer.Publish("stateChange", "shouldWash")
		return nil
	}
	return nil
}

// UnsetShouldWash unset should wash state
func (h *dfpRepository) UnsetShouldWash() error {
	if h.state.ShouldWash {
		h.state.ShouldWash = false
		h.eventer.Publish("stateChange", "shouldNotWash")
		return nil
	}
	return nil
}

// UnsetWashed set washed state to false
func (h *dfpRepository) UnsetWashed() error {
	if h.state.IsWashed {
		h.state.IsWashed = false
		h.eventer.Publish("stateChange", "isNotWashed")
		return nil
	}
	return nil
}

// SetSecurity set security state to true
func (h *dfpRepository) SetSecurity() error {
	if !h.state.IsSecurity {
		h.state.IsSecurity = true
		h.state.IsWashed = false
		h.eventer.Publish("stateChange", "isSecurity")

		// Save config
		err := h.updateConfig()

		return err
	}
	return nil
}

// UnsetSecurity set security state to false
func (h *dfpRepository) UnsetSecurity() error {
	if h.state.IsSecurity {
		h.state.IsSecurity = false
		h.state.ShouldWash = false
		h.state.IsWashed = false
		h.eventer.Publish("stateChange", "isNotSecurity")

		// Save config
		err := h.updateConfig()

		return err
	}
	return nil
}

// SetAuto set auto state to true
func (h *dfpRepository) SetAuto() error {
	if !h.state.IsAuto {
		h.state.IsAuto = true
		h.eventer.Publish("stateChange", "isAuto")

		// Save config
		err := h.updateConfig()

		return err
	}
	return nil
}

// UnsetAuto set auto state to false
func (h *dfpRepository) UnsetAuto() error {
	if h.state.IsAuto {
		h.state.IsAuto = false
		h.eventer.Publish("stateChange", "isNotAuto")

		// Save config
		err := h.updateConfig()

		return err
	}
	return nil
}

// SetStop set stop state to true
func (h *dfpRepository) SetStop() error {
	if !h.state.IsStopped {
		h.state.IsWashed = false
		h.state.IsStopped = true
		h.eventer.Publish("stateChange", "isStop")

		// Save config
		err := h.updateConfig()

		return err
	}
	return nil
}

// UnsetStop set stop state to false
func (h *dfpRepository) UnsetStop() error {
	if h.state.IsStopped {
		h.state.IsStopped = false
		h.state.ShouldWash = false
		h.state.IsWashed = false
		h.eventer.Publish("stateChange", "isNotStop")

		// Save config
		err := h.updateConfig()

		return err
	}
	return nil
}

// SetEmergencyStop set emergency state to true
func (h *dfpRepository) SetEmergencyStop() error {
	if !h.state.IsEmergencyStopped {
		h.state.IsWashed = false
		h.state.IsEmergencyStopped = true
		h.eventer.Publish("stateChange", "isEmergencyStop")

		// Save config
		err := h.updateConfig()

		return err
	}
	return nil
}

// UnsetEmergencyStop set emergency state to false
func (h *dfpRepository) UnsetEmergencyStop() error {
	if h.state.IsEmergencyStopped {
		h.state.ShouldWash = false
		h.state.IsWashed = false
		h.state.IsEmergencyStopped = false
		h.eventer.Publish("stateChange", "isNotEmergencyStop")

		// Save config
		err := h.updateConfig()

		return err
	}
	return nil
}

// SetDisableSecurity set disable security state to true
func (h *dfpRepository) SetDisableSecurity() error {
	if !h.state.IsDisableSecurity {
		h.state.IsDisableSecurity = true
		h.eventer.Publish("stateChange", "isDisableSecurity")

		// Save config
		err := h.updateConfig()

		return err
	}
	return nil
}

// UnsetDisableSecurity set disable security state to false
func (h *dfpRepository) UnsetDisableSecurity() error {
	if h.state.IsDisableSecurity {
		h.state.IsDisableSecurity = false
		h.eventer.Publish("stateChange", "isNotDisableSecurity")

		// Save config
		err := h.updateConfig()

		return err
	}
	return nil
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

	// Save config
	err := h.updateConfig()

	return err
}

func (h *dfpRepository) String() string {
	return h.state.String()
}

func (h *dfpRepository) State() *models.DFPState {
	return h.state
}

func (h *dfpRepository) updateConfig() error {
	ctx := context.Background()
	config, err := h.config.Get(ctx)
	if err != nil {
		return err
	}

	config.SecurityDisabled = h.State().IsDisableSecurity
	config.Auto = h.State().IsAuto
	config.EmergencyStopped = h.State().IsEmergencyStopped
	config.Stopped = h.State().IsStopped
	config.LastWashing = h.State().LastWashing

	err = h.config.Update(ctx, config)
	if err != nil {
		return err
	}

	return nil
}
