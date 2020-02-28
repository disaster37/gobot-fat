package models

/*
// AverageDurationSecond compute average duration
func (h *DFP) AverageDurationSecond() uint64 {
	if len(h.WashingHistory) > 0 {
		totalDuration := uint64(0)
		for i := 0; i < len(h.WashingHistory); i++ {
			if (i + 1) < len(h.WashingHistory) {
				totalDuration = totalDuration + uint64(h.WashingHistory[i+1].Sub(h.WashingHistory[i]).Seconds())
			} else {
				totalDuration = totalDuration + uint64(h.WashingHistory[0].Sub(h.WashingHistory[i]).Seconds())
			}
		}

		return (totalDuration / uint64(len(h.WashingHistory)))
	}

	return 0

}
*/
