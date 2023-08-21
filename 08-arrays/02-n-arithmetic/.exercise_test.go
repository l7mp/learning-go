package narithmetic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNArithmetic(t *testing.T) {
	input := [10]int{1, 2, 3, 4, 5, 6, 7, 8, 9, 10}

	{{if eq (index . "text") "Addition"}}
	assert.Equal(t, 55, nArithmetic(input))
	{{end}}

	{{if eq (index . "text") "Multiplication"}}
	assert.Equal(t, 3628800, nArithmetic(input))
	{{end}}

	{{if eq (index . "text") "Subtraction"}}
	assert.Equal(t, -53, nArithmetic(input))
	{{end}}
}
