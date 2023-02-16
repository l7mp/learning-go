package lib

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// findFile walks up the directory tree from `dir` and returns the full path for the first file
// named `name` or an error is no such file was found.
func findFile(dir, name string) (string, error) {
	for {
		// check if file exists in current directory
		path := filepath.Join(dir, name)
		if _, err := os.Stat(path); err == nil {
			return path, nil
		}

		// move up one directory level
		parentDir := filepath.Dir(dir)
		// if we've reached the root directory, return error
		if parentDir == "/" {
			return "", fmt.Errorf("file %s not found in directory tree", name)
		}
		dir = parentDir
	}
}

// findStudentId returns the raw student id given in the argument `id`, or in the STUDENT_ID file
// searched upwards from the current directory, or from the environment variable `STUDENT_ID`, or
// an error is no id was found or the student id is set as the default
func findStudentId(id *string) (string, error) {
	if id != nil && *id != "" {
		return *id, nil
	}

	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	if file, err := findFile(dir, StudentIdFile); err == nil {
		dat, err := os.ReadFile(file)
		if err != nil {
			return "", fmt.Errorf("internal error: could not find student-id file %q returned by findFile",
				StudentIdFile)
		}
		return strings.TrimSpace(string(dat)), nil
	}

	env := os.Getenv(StudentEnvVar)
	if env != "" {
		return strings.TrimSpace(env), nil
	}

	return "", fmt.Errorf("no student id file found searhing up from %q", dir)

}
