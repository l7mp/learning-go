// Package resilient is a generic implementation of some basic resilience patterns, like timeout, retry, circuit breaking, etc.
package resilient

import "context"

// Closure is the type of the function we want to decorate with a resilience policy.
type Closure func() error

// ClosureContext is a Closure that also accepts a Context.
type ClosureContext func(context.Context) error

// WithTimeoutRetry returns a closure decorated with a context that will do a configurable number
// of retries but times out according to the passed in context.
func WithTimeoutRetry(f Closure, wait Backoff) ClosureContext {
	return WithTimeout(WithRetry(f, wait))
}
