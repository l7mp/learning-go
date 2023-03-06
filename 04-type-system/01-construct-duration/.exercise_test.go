package constructduration

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func setConverter(s string) time.Duration {
	switch s {
	case "hours":
		return time.Hour
	case "minutes":
		return time.Minute
	case "seconds":
		return time.Second
	case "milliseconds":
		return time.Millisecond
	case "microsecond":
		return time.Microsecond
	default:
		return 0 // Will alert for checks in input
 	}
}

func TestConstructDuration(t *testing.T) {
	arg1Conv := setConverter ("{{index . "type1"}}")
	arg2Conv := setConverter ("{{index . "type2"}}")
	assert.Equal(t, time.Duration({{index . "arg1"}}*arg1Conv+{{index . "arg2"}}*arg2Conv), constructDuration({{index . "arg1"}}, {{index . "arg2"}}))
}
