package dfp

import (
	"github.com/disaster37/gobot-fat/models"
)

// Repository is the interface to manage the state of DFP
type Repository interface {
	SetWashed() error
	SetShouldWash() error
	UnsetShouldWash() error
	UnsetWashed() error
	SetSecurity() error
	UnsetSecurity() error
	SetAuto() error
	UnsetAuto() error
	SetStop() error
	UnsetStop() error
	SetEmergencyStop() error
	UnsetEmergencyStop() error
	SetDisableSecurity() error
	UnsetDisableSecurity() error
	UpdateLastWashing() error
	CanWash() bool
	CanStartMotor() bool
	LastWashDurationSecond() uint64
	String() string
	State() *models.DFPState
}
