package scanning

import (
	"testing"
	"bytes"

	"github.com/stretchr/testify/assert"
)

func TestCounter(t *testing.T) {
	var stdin bytes.Buffer
	stdin.WriteString("{{index . "text"}}")
	assert.Equal(t, {{index . "number"}}, counter(&stdin), "The two numbers should be equal.")
}
