package resilient

import (
	"time"
)

// WithDebounceFirst receives a Closure and returns the same Closure decorated with a
// function-first debouncer that suppresses calls within the time interval d.
func WithDebounceFirst(f Closure, d time.Duration) Closure {
	var threshold time.Time
	var err error
	return func() error {
		defer func() {
			threshold = time.Now().Add(d)
		}()

		if time.Now().Before(threshold) {
			return err
		}

		err = f()

		return err
	}
}
