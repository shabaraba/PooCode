package evaluator

import (
	"testing"
)

func TestAssignmentStatements(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"a = 5; a;", 5},
		{"a = 5 * 5; a;", 25},
		{"a = 5; b = a; b;", 5},
		{"a = 5; b = a; c = a + b + 5; c;", 15},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestFunctionApplication(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"let identity = fn(x) { x; }; identity(5);", 5},
		{"let identity = fn(x) { return x; }; identity(5);", 5},
		{"let double = fn(x) { x * 2; }; double(5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5, 5);", 10},
		{"let add = fn(x, y) { x + y; }; add(5 + 5, add(5, 5));", 20},
		{"fn(x) { x; }(5)", 5},
		{
			`
			let factorial = fn(n) { 
				if (n == 0) {
					return 1;
				}
				return n * factorial(n - 1);
			};
			factorial(5);
			`,
			120,
		},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}

func TestPipelineOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected int64
	}{
		{"5 | fn(x) { x + 1; }", 6},
		{"5 | fn(x) { x * 2; }", 10},
		{"5 | fn(x) { x + 1; } | fn(x) { x * 2; }", 12},
		{"5 | fn(x) { x + 1; } | fn(x) { x * 2; } | fn(x) { x + 1; }", 13},
	}

	for _, tt := range tests {
		testIntegerObject(t, testEval(tt.input), tt.expected)
	}
}
