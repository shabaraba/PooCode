package parser

import (
	"testing"

	"github.com/uncode/ast"
	"github.com/uncode/lexer"
)

// TestGlobalStatement はグローバル変数宣言の解析をテストする
func TestGlobalStatement(t *testing.T) {
	tests := []struct {
		input    string
		name     string
		dataType string
	}{
		{"global x;", "x", ""},
		{"global int y;", "y", "int"},
		{"global string name;", "name", "string"},
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
			t.Fatalf("program.Statements does not contain 1 statement. got=%d", len(program.Statements))
		}

		stmt, ok := program.Statements[0].(*ast.GlobalStatement)
		if !ok {
			t.Fatalf("program.Statements[0] is not ast.GlobalStatement. got=%T", program.Statements[0])
		}

		if stmt.Name.Value != tt.name {
			t.Errorf("stmt.Name.Value not '%s'. got=%s", tt.name, stmt.Name.Value)
		}

		if stmt.Type != tt.dataType {
			t.Errorf("stmt.Type not '%s'. got=%s", tt.dataType, stmt.Type)
		}

		if stmt.TokenLiteral() != "global" {
			t.Errorf("stmt.TokenLiteral not 'global'. got=%s", stmt.TokenLiteral())
		}
	}
}

// TestBlockStatement はブロック文の解析をテストする
func TestBlockStatement(t *testing.T) {
	input := `
	{
		x;
		y;
	}
	`

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

	blockStmt, ok := program.Statements[0].(*ast.BlockStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.BlockStatement. got=%T", program.Statements[0])
	}

	if len(blockStmt.Statements) != 2 {
		t.Fatalf("blockStmt.Statements does not contain 2 statements. got=%d", len(blockStmt.Statements))
	}

	// ブロック内のステートメントを検証
	testIdentifier(t, blockStmt.Statements[0].(*ast.ExpressionStatement).Expression, "x")
	testIdentifier(t, blockStmt.Statements[1].(*ast.ExpressionStatement).Expression, "y")
}

// TestAssignStatement は代入文の解析をテストする
func TestAssignStatement(t *testing.T) {
	input := "x >> 5;"

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

	exp, ok := stmt.Expression.(*ast.InfixExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.InfixExpression. got=%T", stmt.Expression)
	}

	if !testIdentifier(t, exp.Left, "x") {
		return
	}

	if exp.Operator != ">>" {
		t.Fatalf("exp.Operator is not '>>'. got=%s", exp.Operator)
	}

	if !testLiteralExpression(t, exp.Right, 5) {
		return
	}
}

// TestEqualAssignStatement は等号を使った代入文の解析をテストする
func TestEqualAssignStatement(t *testing.T) {
	input := "x = 5;"

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

	exp, ok := stmt.Expression.(*ast.InfixExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.InfixExpression. got=%T", stmt.Expression)
	}

	if !testIdentifier(t, exp.Left, "x") {
		return
	}

	if exp.Operator != "=" {
		t.Fatalf("exp.Operator is not '='. got=%s", exp.Operator)
	}

	if !testLiteralExpression(t, exp.Right, 5) {
		return
	}
}

// TestPropertyAssignment はプロパティへの代入をテストする
func TestPropertyAssignment(t *testing.T) {
	input := "obj.prop = 42;"

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

	assignExp, ok := stmt.Expression.(*ast.InfixExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.InfixExpression. got=%T", stmt.Expression)
	}

	propExp, ok := assignExp.Left.(*ast.PropertyAccessExpression)
	if !ok {
		t.Fatalf("assignExp.Left is not ast.PropertyAccessExpression. got=%T", assignExp.Left)
	}

	if !testIdentifier(t, propExp.Object, "obj") {
		return
	}

	if !testIdentifier(t, propExp.Property, "prop") {
		return
	}

	if assignExp.Operator != "=" {
		t.Fatalf("assignExp.Operator is not '='. got=%s", assignExp.Operator)
	}

	if !testLiteralExpression(t, assignExp.Right, 42) {
		return
	}
}

// TestIndexAssignment は配列要素への代入をテストする
func TestIndexAssignment(t *testing.T) {
	input := "arr[0] = 99;"

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

	assignExp, ok := stmt.Expression.(*ast.InfixExpression)
	if !ok {
		t.Fatalf("stmt.Expression is not ast.InfixExpression. got=%T", stmt.Expression)
	}

	indexExp, ok := assignExp.Left.(*ast.IndexExpression)
	if !ok {
		t.Fatalf("assignExp.Left is not ast.IndexExpression. got=%T", assignExp.Left)
	}

	if !testIdentifier(t, indexExp.Left, "arr") {
		return
	}

	if !testLiteralExpression(t, indexExp.Index, 0) {
		return
	}

	if assignExp.Operator != "=" {
		t.Fatalf("assignExp.Operator is not '='. got=%s", assignExp.Operator)
	}

	if !testLiteralExpression(t, assignExp.Right, 99) {
		return
	}
}
