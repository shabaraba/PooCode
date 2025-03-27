package evaluator

import (
	"testing"
	
	"github.com/uncode/object"
)

func TestErrorHandling(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"5 + true;",
			"型の不一致: INTEGER + BOOLEAN",
		},
		{
			"5 + true; 5;",
			"型の不一致: INTEGER + BOOLEAN",
		},
		{
			"-true",
			"未知の演算子: -BOOLEAN",
		},
		{
			"true + false;",
			"未知の演算子: BOOLEAN + BOOLEAN",
		},
		{
			"5; true + false; 5",
			"未知の演算子: BOOLEAN + BOOLEAN",
		},
		{
			"if (10 > 1) { true + false; }",
			"未知の演算子: BOOLEAN + BOOLEAN",
		},
		{
			`
			if (10 > 1) {
				if (10 > 1) {
					return true + false;
				}
				return 1;
			}
			`,
			"未知の演算子: BOOLEAN + BOOLEAN",
		},
		{
			"foobar",
			"識別子が見つかりません: foobar",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned. got=%T(%+v)",
				evaluated, evaluated)
			continue
		}

		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message. expected=%q, got=%q",
				tt.expectedMessage, errObj.Message)
		}
	}
}

// TestIndexAccessErrors は配列インデックスアクセスに関するエラーをテストします
func TestIndexAccessErrors(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"let a = [1, 2, 3]; a[10]",
			"インデックスが範囲外です: インデックス=10, 長さ=3",
		},
		{
			"let a = [1, 2, 3]; a[-1]",
			"インデックスが不正です: -1",
		},
		{
			"let a = 5; a[0]",
			"インデックス演算子はハッシュまたは配列にのみ使用できます: INTEGER",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned for index access. got=%T(%+v)",
				evaluated, evaluated)
			continue
		}

		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message. expected=%q, got=%q",
				tt.expectedMessage, errObj.Message)
		}
	}
}

// TestTypeCheckErrors は型チェック関連のエラーをテストします
func TestTypeCheckErrors(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"fn(x: int) { x * 2 } |> \"hello\"",
			"🍕の型が不正です: 期待=int, 実際=str",
		},
		{
			"fn(x: int) -> str { x * 2 }(5)",
			"💩の型が不正です: 期待=str, 実際=int",
		},
		{
			"fn process(data: array) { data[0] } |> 5",
			"🍕の型が不正です: 期待=array, 実際=int",
		},
		{
			"fn isEmpty(s: str) -> bool { s == \"\" } |> 42",
			"🍕の型が不正です: 期待=str, 実際=int",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned for type check. got=%T(%+v)",
				evaluated, evaluated)
			continue
		}

		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message. expected=%q, got=%q",
				tt.expectedMessage, errObj.Message)
		}
	}
}

// TestFunctionCallErrors は関数呼び出し関連のエラーをテストします
func TestFunctionCallErrors(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"fn add(a, b, c) { a + b + c }; add(1, 2)",
			"引数の数が一致しません: 期待=3, 実際=2",
		},
		{
			"fn() { 5 }(1, 2, 3)",
			"引数の数が一致しません: 期待=0, 実際=3",
		},
		{
			"5()",
			"関数ではありません: INTEGER",
		},
		{
			"let x = 10; x(5)",
			"関数ではありません: INTEGER",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned for function call. got=%T(%+v)",
				evaluated, evaluated)
			continue
		}

		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message. expected=%q, got=%q",
				tt.expectedMessage, errObj.Message)
		}
	}
}

// TestDivisionByZeroError はゼロ除算エラーをテストします
func TestDivisionByZeroError(t *testing.T) {
	input := "10 / 0"
	expected := "ゼロ除算エラー: 0で割ることはできません"

	evaluated := testEval(input)

	errObj, ok := evaluated.(*object.Error)
	if !ok {
		t.Errorf("no error object returned for division by zero. got=%T(%+v)",
			evaluated, evaluated)
		return
	}

	if errObj.Message != expected {
		t.Errorf("wrong error message. expected=%q, got=%q",
			expected, errObj.Message)
	}
}

// TestPropertyAccessErrors はプロパティアクセスエラーをテストします
func TestPropertyAccessErrors(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"let a = null; a.something",
			"プロパティアクセスエラー: NULL型にはプロパティがありません",
		},
		{
			"5.length",
			"プロパティアクセスエラー: INTEGER型にはプロパティがありません",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)

		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned for property access. got=%T(%+v)",
				evaluated, evaluated)
			continue
		}

		if errObj.Message != tt.expectedMessage {
			t.Errorf("wrong error message. expected=%q, got=%q",
				tt.expectedMessage, errObj.Message)
		}
	}
}
