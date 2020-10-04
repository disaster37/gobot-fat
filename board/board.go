package board

import (
	"context"
	"time"

	"github.com/disaster37/gobot-fat/models"
	log "github.com/sirupsen/logrus"
)

// Board represent generic board
type Board interface {
	// IsOnline permit to know if board is online
	IsOnline() bool

	// Start start the main function
	Start(ctx context.Context) error

	// Stop interrupt the main function
	Stop(ctx context.Context) error

	// Name return the board name
	Name() string

	// Board return the board data
	Board() *models.Board
}

// NewHandler is a generic handler that run in background
func NewHandler(ctx context.Context, loopDuration time.Duration, chStop chan bool, process func(ctx context.Context)) {

	go func() {

		timer := time.NewTimer(loopDuration)
		ctx, cancel := context.WithCancel(ctx)
		defer cancel()

		for {
			select {
			case <-ctx.Done():
				log.Debug("Context timeout on handler")
				cancel()
				return
			case <-chStop:
				log.Debug("Handler stopped")
				cancel()
				return
			case <-timer.C:
				process(ctx)
				timer.Reset(loopDuration)
			}
		}
	}()

}
