# Channel Broadcaster

Write a function that broadcasts the content of an input channel to multiple output channels. The function should have the following signature:

```go
func channelBroadcast(ctx context.Context, input <-chan any, outputs []chan<- any)
```

## Requirements
1. The function should start its own goroutine to handle the broadcasting.
2. The function should send each value from the input channel to all output channels.
3. The function should use the provided context for cancellation.
4. When the context is cancelled or the input channel is closed, the function should stop broadcasting and close all output channels.
5. The input channel is unbuffered, whereas output channels are buffered. You do not have to handle the case when an output channel is not immediately writable, so it's acceptable to potentially miss sends if an output channel is full.
6. The function should not leave any goroutines running after it completes.

## Example Usage
```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

input := make(chan any)
output1 := make(chan any, 10) // Buffered channel
output2 := make(chan any, 10) // Buffered channel

channelBroadcast(ctx, input, []chan<- any{output1, output2})

// Use the channels...

close(input) // This will cause the broadcaster to close all output channels
```

## Notes
- The function should be concurrent-safe.
- Error handling should be considered, especially around channel operations.
- Pay attention to potential race conditions and deadlocks.

Insert the code into the file `exercise.go` at the placeholder `// INSERT YOUR CODE HERE`.
