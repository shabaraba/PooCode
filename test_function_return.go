package main

import (
	"fmt"
	
	"github.com/uncode/evaluator"
	"github.com/uncode/lexer"
	"github.com/uncode/object"
	"github.com/uncode/parser"
)

// 関数の戻り値処理テスト
func main() {
	// テストケース
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name: "💩による明示的な戻り値",
			input: `
				func add(x, y) {
					x + y >> 💩
				}
				add(5, 3)
			`,
			expected: "8",
		},
		{
			name: "最後の式が暗黙の戻り値",
			input: `
				func multiply(x, y) {
					x * y
				}
				multiply(4, 5)
			`,
			expected: "20",
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
			expected: "10",
		},
	}

	// 各テストケースを実行
	for _, tt := range tests {
		fmt.Printf("=== テスト: %s ===\n", tt.name)
		
		env := object.NewEnvironment()
		l := lexer.New(tt.input)
		p := parser.New(l)

		program := p.ParseProgram()
		if len(p.Errors()) != 0 {
			fmt.Printf("パースエラー:\n%v\n", p.Errors())
			continue
		}

		evaluated := evaluator.Eval(program, env)
		if evaluated.Inspect() == tt.expected {
			fmt.Printf("成功: %s\n", evaluated.Inspect())
		} else {
			fmt.Printf("失敗: 期待値=%s, 実際=%s\n", tt.expected, evaluated.Inspect())
		}
		fmt.Println()
	}
}
