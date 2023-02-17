package helloworld

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHelloWorld(t *testing.T) {
	assert.Equal(t, helloWorld(), "{{index . "text"}}", "hello world in {{index . "language"}}")
}
