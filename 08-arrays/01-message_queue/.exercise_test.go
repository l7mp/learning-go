package messagequeue

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMessageQueue(t *testing.T) {
	x := messageQueue("Hello", "World", "Again")
	assert.Equal(t, "{{index . "text" 0}}", x[0])
	assert.Equal(t, "{{index . "text" 1}}", x[1])
	assert.Equal(t, "{{index . "text" 2}}", x[2])
}
