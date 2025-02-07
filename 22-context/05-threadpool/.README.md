# Threadpools

In this exercise, you'll implement a simple threadpool in Go. The threadpool should be able to run tasks concurrently, manage errors, and handle graceful shutdown.

Implement a `ThreadPool` type with the following interface:

```go
type Runnable interface {
    Run(context.Context) error
}

type ThreadPool interface {
    Run(Runnable)
    Close()
}

func NewThreadPool(n int) (ThreadPool, chan error)
```

Requirements:

1. `NewThreadPool(n int)` should create a new threadpool with `n` worker goroutines.
2. The returned `chan error` should be used to watch for errors from the tasks.
3. `Run(Runnable)` should submit a task to the threadpool (note: it is non-blocking, which means it does not wait for the task to finish).
4. `Close()` should immediately stop all running threads and close the error channel.
5. At most `n` tasks should run simultaneously.
6. If a task returns an error, it should be sent on the error channel.

Hints:

- Use a `sync.WaitGroup` to keep track of running threads and ensure they all finish when closing the threadpool (and before closing the error channel).
- Use `context.Context` to handle cancellation of running tasks when `Close()` is called.
- Do not export your threadpool implementation, but make sure to implement the exported ThreadPool interface.
- Use `sync.Once` in the `Close()` to ensure it's only called once, preventing potential race conditions or panics from multiple closes.
- In the `Run()` method, use a select statement to either submit the task or detect if the pool has been closed.
- Make sure the error channel is buffered in order to prevent blocking when sending errors. If the channel becomes full, errors should be logged to the standard output instead.
  
Warning: You may find this exercise fairly difficult to get right. Make sure you understand the semantics of channels and contexts and, if something is not working, insert printfs at strategic places in the code and the tests for debugging.

Insert the code into the file `exercise.go` at the placeholder `// INSERT YOUR CODE HERE`.
