# Describe an Earthquake on the Richter scale

Write a small function that returns the textual description of a given earthquake magnitude value
on the [Richter scale](https://en.wikipedia.org/wiki/Richter_magnitude_scale).

The below table gives the textual descriptions for each Richter scale magnitude:

| MAGNITUDE              | DESCRIPTION |
| ---------------------- | ----------- |
| Less than 2.0          | micro       |
| 2.0 and less than 3.0  | very minor  |
| 3.0 and less than 4.0  | minor       |
| 4.0 and less than 5.0  | light       |
| 5.0 and less than 6.0  | moderate    |
| 6.0 and less than 7.0  | strong      |
| 7.0 and less than 8.0  | major       |
| 8.0 and less than 10.0 | great       |
| 10.0 or more           | massive     |

For instance, the function should return "massive" for a magnitude of "11.0".

Insert your code into the file `exercise.go` near the placeholder `// INSERT YOUR CODE HERE`.

HINT: use the [`switch` statement](https://go.dev/tour/flowcontrol/9).
