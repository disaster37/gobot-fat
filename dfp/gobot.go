package dfp

// Gobot is the interface to handle I/O
type Gobot interface {
	StartWashingPump()
	StopWashingPump()
	StartBarrelMotor()
	StopBarrelMotor()
	StopMotors()
	Start()
	Stop() error
}
