package messagequeue

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessageQueue(t *testing.T) {
	x := messageQueue("Hello", "World", "Again")
	assert.Equal(t, x[0], "{{index . "text" 0}}")
	assert.Equal(t, x[1], "{{index . "text" 1}}")
	assert.Equal(t, x[2], "{{index . "text" 2}}")
}
