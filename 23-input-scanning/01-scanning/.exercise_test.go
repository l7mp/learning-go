package scanning

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"bytes"
)

func TestCounter(t *testing.T) {
	var stdin bytes.Buffer
	stdin.WriteString("{{index . "text"}}")

	assert.Equal(t, {{index . "number"}}, counter(&stdin), "The two numbers should be equal.")
}
