package factorial

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCalcFactorial(t *testing.T) {
	assert.Equal(t, calcFactorial({{index . "n"}}), {{index . "nfactor"}})
}
