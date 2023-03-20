package strings

//go:generate go run ../../exercises-cli.go -student-id=$STUDENT_ID generate

// replaceNetMask will replace the subnet part of a CIDR address
// Example -> 10.0.0.1/10 -> 10.0.0.1/14
func replaceNetMask(address, newSubnet string) string {
	// INSERT YOUR CODE HERE
}
