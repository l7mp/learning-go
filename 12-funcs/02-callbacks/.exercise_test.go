package calculator

import (
	"testing"
)

func TestCalculate(t *testing.T) {
	testCases := []struct {
		op       OperationType
		a, b     float64
		expected float64
	}{
		{Add, 10, 5, 15},
		{Subtract, 20, 7, 13},
		{Add, -5, 3, -2},
		{Subtract, 15, 8, 7},
		{Multiply, 2, 3, 6},
	}

	for _, tc := range testCases {
		result := Calculate(tc.op, tc.a, tc.b)
		if result != tc.expected {
			t.Errorf("Calculate(%s, %.2f, %.2f) = %.2f; want %.2f", tc.op, tc.a, tc.b, result, tc.expected)
		}
	}
}

func TestOperationTypeString(t *testing.T) {
	testCases := []struct {
		opType   OperationType
		expected string
	}{
		{Add, "Add"},
		{Subtract, "Subtract"},
		{Multiply, "Multiply"},
	}

	for _, tc := range testCases {
		result := tc.opType.String()
		if result != tc.expected {
			t.Errorf("%d.String() = %s; want %s", tc.opType, result, tc.expected)
		}
	}
}
