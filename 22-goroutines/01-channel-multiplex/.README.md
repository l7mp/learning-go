# Channel Multiplexer

Write a function that multiplexes multiple input channels into a single output channel. The function should have the following signature:

```go
func ChannelMultiplex(ctx context.Context, inputs []chan any) chan any
```

## Requirements:
1. The function should return a new channel that receives values from all input channels.
2. The function should use the provided context for cancellation.
3. When the context is cancelled, the function should close the output channel and stop listening to input channels.
4. The function should not block if no one is receiving from the output channel.
5. The function should not leave any goroutines running after the output channel is closed.

## Example Usage:
```go
ctx, cancel := context.WithCancel(context.Background())
defer cancel()

ch1 := make(chan any)
ch2 := make(chan any)
ch3 := make(chan any)

outputCh := channelMultiplex(ctx, []chan any{ch1, ch2, ch3})

// Now you can receive from outputCh, which will contain values from ch1, ch2, and ch3
```

Insert the code into the file `exercise.go` at the placeholder `// INSERT YOUR CODE HERE`.

