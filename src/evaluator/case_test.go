package evaluator

import (
	"strings"
	"testing"

	"github.com/uncode/object"
)

// TestCaseStatement はcase文の基本的な機能をテストする
func TestCaseStatement(t *testing.T) {
	tests := []struct {
		input    string
		expected interface{}
	}{
		{
			// 基本的なcase文のテスト
			`
			func test() {
				case 🍕 % 3 == 0:
					"Divisible by 3" >> 💩
				case 🍕 % 5 == 0:
					"Divisible by 5" >> 💩
				case default:
					"Not divisible by 3 or 5" >> 💩
			}
			3 |> test
			`,
			"Divisible by 3",
		},
		{
			// 複数条件のテスト
			`
			func test() {
				case 🍕 % 3 == 0:
					"Divisible by 3" >> 💩
				case 🍕 % 5 == 0:
					"Divisible by 5" >> 💩
				case default:
					"Not divisible by 3 or 5" >> 💩
			}
			5 |> test
			`,
			"Divisible by 5",
		},
		{
			// defaultケースのテスト
			`
			func test() {
				case 🍕 % 3 == 0:
					"Divisible by 3" >> 💩
				case 🍕 % 5 == 0:
					"Divisible by 5" >> 💩
				case default:
					"Not divisible by 3 or 5" >> 💩
			}
			7 |> test
			`,
			"Not divisible by 3 or 5",
		},
		{
			// 複雑な条件式のテスト
			`
			func test() {
				case 🍕 > 10 && 🍕 < 20:
					"Between 10 and 20" >> 💩
				case 🍕 <= 10:
					"10 or less" >> 💩
				case default:
					"Greater than 20" >> 💩
			}
			15 |> test
			`,
			"Between 10 and 20",
		},
		{
			// defaultなしのテスト
			`
			func test() {
				case 🍕 > 10:
					"Greater than 10" >> 💩
				case 🍕 < 5:
					"Less than 5" >> 💩
			}
			7 |> test
			`,
			NULL,
		},
		{
			// FizzBuzzテスト
			`
			func fizzbuzz() {
				case 🍕 % 15 == 0:
					"FizzBuzz" >> 💩
				case 🍕 % 3 == 0:
					"Fizz" >> 💩
				case 🍕 % 5 == 0:
					"Buzz" >> 💩
				case default:
					🍕 >> 💩
			}
			15 |> fizzbuzz
			`,
			"FizzBuzz",
		},
	}

	for i, tt := range tests {
		evaluated := testEval(tt.input)
		
		switch expected := tt.expected.(type) {
		case string:
			strObj, ok := evaluated.(*object.String)
			if !ok {
				t.Errorf("テスト %d: object is not String. got=%T (%+v)", i, evaluated, evaluated)
				continue
			}
			if strObj.Value != expected {
				t.Errorf("テスト %d: wrong string value. got=%q, want=%q", i, strObj.Value, expected)
			}
		case int:
			testIntegerObject(t, evaluated, int64(expected))
		case bool:
			testBooleanObject(t, evaluated, expected)
		case nil:
			if evaluated != NULL {
				t.Errorf("テスト %d: object is not NULL. got=%T (%+v)", i, evaluated, evaluated)
			}
		}
	}
}

// TestCaseStatementErrors はcase文のエラーケースをテストする
func TestCaseStatementErrors(t *testing.T) {
	tests := []struct {
		input    string
		expected string // エラーメッセージの一部
	}{
		{
			// 関数外でのcase文使用（エラー）
			`
			case 1 == 1:
				"True" >> 💩
			`,
			"関数ブロック内のトップレベルでのみ使用できます",
		},
		{
			// caseの後に条件式がない
			`
			func test() {
				case:
					"Missing condition" >> 💩
			}
			`,
			"条件式が必要です",
		},
	}

	for i, tt := range tests {
		evaluated := testEval(tt.input)
		
		errObj, ok := evaluated.(*object.Error)
		if !ok {
			t.Errorf("テスト %d: object.Error ではありません。got=%T (%+v)", i, evaluated, evaluated)
			continue
		}
		
		if !strings.Contains(errObj.Message, tt.expected) {
			t.Errorf("テスト %d: 間違ったエラーメッセージ。expected=%q, got=%q", i, tt.expected, errObj.Message)
		}
	}
}

