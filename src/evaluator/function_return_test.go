package evaluator

import (
	"testing"

	"github.com/uncode/object"
)

// TestFunctionReturnValues は関数の戻り値処理をテストする
func TestFunctionReturnValues(t *testing.T) {
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
				testArrayObject(t, evaluated, expected)
			case string:
				testStringObject(t, evaluated, expected)
			case nil:
				testNullObject(t, evaluated)
			default:
				t.Fatalf("未対応の期待値の型: %T", expected)
			}
		})
	}
}

func testArrayObject(t *testing.T, obj object.Object, expected []int64) {
	t.Helper()
	arrayObj, ok := obj.(*object.Array)
	if !ok {
		t.Errorf("オブジェクトが配列ではありません。got=%T (%+v)", obj, obj)
		return
	}
	if len(arrayObj.Elements) != len(expected) {
		t.Errorf("配列の長さが不正です。期待=%d, 実際=%d", len(expected), len(arrayObj.Elements))
		return
	}
	for i, expectedElement := range expected {
		testIntegerObject(t, arrayObj.Elements[i], expectedElement)
	}
}
