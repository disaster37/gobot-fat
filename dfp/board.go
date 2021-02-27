package dfp

import (
	"context"

	"github.com/disaster37/gobot-fat/board"
	"github.com/disaster37/gobot-fat/models"
)

type Board interface {
	// StartDFP put dfp on auto
	StartDFP(ctx context.Context) error

	// StopDFP stop dfp and disable auto
	StopDFP(ctx context.Context) error

	// ForceWashing start a washing cycle
	ForceWashing(ctx context.Context) error

	// StartManualDrum force start drum motor
	StartManualDrum(ctx context.Context) error

	// StopManualDrum force stop drum motor
	StopManualDrum(ctx context.Context) error

	// StartManualPump force start pump
	StartManualPump(ctx context.Context) error

	// StopManualPump force stop pump
	StopManualPump(ctx context.Context) error

	State() models.DFPState

	IO() models.DFPIO

	board.Board
}
