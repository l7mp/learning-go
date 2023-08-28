package resilient

import (
	"errors"
	"testing"
	"time"
)

var errFailed error = errors.New("failed")

var testBreaker = Breaker{
	FailureThreshold: 3,
	CloseInterval:    100 * time.Millisecond,
}

func errFunc() error {
	return errFailed
}

func TestCircuitBreaker(t *testing.T) {
	breaker := WithCircuitBreaker(errFunc, testBreaker)

	// fail with "failed"
	for i := 0; i < 3; i++ {
		err := breaker()
		if err == nil || !errors.Is(err, errFailed) {
			t.Error("expecting errFailed")
		}
	}

	// fail with "service unavailable"
	err := breaker()
	if err == nil || !errors.Is(err, errUnavailable) {
		t.Error("expecting errUnavailable")
	}

	// wait until the breaker reopens
	time.Sleep(2 * testBreaker.CloseInterval)

	err = breaker()
	if err == nil || !errors.Is(err, errFailed) {
		t.Error("expecting errFailed")
	}
}
