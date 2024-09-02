package wordcount

import (
	"reflect"
	"testing"
)

func TestCountWords(t *testing.T) {
	testCases := []struct {
		name     string
		input    []string
		expected map[string]int
	}{
		{"Single text", []string{"the quick brown fox jumps over the lazy dog"},
			map[string]int{"the": 2, "quick": 1, "brown": 1, "fox": 1, "jumps": 1, "over": 1, "lazy": 1, "dog": 1}},
		{"Multiple texts", []string{"hello world", "world hello", "hello hello"},
			map[string]int{"hello": 3, "world": 2}},
		{"Empty input", []string{}, map[string]int{}},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := CountWords(tc.input)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("CountWords(%v) = %v, want %v", tc.input, result, tc.expected)
			}
		})
	}
}
