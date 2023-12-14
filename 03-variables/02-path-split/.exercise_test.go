package pathsplit

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSplitPath(t *testing.T) {
	assert.Equal(t, "{{index . "tests" 0 "output"}}", splitPath("{{index . "tests" 0 "fullpath"}}"))
	assert.Equal(t, "{{index . "tests" 1 "output"}}", splitPath("{{index . "tests" 1 "fullpath"}}"))
	assert.Equal(t, "{{index . "tests" 2 "output"}}", splitPath("{{index . "tests" 2 "fullpath"}}"))
	assert.Equal(t, "{{index . "tests" 3 "output"}}", splitPath("{{index . "tests" 3 "fullpath"}}"))
	assert.Equal(t, "{{index . "tests" 4 "output"}}", splitPath("{{index . "tests" 4 "fullpath"}}"))
}
