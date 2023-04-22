package printer

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func frmt(tp, val string) string {
	return fmt.Sprintf(`{{index . "op1"}} %s {{index . "op2"}} %s`, tp, val)
}

func TestPrinter(t *testing.T) {
	assert.Equal(t, frmt("boolean", "true"), printBool(true), "bool")
	assert.Equal(t, frmt("boolean", "false"), printBool(false), "bool")

	assert.Equal(t, frmt("integer", "12"), printInt(12), "int")
	assert.Equal(t, frmt("integer", "33"), printInt(33), "int")

	assert.Equal(t, frmt("integer in hexadecimal form", "10"), printHex(16), "hex")
	assert.Equal(t, frmt("integer in hexadecimal form", "9"), printHex(9), "hex")
	assert.Equal(t, frmt("integer in hexadecimal form", "ff"), printHex(255), "hex")

	assert.Equal(t, frmt("float", "1.12"), printFloat(1.1234), "float 1")
	assert.Equal(t, frmt("float", "9.97"), printFloat(9.971), "float 2")
	assert.Equal(t, frmt("float", "0.01"), printFloat(0.0123), "float 3")

	assert.Equal(t, "ab", concatStrings("a", "b"), "concat 1")
	assert.Equal(t, "aaabb", concatStrings("aaa", "bb"), "concat 2")

	assert.Equal(t, frmt("string", `"ab"`), printConcatStrings("a", "b"), "print-concat 1")
	assert.Equal(t, frmt("string", `"aaabb"`), printConcatStrings("aaa", "bb"), "print concat 2")
}
