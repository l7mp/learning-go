package resilient

import (
	"context"
	"errors"
	"sync/atomic"
	"time"
)

// RateLimiter defines the basic config for the token-bucket rate limiter. The RateLimiter admits
// bursts of size Capacity and a steady state rate of MaxRate=Fill/Period.
type RateLimiter struct {
	// Capacity defines the maximum number of tokens in the bucket (burst size).
	Capacity int32
	// Fill defines the number of tokens to add per every period.
	Fill int32
	// In every Period exatcly Fill tokens will be added to the token bucket.
	Period time.Duration
}

// WithRateLimiter receives a Closure and returns the same Closure decorated with a token bucket
// based rate-limiter that rejects each request that comes in excess of the tokens actually in the
// bucket.
func WithRateLimiter(f Closure, ctx context.Context, r RateLimiter) Closure {
	var tokens atomic.Int32
	tokens.Store(r.Capacity)

	ticker := time.NewTicker(r.Period)
	go func() {
		defer ticker.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-ticker.C:
				t := tokens.Load() + r.Fill
				if t > r.Capacity {
					t = r.Capacity
				}
				tokens.Store(t)
			}
		}
	}()

	return func() error {
		if tokens.Load() <= 0 {
			return errors.New("Too many calls")
		}

		tokens.Store(tokens.Load() - 1)
		return f()
	}
}

// func WithRateLimiter(f Closure, ctx context.Context, r RateLimiter) Closure {
// 	tokens := r.Capacity
// 	ticker := time.NewTicker(r.Period)
// 	go func() {
// 		defer ticker.Stop()
// 		for {
// 			select {
// 			case <-ctx.Done():
// 				return
// 			case <-ticker.C:
// 				t := tokens + r.Fill
// 				if t > r.Capacity {
// 					t = r.Capacity
// 				}
// 				tokens = t
// 			}
// 		}
// 	}()

// 	return func() error {
// 		if tokens <= 0 {
// 			return errors.New("Too many calls")
// 		}

// 		tokens--
// 		return f()
// 	}
// }
