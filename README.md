# Learning Go: Basic Go exercises

A series of basic Go exercises inspired from the [Go: Bootcamp
Course](https://github.com/inancgumus/learngo). Solving all exercises will not make you a Go ninja
but it is enough to get the basics and start to learn the advanced stuff on your own.

## Getting started

Change to the root of the git repo and generate the exercises:

``` console
echo <MY-STUDENT-ID> > STUDENT_ID
make generate
```

Your student id should always be available in the file named `STUDENT_ID` in the main
directory. You can override this by setting the id in the `STUDENT_ID` environment variable.

``` console
STUDENT_ID=MY-STUDENT-ID> make generate
```

## Write code

Go through the subdirectories, preferably in the same order as the file names, understand the
exercise specified in the README, and insert your solution into `exercise.go` near the
pleaceholder. You can divide your code to as many files as you want but don't forget to add each
file to your git repo.

## Test

At any point in time you can test your solutions as follows.

``` console
make test
```

## Add a new exercise

Add a new subdirectory and add the following files:
- `exercise.yaml`: An exercise definition with a set of inputs, from which `make generate` will
  choose one by hashing on the student id to generate the exercise.
- `exercise.go`: Placeholder for the solution.
- `.README.md`: a README template with instructions.
- `.exercise_test.go`: the test file to check your solutions.

If you add a new top-level directory, don't forget to include it in the `EXERCISE_DIRS` in the
Makefile.

Then run `make clean`, this will add the placeholders for the exercise (these will be overwritten
by `make generate`), add all files in the exercise dir to the git repo, and git-push.

## Clean up

Reset all generated files to the default placeholder.

``` console
make clean
```

Also reset the student id: this is required before pushing any modification to the exercises.

``` console
make realclean
```

## Help

Ask any questions or file an issue at the original [GitHub repo](https://github.com/l7mp/learning-go).

## License

Copyright 2021-2023 by its authors. Some rights reserved. See [AUTHORS](AUTHORS).

[Creative Commons Attribution-NonCommercial-ShareAlike](https://creativecommons.org/licenses/by-nc-sa/4.0/) - see [LICENSE](LICENSE) for full text.

## Acknowledgments

Examples adopted from the [[Go: Bootcamp Course](https://github.com/inancgumus/learngo).
