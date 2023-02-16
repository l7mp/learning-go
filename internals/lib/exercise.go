package lib

import (
	"os"
	"path/filepath"

	yaml "gopkg.in/yaml.v3"
)

const (
	ExerciseFileName = "exercise.yaml"
)

type Exercise struct {
	Name  string    `yaml:"name"`
	Input [](Input) `yaml:"input"`
}

// NewExercise reads the exercise definition in the given directory
func NewExercise(dir string) (*Exercise, error) {
	data, err := os.ReadFile(filepath.Join(dir, ExerciseFileName))
	if err != nil {
		return nil, err
	}

	ex := Exercise{}
	if err := yaml.Unmarshal(data, &ex); err != nil {
		return nil, err
	}

	return &ex, nil
}

// GetInput chooses an input from the given exercise using the student id
func (ex *Exercise) GetInput(id string) Input {
	return ex.Input[GetStudentHash(id)%len(ex.Input)]
}
