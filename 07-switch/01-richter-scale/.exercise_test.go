package richterscale

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDescribeEarthquake(t *testing.T) {
{{range index . "tests"}}
	assert.Equal(t, describeEarthquake({{index . "magnitude"}}), "{{index . "description"}}")
{{end}}
}
