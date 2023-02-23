package repaint

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepaintColor(t *testing.T) {
	assert.Equal(t, repaintColor("{{index . "primary"}}"), "{{index . "complementary"}}")
}
