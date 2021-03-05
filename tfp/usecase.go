package tfp

import (
	"context"

	"github.com/disaster37/gobot-fat/models"
)

// Usecase represent the tfp usecase
type Usecase interface {
	PondPump(ctx context.Context, status bool) error
	WaterfallPump(ctx context.Context, status bool) error
	UVC1(ctx context.Context, status bool) error
	UVC2(ctx context.Context, status bool) error
	PondBubble(ctx context.Context, status bool) error
	FilterBubble(ctx context.Context, status bool) error
	UVC1BlisterNew(ctx context.Context) error
	UVC2BlisterNew(ctx context.Context) error
	OzoneBlisterNew(ctx context.Context) error
	GetState(ctx context.Context) (models.TFPState, error)
	GetIO(ctx context.Context) (models.TFPIO, error)
}
