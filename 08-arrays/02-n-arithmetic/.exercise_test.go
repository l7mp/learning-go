package narithmetic

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNArithmetic(t *testing.T) {
	assert.EqualValues(t, nArithmetic(), mySolution())
}
