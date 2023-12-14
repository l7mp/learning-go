package richterscale

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDescribeEarthquake(t *testing.T) {
{{range index . "tests"}}
	assert.Equal(t, "{{index . "description"}}", describeEarthquake({{index . "magnitude"}}))
{{end}}
}
