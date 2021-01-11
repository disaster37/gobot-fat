package usecase

import (
	"context"
)

type MockUsecasetBase struct{}

func (m *MockUsecasetBase) Get(ctx context.Context, id uint, data interface{}) error { return nil }
func (m *MockUsecasetBase) List(ctx context.Context, listData interface{}) error     { return nil }
func (m *MockUsecasetBase) Update(ctx context.Context, data interface{}) error       { return nil }
func (m *MockUsecasetBase) Create(ctx context.Context, data interface{}) error       { return nil }
func (m *MockUsecasetBase) Init(ctx context.Context, data interface{}) error         { return nil }

func NewMockUsecasetBase() UsecaseCRUD {
	return &MockUsecasetBase{}
}
