package parser

import (
	"fmt"
	"testing"

	"github.com/uncode/ast"
	"github.com/uncode/lexer"
)

// TestIdentifierExpression は識別子の解析をテストする
func TestIdentifierExpression(t *testing.T) {
	input := "foobar;"
	l := lexer.NewLexer(input)
	tokens, _ := l.Tokenize()
	p := NewParser(tokens)
	program, err := p.ParseProgram()
	
	if err != nil {
		t.Fatalf("Parser error: %v", err)
	}
	
	if len(program.Statements) != 1 {
		t.Fatalf("Program has not enough statements. got=%d", len(program.Statements))
	}
	
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}
	
	ident, ok := stmt.Expression.(*ast.Identifier)
	if !ok {
		t.Fatalf("exp not *ast.Identifier. got=%T", stmt.Expression)
	}
	
	if ident.Value != "foobar" {
		t.Errorf("ident.Value not %s. got=%s", "foobar", ident.Value)
	}
	
	if ident.TokenLiteral() != "foobar" {
		t.Errorf("ident.TokenLiteral not %s. got=%s", "foobar", ident.TokenLiteral())
	}
}

// TestIntegerLiteral は整数リテラルの解析をテストする
func TestIntegerLiteral(t *testing.T) {
	input := "5;"
	l := lexer.NewLexer(input)
	tokens, _ := l.Tokenize()
	p := NewParser(tokens)
	program, err := p.ParseProgram()
	
	if err != nil {
		t.Fatalf("Parser error: %v", err)
	}
	
	if len(program.Statements) != 1 {
		t.Fatalf("Program has not enough statements. got=%d", len(program.Statements))
	}
	
	stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}
	
	literal, ok := stmt.Expression.(*ast.IntegerLiteral)
	if !ok {
		t.Fatalf("exp not *ast.IntegerLiteral. got=%T", stmt.Expression)
	}
	
	if literal.Value != 5 {
		t.Errorf("literal.Value not %d. got=%d", 5, literal.Value)
	}
	
	if literal.TokenLiteral() != "5" {
		t.Errorf("literal.TokenLiteral not %s. got=%s", "5", literal.TokenLiteral())
	}
}

// TestStringLiteral は文字列リテラルの解析をテストする
func TestStringLiteral(t *testing.T) {
	input := `"hello world";`
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
	
	literal, ok := stmt.Expression.(*ast.StringLiteral)
	if !ok {
		t.Fatalf("exp not *ast.StringLiteral. got=%T", stmt.Expression)
	}
	
	if literal.Value != "hello world" {
		t.Errorf("literal.Value not %s. got=%s", "hello world", literal.Value)
	}
}

// TestBooleanLiteral はブール値リテラルの解析をテストする
func TestBooleanLiteral(t *testing.T) {
	tests := []struct {
		input    string
		expected bool
	}{
		{"true;", true},
		{"false;", false},
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
		
		literal, ok := stmt.Expression.(*ast.BooleanLiteral)
		if !ok {
			t.Fatalf("exp not *ast.BooleanLiteral. got=%T", stmt.Expression)
		}
		
		if literal.Value != tt.expected {
			t.Errorf("literal.Value not %t. got=%t", tt.expected, literal.Value)
		}
	}
}

// TestPrefixExpressions は前置演算子式の解析をテストする
func TestPrefixExpressions(t *testing.T) {
	prefixTests := []struct {
		input    string
		operator string
		value    interface{}
	}{
		{"!5;", "!", 5},
		{"-15;", "-", 15},
		{"!true;", "!", true},
		{"!false;", "!", false},
	}
	
	for _, tt := range prefixTests {
		l := lexer.NewLexer(tt.input)
		tokens, _ := l.Tokenize()
		p := NewParser(tokens)
		program, err := p.ParseProgram()
		
		if err != nil {
			t.Fatalf("Parser error: %v", err)
		}
		
		if len(program.Statements) != 1 {
			t.Fatalf("Program has not enough statements. got=%d", len(program.Statements))
		}
		
		stmt, ok := program.Statements[0].(*ast.ExpressionStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
		}
		
		exp, ok := stmt.Expression.(*ast.PrefixExpression)
		if !ok {
			t.Fatalf("stmt is not ast.PrefixExpression. got=%T", stmt.Expression)
		}
		
		if exp.Operator != tt.operator {
			t.Fatalf("exp.Operator is not '%s'. got=%s", tt.operator, exp.Operator)
		}
		
		switch v := tt.value.(type) {
		case int:
			testIntegerLiteral(t, exp.Right, int64(v))
		case bool:
			testBooleanLiteral(t, exp.Right, v)
		}
	}
}

