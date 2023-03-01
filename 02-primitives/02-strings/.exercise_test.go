package strings

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReplaceSubnet(t *testing.T) {
	assert.Equal(t, replaceSubnet("{{index . "cidr"}}", "{{index . "subnet"}}"), "{{index . "address"}}/{{index . "subnet"}}")
}
