package filteringdata

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFilteringData(t *testing.T) {
	assert.EqualValues(t, filteringData({{index . "keys"}}, {{index . "values"}}), mySolution({{index . "keys"}}, {{index . "values"}}))
}
