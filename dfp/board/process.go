package dfpboard

func (h *DFPBoard) wash() {

}

func (h *DFPBoard) updateState() {

	h.mutexState.Lock()
	defer h.mutexState.Unlock()

}

func (h *DFPBoard) sendEvent(kind string, name string, args ...interface{}) {

}
