package resilient

import (
	"testing"
	"time"
)

func TestDebounce(t *testing.T) {
	counter := 0
	counterFunc := func() error {
		counter++
		return nil
	}

	debounce := WithDebounceFirst(counterFunc, 200*time.Millisecond)
	// 3 quick calls
	debounce()
	debounce()
	debounce()

	time.Sleep(200 * time.Millisecond)
	debounce()

	if counter != 2 {
		t.Error("expecting only 2 calls")
	}
}
