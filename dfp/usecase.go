package dfp

import (
	"context"

	"github.com/disaster37/gobot-fat/models"
)

// Usecase represent the dfp usecase
type Usecase interface {
	Wash(ctx context.Context) error
	Stop(ctx context.Context) error
	Start(ctx context.Context) error
	ManualDrum(ctx context.Context, status bool) error
	ManualPump(ctx context.Context, status bool) error
	Security(ctx context.Context, status bool) error
	EmergencyStop(ctx context.Context, status bool) error
	GetState(ctx context.Context) (models.DFPState, error)
	GetIO(ctx context.Context) (models.DFPIO, error)
}
