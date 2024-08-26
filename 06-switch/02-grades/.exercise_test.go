package grades

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGradeExam(t *testing.T) {
	assert.Equal(t, {{index . "grades" 0}}, gradeExam(100))
	assert.Equal(t, {{index . "grades" 0}}, gradeExam({{index . "percents" 0}}))
	assert.Equal(t, {{index . "grades" 1}}, gradeExam({{index . "percents" 1}}))
	assert.Equal(t, {{index . "grades" 2}}, gradeExam({{index . "percents" 2}}))
	assert.Equal(t, {{index . "grades" 3}}, gradeExam({{index . "percents" 3}}))
	assert.Equal(t, {{index . "grades" 4}}, gradeExam(0))
}
