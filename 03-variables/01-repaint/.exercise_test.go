package repaint

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepaintColor(t *testing.T) {
	c, err := repaintColor("{{index . "primary"}}")
	assert.Equal(t, c, "{{index . "complementary"}}")
	assert.NoError(t, err)

	c, err = repaintColor("{{index . "complementary"}}")
	assert.Equal(t, c, "{{index . "primary"}}")
	assert.NoError(t, err)

	c, err = repaintColor("xxx")
	assert.Equal(t, c, "")
	assert.Error(t, err)
}
