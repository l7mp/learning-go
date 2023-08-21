package calculator

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalculator(t *testing.T) {
	assert.Equal(t, {{index . "sol0"}}, {{index . "funcName"}}(1.2, 2.4), "basic test 0")
	assert.Equal(t, {{index . "sol1"}}, {{index . "funcName"}}(3.67, 100.0), "basic test 1")
	assert.Equal(t, {{index . "sol2"}}, {{index . "funcName"}}(34.91, 144.02), "basic test 2")

	val, err := {{index . "funcName"}}String("3.67", "100.0")
	assert.NoError(t, err, "string parse test 1")
	assert.Equal(t, {{index . "sol1"}}, val, "string test 1")

	val, err = {{index . "funcName"}}String("34.91", "144.02")
	assert.NoError(t, err, "string parse test 2")
	assert.Equal(t, {{index . "sol2"}}, val, "string test 2")

	val, err = {{index . "funcName"}}String("dummmy", "100.0")
	assert.Error(t, err, "string parse test x")
	assert.Equal(t, 0, val, "string test 1")

	val, err = {{index . "funcName"}}String("1.0", "dummmy")
	assert.Error(t, err, "string parse test y")
	assert.Equal(t, 0, val, "string test 2")
}
