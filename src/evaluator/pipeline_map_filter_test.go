package evaluator

import (
	"testing"
)

// TestMapOperator は +> 演算子（map）をテストする
func TestMapOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		// 引数なしの+>演算子（シンプルなmap操作）
		{
			"let double = fn() { 🍕 * 2 }; [1, 2, 3] +> double",
			[]int64{2, 4, 6},
		},
		// 引数を持つ関数を使った+>演算子
		{
			"let addNum = fn(n) { 🍕 + n }; [1, 2, 3] +> addNum(10)",
			[]int64{11, 12, 13},
		},
		// 複数の配列操作の組み合わせ
		{
			"let double = fn() { 🍕 * 2 }; [1..5] +> double",
			[]int64{2, 4, 6, 8, 10},
		},
		/* 文字列配列に対する操作は別途テスト
		{
			"let addExclamation = fn() { 🍕 + \"!\" }; [\"hello\", \"world\"] +> addExclamation",
			[]string{"hello!", "world!"},
		},
		*/
		// パイプラインとの組み合わせ
		{
			"let double = fn() { 🍕 * 2 }; [1, 2, 3] |> map double",
			[]int64{2, 4, 6},
		},
		// +>演算子同士の連結
		{
			"let double = fn() { 🍕 * 2 }; let addOne = fn() { 🍕 + 1 }; [1, 2, 3] +> double +> addOne",
			[]int64{3, 5, 7},
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		
		switch expected := tt.expected.(type) {
		case []int64:
			testIntegerArray(t, evaluated, expected)
		}
	}
}

// TestFilterOperator は ?> 演算子（filter）をテストする
func TestFilterOperator(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		// 引数なしのfilter演算子
		{
			"let isEven = fn() { 🍕 % 2 == 0 }; [1, 2, 3, 4, 5] ?> isEven",
			[]int64{2, 4},
		},
		// 引数を持つ関数を使ったfilter
		{
			"let greaterThan = fn(n) { 🍕 > n }; [1, 2, 3, 4, 5] ?> greaterThan(2)",
			[]int64{3, 4, 5},
		},
		// 複数の配列操作の組み合わせ
		{
			"let isEven = fn() { 🍕 % 2 == 0 }; [1..10] ?> isEven",
			[]int64{2, 4, 6, 8, 10},
		},
		/* 文字列配列に対する操作は別途テスト
		{
			"let isLong = fn() { len(🍕) > 3 }; [\"a\", \"ab\", \"abc\", \"abcd\"] ?> isLong",
			[]string{"abcd"},
		},
		*/
		// パイプラインとの組み合わせ
		{
			"let isEven = fn() { 🍕 % 2 == 0 }; [1, 2, 3, 4, 5] |> filter isEven",
			[]int64{2, 4},
		},
		// ?> (filter) と +> (map) の連結
		{
			"let isEven = fn() { 🍕 % 2 == 0 }; let double = fn() { 🍕 * 2 }; [1, 2, 3, 4, 5] ?> isEven +> double",
			[]int64{4, 8},
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		
		switch expected := tt.expected.(type) {
		case []int64:
			testIntegerArray(t, evaluated, expected)
		}
	}
}

// TestMapFilterOperatorsCombined は +> と ?> 演算子の組み合わせをテストする
func TestMapFilterOperatorsCombined(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		// +>と?>の連結
		{
			"let double = fn() { 🍕 * 2 }; let isEven = fn() { 🍕 % 2 == 0 }; [1, 2, 3, 4, 5] +> double ?> isEven",
			[]int64{2, 4, 6, 8, 10},
		},
		// ?>と+>の連結
		{
			"let double = fn() { 🍕 * 2 }; let isEven = fn() { 🍕 % 2 == 0 }; [1, 2, 3, 4, 5] ?> isEven +> double",
			[]int64{4, 8},
		},
		// 複雑な連結
		{
			"let double = fn() { 🍕 * 2 }; let addOne = fn() { 🍕 + 1 }; let isEven = fn() { 🍕 % 2 == 0 }; [1, 2, 3, 4, 5] +> double ?> isEven +> addOne",
			[]int64{5, 9},
		},
		// パイプラインとの組み合わせ
		{
			"let double = fn() { 🍕 * 2 }; let isGreaterThan5 = fn() { 🍕 > 5 }; [1, 2, 3, 4, 5] |> map double |> filter isGreaterThan5",
			[]int64{6, 8, 10},
		},
		// +>, ?>, |>の混合
		{
			"let double = fn() { 🍕 * 2 }; let isEven = fn() { 🍕 % 2 == 0 }; let addOne = fn() { 🍕 + 1 }; [1, 2, 3, 4, 5] +> double ?> isEven |> map addOne",
			[]int64{5, 9},
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		switch expected := tt.expected.(type) {
		case []int64:
			testIntegerArray(t, evaluated, expected)
		}
	}
}

// テスト用ヘルパー関数の参照
// テストヘルパー関数（testIntegerArray, testStringArray）は
// array_features_test.goですでに定義されています