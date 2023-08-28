package resilient

import (
	"errors"
	"time"
)

var errUnavailable = errors.New("service unreachable")

type breakerState int

const (
	stateOpened breakerState = iota
	stateClosed
)

// Breaker defines the circuit breaker parameters.
type Breaker struct {
	// Number of consecutive failures after which the circuit is opened.
	FailureThreshold int
	// Time interval after which the circuit is spontaneously closed.
	CloseInterval time.Duration
}

var DefaultCircuitBreaker = Breaker{
	FailureThreshold: 3,
	CloseInterval:    time.Second,
}

// WithCircuitBreaker takes a Closure and returns a function that, when run, opens the circuit after the failureThreshold errors and spontaneously closes an open open circuit after closeInterval.
func WithCircuitBreaker(f Closure, breaker Breaker) Closure {
	consecutiveFailures := 0
	lastAttempt := time.Now()
	state := stateClosed

	return func() error {
		if state == stateOpened && time.Since(lastAttempt) > breaker.CloseInterval {
			state = stateClosed
			consecutiveFailures = 0
		}

		if state == stateOpened {
			return errUnavailable
		}

		lastAttempt = time.Now()
		err := f()
		if err != nil {
			consecutiveFailures++
			if consecutiveFailures >= breaker.FailureThreshold {
				state = stateOpened
			}
		} else {
			consecutiveFailures = 0
		}

		return err
	}
}
