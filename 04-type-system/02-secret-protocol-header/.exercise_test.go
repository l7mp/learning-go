package secretprotocolheader

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

const PACKET_TYPE_LSB = 5
const FIRST_ATTEMPT_LSB = 4
const QOS_LSB = 3
const BROADCST_LSB = 1
const SECURE_LSB = 0

func mySolution(isFirstAttempt byte, isBroadcasted byte, isSecure byte, qos byte) byte {
	return byte(0x02 << PACKET_TYPE_LSB | isFirstAttempt << FIRST_ATTEMPT_LSB | qos << QOS_LSB | isBroadcasted << BROADCST_LSB | isSecure << SECURE_LSB) 
	
}

func TestCreatePublishFixHeader(t *testing.T) {
	assert.Equal(t, createPublishFixHeader({{index . "fa" "val"}}, {{index . "bc" "val"}}, {{index . "sc" "val"}}), mySolution({{index . "fa" "val"}}, {{index . "bc" "val"}}, {{index . "sc" "val"}}, {{index . "qos" "val"}}))
	fmt.Println(mySolution({{index . "fa" "val"}}, {{index . "bc" "val"}}, {{index . "sc" "val"}}, {{index . "qos" "val"}}))
}