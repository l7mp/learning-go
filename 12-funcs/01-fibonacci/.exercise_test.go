package fibonacci

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFibonacci(t *testing.T) {
	testCases := []struct {
		input    int
		expected int
	}{
		{0, 0},
		{1, 1},
		{2, 1},
		{3, 2},
		{7, 13},
		{10, 55},
	}

	for _, tc := range testCases {
		assert.Equal(t, tc.expected, fibonacci(tc.input))
	}
}
