package resilient

import (
	"context"
	// "fmt"
	"testing"
	"time"
)

func slowFunction() error { time.Sleep(40 * time.Second); return nil }
func fastFunction() error { return nil }

func TestTimeoutTimedOut(t *testing.T) {
	bground := context.Background()
	ctx, _ := context.WithTimeout(bground, 1*time.Second)

	timeout := WithTimeout(slowFunction)
	now := time.Now()
	err := timeout(ctx)
	interval := time.Since(now)

	// times out: should return "context deadline exceeded" error
	if err == nil {
		t.Error("expecting non-nil error on timeout")
	}

	if interval < 950*time.Millisecond || interval < 105*time.Millisecond {
		t.Error("expecting timeout in 1sec")
	}
}

func TestTimeoutNormalReturn(t *testing.T) {
	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)

	timeout := WithTimeout(fastFunction)
	now := time.Now()
	err := timeout(ctx)
	interval := time.Since(now)

	// normal return
	if err != nil {
		t.Error("expecting nil error on normal return")
	}

	if interval > 100*time.Millisecond {
		t.Error("expecting immediate return")
	}
}
