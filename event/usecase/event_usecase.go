package usecase

import (
	"context"
	"time"

	"github.com/disaster37/gobot-fat/event"
	"github.com/disaster37/gobot-fat/models"
)

type eventUsecase struct {
	eventRepository event.Repository
	contextTimeout  time.Duration
}

// NewEventUsecase return implementation of event usecase
func NewEventUsecase(eventRepository event.Repository, timeout time.Duration) event.Usecase {
	return &eventUsecase{
		eventRepository: eventRepository,
		contextTimeout:  timeout,
	}
}

func (h *eventUsecase) Fetch(c context.Context, from int, size int) (res []*models.Event, nextFrom int, err error) {
	ctx, cancel := context.WithTimeout(c, h.contextTimeout)
	defer cancel()

	return h.eventRepository.Fetch(ctx, from, size)
}

func (h *eventUsecase) GetByID(c context.Context, id string) (*models.Event, error) {
	ctx, cancel := context.WithTimeout(c, h.contextTimeout)
	defer cancel()

	return h.eventRepository.GetByID(ctx, id)
}

func (h *eventUsecase) Search(c context.Context, query map[string]interface{}, minimalScoring float64) ([]*models.Event, error) {
	ctx, cancel := context.WithTimeout(c, h.contextTimeout)
	defer cancel()

	return h.eventRepository.Search(ctx, query, minimalScoring)
}

func (h *eventUsecase) Update(c context.Context, object *models.Event) error {
	ctx, cancel := context.WithTimeout(c, h.contextTimeout)
	defer cancel()

	return h.eventRepository.Update(ctx, object)
}

func (h *eventUsecase) Store(c context.Context, object *models.Event) error {
	ctx, cancel := context.WithTimeout(c, h.contextTimeout)
	defer cancel()

	return h.eventRepository.Store(ctx, object)
}

func (h *eventUsecase) Delete(c context.Context, id string) error {
	ctx, cancel := context.WithTimeout(c, h.contextTimeout)
	defer cancel()

	return h.eventRepository.Delete(ctx, id)
}
