package dfpState

import(
	"context"

	"github.com/disaster37/gobot-fat/models"
)

// Usecase represent the DFP state usecase
type Usecase interface {
	Fetch(ctx context.Context) (*models.DFPState, error)
}