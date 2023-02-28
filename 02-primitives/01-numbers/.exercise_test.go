package numbers

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const r float64 = 3.4
const a float64 = 3.6
const b float64 = 4.2
const h float64 = 6.9


func TestArea(t *testing.T) {
	switch {{index . "shape"}} {
	case "Circle":
		assert.Equal(t, shape(r), 36.32)
	case "Ellipse":
		assert.Equal(t, shape(a, b), 47.5)
	case "Trapezoid":
		assert.Equal(t, shape(a, b, h), 26.91)
	case "Pentagon":
		assert.Equal(t, shape(a), 22.297)
	case "Hexagon":
		assert.Equal(t, shape(a), 33.67)
	case "Octagon":
		assert.Equal(t, shape(a), 62.58)
	}
}
