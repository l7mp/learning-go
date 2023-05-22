package secretprotocolheader

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestCreatePublishFixHeader(t *testing.T) {
	assert.Equal(t, createPublishFixHeader({{index . "fa" "val"}}, {{index . "bc" "val"}}, {{index . "sc" "val"}}), mySolution({{index . "fa" "val"}}, {{index . "bc" "val"}}, {{index . "sc" "val"}}, {{index . "qos" "val"}}))
}