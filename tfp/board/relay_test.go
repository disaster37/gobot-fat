package tfpboard

import (
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStartStopPondPump(t *testing.T) {
	board, adaptor := initTestBoard()
	err := board.Start(context.Background())
	if err != nil {
		panic(err)
	}

	// Start pond pomp when already started
	adaptor.DigitalPinState["1"] = 0
	err = board.StartPondPump(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 0, adaptor.DigitalPinState["1"])

	// Start pond pomp when stopped
	adaptor.DigitalPinState["1"] = 1
	err = board.StartPondPump(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 0, adaptor.DigitalPinState["1"])

	// Start pond pomp when error
	adaptor.SetError(errors.New("test"))
	err = board.StartPondPump(context.Background())
	assert.Error(t, err)
	adaptor.SetError(nil)

	// Stop pond pomp when already stopped
	adaptor.DigitalPinState["1"] = 1
	err = board.StopPondPump(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 1, adaptor.DigitalPinState["1"])

	// Stop pond pomp when started
	adaptor.DigitalPinState["1"] = 0
	err = board.StopPondPump(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 1, adaptor.DigitalPinState["1"])

	// Stop pond pomp when error
	adaptor.SetError(errors.New("test"))
	err = board.StopPondPump(context.Background())
	assert.Error(t, err)
	adaptor.SetError(nil)

	// Stop UVC1 and UVC2 when stop pond pump
	board.state.UVC1Running = true
	board.state.UVC2Running = true
	err = board.StopPondPump(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 1, adaptor.DigitalPinState["2"])
	assert.Equal(t, 1, adaptor.DigitalPinState["3"])

	// Can't start pond pomp when emergency stop
	board.state.IsEmergencyStopped = true
	adaptor.DigitalPinState["1"] = 1
	err = board.StartPondPump(context.Background())
	assert.Error(t, err)
	assert.Equal(t, 1, adaptor.DigitalPinState["1"])

	// Can't start pond pomp when security
	board.state.IsEmergencyStopped = false
	board.state.IsSecurity = true
	adaptor.DigitalPinState["1"] = 1
	err = board.StartPondPump(context.Background())
	assert.Error(t, err)
	assert.Equal(t, 1, adaptor.DigitalPinState["1"])

	// Can start  pond pomp when security and disable security
	board.state.IsEmergencyStopped = false
	board.state.IsSecurity = true
	board.state.IsDisableSecurity = true
	adaptor.DigitalPinState["1"] = 1
	err = board.StartPondPump(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 0, adaptor.DigitalPinState["1"])

}

func TestStartStopUVC1(t *testing.T) {

	board, adaptor := initTestBoard()
	err := board.Start(context.Background())
	if err != nil {
		panic(err)
	}
	board.state.PondPumpRunning = true

	// Start UVC1 when already started
	adaptor.DigitalPinState["2"] = 0
	err = board.StartUVC1(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 0, adaptor.DigitalPinState["2"])

	// Start UVC1 when stopped
	adaptor.DigitalPinState["2"] = 1
	err = board.StartUVC1(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 0, adaptor.DigitalPinState["2"])

	// Start UVC1 when error
	adaptor.SetError(errors.New("test"))
	err = board.StartUVC1(context.Background())
	assert.Error(t, err)
	adaptor.SetError(nil)

	// Stop UVC1 when already stopped
	adaptor.DigitalPinState["2"] = 1
	err = board.StopUVC1(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 1, adaptor.DigitalPinState["2"])

	// Stop UVC1 when started
	adaptor.DigitalPinState["2"] = 0
	err = board.StopUVC1(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 1, adaptor.DigitalPinState["2"])

	// Stop UVC1 when error
	adaptor.SetError(errors.New("test"))
	err = board.StopUVC1(context.Background())
	assert.Error(t, err)
	adaptor.SetError(nil)

	// Can't start UVC1 when pond pump is stopped
	board.state.PondPumpRunning = false
	adaptor.DigitalPinState["2"] = 1
	err = board.StartUVC1(context.Background())
	assert.Error(t, err)
	assert.Equal(t, 1, adaptor.DigitalPinState["2"])

	// Can stop UVC1 when pond pump is started
	board.state.PondPumpRunning = false
	adaptor.DigitalPinState["2"] = 0
	err = board.StopUVC1(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 1, adaptor.DigitalPinState["2"])

	// Can't start UVC1 when emergency stop
	board.state.PondPumpRunning = true
	board.state.IsEmergencyStopped = true
	adaptor.DigitalPinState["2"] = 1
	err = board.StartUVC1(context.Background())
	assert.Error(t, err)
	assert.Equal(t, 1, adaptor.DigitalPinState["2"])

	// Can't start UVC1 when security
	board.state.IsEmergencyStopped = false
	board.state.IsSecurity = true
	adaptor.DigitalPinState["2"] = 1
	err = board.StartUVC1(context.Background())
	assert.Error(t, err)
	assert.Equal(t, 1, adaptor.DigitalPinState["2"])

	// Can start  UVC1 when security and disable security
	board.state.IsEmergencyStopped = false
	board.state.IsSecurity = true
	board.state.IsDisableSecurity = true
	adaptor.DigitalPinState["2"] = 1
	err = board.StartUVC1(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 0, adaptor.DigitalPinState["2"])
}

func TestStartStopUVC2(t *testing.T) {
	board, adaptor := initTestBoard()
	err := board.Start(context.Background())
	if err != nil {
		panic(err)
	}
	board.state.PondPumpRunning = true

	// Start UVC2 when already started
	adaptor.DigitalPinState["3"] = 0
	err = board.StartUVC2(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 0, adaptor.DigitalPinState["3"])

	// Start UVC2 when stopped
	adaptor.DigitalPinState["3"] = 1
	err = board.StartUVC2(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 0, adaptor.DigitalPinState["3"])

	// Start UVC2 when error
	adaptor.SetError(errors.New("test"))
	err = board.StartUVC2(context.Background())
	assert.Error(t, err)
	adaptor.SetError(nil)

	// Stop UVC2 when already stopped
	adaptor.DigitalPinState["3"] = 1
	err = board.StopUVC2(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 1, adaptor.DigitalPinState["3"])

	// Stop UVC2 when started
	adaptor.DigitalPinState["3"] = 0
	err = board.StopUVC2(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 1, adaptor.DigitalPinState["3"])

	// Stop UVC2 when error
	adaptor.SetError(errors.New("test"))
	err = board.StopUVC2(context.Background())
	assert.Error(t, err)
	adaptor.SetError(nil)

	// Can't start UVC2 when pond pump is stopped
	board.state.PondPumpRunning = false
	adaptor.DigitalPinState["3"] = 1
	err = board.StartUVC2(context.Background())
	assert.Error(t, err)
	assert.Equal(t, 1, adaptor.DigitalPinState["3"])

	// Can stop UVC2 when pond pump is started
	board.state.PondPumpRunning = false
	adaptor.DigitalPinState["3"] = 0
	err = board.StopUVC2(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 1, adaptor.DigitalPinState["3"])

	// Can't start UVC2 when emergency stop
	board.state.PondPumpRunning = true
	board.state.IsEmergencyStopped = true
	adaptor.DigitalPinState["3"] = 1
	err = board.StartUVC2(context.Background())
	assert.Error(t, err)
	assert.Equal(t, 1, adaptor.DigitalPinState["3"])

	// Can't start UVC2 when security
	board.state.IsEmergencyStopped = false
	board.state.IsSecurity = true
	adaptor.DigitalPinState["3"] = 1
	err = board.StartUVC2(context.Background())
	assert.Error(t, err)
	assert.Equal(t, 1, adaptor.DigitalPinState["3"])

	// Can start  UVC2 when security and disable security
	board.state.IsEmergencyStopped = false
	board.state.IsSecurity = true
	board.state.IsDisableSecurity = true
	adaptor.DigitalPinState["3"] = 1
	err = board.StartUVC2(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 0, adaptor.DigitalPinState["3"])
}

func TestStartPondPumpWithUVC(t *testing.T) {
	board, adaptor := initTestBoard()
	err := board.Start(context.Background())
	if err != nil {
		panic(err)
	}

	// When already started
	err = board.StartPondPumpWithUVC(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 0, adaptor.DigitalPinState["1"])
	assert.Equal(t, 0, adaptor.DigitalPinState["2"])
	assert.Equal(t, 0, adaptor.DigitalPinState["3"])

	// When all stopped
	adaptor.DigitalPinState["1"] = 1
	adaptor.DigitalPinState["2"] = 1
	adaptor.DigitalPinState["3"] = 1
	err = board.StartPondPumpWithUVC(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 0, adaptor.DigitalPinState["1"])
	assert.Equal(t, 0, adaptor.DigitalPinState["2"])
	assert.Equal(t, 0, adaptor.DigitalPinState["3"])
}

func TestStartStopPondBubble(t *testing.T) {
	board, adaptor := initTestBoard()
	err := board.Start(context.Background())
	if err != nil {
		panic(err)
	}

	// Start pond bubble when already started
	adaptor.DigitalPinState["4"] = 0
	err = board.StartPondBubble(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 0, adaptor.DigitalPinState["4"])

	// Start pond bubble when stopped
	adaptor.DigitalPinState["4"] = 1
	err = board.StartPondBubble(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 0, adaptor.DigitalPinState["4"])

	// Start pond bubble when error
	adaptor.SetError(errors.New("test"))
	err = board.StartPondBubble(context.Background())
	assert.Error(t, err)
	adaptor.SetError(nil)

	// Stop pond bubble when already stopped
	adaptor.DigitalPinState["4"] = 1
	err = board.StopPondBubble(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 1, adaptor.DigitalPinState["4"])

	// Stop pond bubble when started
	adaptor.DigitalPinState["4"] = 0
	err = board.StopPondBubble(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 1, adaptor.DigitalPinState["4"])

	// Stop pond bubble when error
	adaptor.SetError(errors.New("test"))
	err = board.StopPondBubble(context.Background())
	assert.Error(t, err)
	adaptor.SetError(nil)

	// Can't start pond bubble when emergency stop
	board.state.IsEmergencyStopped = true
	adaptor.DigitalPinState["4"] = 1
	err = board.StartPondBubble(context.Background())
	assert.Error(t, err)
	assert.Equal(t, 1, adaptor.DigitalPinState["4"])

	// Can't start pond bubble when security
	board.state.IsEmergencyStopped = false
	board.state.IsSecurity = true
	adaptor.DigitalPinState["4"] = 1
	err = board.StartPondBubble(context.Background())
	assert.Error(t, err)
	assert.Equal(t, 1, adaptor.DigitalPinState["4"])

	// Can start  pond pomp when security and disable security
	board.state.IsEmergencyStopped = false
	board.state.IsSecurity = true
	board.state.IsDisableSecurity = true
	adaptor.DigitalPinState["4"] = 1
	err = board.StartPondBubble(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 0, adaptor.DigitalPinState["4"])
}

func TestStartStopFilterBubble(t *testing.T) {
	board, adaptor := initTestBoard()
	err := board.Start(context.Background())
	if err != nil {
		panic(err)
	}

	// Start filter bubble when already started
	adaptor.DigitalPinState["5"] = 0
	err = board.StartFilterBubble(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 0, adaptor.DigitalPinState["5"])

	// Start filter bubble when stopped
	adaptor.DigitalPinState["5"] = 1
	err = board.StartFilterBubble(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 0, adaptor.DigitalPinState["5"])

	// Start filter bubble when error
	adaptor.SetError(errors.New("test"))
	err = board.StartFilterBubble(context.Background())
	assert.Error(t, err)
	adaptor.SetError(nil)

	// Stop filter bubble when already stopped
	adaptor.DigitalPinState["5"] = 1
	err = board.StopFilterBubble(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 1, adaptor.DigitalPinState["5"])

	// Stop filter bubble when started
	adaptor.DigitalPinState["5"] = 0
	err = board.StopFilterBubble(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 1, adaptor.DigitalPinState["5"])

	// Stop filter bubble when error
	adaptor.SetError(errors.New("test"))
	err = board.StopFilterBubble(context.Background())
	assert.Error(t, err)
	adaptor.SetError(nil)

	// Can't start filter bubble when emergency stop
	board.state.IsEmergencyStopped = true
	adaptor.DigitalPinState["5"] = 1
	err = board.StartFilterBubble(context.Background())
	assert.Error(t, err)
	assert.Equal(t, 1, adaptor.DigitalPinState["5"])

	// Can't start filter bubble when security
	board.state.IsEmergencyStopped = false
	board.state.IsSecurity = true
	adaptor.DigitalPinState["5"] = 1
	err = board.StartFilterBubble(context.Background())
	assert.Error(t, err)
	assert.Equal(t, 1, adaptor.DigitalPinState["5"])

	// Can start  filter bubble when security and disable security
	board.state.IsEmergencyStopped = false
	board.state.IsSecurity = true
	board.state.IsDisableSecurity = true
	adaptor.DigitalPinState["5"] = 1
	err = board.StartFilterBubble(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 0, adaptor.DigitalPinState["5"])
}

func TestStartStopWaterfallPomp(t *testing.T) {
	board, adaptor := initTestBoard()
	err := board.Start(context.Background())
	if err != nil {
		panic(err)
	}

	// Start waterfall pomp when already started
	adaptor.DigitalPinState["6"] = 1
	err = board.StartWaterfallPump(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 1, adaptor.DigitalPinState["6"])

	// Start waterfall pomp when stopped
	adaptor.DigitalPinState["6"] = 0
	err = board.StartWaterfallPump(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 1, adaptor.DigitalPinState["6"])

	// Start waterfall pomp when error
	adaptor.SetError(errors.New("test"))
	err = board.StartWaterfallPump(context.Background())
	assert.Error(t, err)
	adaptor.SetError(nil)

	// Stop waterfall pomp when already stopped
	adaptor.DigitalPinState["6"] = 0
	err = board.StopWaterfallPump(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 0, adaptor.DigitalPinState["6"])

	// Stop waterfall pomp when started
	adaptor.DigitalPinState["6"] = 1
	err = board.StopWaterfallPump(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 0, adaptor.DigitalPinState["6"])

	// Stop waterfall pomp when error
	adaptor.SetError(errors.New("test"))
	err = board.StopWaterfallPump(context.Background())
	assert.Error(t, err)
	adaptor.SetError(nil)

	// Can't start waterfall pomp when emergency stop
	board.state.IsEmergencyStopped = true
	adaptor.DigitalPinState["6"] = 0
	err = board.StartWaterfallPump(context.Background())
	assert.Error(t, err)
	assert.Equal(t, 0, adaptor.DigitalPinState["6"])

	// Can't start waterfall pomp when security
	board.state.IsEmergencyStopped = false
	board.state.IsSecurity = true
	adaptor.DigitalPinState["6"] = 0
	err = board.StartWaterfallPump(context.Background())
	assert.Error(t, err)
	assert.Equal(t, 0, adaptor.DigitalPinState["6"])

	// Can start waterfall pomp when security and disable security
	board.state.IsEmergencyStopped = false
	board.state.IsSecurity = true
	board.state.IsDisableSecurity = true
	adaptor.DigitalPinState["6"] = 0
	err = board.StartWaterfallPump(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 1, adaptor.DigitalPinState["6"])
}
func TestStopRelais(t *testing.T) {
	board, adaptor := initTestBoard()
	err := board.Start(context.Background())
	if err != nil {
		panic(err)
	}

	adaptor.DigitalPinState["1"] = 0
	adaptor.DigitalPinState["2"] = 0
	adaptor.DigitalPinState["3"] = 0
	adaptor.DigitalPinState["4"] = 0
	adaptor.DigitalPinState["5"] = 0
	adaptor.DigitalPinState["6"] = 1

	err = board.StopRelais(context.Background())
	assert.NoError(t, err)
	assert.Equal(t, 1, adaptor.DigitalPinState["1"])
	assert.Equal(t, 1, adaptor.DigitalPinState["2"])
	assert.Equal(t, 1, adaptor.DigitalPinState["3"])
	assert.Equal(t, 1, adaptor.DigitalPinState["4"])
	assert.Equal(t, 1, adaptor.DigitalPinState["5"])
	assert.Equal(t, 0, adaptor.DigitalPinState["6"])

}
