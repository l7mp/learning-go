package pipeline

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestPipeline(t *testing.T) {
	input := []int{{index . "input"}}
	result := []float32{{index . "output"}}
	assert.Equal(t, result, collector({{index . "type"}}(generator(input))))
}
