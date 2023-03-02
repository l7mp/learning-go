package strings

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReplaceNetMask(t *testing.T) {
	assert.Equal(t, replaceNetMask("{{index . "cidr"}}", "{{index . "subnet"}}"), "{{index . "address"}}/{{index . "subnet"}}")
}
