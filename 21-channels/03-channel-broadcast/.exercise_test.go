package channelbroadcaster

import (
	"context"
	"testing"
	"time"
)

func TestChannelBroadcast(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	input := make(chan any)
	output1 := make(chan any, 3) // Buffered channel
	output2 := make(chan any, 3) // Buffered channel

	channelBroadcast(ctx, input, []chan<- any{output1, output2})

	// Send values to input channel
	go func() {
		input <- 1
		input <- "two"
		input <- 3.0
		close(input)
	}()

	// Wait for a short time to ensure broadcasting is complete
	time.Sleep(time.Millisecond * 100)

	// Check if both output channels received all values
	expected := []any{1, "two", 3.0}
	outputs := []<-chan any{output1, output2}

	for i, out := range outputs {
		var received []any
		for j := 0; j < len(expected); j++ {
			select {
			case v := <-out:
				received = append(received, v)
			case <-time.After(time.Millisecond * 100):
				t.Errorf("Timeout waiting for value from output channel %d", i+1)
			}
		}

		if len(received) != len(expected) {
			t.Errorf("Output channel %d: expected %d values, got %d", i+1, len(expected), len(received))
		}
		for j, v := range expected {
			if received[j] != v {
				t.Errorf("Output channel %d: expected %v at index %d, got %v", i+1, v, j, received[j])
			}
		}
	}

	// Check if the output channels are closed
	for i, out := range outputs {
		select {
		case _, ok := <-out:
			if ok {
				t.Errorf("Output channel %d should be closed", i+1)
			}
		case <-time.After(time.Millisecond * 100):
			t.Errorf("Timeout waiting for output channel %d to close", i+1)
		}
	}
}

func TestChannelBroadcastCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	input := make(chan any)
	output1 := make(chan any, 1)
	output2 := make(chan any, 1)

	channelBroadcast(ctx, input, []chan<- any{output1, output2})

	// Send one value
	input <- "test"

	// Wait a bit to ensure the value is broadcasted
	time.Sleep(time.Millisecond * 50)

	// Cancel the context
	cancel()

	// Wait a bit for cancellation to take effect
	time.Sleep(time.Millisecond * 50)

	// Consume content
	for i, ch := range []<-chan any{output1, output2} {
		select {
		case v, ok := <-ch:
			if !ok {
				t.Errorf("Output channel %d should be opened before drain", i+1)
			}
			if v != "test" {
				t.Errorf("Expected different input on output channel %d", i+1)
			}
		default:
			t.Errorf("Output channel %d should be opened before drain", i+1)
		}
	}

	// Check if the output channels are closed
	for i, ch := range []<-chan any{output1, output2} {
		select {
		case _, ok := <-ch:
			if ok {
				t.Errorf("Output channel %d should be closed after cancellation", i+1)
			}
		default:
			t.Errorf("Output channel %d should be closed after cancellation", i+1)
		}
	}
}
