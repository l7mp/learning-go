package threadpool

import (
	"context"
	"errors"
	"sync"
	"testing"
	"time"
)

var mu sync.Mutex
var counter = 0

type mockTask struct {
	delay time.Duration
	err   error
}

func (m mockTask) Run(ctx context.Context) error {
	mu.Lock()
	counter++
	mu.Unlock()

	select {
	case <-time.After(m.delay):
		return m.err
	case <-ctx.Done():
		return ctx.Err()
	}
}

func TestThreadPool(t *testing.T) {
	t.Run("creates correct number of workers", func(t *testing.T) {
		counter = 0
		pool, _ := NewThreadPool(5)

		for i := 0; i < 10; i++ {
			pool.Run(mockTask{
				delay: 50 * time.Millisecond,
				err:   nil,
			})
		}

		time.Sleep(10 * time.Millisecond)

		mu.Lock()
		if counter != 5 {
			t.Errorf("Expected at most 5 concurrent tasks, got %d", counter)
		}
		mu.Unlock()

		time.Sleep(100 * time.Millisecond)

		pool.Close()
	})

	t.Run("reports errors correctly", func(t *testing.T) {
		pool, errChan := NewThreadPool(1)

		expectedErr := errors.New("test error")
		pool.Run(mockTask{
			delay: 10 * time.Millisecond,
			err:   expectedErr,
		})

		select {
		case err := <-errChan:
			if err != expectedErr {
				t.Errorf("Expected error %v, got %v", expectedErr, err)
			}
		case <-time.After(20 * time.Millisecond):
			t.Error("Timeout waiting for error")
		}

		pool.Close()
	})

	t.Run("waits for threads to finish when closing", func(t *testing.T) {
		pool, errChan := NewThreadPool(1)

		expectedErr := errors.New("test error")
		pool.Run(mockTask{
			delay: 50 * time.Millisecond,
			err:   expectedErr,
		})

		// make sure task has enough time to start
		time.Sleep(10 * time.Millisecond)
		pool.Close()

		select {
		case err := <-errChan:
			if err != expectedErr && err != context.Canceled {
				t.Errorf("Received unexpected error %v", err)
			}
		case <-time.After(20 * time.Millisecond):
			t.Error("Timeout waiting for error")
		}
	})

	t.Run("closes gracefully", func(t *testing.T) {
		pool, errChan := NewThreadPool(1)

		done := make(chan struct{})
		go func() {
			pool.Run(mockTask{
				delay: 100 * time.Millisecond,
				err:   nil,
			})
			close(done)
		}()

		time.Sleep(10 * time.Millisecond)
		pool.Close()

		select {
		case <-done:
			// Task was cancelled successfully
		case <-time.After(50 * time.Millisecond):
			t.Error("Close did not cancel running task")
		}

		time.Sleep(20 * time.Millisecond)

		err, ok := <-errChan
		if !ok || err == nil {
			t.Error("Error channel should contain a context cancelled error")
		}

		_, ok = <-errChan
		if ok {
			t.Error("Error channel was not closed")
		}
	})
}
