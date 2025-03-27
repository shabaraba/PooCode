package evaluator

import (
	"testing"
	
	"github.com/uncode/object"
)

// TestRangeExpressions は配列範囲式をテストする
func TestRangeExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected []int64
	}{
		{
			"[1..5]",
			[]int64{1, 2, 3, 4, 5},
		},
		{
			"[5..1]",
			[]int64{5, 4, 3, 2, 1},
		},
		{
			"let start = 2; let end = 6; [start..end]",
			[]int64{2, 3, 4, 5, 6},
		},
		{
			"let a = [1..3]; a",
			[]int64{1, 2, 3},
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerArray(t, evaluated, tt.expected)
	}
}

// TestCharRangeExpressions は文字範囲式をテストする
func TestCharRangeExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{
			"[\"a\"..\"e\"]",
			[]string{"a", "b", "c", "d", "e"},
		},
		{
			"[\"z\"..\"v\"]",
			[]string{"z", "y", "x", "w", "v"},
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testStringArray(t, evaluated, tt.expected)
	}
}

// TestArraySlicing は配列スライシングをテストする
func TestArraySlicing(t *testing.T) {
	tests := []struct {
		input    string
		expected []int64
	}{
		{
			"[1, 2, 3, 4, 5][1..3]",
			[]int64{2, 3},
		},
		{
			"[1, 2, 3, 4, 5][..2]",
			[]int64{1, 2},
		},
		{
			"[1, 2, 3, 4, 5][2..]",
			[]int64{3, 4, 5},
		},
		{
			"[1, 2, 3, 4, 5][..]",
			[]int64{1, 2, 3, 4, 5},
		},
		{
			"let a = [1, 2, 3, 4, 5]; a[1..3]",
			[]int64{2, 3},
		},
		{
			"[1, 2, 3, 4, 5][-2..]",
			[]int64{4, 5},
		},
		{
			"[1, 2, 3, 4, 5][..-2]",
			[]int64{1, 2, 3},
		},
		{
			"[1, 2, 3, 4, 5][-3..-1]",
			[]int64{3, 4},
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		testIntegerArray(t, evaluated, tt.expected)
	}
}

// TestArrayHigherOrderFunctions は配列の高階関数をテストする
func TestArrayHigherOrderFunctions(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		// map関数のテスト（括弧あり）
		{
			"let double = fn() { 🍕 * 2 }; [1, 2, 3] |> map(double)",
			[]int64{2, 4, 6},
		},
		// map関数のテスト（括弧なし）
		{
			"let double = fn() { 🍕 * 2 }; [1, 2, 3] |> map double",
			[]int64{2, 4, 6},
		},
		{
			"let getLength = fn() { len(🍕) }; [\"\", \"hello\", \"world\"] |> map getLength",
			[]int64{0, 5, 5},
		},
		
		// filter関数のテスト（括弧あり）
		{
			"let isEven = fn() { 🍕 % 2 == 0 }; [1, 2, 3, 4, 5] |> filter(isEven)",
			[]int64{2, 4},
		},
		// filter関数のテスト（括弧なし）
		{
			"let isEven = fn() { 🍕 % 2 == 0 }; [1, 2, 3, 4, 5] |> filter isEven",
			[]int64{2, 4},
		},
		{
			"let isLong = fn() { len(🍕) > 1 }; [\"a\", \"ab\", \"abc\"] |> filter isLong",
			[]string{"ab", "abc"},
		},
		
		// 複数パイプラインの連鎖
		{
			"let double = fn() { 🍕 * 2 }; let isEven = fn() { 🍕 % 2 == 0 }; [1, 2, 3, 4, 5] |> filter isEven |> map double",
			[]int64{4, 8},
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		
		switch expected := tt.expected.(type) {
		case []int64:
			testIntegerArray(t, evaluated, expected)
		case []string:
			testStringArray(t, evaluated, expected)
		}
	}
}

// testIntegerArray はオブジェクトが期待する整数配列かをテストする
func testIntegerArray(t *testing.T, obj object.Object, expected []int64) {
	array, ok := obj.(*object.Array)
	if !ok {
		t.Fatalf("object is not Array. got=%T (%+v)", obj, obj)
	}

	if len(array.Elements) != len(expected) {
		t.Fatalf("wrong num of elements. expected=%d, got=%d",
			len(expected), len(array.Elements))
	}

	for i, expectedElem := range expected {
		testIntegerObject(t, array.Elements[i], expectedElem)
	}
}

// testStringArray はオブジェクトが期待する文字列配列かをテストする
func testStringArray(t *testing.T, obj object.Object, expected []string) {
	array, ok := obj.(*object.Array)
	if !ok {
		t.Fatalf("object is not Array. got=%T (%+v)", obj, obj)
	}

	if len(array.Elements) != len(expected) {
		t.Fatalf("wrong num of elements. expected=%d, got=%d",
			len(expected), len(array.Elements))
	}

	for i, expectedElem := range expected {
		testStringObject(t, array.Elements[i], expectedElem)
	}
}

