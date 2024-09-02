# Nested Goroutines

In this exercise, you'll practice using contexts to manage nested goroutines in Go. Your task is to
implement two functions: one that starts a long-running operation in a separate goroutine, and
another that performs the actual long-running operation. You'll use contexts to manage the lifetime
of these goroutines and implement a cancellation mechanism that allows the caller to stop both
goroutines.

The two function signatures:

```go
func StartTask(ctx context.Context) (result string, err error)
func SubTask(ctx context.Context) (result string, err error)
```

Requirements:

1. `StartTask` should:
   - Create a new context with a 1-second timeout, derived from the input context.
   - Start `SubTask` in a new goroutine using the derived context.
   - If the main context is canceled, cancel the subtask and return with an empty string and the
     error provided by the main context (use [`ctx.Err()`](https://pkg.go.dev/context#Context)).
   - If `SubTask` finishes with a non-empty error, return that error with an empty string.
   - Otherwise return whatever string `SubTask` returns prepended with the string `"Main task
     status:"` and a `nil` error.

2. `SubTask` should:
   - Simulate a long-running task by attempting to run for 200 milliseconds.
   - If the provided context is canceled before the task completes, return immediately with the
     error provided by the context (use [`ctx.Err()`](https://pkg.go.dev/context#Context)).
   - Otherwise, return the string `"Subtask completed successfully"` with an empty error.

3. Use channels to communicate the subtask result and errors between the goroutines. Make sure to
   close the channels in the goroutine that actually writes the channels otherwise you will see
   ugly race conditions and random panics.

4. Ensure proper resource cleanup by using `defer` statements where appropriate.

Insert the code into the file `exercise.go` at the placeholder `// INSERT YOUR CODE HERE`.

