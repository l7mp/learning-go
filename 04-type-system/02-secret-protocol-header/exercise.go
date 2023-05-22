package secretprotocolheader

// DO NOT REMOVE THIS COMMENT
//go:generate go run ../../exercises-cli.go -student-id=$STUDENT_ID generate

// createPublishFixHeader constructs an octet (8-bit long byte) based on its three arguments and the fix QoS setting
func createPublishFixHeader(isFirstAttempt byte, isBroadcasted byte, isSecure byte) byte {
	// INSERT YOUR CODE HERE
}
