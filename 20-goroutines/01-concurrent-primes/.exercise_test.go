package concurrentprimes

import (
	"reflect"
	"testing"
)

func TestGeneratePrimes(t *testing.T) {
	testCases := []struct {
		name     string
		input    int
		expected []int
	}{
		{"Primes below 32", 32, []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31}},
		{"Primes below 70", 70, []int{2, 3, 5, 7, 11, 13, 17, 19, 23, 29, 31, 37, 41, 43, 47, 53, 59, 61, 67}},
		{"Negative input", -1, []int{}},
		{"Zero input", 0, []int{}},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			result := GeneratePrimes(tc.input)
			if !reflect.DeepEqual(result, tc.expected) {
				t.Errorf("GeneratePrimes(%d) = %v, want %v", tc.input, result, tc.expected)
			}
		})
	}
}
