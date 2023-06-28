package search

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"bytes"
)

func TestContain(t *testing.T) {
	var stdin bytes.Buffer
	stdin.WriteString("{{index . "text"}}")

	assert.Equal(t, {{index . "result"}}, contain(&stdin, "{{index . "word"}}"), "The result should be: {{index . "result"}}")
}
