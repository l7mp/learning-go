# Learning Go: Basic Go exercises

A series of basic Go exercises inspired from the [Go: Bootcamp Course](https://github.com/inancgumus/learngo). Solving all exercises will not make you a Go ninja, but it should be enough to get you started.

## Getting started

Change to the root of the git repo and generate the exercises:

``` console
echo <MY-STUDENT-ID> > STUDENT_ID
git add STUDENT_ID
git commit -m 'student id added'
make generate
```

Your student id should always be available in the file named `STUDENT_ID` in the main directory. You can override this by setting the id in the `STUDENT_ID` environment variable.

``` console
STUDENT_ID=<MY-STUDENT-ID> make generate
```

## Solve the exercises

### Write code

Go through the subdirectories, preferably in the same order as the file names, understand the exercise specified in the `README.md` file, and insert your solution into `exercise.go` near the placeholder.

Once done with the exercise in the directory `<exercise-directory>`, make sure to git-add and git-commit your solution.

``` console
git add <exercise-directory>/exercise.go
git commit -m 'solved <exercise-directory>'
```

You can divide your code into as many files as you like but don't forget to add each file to your git repo.

### Test

At any point in time you can test your solutions as follows.

``` console
make test
```

### Keep track of repo updates

Sometimes we update the main git repo to fix bugs or add new exercises. The below workflow shows how to update your local working copy from the master without overwriting your solutions already written.

> **Warning**
> Make sure to commit all your code into git: this will guarantee that you will never lose your solutions even if some of the below steps go wrong. You can also back up your git repo, but please make sure your solutions are kept private (a private GitHUb repo will do it).

#### With all local changes committed to git

If you have a clear git tree with all your changes committed, the below should update only the files that change in the master.

``` console
git pull --rebase
```

#### With local changes in the working directory

If you have uncommitted changes, follow the below steps.

1. Store your changes temporarily.

   ``` console
   git stash
   ```

2. Pull updates.

   ``` console
   git pull
   ```

3. Restore your changes.

   ``` console
   git stash
   ```

4. Optional: fix conflicts and re-run tests.

### Add a new exercise

Add a new subdirectory and add the following files:

- `exercise.yaml`: An exercise definition with a set of inputs, from which `make generate` will choose one by hashing on the student id to generate the exercise.
- `exercise.go`: Placeholder for the solution.
- `.README.md`: a README template with instructions.
- `.exercise_test.go`: the test file to check your solutions.

If you add a new top-level directory, don't forget to include it in the `EXERCISE_DIRS` in the Makefile.

Then run `make clean`, this will add the placeholders for the exercise (these will be overwritten by `make generate`), add all files in the exercise dir to the git repo, and git-push.

## Complete the labs

Labs tasks are located in the [99-labs](99-labs/) folder. The labs give you hands-on development and deployment experience. You will learn how to build and containerize Go programs as well as how to run them in Kubernetes. Each lab contains a README that gives you context and specifies the lab exercises. The labs depend on each other, so it is recommended to complete them one after the other.

## Clean up

Reset all generated files to the default placeholder. Warning: this will drop all your uncommitted
local changes, use it at your own risk.

``` console
make clean
```

Also reset the student id: this is required before pushing any modification to any of the exercises.

``` console
make realclean
```

## Help

Ask any questions or file an issue at the original [GitHub repo](https://github.com/l7mp/learning-go).

## License

Copyright 2021-2023 by its authors. Some rights reserved. See [AUTHORS](AUTHORS).

[Creative Commons Attribution-NonCommercial-ShareAlike](https://creativecommons.org/licenses/by-nc-sa/4.0/) - see [LICENSE](LICENSE) for full text.

## Acknowledgments

Examples adopted from the [Go: Bootcamp Course](https://github.com/inancgumus/learngo).

<!-- Local Variables: -->
<!-- mode: markdown; coding: utf-8 -->
<!-- auto-fill-mode: nil -->
<!-- visual-line-mode: 1 -->
<!-- markdown-enable-math: t -->
<!-- End: -->
