package channelmultiplexer

import (
	"context"
	"fmt"
	"sort"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestChannelMultiplex(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	ch1 := make(chan any)
	ch2 := make(chan any)
	ch3 := make(chan any)

	outputCh := channelMultiplex(ctx, []chan any{ch1, ch2, ch3})

	// Send values to input channels
	go func() {
		ch1 <- 1
		ch2 <- "two"
		ch3 <- 3.0
		ch1 <- 4
		ch2 <- "five"
		ch3 <- 6.0
		close(ch1)
		close(ch2)
		close(ch3)
	}()

	// Collect results
	var results []any
	for v := range outputCh {
		results = append(results, v)
	}

	// Sort results for consistent comparison
	sort.Slice(results, func(i, j int) bool {
		return fmt.Sprintf("%v", results[i]) < fmt.Sprintf("%v", results[j])
	})

	// Check if we received all values
	expected := []any{1, "two", 3.0, 4, "five", 6.0}
	sort.Slice(expected, func(i, j int) bool {
		return fmt.Sprintf("%v", expected[i]) < fmt.Sprintf("%v", expected[j])
	})

	assert.Len(t, results, len(expected))
	assert.Equal(t, expected, results)
}

func TestChannelMultiplexCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	ch1 := make(chan any)
	ch2 := make(chan any)

	outputCh := channelMultiplex(ctx, []chan any{ch1, ch2})

	// Send one value
	ch1 <- "test"

	// Receive one value
	select {
	case v := <-outputCh:
		if v != "test" {
			t.Errorf("Expected 'test', got %v", v)
		}
	case <-time.After(time.Second):
		t.Error("Timed out waiting for value")
	}

	// Cancel the context
	cancel()

	// The output channel should be closed
	select {
	case _, ok := <-outputCh:
		if ok {
			t.Error("Output channel should be closed after cancellation")
		}
	case <-time.After(time.Second):
		t.Error("Timed out waiting for channel to close")
	}
}
