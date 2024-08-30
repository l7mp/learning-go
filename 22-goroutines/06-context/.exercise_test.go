package subtask

import (
	"context"
	"testing"
	"time"
)

func TestStartLongRunningTask(t *testing.T) {
	t.Run("Successful completion", func(t *testing.T) {
		ctx := context.Background()
		result, err := StartTask(ctx)
		if err != nil {
			t.Errorf("Expected no error, got %v", err)
		}
		if result != "Main task status: Subtask completed successfully" &&
			result != "Main task status:Subtask completed successfully" {
			t.Errorf("Expected 'Main task status:Subtask completed successfully', got %s", result)
		}
	})

	t.Run("Cancelled by main context", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		go func() {
			time.Sleep(100 * time.Millisecond)
			cancel()
		}()

		result, err := StartTask(ctx)
		if err == nil {
			t.Error("Expected an error, got nil")
		}
		if result != "" {
			t.Errorf("Expected empty result, got %s", result)
		}
	})

	t.Run("Immediate cancellation", func(t *testing.T) {
		ctx, cancel := context.WithCancel(context.Background())
		cancel() // Cancel immediately

		result, err := StartTask(ctx)
		if err == nil {
			t.Error("Expected an error, got nil")
		}
		if result != "" {
			t.Errorf("Expected empty result, got %s", result)
		}
	})
}
