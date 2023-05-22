package readsecretregister

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseChannelControlRegister(t *testing.T) {
	TX_CHAN, RX_CHAN, RX_PCODE, TX_PCODE := parseChannelControlRegister({{index . "CHAN_CTRL"}})
	TX_CHAN_REFERNCE, RX_CHAN_REFERNCE, RX_PCODE_REFERNCE, TX_PCODE_REFERNCE := mySolution({{index . "CHAN_CTRL"}})
	assert.Equal(t, TX_CHAN, TX_CHAN_REFERNCE)
	assert.Equal(t, RX_CHAN, RX_CHAN_REFERNCE)
	assert.Equal(t, RX_PCODE, RX_PCODE_REFERNCE)
	assert.Equal(t, TX_PCODE, TX_PCODE_REFERNCE)
}