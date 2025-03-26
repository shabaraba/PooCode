package parser

import (
	"testing"

	"github.com/uncode/ast"
	"github.com/uncode/lexer"
)

// TestArrayLiteral は配列リテラルの解析をテストする
func TestArrayLiteral(t *testing.T) {
	input := "[1, 2 * 2, 3 + 3];"
	
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
	
	array, ok := stmt.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("exp not *ast.ArrayLiteral. got=%T", stmt.Expression)
	}
	
	if len(array.Elements) != 3 {
		t.Fatalf("len(array.Elements) not 3. got=%d", len(array.Elements))
	}
	
	testIntegerLiteral(t, array.Elements[0], 1)
	testInfixExpression(t, array.Elements[1], 2, "*", 2)
	testInfixExpression(t, array.Elements[2], 3, "+", 3)
}

// TestEmptyArrayLiteral は空の配列リテラルの解析をテストする
func TestEmptyArrayLiteral(t *testing.T) {
	input := "[];"
	
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
	
	array, ok := stmt.Expression.(*ast.ArrayLiteral)
	if !ok {
		t.Fatalf("exp not *ast.ArrayLiteral. got=%T", stmt.Expression)
	}
	
	if len(array.Elements) != 0 {
		t.Fatalf("len(array.Elements) not 0. got=%d", len(array.Elements))
	}
}

// TestCallExpression は関数呼び出し式の解析をテストする
func TestCallExpression(t *testing.T) {
	input := "add(1, 2 * 3, 4 + 5);"
	
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
	
	callExpr, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T", stmt.Expression)
	}
	
	if !testIdentifier(t, callExpr.Function, "add") {
		return
	}
	
	if len(callExpr.Arguments) != 3 {
		t.Fatalf("wrong length of arguments. got=%d", len(callExpr.Arguments))
	}
	
	testIntegerLiteral(t, callExpr.Arguments[0], 1)
	testInfixExpression(t, callExpr.Arguments[1], 2, "*", 3)
	testInfixExpression(t, callExpr.Arguments[2], 4, "+", 5)
}

// TestCallExpressionWithoutParentheses は括弧なしの関数呼び出し式の解析をテストする
func TestCallExpressionWithoutParentheses(t *testing.T) {
	input := "print 42;"
	
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
	
	callExpr, ok := stmt.Expression.(*ast.CallExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.CallExpression. got=%T", stmt.Expression)
	}
	
	if !testIdentifier(t, callExpr.Function, "print") {
		return
	}
	
	if len(callExpr.Arguments) != 1 {
		t.Fatalf("wrong length of arguments. got=%d", len(callExpr.Arguments))
	}
	
	testIntegerLiteral(t, callExpr.Arguments[0], 42)
}

// TestFunctionLiteral は関数リテラルの解析をテストする
func TestFunctionLiteral(t *testing.T) {
	input := "function(x, y) { x + y; };"
	
	l := lexer.NewLexer(input)
	tokens, _ := l.Tokenize()
	p := NewParser(tokens)
	program, err := p.ParseProgram()
	
	if err != nil {
		t.Fatalf("Parser error: %v", err)
	}
	
	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}
	
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}
	
	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.FunctionLiteral. got=%T", stmt.Expression)
	}
	
	if len(function.Parameters) != 2 {
		t.Fatalf("function literal parameters wrong. want 2, got=%d", len(function.Parameters))
	}
	
	testLiteralExpression(t, function.Parameters[0], "x")
	testLiteralExpression(t, function.Parameters[1], "y")
	
	if len(function.Body.Statements) != 1 {
		t.Fatalf("function.Body.Statements has not 1 statement. got=%d", len(function.Body.Statements))
	}
	
	bodyStmt, ok := function.Body.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("function body stmt is not ast.ExpressionStatement. got=%T", function.Body.Statements[0])
	}
	
	testInfixExpression(t, bodyStmt.Expression, "x", "+", "y")
}

// TestFunctionLiteralWithName は名前付き関数リテラルの解析をテストする
func TestFunctionLiteralWithName(t *testing.T) {
	input := "function add(x, y) { x + y; };"
	
	l := lexer.NewLexer(input)
	tokens, _ := l.Tokenize()
	p := NewParser(tokens)
	program, err := p.ParseProgram()
	
	if err != nil {
		t.Fatalf("Parser error: %v", err)
	}
	
	if len(program.Statements) != 1 {
		t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
	}
	
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}
	
	function, ok := stmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.FunctionLiteral. got=%T", stmt.Expression)
	}
	
	if function.Name == nil {
		t.Fatalf("function literal name is nil")
	}
	
	if function.Name.Value != "add" {
		t.Fatalf("function name is not 'add'. got=%q", function.Name.Value)
	}
	
	if len(function.Parameters) != 2 {
		t.Fatalf("function literal parameters wrong. want 2, got=%d", len(function.Parameters))
	}
	
	testLiteralExpression(t, function.Parameters[0], "x")
	testLiteralExpression(t, function.Parameters[1], "y")
	
	if len(function.Body.Statements) != 1 {
		t.Fatalf("function.Body.Statements has not 1 statement. got=%d", len(function.Body.Statements))
	}
}

