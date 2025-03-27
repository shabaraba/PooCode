package parser

import (
	"testing"

	"github.com/uncode/ast"
	"github.com/uncode/lexer"
)

// TestMapFilterOperators は+>と?>演算子の解析をテストします
func TestMapFilterOperators(t *testing.T) {
	tests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		// map演算子(+>)のテスト
		{"[1, 2, 3] +> double;", "[1, 2, 3]", "+>", "double"},
		{"arr +> map;", "arr", "+>", "map"},
		{"getData() +> process;", "getData()", "+>", "process"},
		
		// filter演算子(?>)のテスト
		{"[1, 2, 3] ?> isEven;", "[1, 2, 3]", "?>", "isEven"},
		{"arr ?> filter;", "arr", "?>", "filter"},
		{"getData() ?> validate;", "getData()", "?>", "validate"},
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

// TestMapFilterWithArguments はmap/filter演算子で引数を取る関数へのパイプをテストします
func TestMapFilterWithArguments(t *testing.T) {
	tests := []struct {
		input     string
		operator  string
		function  string
		argValues []interface{}
	}{
		// 引数付きのmap演算子(+>)
		{"[1, 2, 3] +> addNum(5);", "+>", "addNum", []interface{}{5}},
		{"data +> transform(\"option\", true);", "+>", "transform", []interface{}{"option", true}},
		
		// 引数付きのfilter演算子(?>)
		{"[1, 2, 3] ?> greaterThan(2);", "?>", "greaterThan", []interface{}{2}},
		{"users ?> hasRole(\"admin\");", "?>", "hasRole", []interface{}{"admin"}},
	}

	for _, tt := range tests {
		l := lexer.NewLexer(tt.input)
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

		if pipeExp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s", tt.operator, pipeExp.Operator)
		}

		// 右辺は関数呼び出し
		rightCall, ok := pipeExp.Right.(*ast.CallExpression)
		if !ok {
			t.Fatalf("pipeExp.Right is not ast.CallExpression. got=%T", pipeExp.Right)
		}

		if !testIdentifier(t, rightCall.Function, tt.function) {
			return
		}

		if len(rightCall.Arguments) != len(tt.argValues) {
			t.Fatalf("関数の引数の数が正しくありません。期待値=%d, 実際=%d", len(tt.argValues), len(rightCall.Arguments))
		}

		// 引数のテスト
		for i, expectedValue := range tt.argValues {
			switch expected := expectedValue.(type) {
			case int:
				testIntegerLiteral(t, rightCall.Arguments[i], int64(expected))
			case string:
				testStringLiteral(t, rightCall.Arguments[i], expected)
			case bool:
				testBooleanLiteral(t, rightCall.Arguments[i], expected)
			}
		}
	}
}

// TestChainedMapFilterOperators は連鎖したmap/filter演算子の解析をテストします
func TestChainedMapFilterOperators(t *testing.T) {
	tests := []struct {
		input           string
		operator1       string
		rightValue1     string
		operator2       string
		rightValue2     string
	}{
		// 連鎖したmap/filter演算子
		{"data +> double +> addOne;", "+>", "double", "+>", "addOne"},
		{"data ?> isEven ?> isPositive;", "?>", "isEven", "?>", "isPositive"},
		{"data +> double ?> isEven;", "+>", "double", "?>", "isEven"},
		{"data ?> isEven +> double;", "?>", "isEven", "+>", "double"},
	}

	for _, tt := range tests {
		l := lexer.NewLexer(tt.input)
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

		// (data +> double) +> addOne のような形になるはず
		topPipe, ok := stmt.Expression.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("exp is not ast.InfixExpression. got=%T", stmt.Expression)
		}

		if topPipe.Operator != tt.operator2 {
			t.Fatalf("topPipe.Operator is not '%s'. got=%s", tt.operator2, topPipe.Operator)
		}

		// 右側は第2の関数
		if !testIdentifier(t, topPipe.Right, tt.rightValue2) {
			return
		}

		// 左側は data +> double
		leftPipe, ok := topPipe.Left.(*ast.InfixExpression)
		if !ok {
			t.Fatalf("topPipe.Left is not ast.InfixExpression. got=%T", topPipe.Left)
		}

		if leftPipe.Operator != tt.operator1 {
			t.Fatalf("leftPipe.Operator is not '%s'. got=%s", tt.operator1, leftPipe.Operator)
		}

		if !testIdentifier(t, leftPipe.Left, "data") {
			return
		}

		if !testIdentifier(t, leftPipe.Right, tt.rightValue1) {
			return
		}
	}
}

// testStringLiteral は式が期待する文字列リテラルかをテストする
func testStringLiteral(t *testing.T, exp ast.Expression, expected string) bool {
	strLit, ok := exp.(*ast.StringLiteral)
	if !ok {
		t.Errorf("exp not *ast.StringLiteral. got=%T", exp)
		return false
	}

	if strLit.Value != expected {
		t.Errorf("strLit.Value not %s. got=%s", expected, strLit.Value)
		return false
	}

	return true
}

// testBooleanLiteral は式が期待する真偽値リテラルかをテストする
func testBooleanLiteral(t *testing.T, exp ast.Expression, expected bool) bool {
	bo, ok := exp.(*ast.Boolean)
	if !ok {
		t.Errorf("exp not *ast.Boolean. got=%T", exp)
		return false
	}

	if bo.Value != expected {
		t.Errorf("bo.Value not %t. got=%t", expected, bo.Value)
		return false
	}

	return true
}
