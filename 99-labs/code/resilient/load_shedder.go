package resilient

import (
	"fmt"
	"sync"
)

// WithLoadShedding receives a Closure and returns the same Closure decorated with a load-shedding
// pattern that returns an error when the number of active calls surpasses a threashold.
func WithLoadShedding(f Closure, threshold int) Closure {
	counter := 0
	var lock sync.RWMutex // Mutual exclusion lock to protect counter
	var err error
	return func() error {
		lock.RLock()
		c := counter
		lock.RUnlock() // Read counter to 'c'

		if c >= threshold {
			return fmt.Errorf("Too many requests: %d > %d", c, threshold)
		}

		lock.Lock()
		counter++
		lock.Unlock() // Increment req counter
		err = f()
		lock.Lock()
		counter--
		lock.Unlock() // Decrement req counter

		return err
	}
}
