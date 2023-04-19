package lib

import (
	"fmt"
	"hash/fnv"
	"log"
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
	// SolutionTemplateFile  is the name of the template file we use to generate the solutions
	SolutionTemplateFile = ".exercise.go"
	// SolutionFile is the name of the generated solution
	SolutionFile = "exercise.go"
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
func Generate(id string, verbose bool) error {
	cwd, err := os.Getwd()
	if err != nil {
		return err
	}

	ex, err := NewExercise(cwd)
	if err != nil {
		return err
	}

	log.Printf("Generating exercise %q in dir %q\n", ex.Name, cwd)

	input := ex.GetInput(id)
	if verbose {
		log.Printf("Using input %#v\n", input)
	}

	if err := GenerateReadme(cwd, input, verbose); err != nil {
		return fmt.Errorf("error generating README in dir %q: %s", cwd, err)
	}

	if err := GenerateTest(cwd, input, verbose); err != nil {
		return fmt.Errorf("error generating tests in dir %q: %s", cwd, err)
	}

	if err := GenerateSolution(cwd, input, verbose); err != nil {
		return fmt.Errorf("error generating solution in dir %q: %s", cwd, err)
	}

	return nil
}

// GenerateReadme generates the README.
func GenerateReadme(dir string, input Input, verbose bool) error {
	if verbose {
		log.Printf("Generating README in dir %q\n", dir)
	}

	templateFile := filepath.Join(dir, ReadmeTemplateFile)
	t, err := template.ParseFiles(templateFile)
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
		return fmt.Errorf("template %q: %w", templateFile, err)
	}

	return nil
}

// GenerateTest generates the test file.
func GenerateTest(dir string, input Input, verbose bool) error {
	if verbose {
		log.Printf("Generating tests in dir %q\n", dir)
	}

	templateFile := filepath.Join(dir, TestTemplateFile)
	t, err := template.ParseFiles(templateFile)
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
		return fmt.Errorf("template %q: %w", templateFile, err)
	}

	return nil
}

// GenerateSolution generates the solution file.
func GenerateSolution(dir string, input Input, verbose bool) error {
	// the solution template exists only in the solution repo, so of this is missing it is not
	// critical
	path := filepath.Join(dir, SolutionTemplateFile)
	if _, err := os.Stat(path); err != nil {
		return nil
	}

	if verbose {
		log.Printf("Generating solution in dir %q\n", dir)
	}

	t, err := template.ParseFiles(path)
	if err != nil {
		return err
	}

	r, err := os.Create(filepath.Join(dir, SolutionFile))
	if err != nil {
		return err
	}
	defer r.Close()

	err = t.Execute(r, input)
	if err != nil {
		return fmt.Errorf("template %q: %w", path, err)
	}

	return nil
}
