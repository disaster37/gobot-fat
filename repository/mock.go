package repository

import (
	"context"
	"reflect"
)

type Mock struct {
	Result      interface{}
	Err         error
	ShouldError bool
	IsUpdate    bool
	IsCreate    bool
	IsGet       bool
	IsList      bool
}

func NewMock() *Mock {
	return &Mock{
		ShouldError: false,
	}
}

func (h *Mock) Reset() {
	h.Result = nil
	h.Err = nil
	h.ShouldError = false
	h.IsUpdate = false
	h.IsCreate = false
	h.IsGet = false
	h.IsList = false
}

func (h *Mock) Get(ctx context.Context, id uint, data interface{}) error {

	h.IsGet = true

	if h.Result == nil {
		return ErrRecordNotFoundError
	} else {
		reflect.ValueOf(data).Elem().Set(reflect.ValueOf(h.Result).Elem())
	}

	if h.ShouldError {
		return h.Err
	}
	return nil
}

func (h *Mock) List(ctx context.Context, listData interface{}) error {

	h.IsList = true

	if h.Result == nil {
		return nil
	}

	rt := reflect.ValueOf(listData).Elem()

	rt.Set(reflect.ValueOf(h.Result))

	if h.ShouldError {
		return h.Err
	}
	return nil
}

func (h *Mock) Update(ctx context.Context, data interface{}) error {

	h.IsUpdate = true

	if h.ShouldError {
		return h.Err
	}

	h.Result = data
	return nil
}

func (h *Mock) Create(ctx context.Context, data interface{}) error {

	h.IsCreate = true

	if h.ShouldError {
		return h.Err
	}
	h.Result = data
	return nil
}
