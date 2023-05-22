package secretprotocolheader

// DO NOT REMOVE THIS COMMENT
//go:generate go run ../../exercises-cli.go -student-id=$STUDENT_ID generate

const PACKET_TYPE_LSB = 5
const FIRST_ATTEMPT_LSB = 4
const QOS_LSB = 3
const BROADCST_LSB = 1
const SECURE_LSB = 0

func createPublishFixHeader(isFirstAttempt bool, isBroadcasted bool, isSecure bool) uint8 {
	return uint8(0x02 << PACKET_TYPE_LSB | isFirstAttempt << FIRST_ATTEMPT_LSB | {{index . "qos.val"}} << QOS_QOS_LSB | isBroadcasted << BROADCST_LSB | isSecure << SECURE_LSB) 
}