package resilient

import (
	"context"
	"testing"
	"time"
)

func TestRateLimiter(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	rateLimiter := WithRateLimiter(func() error { return nil }, ctx,
		RateLimiter{Capacity: 4, Fill: 1, Period: 250 * time.Millisecond})

	errorCounter := 0
	countSixErrors := func() {
		for i := 0; i < 6; i++ {
			if err := rateLimiter(); err != nil {
				errorCounter++
			}
		}
	}

	// 6 calls in quick succession, only 4 should succeed
	countSixErrors()

	// Let the bucket refill
	time.Sleep(1100 * time.Millisecond)

	// Another 6 calls in quick succession, only 4 should succeed
	countSixErrors()

	if errorCounter != 4 {
		t.Error("expecting 4 suppressed calls")
	}
}
