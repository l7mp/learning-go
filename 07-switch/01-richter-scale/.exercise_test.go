package richterscale

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDescribeEarthquake(t *testing.T) {
	assert.Equal(t, describeEarthquake({{index . "magnitude"}}), "{{index . "description"}}")
}
