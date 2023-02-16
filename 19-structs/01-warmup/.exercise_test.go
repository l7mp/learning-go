package structs

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewGame(t *testing.T) {
	assert.Equal(t, newGame().id, {{index . "id"}}, "id")
	assert.Equal(t, newGame().name, "{{index . "name"}}", "name")
	assert.Equal(t, newGame().price, {{index . "price"}}, "price")
	assert.Equal(t, newGame().genre, "{{index . "genre"}}", "genre")
}
