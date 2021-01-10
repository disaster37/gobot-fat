package repository

import (
	"context"
	"fmt"
	"reflect"
)

type MockBase struct{}

func (m *MockBase) Get(ctx context.Context, id uint, data interface{}) error { return nil }
func (m *MockBase) List(ctx context.Context, listData interface{}) error     { return nil }
func (m *MockBase) Update(ctx context.Context, data interface{}) error       { return nil }
func (m *MockBase) Create(ctx context.Context, data interface{}) error       { return nil }

type Mock struct {
	expectedData   interface{}
	expectedError  error
	callMethod     string
	callParameters []interface{}
	testGet        func(ctx context.Context, id uint, data interface{}) error
	testList       func(ctx context.Context, listData interface{}) error
	testUpdate     func(ctx context.Context, data interface{}) error
	testCreate     func(ctx context.Context, data interface{}) error
}

func NewMockBase() *MockBase {
	return &MockBase{}
}

func NewMock() *Mock {
	m := &Mock{}
	m.init()

	return m
}

func (h *Mock) init() {
	h.testGet = func(ctx context.Context, id uint, data interface{}) error {
		h.callMethod = "Get"
		h.callParameters = make([]interface{}, 0)
		h.callParameters = append(h.callParameters, ctx)
		h.callParameters = append(h.callParameters, id)
		h.callParameters = append(h.callParameters, data)
		if h.expectedError != nil {
			return h.expectedError
		}

		if h.expectedData == nil {
			return ErrRecordNotFoundError
		}

		reflect.ValueOf(data).Elem().Set(reflect.ValueOf(h.expectedData).Elem())

		return nil
	}

	h.testList = func(ctx context.Context, listData interface{}) error {
		h.callMethod = "List"
		h.callParameters = make([]interface{}, 0)
		h.callParameters = append(h.callParameters, ctx)
		h.callParameters = append(h.callParameters, listData)
		if h.expectedError != nil {
			return h.expectedError
		}

		if h.expectedData == nil {
			return nil
		}

		rt := reflect.ValueOf(listData).Elem()
		rt.Set(reflect.ValueOf(h.expectedData))

		return nil
	}

	h.testCreate = func(ctx context.Context, data interface{}) error {
		h.callMethod = "Create"
		h.callParameters = make([]interface{}, 0)
		h.callParameters = append(h.callParameters, ctx)
		h.callParameters = append(h.callParameters, data)
		if h.expectedError != nil {
			return h.expectedError
		}

		h.expectedData = data
		return nil
	}

	h.testUpdate = func(ctx context.Context, data interface{}) error {
		h.callMethod = "Update"
		h.callParameters = make([]interface{}, 0)
		h.callParameters = append(h.callParameters, ctx)
		h.callParameters = append(h.callParameters, data)
		if h.expectedError != nil {
			return h.expectedError
		}

		h.expectedData = data
		return nil
	}

}

func (h *Mock) TestGet(f func(ctx context.Context, id uint, data interface{}) error) {
	h.testGet = f
}

func (h *Mock) TestList(f func(ctx context.Context, listData interface{}) error) {
	h.testList = f
}

func (h *Mock) TestUpdate(f func(ctx context.Context, data interface{}) error) {
	h.testUpdate = f
}

func (h *Mock) TestCreate(f func(ctx context.Context, data interface{}) error) {
	h.testCreate = f
}

func (h *Mock) Get(ctx context.Context, id uint, data interface{}) error {
	return h.testGet(ctx, id, data)
}

func (h *Mock) List(ctx context.Context, listData interface{}) error {

	return h.testList(ctx, listData)
}

func (h *Mock) Update(ctx context.Context, data interface{}) error {

	return h.testUpdate(ctx, data)
}

func (h *Mock) Create(ctx context.Context, data interface{}) error {

	return h.testCreate(ctx, data)
}

func (h *Mock) ExpectCall(functionName string, params ...interface{}) (isExpect bool, reason string) {

	if functionName != h.callMethod {
		return false, fmt.Sprintf("Call method '%s', expect method '%s'", h.callMethod, functionName)
	}

	isMatch := true
	for _, param := range params {
		isFound := false
		for _, calledParam := range h.callParameters {
			if reflect.DeepEqual(param, calledParam) {
				isFound = true
				break
			}
		}

		if !isFound {
			isMatch = false
			break
		}

	}

	if !isMatch {
		return false, "Parameters mistmatch"
	}

	return isMatch, ""

}

func (h *Mock) ExpectResult() interface{} {
	return h.expectedData
}

func (h *Mock) ExpectError() error {
	return h.expectedError
}

func (h *Mock) SetData(data interface{}) {
	h.expectedData = data
}

func (h *Mock) SetError(err error) {
	h.expectedError = err
}
