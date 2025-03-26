package parser

import (
	"testing"

	"github.com/uncode/ast"
	"github.com/uncode/lexer"
)

// TestPipelineExpression はパイプライン式の解析をテストする
func TestPipelineExpression(t *testing.T) {
	tests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 |> add;", 5, "|>", "add"},
		{"x |> print;", "x", "|>", "print"},
		{"getData() |> process;", "getData()", "|>", "process"},
		{"[1, 2, 3] |> map;", "[1, 2, 3]", "|>", "map"},
	}

	for _, tt := range tests {
		l := lexer.NewLexer(tt.input)
		tokens, _ := l.Tokenize()
		p := NewParser(tokens)
		program, err := p.ParseProgram()

		if err != nil {
			t.Fatalf("Parser error: %v", err)
		}

		if len(program.Statements) != 1 {
			t.Fatalf("Program has wrong number of statements. got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}

		exp, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("exp is not ast.InfixExpression. got=%T", stmt.Expression)
		}

		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s", tt.operator, exp.Operator)
		}

		// 左辺値のテスト
		switch left := tt.leftValue.(type) {
		case int:
			testIntegerLiteral(t, exp.Left, int64(left))
		case string:
			// 文字列がリテラルか式かを判断
			if left == "[1, 2, 3]" {
				_, ok := exp.Left.(*ast.ArrayLiteral)
				if !ok {
					t.Fatalf("exp.Left is not ast.ArrayLiteral. got=%T", exp.Left)
				}
			} else if left == "getData()" {
				_, ok := exp.Left.(*ast.CallExpression)
				if !ok {
					t.Fatalf("exp.Left is not ast.CallExpression. got=%T", exp.Left)
				}
			} else {
				testIdentifier(t, exp.Left, left)
			}
		}

		// 右辺値のテスト
		if rightIdent, ok := tt.rightValue.(string); ok {
			testIdentifier(t, exp.Right, rightIdent)
		}
	}
}

// TestPipelineWithArgument はパイプラインで引数を取る関数へのパイプをテストする
func TestPipelineWithArgument(t *testing.T) {
	input := "5 |> add 3;"

	l := lexer.NewLexer(input)
	tokens, _ := l.Tokenize()
	p := NewParser(tokens)
	program, err := p.ParseProgram()

	if err != nil {
		t.Fatalf("Parser error: %v", err)
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	pipeExp, ok := stmt.Expression.(*ast.InfixExpression)
	if !ok {
		t.Fatalf("exp is not ast.InfixExpression. got=%T", stmt.Expression)
	}

	if pipeExp.Operator != "|>" {
		t.Fatalf("exp.Operator is not '|>'. got=%s", pipeExp.Operator)
	}

	// 左辺は整数5
	testIntegerLiteral(t, pipeExp.Left, 5)

	// 右辺は関数呼び出し add(3)
	rightCall, ok := pipeExp.Right.(*ast.CallExpression)
	if !ok {
		t.Fatalf("pipeExp.Right is not ast.CallExpression. got=%T", pipeExp.Right)
	}

	if !testIdentifier(t, rightCall.Function, "add") {
		return
	}

	if len(rightCall.Arguments) != 1 {
		t.Fatalf("add関数の引数の数が正しくありません。期待値=1, 実際=%d", len(rightCall.Arguments))
	}

	testIntegerLiteral(t, rightCall.Arguments[0], 3)
}

// TestParallelPipeExpression は並列パイプ演算子の解析をテストする
func TestParallelPipeExpression(t *testing.T) {
	input := "data | process;"

	l := lexer.NewLexer(input)
	tokens, _ := l.Tokenize()
	p := NewParser(tokens)
	program, err := p.ParseProgram()

	if err != nil {
		t.Fatalf("Parser error: %v", err)
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	pipeExp, ok := stmt.Expression.(*ast.InfixExpression)
	if !ok {
		t.Fatalf("exp is not ast.InfixExpression. got=%T", stmt.Expression)
	}

	if pipeExp.Operator != "|" {
		t.Fatalf("exp.Operator is not '|'. got=%s", pipeExp.Operator)
	}

	// 左辺と右辺が正しく解析されていることを確認
	if !testIdentifier(t, pipeExp.Left, "data") {
		return
	}

	if !testIdentifier(t, pipeExp.Right, "process") {
		return
	}
}

// TestChainedPipelineExpression は連鎖したパイプライン式の解析をテストする
func TestChainedPipelineExpression(t *testing.T) {
	input := "data |> process |> display;"

	l := lexer.NewLexer(input)
	tokens, _ := l.Tokenize()
	p := NewParser(tokens)
	program, err := p.ParseProgram()

	if err != nil {
		t.Fatalf("Parser error: %v", err)
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	// (data |> process) |> display のような形になるはず
	topPipe, ok := stmt.Expression.(*ast.InfixExpression)
	if !ok {
		t.Fatalf("exp is not ast.InfixExpression. got=%T", stmt.Expression)
	}

	if topPipe.Operator != "|>" {
		t.Fatalf("topPipe.Operator is not '|>'. got=%s", topPipe.Operator)
	}

	// 右側はdisplay
	if !testIdentifier(t, topPipe.Right, "display") {
		return
	}

	// 左側は data |> process
	leftPipe, ok := topPipe.Left.(*ast.InfixExpression)
	if !ok {
		t.Fatalf("topPipe.Left is not ast.InfixExpression. got=%T", topPipe.Left)
	}

	if leftPipe.Operator != "|>" {
		t.Fatalf("leftPipe.Operator is not '|>'. got=%s", leftPipe.Operator)
	}

	if !testIdentifier(t, leftPipe.Left, "data") {
		return
	}

	if !testIdentifier(t, leftPipe.Right, "process") {
		return
	}
}

// TestMixedPipeExpressions は異なる種類のパイプを混合した式のテスト
func TestMixedPipeExpressions(t *testing.T) {
	input := "data |> process | display;"

	l := lexer.NewLexer(input)
	tokens, _ := l.Tokenize()
	p := NewParser(tokens)
	program, err := p.ParseProgram()

	if err != nil {
		t.Fatalf("Parser error: %v", err)
	}

	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}

	// (data |> process) | display のような形になるはず
	topPipe, ok := stmt.Expression.(*ast.InfixExpression)
	if !ok {
		t.Fatalf("exp is not ast.InfixExpression. got=%T", stmt.Expression)
	}

	if topPipe.Operator != "|" {
		t.Fatalf("topPipe.Operator is not '|'. got=%s", topPipe.Operator)
	}

	// 右側はdisplay
	if !testIdentifier(t, topPipe.Right, "display") {
		return
	}

	// 左側は data |> process
	leftPipe, ok := topPipe.Left.(*ast.InfixExpression)
	if !ok {
		t.Fatalf("topPipe.Left is not ast.InfixExpression. got=%T", topPipe.Left)
	}

	if leftPipe.Operator != "|>" {
		t.Fatalf("leftPipe.Operator is not '|>'. got=%s", leftPipe.Operator)
	}

	if !testIdentifier(t, leftPipe.Left, "data") {
		return
	}

	if !testIdentifier(t, leftPipe.Right, "process") {
		return
	}
}
