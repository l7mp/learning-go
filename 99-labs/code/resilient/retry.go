package resilient

import (
	"math/rand"
	"time"
)

// Backoff allows to configure the parameters used for the random exponential backoff algorithm.
type Backoff struct {
	// The initial duration.
	Base time.Duration
	// An upper limit on the delay between retries.
	Cap time.Duration
	// The sleep time is the base plus an additional random jitter.
	Jitter float64
	// NumTrials is the number of times the function is run.
	NumTrials int
}

var DefaultBackoff = Backoff{
	Base:      250 * time.Millisecond,
	Cap:       5 * time.Second,
	Jitter:    3,
	NumTrials: 4,
}

// WithRetry takes a Closure and a set of backoff parameters and returns a function that, when run,
// retries the Closure using random exponential backoff on failure.
func WithRetry(f Closure, wait Backoff) Closure {
	return func() error {
		err := f()
		base, cap := wait.Base, wait.Cap
		for backoff, step := base, wait.NumTrials-1; err != nil && step > 0; backoff, step = backoff<<1, step-1 {
			if backoff > cap {
				backoff = cap
			}
			jitter := rand.Float64() * float64(backoff) * wait.Jitter
			sleep := base + time.Duration(jitter)
			time.Sleep(sleep)
			err = f()
		}
		return err
	}
}
