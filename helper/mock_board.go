package helper

import (
	"sync"

	"gobot.io/x/gobot/drivers/gpio"
)

type MockPlateform struct {
	name                    string
	mtx                     sync.Mutex
	testAdaptorReconnect    func() error
	testAdaptorDigitalWrite func(pin string, val byte) (err error)
	testAdaptorServoWrite   func(pin string, val byte) (err error)
	testAdaptorPwmWrite     func(pin string, val byte) (err error)
	testAdaptorAnalogRead   func(ping string) (val int, err error)
	testAdaptorDigitalRead  func(ping string) (val int, err error)
	testAdaptorValueRead    func(name string) (val interface{}, err error)
	testAdaptorValuesRead   func() (vals map[string]interface{}, err error)
	testAdaptorFunctionCall func(name string, parameters string) (val int, err error)
	DigitalPinState         map[string]int
	ValueReadState          map[string]interface{}
	CallFunctionState       map[string]int
	invertedInitialState    map[string]bool
	expectedError           error
}

func (m *MockPlateform) SetInvertInitialPinState(pin string) { m.invertedInitialState[pin] = true }

// Adaptor interface
func (m *MockPlateform) Name() string     { return "test" }
func (m *MockPlateform) SetName(n string) { return }
func (m *MockPlateform) Connect() error   { return nil }
func (m *MockPlateform) Finalize() error  { return nil }
func (m *MockPlateform) SetInputPullup(listPins []*gpio.ButtonDriver) (err error) {

	for _, button := range listPins {

		if !m.invertedInitialState[button.Pin()] {
			button.DefaultState = 1

			// When InputPullup, the default button state is 1
			m.DigitalPinState[button.Pin()] = 1
		}
	}

	return

}

// Arest interface
func (t *MockPlateform) TestReconnect(f func() error) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	t.testAdaptorReconnect = f
}
func (t *MockPlateform) Reconnect() error {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	return t.testAdaptorReconnect()
}

// gpio interface
func (t *MockPlateform) TestAdaptorDigitalWrite(f func(pin string, val byte) (err error)) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	t.testAdaptorDigitalWrite = f
}
func (t *MockPlateform) TestAdaptorServoWrite(f func(pin string, val byte) (err error)) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	t.testAdaptorServoWrite = f
}
func (t *MockPlateform) TestAdaptorPwmWrite(f func(pin string, val byte) (err error)) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	t.testAdaptorPwmWrite = f
}
func (t *MockPlateform) TestAdaptorAnalogRead(f func(pin string) (val int, err error)) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	t.testAdaptorAnalogRead = f
}
func (t *MockPlateform) TestAdaptorDigitalRead(f func(pin string) (val int, err error)) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	t.testAdaptorDigitalRead = f
}
func (t *MockPlateform) TestAdaptorReadValue(f func(pin string) (val int, err error)) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	t.testAdaptorDigitalRead = f
}

func (t *MockPlateform) ServoWrite(pin string, val byte) (err error) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	return t.testAdaptorServoWrite(pin, val)
}
func (t *MockPlateform) PwmWrite(pin string, val byte) (err error) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	return t.testAdaptorPwmWrite(pin, val)
}
func (t *MockPlateform) AnalogRead(pin string) (val int, err error) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	return t.testAdaptorAnalogRead(pin)
}
func (t *MockPlateform) DigitalRead(pin string) (val int, err error) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	return t.testAdaptorDigitalRead(pin)
}
func (t *MockPlateform) DigitalWrite(pin string, val byte) (err error) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	return t.testAdaptorDigitalWrite(pin, val)
}

// Extra interface
func (t *MockPlateform) TestAdaptorValueRead(f func(name string) (val interface{}, err error)) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	t.testAdaptorValueRead = f
}
func (t *MockPlateform) TestAdaptorValuesRead(f func() (vals map[string]interface{}, err error)) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	t.testAdaptorValuesRead = f
}
func (t *MockPlateform) TestAdaptorFunctionCall(f func(name string, parameters string) (val int, err error)) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	t.testAdaptorFunctionCall = f
}

func (t *MockPlateform) ValueRead(name string) (val interface{}, err error) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	return t.testAdaptorValueRead(name)
}
func (t *MockPlateform) ValuesRead() (vals map[string]interface{}, err error) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	return t.testAdaptorValuesRead()
}
func (t *MockPlateform) FunctionCall(name string, parameters string) (val int, err error) {
	t.mtx.Lock()
	defer t.mtx.Unlock()
	return t.testAdaptorFunctionCall(name, parameters)
}

func (t *MockPlateform) SetError(err error) {
	t.expectedError = err
}

func (t *MockPlateform) init() {
	t.testAdaptorDigitalRead = func(pin string) (val int, err error) {
		if t.expectedError != nil {
			return 0, t.expectedError
		}

		return t.DigitalPinState[pin], nil
	}

	t.testAdaptorValueRead = func(name string) (val interface{}, err error) {
		if t.expectedError != nil {
			return 0, t.expectedError
		}
		if t.ValueReadState[name] != nil {
			return t.ValueReadState[name], nil
		}
		return nil, nil
	}

	t.testAdaptorFunctionCall = func(name string, parameters string) (val int, err error) {
		if t.expectedError != nil {
			return 0, t.expectedError
		}
		return t.CallFunctionState[name], nil
	}

	t.testAdaptorDigitalWrite = func(pin string, val byte) (err error) {
		if t.expectedError != nil {
			return t.expectedError
		}
		t.DigitalPinState[pin] = int(val)
		return nil
	}

}

func NewMockPlateform() *MockPlateform {
	m := &MockPlateform{
		DigitalPinState:      make(map[string]int),
		ValueReadState:       make(map[string]interface{}),
		CallFunctionState:    make(map[string]int),
		invertedInitialState: make(map[string]bool),

		testAdaptorServoWrite: func(pin string, val byte) (err error) {
			return nil
		},
		testAdaptorPwmWrite: func(pin string, val byte) (err error) {
			return nil
		},
		testAdaptorAnalogRead: func(pin string) (val int, err error) {
			return 0, nil
		},
		testAdaptorValuesRead: func() (vals map[string]interface{}, err error) {
			return map[string]interface{}{
				"test": 99,
			}, nil
		},
		testAdaptorReconnect: func() (err error) { return nil },
	}

	m.init()

	return m
}
