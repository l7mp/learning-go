package secretprotocolheader

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreatePublishFixHeader(t *testing.T) {
	assert.Equal(t, byte({{index . "sol0"}}), createPublishFixHeader(false, false, false))
	assert.Equal(t, byte({{index . "sol1"}}), createPublishFixHeader(false, false, true))
	assert.Equal(t, byte({{index . "sol2"}}), createPublishFixHeader(false, true, false))
	assert.Equal(t, byte({{index . "sol3"}}), createPublishFixHeader(false, true, true))
	assert.Equal(t, byte({{index . "sol4"}}), createPublishFixHeader(true, false, false))
	assert.Equal(t, byte({{index . "sol5"}}), createPublishFixHeader(true, false, true))
	assert.Equal(t, byte({{index . "sol6"}}), createPublishFixHeader(true, true, false))
	assert.Equal(t, byte({{index . "sol7"}}), createPublishFixHeader(true, true, true))
}
