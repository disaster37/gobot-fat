package usecase

import (
	"context"

	"github.com/disaster37/gobot-fat/event"
	"github.com/disaster37/gobot-fat/models"
)

type MockEventBase struct{}

func (m *MockEventBase) Fetch(ctx context.Context, from int, size int) (res []*models.Event, nextFrom int, err error) {
	return
}
func (m *MockEventBase) GetByID(ctx context.Context, id string) (event *models.Event, err error) {
	return
}
func (m *MockEventBase) Search(ctx context.Context, query map[string]interface{}, minimalScoring float64) (listEvents []*models.Event, err error) {
	return
}
func (m *MockEventBase) Update(ctx context.Context, object *models.Event) (err error) { return }
func (m *MockEventBase) Store(ctx context.Context, object *models.Event) (err error)  { return }
func (m *MockEventBase) Delete(ctx context.Context, id string) (err error)            { return }

func NewMockEventBase() event.Usecase {
	return &MockEventBase{}
}
