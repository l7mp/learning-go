package messagequeue

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestmessageQueue(t *testing.T) {
	assert.EqualValues(t, messageQueue("{{index . "text" 0}}", "{{index . "text" 1}}", "{{index . "text" 2}}"), mySolution("{{index . "text" 0}}", "{{index . "text" 1}}", "{{index . "text" 2}}"))
}