// TestInfixExpressions は中置演算子式の解析をテストする
func TestInfixExpressions(t *testing.T) {
	infixTests := []struct {
		input      string
		leftValue  interface{}
		operator   string
		rightValue interface{}
	}{
		{"5 + 5;", 5, "+", 5},
		{"5 - 5;", 5, "-", 5},
		{"5 * 5;", 5, "*", 5},
		{"5 / 5;", 5, "/", 5},
		{"5 > 5;", 5, ">", 5},
		{"5 < 5;", 5, "<", 5},
		{"5 == 5;", 5, "==", 5},
		{"5 != 5;", 5, "!=", 5},
		{"true == true", true, "==", true},
		{"true != false", true, "!=", false},
		{"false == false", false, "==", false},
	}
	
	for _, tt := range infixTests {
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
		
		switch left := tt.leftValue.(type) {
		case int:
			testIntegerLiteral(t, exp.Left, int64(left))
		case bool:
			testBooleanLiteral(t, exp.Left, left)
		}
		
		switch right := tt.rightValue.(type) {
		case int:
			testIntegerLiteral(t, exp.Right, int64(right))
		case bool:
			testBooleanLiteral(t, exp.Right, right)
		}
	}
}

// TestGroupedExpressions はグループ化された式の解析をテストする
func TestGroupedExpressions(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"(5 + 5) * 2",
			"((5 + 5) * 2)",
		},
		{
			"2 * (5 + 5)",
			"(2 * (5 + 5))",
		},
		{
			"(5 + 5) * 2 * (5 + 5)",
			"(((5 + 5) * 2) * (5 + 5))",
		},
		{
			"-(5 + 5)",
			"(-(5 + 5))",
		},
		{
			"!(true == true)",
			"(!(true == true))",
		},
	}
	
	for i, tt := range tests {
		l := lexer.NewLexer(tt.input)
		tokens, _ := l.Tokenize()
		p := NewParser(tokens)
		program, err := p.ParseProgram()
		
		if err != nil {
			t.Fatalf("Test %d: Parser error: %v", i, err)
		}
		
		actual := program.String()
		if actual != tt.expected {
			t.Errorf("Test %d: expected=%q, got=%q", i, tt.expected, actual)
		}
	}
}

// Helper functions
func testIntegerLiteral(t *testing.T, il ast.Expression, value int64) bool {
	integ, ok := il.(*ast.IntegerLiteral)
	if !ok {
		t.Errorf("il not *ast.IntegerLiteral. got=%T", il)
		return false
	}
	
	if integ.Value != value {
		t.Errorf("integ.Value not %d. got=%d", value, integ.Value)
		return false
	}
	
	if integ.TokenLiteral() != fmt.Sprintf("%d", value) {
		t.Errorf("integ.TokenLiteral not %d. got=%s", value, integ.TokenLiteral())
		return false
	}
	
	return true
}

func testBooleanLiteral(t *testing.T, exp ast.Expression, value bool) bool {
	bo, ok := exp.(*ast.BooleanLiteral)
	if !ok {
		t.Errorf("exp not *ast.BooleanLiteral. got=%T", exp)
		return false
	}
	
	if bo.Value != value {
		t.Errorf("bo.Value not %t. got=%t", value, bo.Value)
		return false
	}
	
	if bo.TokenLiteral() != fmt.Sprintf("%t", value) {
		t.Errorf("bo.TokenLiteral not %t. got=%s", value, bo.TokenLiteral())
		return false
	}
	
	return true
}
