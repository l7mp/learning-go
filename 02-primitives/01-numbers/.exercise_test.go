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
	switch "{{index . "shape"}}"{
	case "Circle":
		assert.Equal(t, area(r, a, b, h), 36.32)
	case "Ellipse":
		assert.Equal(t, area(r, a, b, h), 47.5)
	case "Trapezoid":
		assert.Equal(t, area(r, a, b, h), 26.91)
	case "Pentagon":
		assert.Equal(t, area(r, a, b, h), 22.297)
	case "Hexagon":
		assert.Equal(t, area(r, a, b, h), 33.67)
	case "Octagon":
		assert.Equal(t, area(r, a, b, h), 62.58)
	}
}
