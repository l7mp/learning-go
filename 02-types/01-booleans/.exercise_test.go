package logicalops

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

{{index . "demorganimpl"}}

func TestInverse(t *testing.T) {
	assert.Equal(t, false, inverse(true), "inv 1")
	assert.Equal(t, true, inverse(false), "inv 1")
}

func Test{{index . "op"}}(t *testing.T) {
	assert.Equal(t, true, {{index . "name"}}(true, true), "op1")
	assert.Equal(t, true {{index . "impl"}} false, {{index . "name"}}(true, false), "op2")
	assert.Equal(t, false {{index . "impl"}} true, {{index . "name"}}(false, true), "op3")
	assert.Equal(t, false, {{index . "name"}}(false, false), "op4")
}

func TestDeMorgan(t *testing.T) {
	assert.Equal(t, dm(true,  true),  deMorgan(true, true), "op1")
	assert.Equal(t, dm(true,  false), deMorgan(true, false), "op1")
	assert.Equal(t, dm(false, true),  deMorgan(false, true), "op1")
	assert.Equal(t, dm(false, false), deMorgan(false, false), "op1")
}
