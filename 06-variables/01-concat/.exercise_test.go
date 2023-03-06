package concat

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConcatVariables(t *testing.T) {
	a, b := concatVariables("{{index . "var1"}}", "{{index . "var2"}}")
	assert.Equal(t, a, "{{index . "var1"}}{{index . "var2"}}")
	assert.Equal(t, b, "{{index . "var2"}}{{index . "var1"}}")
}
