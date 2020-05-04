package tfp

// Gobot is the interface to handle I/O
type Gobot interface {
	StartPondPump() error
	StopPondPump() error
	StartWaterfallPump() error
	StopWaterfallPump() error
	StartUVC1() error
	StopUVC1() error
	StartUVC2() error
	StopUVC2() error
	StartPondBubble() error
	StopPondBubble() error
	StartFilterBubble() error
	StopFilterBubble() error
	StopRelais() error
	Start() error
	Stop() error
	Reconnect() error
}
