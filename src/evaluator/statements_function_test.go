package evaluator

import (
	"testing"

	"github.com/uncode/object"
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

func TestStatementFunctionReturnValues(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
	}{
		{
			name: "直接の戻り値 (💩による戻り値)",
			input: `
				func add(x, y) {
					x + y >> 💩
				}
				add(5, 3)
			`,
			expected: int64(8),
		},
		{
			name: "暗黙の戻り値 (最後の式の結果)",
			input: `
				func multiply(x, y) {
					x * y
				}
				multiply(4, 5)
			`,
			expected: int64(20),
		},
		{
			name: "高階関数と戻り値 (mapに渡す関数)",
			input: `
				func double(x) {
					x * 2
				}
				[1, 2, 3].map(double)
			`,
			expected: []int64{2, 4, 6},
		},
		{
			name: "複数の関数呼び出しと戻り値",
			input: `
				func add(x, y) {
					x + y >> 💩
				}
				func multiply(x, y) {
					x * y >> 💩
				}
				add(multiply(2, 3), multiply(4, 5))
			`,
			expected: int64(26), // (2*3) + (4*5) = 6 + 20 = 26
		},
		{
			name: "関数のパイプライン適用と戻り値",
			input: `
				func double(x) {
					x * 2 >> 💩
				}
				func add5(x) {
					x + 5 >> 💩
				}
				10 |> double |> add5
			`,
			expected: int64(25), // (10*2)+5 = 25
		},
		{
			name: "条件分岐がある関数",
			input: `
				func abs(x) {
					if x < 0 {
						-x >> 💩
					} else {
						x >> 💩
					}
				}
				abs(-10)
			`,
			expected: int64(10),
		},
		{
			name: "入れ子のブロック文",
			input: `
				func complexFunc(x) {
					let result = 0
					if x > 0 {
						{
							let temp = x * 2
							result = temp + 1
						}
					} else {
						{
							let temp = -x
							result = temp * 2
						}
					}
					result >> 💩
				}
				complexFunc(5)
			`,
			expected: int64(11), // (5*2)+1 = 11
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluated := testEval(tt.input)
			
			switch expected := tt.expected.(type) {
			case int64:
				testIntegerObject(t, evaluated, expected)
			case []int64:
				result, ok := evaluated.(*object.Array)
				if !ok {
					t.Fatalf("object is not Array. got=%T (%+v)", evaluated, evaluated)
				}
				if len(result.Elements) != len(expected) {
					t.Fatalf("array has wrong num of elements. got=%d", len(result.Elements))
				}
				for i, expectedElement := range expected {
					testIntegerObject(t, result.Elements[i], expectedElement)
				}
			default:
				t.Fatalf("未対応の期待値の型: %T", expected)
			}
		})
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

func TestCaseStatements(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected interface{}
	}{
		{
			name: "シンプルなcase文",
			input: `
				func test() {
					case 1:
						100 >> 💩
					case 2:
						200 >> 💩
					case default:
						300 >> 💩
				}
				test()
			`,
			expected: int64(300),
		},
		{
			name: "条件付きcase文",
			input: `
				func test() {
					case 🍕 % 2 == 0:
						"偶数" >> 💩
					case 🍕 % 2 != 0:
						"奇数" >> 💩
					case default:
						"不明" >> 💩
				}
				test(4)
			`,
			expected: "偶数",
		},
		{
			name: "複数のcase文",
			input: `
				func test() {
					case 🍕 < 0:
						"負の数" >> 💩
					case 🍕 == 0:
						"ゼロ" >> 💩
					case 🍕 > 0:
						"正の数" >> 💩
					case default:
						"不明" >> 💩
				}
				test(-5)
			`,
			expected: "負の数",
		},
		{
			name: "case文のネスト",
			input: `
				func test() {
					case 🍕 > 0:
						case 🍕 % 2 == 0:
							"正の偶数" >> 💩
						case 🍕 % 2 != 0:
							"正の奇数" >> 💩
					case 🍕 < 0:
						case 🍕 % 2 == 0:
							"負の偶数" >> 💩
						case 🍕 % 2 != 0:
							"負の奇数" >> 💩
					case default:
						"ゼロ" >> 💩
				}
				test(3)
			`,
			expected: "正の奇数",
		},
		{
			name: "case文のデフォルト",
			input: `
				func test() {
					case 🍕 == 1:
						"One" >> 💩
					case 🍕 == 2:
						"Two" >> 💩
					case default:
						"Other" >> 💩
				}
				test(10)
			`,
			expected: "Other",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			evaluated := testEval(tt.input)
			
			switch expected := tt.expected.(type) {
			case int64:
				testIntegerObject(t, evaluated, expected)
			case string:
				result, ok := evaluated.(*object.String)
				if !ok {
					t.Fatalf("object is not String. got=%T (%+v)", evaluated, evaluated)
				}
				if result.Value != expected {
					t.Fatalf("string has wrong value. got=%s, want=%s", result.Value, expected)
				}
			default:
				t.Fatalf("未対応の期待値の型: %T", expected)
			}
		})
	}
}
