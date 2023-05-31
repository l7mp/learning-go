package factorial

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalcFactorialOrSum(t *testing.T) {
	{{if eq (index . "func") "factorial"}}
	assert.Equal(t, calcFactorial({{index . "input"}}), {{index . "output"}})
	{{end}}
	{{if eq (index . "func") "sum"}}
	assert.Equal(t, calcSum({{index . "input"}}), {{index . "output"}})
	{{end}}
}
