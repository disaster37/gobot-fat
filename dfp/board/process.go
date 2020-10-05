package dfpboard

import (
	"context"
	"time"

	"github.com/disaster37/gobot-fat/models"
	log "github.com/sirupsen/logrus"
)

func (h *DFPBoard) wash() {

}

func (h *DFPBoard) updateState() {

	h.mutexState.Lock()
	defer h.mutexState.Unlock()

}

func (h *DFPBoard) sendEvent(ctx context.Context, kind string, name string, args ...interface{}) {
	event := &models.Event{
		SourceID:   h.state.Name,
		SourceName: h.state.Name,
		Timestamp:  time.Now(),
		EventType:  name,
		EventKind:  kind,
	}
	err := h.eventUsecase.Store(ctx, event)
	if err != nil {
		log.Errorf("Error when store new event: %s", err.Error())
	}
}
