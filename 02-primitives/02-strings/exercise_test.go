package strings

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReplaceSubnet(t *testing.T) {
	assert.Equal(t, replaceSubnet("192.168.0.1/10", "11"), "192.168.0.1/11")
}
