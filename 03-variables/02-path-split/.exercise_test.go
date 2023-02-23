package pathsplit

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitPath(t *testing.T) {
	assert.Equal(t, splitPath("{{index . "fullpath"}}", "{{index . "component"}}"), "{{index . "output"}}")
}
