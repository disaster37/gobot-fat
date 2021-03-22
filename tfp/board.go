package tfp

import (
	"context"

	"github.com/disaster37/gobot-fat/board"
	"github.com/disaster37/gobot-fat/models"
)

// Board is the interface to handle I/O
type Board interface {
	StartPondPump(ctx context.Context) error
	StopPondPump(ctx context.Context) error
	StartWaterfallPump(ctx context.Context) error
	StopWaterfallPump(ctx context.Context) error
	StartUVC1(ctx context.Context) error
	StopUVC1(ctx context.Context) error
	StartUVC2(ctx context.Context) error
	StopUVC2(ctx context.Context) error
	StartPondBubble(ctx context.Context) error
	StopPondBubble(ctx context.Context) error
	StartFilterBubble(ctx context.Context) error
	StopFilterBubble(ctx context.Context) error
	StopRelais(ctx context.Context) error
	State() models.TFPState
	IO() models.TFPIO
	Config() models.TFPConfig
	board.Board
}
