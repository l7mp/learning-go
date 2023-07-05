package pointerbasic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRetriveValue(t *testing.T) {
	value := {{index . "value"}}
	assert.Equal(t, {{index . "value"}}, retriveValue(&value), "The result should be: {{index . "value"}}")
}
