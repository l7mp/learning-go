package resilient

import (
	"fmt"
	"testing"
	"time"
)

var FastRetry = Backoff{
	Base:      10 * time.Millisecond,
	Cap:       10000 * time.Millisecond,
	Jitter:    1,
	NumTrials: 10,
}

func getRetriableFunction(times int) Closure {
	return func() error {
		if times > 1 {
			times -= 1
			return fmt.Errorf("failed, remaining fail budget: %d", times)
		}
		return nil
	}
}

func TestRetry(t *testing.T) {
	for retry := 1; retry <= 4; retry += 1 {
		backoff := FastRetry
		backoff.NumTrials = retry
		for fail := 1; fail <= 4; fail += 1 {
			f := getRetriableFunction(fail)
			retrier := WithRetry(f, backoff)
			err := retrier()

			if fail <= retry {
				if err != nil {
					t.Error("expecting nil error")
				}
			} else {
				if err == nil {
					t.Error("expecting non-nil error")
				}
				break
			}
		}
	}
}
