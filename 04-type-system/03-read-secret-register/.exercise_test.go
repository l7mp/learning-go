package readsecretregister

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseChannelControlRegister(t *testing.T) {
	a,b,c,d := parseChannelControlRegister({{index . "CHAN_CTRL"}})
	assert.Equal(t, {{index . "test0"}}, []byte{a,b,c,d})

	a,b,c,d = parseChannelControlRegister(0xdeadbeef)
	assert.Equal(t, {{index . "test1"}}, []byte{a,b,c,d})

	a,b,c,d = parseChannelControlRegister(0x01234567)
	assert.Equal(t, {{index . "test2"}}, []byte{a,b,c,d})
}
