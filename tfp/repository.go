package tfp

import (
	"github.com/disaster37/gobot-fat/models"
)

// Repository is the interface to manage the state of DFP
type Repository interface {
	StartPondPump() error
	StopPondPump() error
	StartWaterfallPump() error
	StopWaterfallPump() error
	StartUVC1() error
	StopUVC1() error
	StartUVC2() error
	StopUVC2() error
	StartPondBubble() error
	StopPondBubble() error
	StartFilterBubble() error
	StopFilterBubble() error
	CanStartRelay() bool
	String() string
	State() *models.TFPState
}
