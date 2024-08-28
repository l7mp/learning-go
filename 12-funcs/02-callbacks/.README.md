# Calculator

In this exercise, you'll implement a simple calculator that uses enum-like constants for operation names and callbacks to perform arithmetic operations. Your task is to implement the necessary types and functions to make the calculator work.

Requirements:

1. Define constants for operation types using iota:
   ```go
   const (
       Add OperationType = iota
       Subtract
       // Add more operations here
   )
   ```

2. Implement a function to get the string representation of an OperationType:
   ```go
   func (op OperationType) String() string
   ```

3. Implement 3 operations:
   - `Add`: Returns the sum of two numbers
   - `Subtract`: Returns the difference between two numbers
   - `Multiply`: Returns the product of two numbers

5. Implement a Calculator function with the following signature:
   ```go
   func Calculate(op OperationType, a, b float64) float64
   ```
   This function should take an operation type and two float64 numbers and return the result of
   applying the corresponding operation to the numbers.

Example usage:
```go
result := Calculate(Add, 10, 5)
// result should be 15
```

Insert your code into the file `exercise.go` at the placeholder `// INSERT YOUR CODE HERE`.
