package constructduration

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestConstructDuration(t *testing.T) {
	assert.Equal(t, `{{index . "sol1"}}`, constructDuration(1, 2).String())
	assert.Equal(t, `{{index . "sol2"}}`, constructDuration(4, 7).String())
	assert.Equal(t, `{{index . "sol3"}}`, constructDuration(12, 23).String())
}