// TestSpecialLiterals は特殊リテラル（Pizza, Poo）の解析をテストする
func TestSpecialLiterals(t *testing.T) {
	tests := []struct {
		input    string
		expected string
		litType  string
	}{
		{"🍕;", "🍕", "PizzaLiteral"},
		{"💩;", "💩", "PooLiteral"},
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
		
		var literal ast.Expression
		switch tt.litType {
		case "PizzaLiteral":
			lit, ok := stmt.Expression.(*ast.PizzaLiteral)
			if !ok {
				t.Fatalf("Expression is not *ast.PizzaLiteral. got=%T", stmt.Expression)
			}
			literal = lit
		case "PooLiteral":
			lit, ok := stmt.Expression.(*ast.PooLiteral)
			if !ok {
				t.Fatalf("Expression is not *ast.PooLiteral. got=%T", stmt.Expression)
			}
			literal = lit
		}
		
		if literal.String() != tt.expected {
			t.Errorf("literal.String() not %q. got=%q", tt.expected, literal.String())
		}
	}
}

// TestIndexExpression は配列の添字アクセス式の解析をテストする
func TestIndexExpression(t *testing.T) {
	input := "myArray[1 + 1];"
	
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
	
	indexExp, ok := stmt.Expression.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("exp not *ast.IndexExpression. got=%T", stmt.Expression)
	}
	
	if !testIdentifier(t, indexExp.Left, "myArray") {
		return
	}
	
	if !testInfixExpression(t, indexExp.Index, 1, "+", 1) {
		return
	}
}

// TestPropertyAccessExpression はプロパティアクセス式の解析をテストする
func TestPropertyAccessExpression(t *testing.T) {
	tests := []struct {
		input    string
		object   string
		property string
		operator string
	}{
		{"foo.bar;", "foo", "bar", "."},
		{"foo's bar;", "foo", "bar", "'s"},
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
		
		propExp, ok := stmt.Expression.(*ast.PropertyAccessExpression)
		if !ok {
			t.Fatalf("exp not *ast.PropertyAccessExpression. got=%T", stmt.Expression)
		}
		
		if !testIdentifier(t, propExp.Object, tt.object) {
			return
		}
		
		if !testIdentifier(t, propExp.Property, tt.property) {
			return
		}
		
		if propExp.Token.Literal != tt.operator {
			t.Errorf("propExp.Token.Literal not %s. got=%s", tt.operator, propExp.Token.Literal)
		}
	}
}

// Helper functions
func testIdentifier(t *testing.T, exp ast.Expression, value string) bool {
	ident, ok := exp.(*ast.Identifier)
	if !ok {
		t.Errorf("exp not *ast.Identifier. got=%T", exp)
		return false
	}
	
	if ident.Value != value {
		t.Errorf("ident.Value not %s. got=%s", value, ident.Value)
		return false
	}
	
	if ident.TokenLiteral() != value {
		t.Errorf("ident.TokenLiteral not %s. got=%s", value, ident.TokenLiteral())
		return false
	}
	
	return true
}

func testLiteralExpression(t *testing.T, exp ast.Expression, expected interface{}) bool {
	switch v := expected.(type) {
	case int:
		return testIntegerLiteral(t, exp, int64(v))
	case int64:
		return testIntegerLiteral(t, exp, v)
	case string:
		return testIdentifier(t, exp, v)
	case bool:
		return testBooleanLiteral(t, exp, v)
	}
	t.Errorf("type of exp not handled. got=%T", exp)
	return false
}

func testInfixExpression(t *testing.T, exp ast.Expression, left interface{}, operator string, right interface{}) bool {
	opExp, ok := exp.(*ast.InfixExpression)
	if !ok {
		t.Errorf("exp is not ast.InfixExpression. got=%T(%s)", exp, exp)
		return false
	}
	
	if !testLiteralExpression(t, opExp.Left, left) {
		return false
	}
	
	if opExp.Operator != operator {
		t.Errorf("exp.Operator is not '%s'. got=%q", operator, opExp.Operator)
		return false
	}
	
	if !testLiteralExpression(t, opExp.Right, right) {
		return false
	}
	
	return true
}
