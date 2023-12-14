package helloworld

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHelloWorld(t *testing.T) {
	assert.Equal(t, "{{index . "text"}}", helloWorld(), "hello world in {{index . "language"}}")
}
