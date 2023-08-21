package factorial

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalcFactorialOrSum(t *testing.T) {
	{{if eq (index . "func") "factorial"}}
	assert.Equal(t, 3628800, calcFactorial(10))
	assert.Equal(t, 24, calcFactorial(4))
	assert.Equal(t, 1, calcFactorial(1))
	{{end}}
	{{if eq (index . "func") "sum"}}
	assert.Equal(t, 55, calcSum(10))
	assert.Equal(t, 10, calcSum(4))
	assert.Equal(t, 1, calcSum(1))
	{{end}}
	{{if eq (index . "func") "abs"}}
	assert.Equal(t, 10, calcSum(10))
	assert.Equal(t, 4, calcSum(-4))
	assert.Equal(t, 1, calcSum(-1))
	{{end}}
}
