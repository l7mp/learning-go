package digits

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDigits(t *testing.T) {
	{{if eq (index . "func") "sum"}}
	{{range index . "tests"}}
		assert.Equal(t, {{index . "output"}}, sumDigits({{index . "input"}}))
	{{end}}
	{{end}}
	{{if eq (index . "func") "product"}}
	{{range index . "tests"}}
		assert.Equal(t, {{index . "output"}}, multiplyDigits({{index . "input"}}))
	{{end}}
	{{end}}
}
