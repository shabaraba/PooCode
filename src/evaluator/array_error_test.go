package evaluator

import (
	"strings"
	"testing"
	
	"github.com/uncode/object"
)

// TestArrayErrorCases は配列関連のエラーケースをテストする
func TestArrayErrorCases(t *testing.T) {
	tests := []struct {
		input           string
		expectedMessage string
	}{
		{
			"[1, 2, 3][\"a\"]",
			"配列のインデックスは整数",
		},
		{
			"[1, 2, 3] |> map(1)",
			"map関数の第2引数は関数",
		},
		{
			"[1, 2, 3] |> filter(fn(x) { x == 1 })",
			"filter関数に渡された関数はパラメーターを取るべきではありません",
		},
		{
			"[1, 2, 3] |> map(fn(x) { x * 2 })",
			"map関数に渡された関数はパラメーターを取るべきではありません",
		},
		{
			"1 |> map(fn() { 🍕 })",
			"map関数の第1引数は配列",
		},
		{
			"[1, 2, 3][100]",
			"インデックスが範囲外",
		},
		{
			"[\"a\"..1]",
			"サポートされていない範囲式の型",
		},
	}

	for _, tt := range tests {
		evaluated := testEval(tt.input)
		
		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("no error object returned. got=%T (%+v)", evaluated, evaluated)
			continue
		}
		
		if !strings.Contains(errObj.Message, tt.expectedMessage) {
			t.Errorf("wrong error message. expected to contain=%q, got=%q",
				tt.expectedMessage, errObj.Message)
		}
	}
}
