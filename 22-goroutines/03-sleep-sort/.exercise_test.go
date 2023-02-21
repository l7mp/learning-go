package sleepSort

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSleepSort(t *testing.T) {
	input := []uint{{index . "input"}}
	result := []uint{{index . "result"}}
	assert.Equal(t, {{index . "funcName"}}(input), result, "{{index . "dir"}} sleep-sort")
}
