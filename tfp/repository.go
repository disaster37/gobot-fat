package tfp

import (
	"github.com/disaster37/gobot-fat/models"
)

// Repository is the interface to manage the state of DFP
type Repository interface {
	StartPondPump() (bool, error)
	StopPondPump() (bool, error)
	StartWaterfallPump() (bool, error)
	StopWaterfallPump() (bool, error)
	StartUVC1() (bool, error)
	StopUVC1() (bool, error)
	StartUVC2() (bool, error)
	StopUVC2() (bool, error)
	StartPondBubble() (bool, error)
	StopPondBubble() (bool, error)
	StartFilterBubble() (bool, error)
	StopFilterBubble() (bool, error)
	CanStartRelay() bool
	String() string
	State() *models.TFPState
}
