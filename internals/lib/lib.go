package lib

import (
	"fmt"
	"hash/fnv"
	"os"
	"path/filepath"
	"text/template"
)

const (
	// StudentIdFile is the name of the file that holds the student id
	StudentIdFile = "STUDENT_ID"
	// StudentEnvVar is the name of the environment variable that holds the student id
	StudentEnvVar = "STUDENT_ID"
	// DefaultStudentId is what we store in the student id file by default
	DefaultStudentId = "PLEASE SET STUDENT ID"
	// ReadmeTemplateFile is the name of the template file we use to generate the README
	ReadmeTemplateFile = ".README.md"
	// ReadmeFile is the name README file
	ReadmeFile = "README.md"
	// TestTemplateFile  is the name of the template file we use to generate the tests
	TestTemplateFile = ".exercise_test.go"
	// TestFile is the name of the generated test
	TestFile = "exercise_test.go"
)

// Input is a type alias to the input field of an exercise.
type Input = map[string]any

// GetStudentId returns the student id given in the argument `id`, or in the STUDENT_ID file
// searched upwards from the current directory, or from the environment variable `STUDENT_ID`, or
// an error is no id was found
func GetStudentId(id *string) (string, error) {
	student, err := findStudentId(id)
	if err != nil {
		return "", err
	}

	if student == DefaultStudentId {
		return "", fmt.Errorf("%s", DefaultStudentId)
	}

	return student, nil
}

// GetStudentHash creates a hash from a student id.
func GetStudentHash(id string) int {
	h := fnv.New32a()
	h.Write([]byte(id))
	return int(h.Sum32())
}

// Generate generates the README and the test file in the given dir.
func Generate(id string) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	ex, err := NewExercise(cwd)
	if err != nil {
		return err
	}

	input := ex.GetInput(id)

	if err := GenerateReadme(cwd, input); err != nil {
		return err
	}

	if err := GenerateTest(cwd, input); err != nil {
		return err
	}

	return nil
}

// GenerateReadme generates the README.
func GenerateReadme(dir string, input Input) error {
	t, err := template.ParseFiles(filepath.Join(dir, ReadmeTemplateFile))
	if err != nil {
		return err
	}

	r, err := os.Create(filepath.Join(dir, ReadmeFile))
	if err != nil {
		return err
	}
	defer r.Close()

	err = t.Execute(r, input)
	if err != nil {
		return err
	}

	return nil
}

// GenerateTest generates the test file.
func GenerateTest(dir string, input Input) error {
	t, err := template.ParseFiles(filepath.Join(dir, TestTemplateFile))
	if err != nil {
		return err
	}

	r, err := os.Create(filepath.Join(dir, TestFile))
	if err != nil {
		return err
	}
	defer r.Close()

	err = t.Execute(r, input)
	if err != nil {
		return err
	}

	return nil
}
