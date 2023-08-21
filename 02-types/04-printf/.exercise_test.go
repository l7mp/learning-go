package printer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPrinter(t *testing.T) {
	assert.Equal(t, `{{index . "op1"}} boolean {{index . "op2"}} true`, printBool(true), "bool")
	assert.Equal(t, `{{index . "op1"}} boolean {{index . "op2"}} false`, printBool(false), "bool")
	assert.Equal(t, `{{index . "op1"}} integer {{index . "op2"}} 12`, printInt(12), "int")
	assert.Equal(t, `{{index . "op1"}} integer {{index . "op2"}} 33`, printInt(33), "int")
	assert.Equal(t, `{{index . "op1"}} integer in hexadecimal form {{index . "op2"}} 10`, printHex(16), "hex")
	assert.Equal(t, `{{index . "op1"}} integer in hexadecimal form {{index . "op2"}} 9`, printHex(9), "hex")
	assert.Equal(t, `{{index . "op1"}} integer in hexadecimal form {{index . "op2"}} ff`, printHex(255), "hex")
	assert.Equal(t, `{{index . "op1"}} float {{index . "op2"}} 1.12`, printFloat(1.1234), "float 1")
	assert.Equal(t, `{{index . "op1"}} float {{index . "op2"}} 9.97`, printFloat(9.971), "float 2")
	assert.Equal(t, `{{index . "op1"}} float {{index . "op2"}} 0.01`, printFloat(0.0123), "float 3")

	assert.Equal(t, "ab", concatStrings("a", "b"), "concat 1")
	assert.Equal(t, "aaabb", concatStrings("aaa", "bb"), "concat 2")

	assert.Equal(t, `{{index . "op1"}} string {{index . "op2"}} "ab"`, printConcatStrings("a", "b"), "print-concat 1")
	assert.Equal(t, `{{index . "op1"}} string {{index . "op2"}} "aaabb"`, printConcatStrings("aaa", "bb"), "print concat 2")
}
