package tfp

// Gobot is the interface to handle I/O
type Gobot interface {
	StartPondPump() error
	StopPondPump()
	StartWaterfallPump() error
	StopWaterfallPump()
	StartUVC1() error
	StopUVC1()
	StartUVC2() error
	StopUVC2()
	StartPondBubble() error
	StopPondBubble()
	StartFilterBubble() error
	StopFilterBubble()
	StopRelais()
	Start()
	Stop() error
}
