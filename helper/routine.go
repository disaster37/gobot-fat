package helper

import "time"

// Every triggers f every t time.Duration until the end of days, or when a Stop()
// is called on the Ticker that is returned by the Every function.
// It does not wait for the previous execution of f to finish before
// it fires the next f.
func Every(t time.Duration, f func()) *time.Ticker {
	ticker := time.NewTicker(t)

	go func() {
		for {
			select {
			case <-ticker.C:
				f()
			}
		}
	}()

	return ticker
}
