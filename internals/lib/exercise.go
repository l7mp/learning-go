package lib

import (
	"fmt"
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v3"
)

var _ = fmt.Sprintf("%d", 1)

const (
	ExerciseFileName = "exercise.yaml"
)

type Exercise struct {
	Name  string    `yaml:"name"`
	Input [](Input) `yaml:"input"`
}

// NewExercise reads the exercise definition in the given directory
func NewExercise(dir string) (*Exercise, error) {
	yml := filepath.Join(dir, ExerciseFileName)
	data, err := os.ReadFile(yml)
	if err != nil {
		return nil, fmt.Errorf("cannot read YAML def for %q: %w", yml, err)
	}

	ex := Exercise{}
	if err := yaml.Unmarshal(data, &ex); err != nil {
		return nil, fmt.Errorf("cannot parse YAML def for %q: %w", yml, err)
	}

	return &ex, nil
}

// GetInput chooses an input from the given exercise using the student id
func (ex *Exercise) GetInput(id string) Input {
	index := GetStudentHash(id) % len(ex.Input)
	return ex.Input[index]
}
