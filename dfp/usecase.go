package dfp

import (
	"context"

	"github.com/disaster37/gobot-fat/models"
)

// Usecase represent the dfp usecase
type Usecase interface {
	Wash(ctx context.Context) error
	Stop(ctx context.Context, status bool) error
	EmergencyStop(ctx context.Context, status bool) error
	Auto(ctx context.Context, status bool) error
	ForceWashingDrum(ctx context.Context, status bool) error
	ForceWashingPump(ctx context.Context, status bool) error
	DisableSecurity(ctx context.Context, status bool) error
	GetState(ctx context.Context) (*models.DFPState, error)
	StartRobot(ctx context.Context) error
	StopRobot(ctx context.Context) error
}
