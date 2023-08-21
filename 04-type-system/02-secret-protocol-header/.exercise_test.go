package secretprotocolheader

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreatePublishFixHeader(t *testing.T) {
	assert.Equal(t, {{index . "sol0"}}, createPublishFixHeader(false, false, false))
	assert.Equal(t, {{index . "sol1"}}, createPublishFixHeader(false, false, true))
	assert.Equal(t, {{index . "sol2"}}, createPublishFixHeader(false, true, false))
	assert.Equal(t, {{index . "sol3"}}, createPublishFixHeader(false, true, true))
	assert.Equal(t, {{index . "sol4"}}, createPublishFixHeader(true, false, false))
	assert.Equal(t, {{index . "sol5"}}, createPublishFixHeader(true, false, true))
	assert.Equal(t, {{index . "sol6"}}, createPublishFixHeader(true, true, false))
	assert.Equal(t, {{index . "sol7"}}, createPublishFixHeader(true, true, true))
}
