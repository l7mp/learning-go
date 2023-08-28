package resilient

import (
	"context"
)

// WithTimeout receives a Closure as argument and returns the same Closure decorated with a context
// that can be used to control its lifetime.
func WithTimeout(f Closure) ClosureContext {
	return func(ctx context.Context) error {
		cherr := make(chan error)
		defer close(cherr)

		go func() {
			cherr <- f()
		}()

		select {
		case err := <-cherr:
			return err
		case <-ctx.Done():
			return ctx.Err()
		}
	}
}
