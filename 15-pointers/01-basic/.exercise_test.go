package pointerbasic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRetrieveValue(t *testing.T) {
	{{if eq (index . "type") "int"}}
	var x int
	x = 100
	assert.Equal(t, x, retrieveValue(&x))
	x = -100
	assert.Equal(t, x, retrieveValue(&x))
	x = 1
	assert.Equal(t, x, retrieveValue(&x))
	x = 11
	assert.Equal(t, x, retrieveValue(&x))
	{{end}}

	{{if eq (index . "type") "string"}}
	var x string
	x = "100"
	assert.Equal(t, x, retrieveValue(&x))
	x = "aaa"
	assert.Equal(t, x, retrieveValue(&x))
	x = "Joe"
	assert.Equal(t, x, retrieveValue(&x))
	x = ""
	assert.Equal(t, x, retrieveValue(&x))
	{{end}}

	{{if eq (index . "type") "bool"}}
	var x bool
	x = true
	assert.Equal(t, x, retrieveValue(&x))
	x = false
	assert.Equal(t, x, retrieveValue(&x))
	{{end}}
}
