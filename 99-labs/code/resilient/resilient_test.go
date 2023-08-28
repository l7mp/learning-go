package resilient

import (
	"context"
	"testing"
	"time"
)

var SlowRetry = Backoff{
	Base:      time.Second,
	Cap:       time.Minute,
	Jitter:    1,
	NumTrials: 3,
}

func TestTimoutRetrySucceed(t *testing.T) {
	// should succeed with fast retries
	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)

	f := getRetriableFunction(3)
	retrier := WithTimeoutRetry(f, FastRetry)
	err := retrier(ctx)

	if err != nil {
		t.Error("expecting nil error")
	}
}

func TestTimoutRetryFail(t *testing.T) {
	// should fail with slow retries
	ctx, _ := context.WithTimeout(context.Background(), 1*time.Second)

	f := getRetriableFunction(3)
	retrier := WithTimeoutRetry(f, SlowRetry)
	err := retrier(ctx)

	if err == nil {
		t.Error("expecting non-nil error")
	}
}
