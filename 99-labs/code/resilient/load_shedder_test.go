package resilient

import (
	"sync"
	"testing"
	"time"
)

func costlyFunction() error { time.Sleep(1 * time.Second); return nil }

func TestLoadShedding(t *testing.T) {
	errorCounter := 0
	var lock sync.RWMutex // Mutual exclusion lock to protect counter
	var wg sync.WaitGroup

	loadShedder := WithLoadShedding(costlyFunction, 2)

	wg.Add(5)
	for i := 0; i < 5; i++ { // run 3 goroutines: one of them should fail
		go func() {
			err := loadShedder()
			if err != nil {
				lock.Lock()
				errorCounter++
				lock.Unlock() // Increment req counter
			}
			wg.Done()
		}()
	}

	wg.Wait() // Wait until all threads finish

	if errorCounter != 3 {
		t.Error("expecting 3 suppressed calls")
	}
}
