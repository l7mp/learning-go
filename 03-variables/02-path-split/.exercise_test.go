package pathsplit

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitPath(t *testing.T) {
	assert.Equal(t, splitPath("{{index . "tests" 0 "fullpath"}}"), "{{index . "tests" 0 "output"}}")
	assert.Equal(t, splitPath("{{index . "tests" 1 "fullpath"}}"), "{{index . "tests" 1 "output"}}")
	assert.Equal(t, splitPath("{{index . "tests" 2 "fullpath"}}"), "{{index . "tests" 2 "output"}}")
	assert.Equal(t, splitPath("{{index . "tests" 3 "fullpath"}}"), "{{index . "tests" 3 "output"}}")
	assert.Equal(t, splitPath("{{index . "tests" 4 "fullpath"}}"), "{{index . "tests" 4 "output"}}")
}
