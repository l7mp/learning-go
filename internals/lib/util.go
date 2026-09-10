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

func findStudentId(id *string) (string, error) {
	if id != nil && *id != "" {
		return strings.ToUpper(*id), nil
	}

	// env overrides stundent-id file
	env := os.Getenv(StudentEnvVar)
	if env != "" {
		return strings.ToUpper(strings.TrimSpace(env)), nil
	}

	dir, err := os.Getwd()
	if err != nil {
		return "", err
	}

	// The override file wins, and a file still holding the placeholder counts as unset,
	// so a repository that carries an override is not held back by the STUDENT_ID file
	// that every repository tracks.
	for _, name := range []string{OverrideStudentIdFile, StudentIdFile} {
		file, err := findFile(dir, name)
		if err != nil {
			continue
		}

		dat, err := os.ReadFile(file)
		if err != nil {
			return "", fmt.Errorf("could not read student-id file %q: %w", file, err)
		}

		student := strings.ToUpper(strings.TrimSpace(string(dat)))
		if student == "" || student == DefaultStudentId {
			continue
		}

		return student, nil
	}

	return "", fmt.Errorf("no student id: searched for %s and %s upwards from %q, "+
		"put your id in %s and commit it", OverrideStudentIdFile, StudentIdFile, dir,
		StudentIdFile)
}
