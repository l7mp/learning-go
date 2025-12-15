package testimony

import (
	"os"
	"regexp"
	"testing"
)

const requiredAsserts = 3

func TestStudentUsesAssertions(t *testing.T) {
	src, err := os.ReadFile("exercise.go")
	if err != nil {
		t.Fatalf("cannot read exercise.go: %v", err)
	}

	source := string(src)

	re := regexp.MustCompile(`assert\.\w+\s*\(`)
	if !re.MatchString(source) {
		t.Fatalf("no assert.* call found in exercise.go")
	}
}

func TestStudentHasEnoughAsserts(t *testing.T) {
	src, err := os.ReadFile("exercise.go")
	if err != nil {
		t.Fatalf("cannot read exercise.go: %v", err)
	}

	source := string(src)

	re := regexp.MustCompile(`assert\.\w+\s*\(`)
	matches := re.FindAllString(source, -1)

	if len(matches) < requiredAsserts {
		t.Fatalf("not enough assert calls: found %d, required %d", len(matches), requiredAsserts)
	}
}

func TestNoTrivialAssertions(t *testing.T) {
	src, err := os.ReadFile("exercise.go")
	if err != nil {
		t.Fatalf("cannot read exercise.go: %v", err)
	}

	source := string(src)

	patterns := []string{
		`assert\.False\s*\(\s*t\s*,\s*false\s*[,)]`,
		`assert\.True\s*\(\s*t\s*,\s*true\s*[,)]`,
		`assert\.Equal\s*\(\s*t\s*,\s*true\s*,\s*true\s*[,)]`,
		`assert\.Equal\s*\(\s*t\s*,\s*false\s*,\s*false\s*[,)]`,
	}

	for _, p := range patterns {
		re := regexp.MustCompile(p)
		if re.MatchString(source) {
			t.Fatalf("trivial assertion detected")
		}
	}
}

func TestRunStudentTests(t *testing.T) {
{{if eq (index . "name") "strstr"}}
	t.Run("StrStrMatch", StrStrMatch)
	t.Run("StrStrNoMatch", StrStrNoMatch)
	t.Run("StrStrEmptySubString", StrStrEmptySubString)
{{end}}
{{if eq (index . "name") "strncat"}}
	t.Run("StrNCatInbounds", StrNCatInbounds)
	t.Run("StrNCatOutOfBounds", StrNCatOutOfBounds)
	t.Run("StrNCatEmptySource", StrNCatEmptySource)
{{end}}
{{if eq (index . "name") "strncpy"}}
	t.Run("StrNCpyINbounds", StrNCpyINbounds)
	t.Run("StrNCpyEmptyDestination", StrNCpyEmptyDestination)
	t.Run("StrNCpyEmptySource", StrNCpyEmptySource)
{{end}}
}
