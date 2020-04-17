package dfp

import (
	"github.com/disaster37/gobot-fat/models"
)

// Repository is the interface to manage the state of DFP
type Repository interface {
	SetWashed() (bool, error)
	SetShouldWash() (bool, error)
	UnsetShouldWash() (bool, error)
	UnsetWashed() (bool, error)
	SetSecurity() (bool, error)
	UnsetSecurity() (bool, error)
	SetAuto() (bool, error)
	UnsetAuto() (bool, error)
	SetStop() (bool, error)
	UnsetStop() (bool, error)
	SetEmergencyStop() (bool, error)
	UnsetEmergencyStop() (bool, error)
	SetDisableSecurity() (bool, error)
	UnsetDisableSecurity() (bool, error)
	UpdateLastWashing() error
	CanWash() bool
	CanStartMotor() bool
	LastWashDurationSecond() uint64
	String() string
	State() *models.DFPState
}
