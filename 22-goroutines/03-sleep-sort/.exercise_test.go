package sleepSort

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSleepSort(t *testing.T) {
	input := []int{{index . "input"}}
	result := []int{{index . "result"}}
	assert.DeepEqual(t, {{index . "funcName"}}(input), result, "{{index . "dir"}} sleep-sort")
}
