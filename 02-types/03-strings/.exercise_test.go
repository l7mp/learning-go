package strings

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStrings(t *testing.T) {
	assert.Equal(t, `some
multiline
string`, multilineString(), "multiline")

	assert.Equal(t, 10, stringLen("1234567890"), "len 1")
	assert.Equal(t, 0, stringLen(""), "len 2")

	assert.Equal(t, "234567890", trimFirstChar("1234567890"), "len 1")
	assert.Equal(t, 9, stringLen(trimFirstChar("1234567890")), "len 1")
	assert.Equal(t, "", trimFirstChar(""), "len 2")
	assert.Equal(t, 0, stringLen(trimFirstChar("")), "len 1")

	assert.Equal(t, "123456789", trimLastChar("1234567890"), "len 1")
	assert.Equal(t, 9, stringLen(trimLastChar("1234567890")), "len 1")
	assert.Equal(t, "", trimLastChar(""), "len 2")
	assert.Equal(t, 0, stringLen(trimLastChar("")), "len 1")

	assert.Equal(t, "A234567890", swapFirstChar("1234567890"), "len 1")
	assert.Equal(t, 10, stringLen(swapFirstChar("1234567890")), "len 1")
	assert.Equal(t, "", swapFirstChar(""), "len 2")
	assert.Equal(t, 0, stringLen(swapFirstChar("")), "len 1")

	assert.Equal(t, "123456789A", swapLastChar("1234567890"), "len 1")
	assert.Equal(t, 10, stringLen(swapLastChar("1234567890")), "len 1")
	assert.Equal(t, "", swapLastChar(""), "len 2")
	assert.Equal(t, 0, stringLen(swapLastChar("")), "len 1")

	assert.Equal(t, "A1234567890", prependChar("1234567890"), "len 1")
	assert.Equal(t, 11, stringLen(prependChar("1234567890")), "len 1")
	assert.Equal(t, "A", prependChar(""), "len 2")
	assert.Equal(t, 1, stringLen(prependChar("")), "len 1")

	assert.Equal(t, "1234567890A", appendChar("1234567890"), "len 1")
	assert.Equal(t, 11, stringLen(appendChar("1234567890")), "len 1")
	assert.Equal(t, "A", appendChar(""), "len 2")
	assert.Equal(t, 1, stringLen(appendChar("")), "len 1")
}
