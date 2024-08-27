package pointernew

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewValue(t *testing.T) {
	assert.Equal(t, {{index . "value"}}, *newValue())
}
