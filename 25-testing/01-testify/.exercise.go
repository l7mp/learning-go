package testimony

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// DO NOT REMOVE THIS COMMENT
//go:generate go run ../../exercises-cli.go -student-id=$STUDENT_ID generate
{{- if eq (index . "name") "strstr"}}

func strstr(str string, substr string) int {
	if substr == "" {
		return 0
	}

	for i := 0; i <= len(str)-len(substr); i++ {
		if str[i:i+len(substr)] == substr {
			return i
		}
	}

	return -1
}
{{- end}}
{{- if eq (index . "name") "strncat"}}

func strncat(dest *string, src string, n int) {
	if n > len(src) {
		n = len(src)
	}
	*dest += src[:n]
}
{{- end}}
{{- if eq (index . "name") "strncpy"}}

func strncpy(dest *string, src string, n int) {
	if n > len(src) {
		n = len(src)
	}
	*dest = src[:n]
}
{{- end}}


// INSERT YOUR CODE HERE
{{- if eq (index . "name") "strstr"}}
func StrStrMatch(t *testing.T) {

}

func StrStrNoMatch(t *testing.T) {

}

func StrStrEmptySubString(t *testing.T) {

}
{{end}}
{{- if eq (index . "name") "strncat"}}
func StrNCatInbounds(t *testing.T) {

}

func StrNCatOutOfBounds(t *testing.T) {

}

func StrNCatEmptySource(t *testing.T) {

}
{{end}}
{{- if eq (index . "name") "strncpy"}}
func StrNCpyINbounds(t *testing.T) {

}

func StrNCpyEmptyDestination(t *testing.T) {

}

func StrNCpyEmptySource(t *testing.T) {

}
{{end}}
