package pbf

// Process action when security is enabled
func (h *FATHandler) doSecurity() {
	h.state.IsSecurity = true

	// stop barrel motor
	h.StopBarrelMotor()

	// stop washing pump
	h.StopWashingPump()
}

// Process action when security is disabled
func (h *FATHandler) undoSecurity() {
	h.state.IsSecurity = false
}
