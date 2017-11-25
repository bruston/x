// Package try provides simple, throttled retrying.
package try

import "time"

// A Throttler regulates the rate of retry attempts.
type Throttler interface {
	Throttle()
}

// Do runs a function a maximum of n times, retrying on failure. If a non-nil
// Throttler is passed, it is called after each failed try. A failed try occurs
// when the given function returns an error. Do returns the last error if no tries
// are successful, or nil on success.
func Do(n int, fn func() error, throt Throttler) error {
	if fn == nil {
		return nil
	}
	var lastErr error
	for i := 0; i < n; i++ {
		if err := fn(); err != nil {
			lastErr = err
			if throt != nil {
				throt.Throttle()
			}
			continue
		}
		return nil
	}
	return lastErr
}

// Delay is a Throttler that pauses for a fixed duration.
type Delay time.Duration

// Throttle pauses for a fixed duration.
func (d Delay) Throttle() {
	time.Sleep(time.Duration(d))
}

// Backoff is a Throttler that pauses for an exponentionally increasing duration.
type Backoff struct {
	Duration   time.Duration
	Multiplier int
}

// Throttle pauses for an increasing duration.
func (b *Backoff) Throttle() {
	time.Sleep(b.Duration)
	b.Duration += b.Duration * time.Duration(b.Multiplier)
}
