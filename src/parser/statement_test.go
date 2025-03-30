package parser

import (
	"strings"
	"testing"

	"github.com/uncode/ast"
	"github.com/uncode/lexer"
	"github.com/uncode/logger"
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

// TestCaseStatement はcase文の解析をテストする
func TestCaseStatement(t *testing.T) {
	// ログレベルを設定
	logger.SetLogLevel(logger.LevelDebug)
	
	// 関数内でのcase文使用のテスト
	input := `
	def test() {
		case 🍕 % 3 == 0:
			"Divisible by 3" >> 💩
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
	
	// 関数定義ステートメントを検証
	exprStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}
	
	funcLit, ok := exprStmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("Expression is not ast.FunctionLiteral. got=%T", exprStmt.Expression)
	}
	
	// 関数本体内のcase文を検証
	if len(funcLit.Body.Statements) != 1 {
		t.Fatalf("function body does not contain 1 statement. got=%d", len(funcLit.Body.Statements))
	}
	
	caseStmt, ok := funcLit.Body.Statements[0].(*ast.CaseStatement)
	if !ok {
		t.Fatalf("function body statement is not ast.CaseStatement. got=%T", funcLit.Body.Statements[0])
	}
	
	// case文のトークンを検証
	if caseStmt.TokenLiteral() != "case" {
		t.Errorf("caseStmt.TokenLiteral not %s. got=%s", "case", caseStmt.TokenLiteral())
	}
	
	// 条件式の検証
	if caseStmt.Condition == nil {
		t.Fatalf("caseStmt.Condition is nil")
	}
	
	// Consequence (結果ブロック) の検証
	if caseStmt.Consequence == nil {
		t.Fatalf("caseStmt.Consequence is nil")
	}
	
	if len(caseStmt.Consequence.Statements) != 1 {
		t.Fatalf("case consequence does not contain 1 statement. got=%d", 
			len(caseStmt.Consequence.Statements))
	}
}

// TestDefaultCaseStatement はdefault case文の解析をテストする
func TestDefaultCaseStatement(t *testing.T) {
	// ログレベルを設定
	logger.SetLogLevel(logger.LevelDebug)
	
	// 関数内でのdefault文使用のテスト
	input := `
	def test() {
		default:
			"Default case" >> 💩
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
	
	// 関数定義ステートメントを検証
	exprStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}
	
	funcLit, ok := exprStmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("Expression is not ast.FunctionLiteral. got=%T", exprStmt.Expression)
	}
	
	// 関数本体内のdefault文を検証
	if len(funcLit.Body.Statements) != 1 {
		t.Fatalf("function body does not contain 1 statement. got=%d", len(funcLit.Body.Statements))
	}
	
	defaultStmt, ok := funcLit.Body.Statements[0].(*ast.DefaultCaseStatement)
	if !ok {
		t.Fatalf("function body statement is not ast.DefaultCaseStatement. got=%T", funcLit.Body.Statements[0])
	}
	
	// default文のトークンを検証
	if defaultStmt.TokenLiteral() != "default" {
		t.Errorf("defaultStmt.TokenLiteral not %s. got=%s", "default", defaultStmt.TokenLiteral())
	}
	
	// Body (結果ブロック) の検証
	if defaultStmt.Body == nil {
		t.Fatalf("defaultStmt.Body is nil")
	}
	
	if len(defaultStmt.Body.Statements) != 1 {
		t.Fatalf("default body does not contain 1 statement. got=%d", 
			len(defaultStmt.Body.Statements))
	}
}

// TestCaseDefaultCombination はcase文とdefault文の組み合わせをテストする
func TestCaseDefaultCombination(t *testing.T) {
	// ログレベルを設定
	logger.SetLogLevel(logger.LevelDebug)
	
	// 関数内でcase文とdefault文の組み合わせをテスト
	input := `
	def fizzbuzz() {
		case 🍕 % 15 == 0:
			"FizzBuzz" >> 💩
		case 🍕 % 3 == 0:
			"Fizz" >> 💩
		case 🍕 % 5 == 0:
			"Buzz" >> 💩
		default:
			🍕 >> 💩
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
	
	// 関数定義ステートメントを検証
	exprStmt, ok := program.Statements[0].(*ast.ExpressionStatement)
	if !ok {
		t.Fatalf("program.Statements[0] is not ast.ExpressionStatement. got=%T", program.Statements[0])
	}
	
	funcLit, ok := exprStmt.Expression.(*ast.FunctionLiteral)
	if !ok {
		t.Fatalf("Expression is not ast.FunctionLiteral. got=%T", exprStmt.Expression)
	}
	
	// 関数名の検証
	if funcLit.Name == nil || funcLit.Name.Value != "fizzbuzz" {
		t.Fatalf("Function name is not 'fizzbuzz'. got=%v", funcLit.Name)
	}
	
	// 関数本体内のステートメント数を検証
	if len(funcLit.Body.Statements) != 4 {
		t.Fatalf("function body does not contain 4 statements. got=%d", len(funcLit.Body.Statements))
	}
	
	// 最初の3つのステートメントがcase文であることを検証
	for i := 0; i < 3; i++ {
		_, ok := funcLit.Body.Statements[i].(*ast.CaseStatement)
		if !ok {
			t.Fatalf("function body statement[%d] is not ast.CaseStatement. got=%T", 
				i, funcLit.Body.Statements[i])
		}
	}
	
	// 最後のステートメントがdefault文であることを検証
	_, ok = funcLit.Body.Statements[3].(*ast.DefaultCaseStatement)
	if !ok {
		t.Fatalf("function body statement[3] is not ast.DefaultCaseStatement. got=%T", 
			funcLit.Body.Statements[3])
	}
}

// TestCaseStatementErrors はcase文のエラーケースをテストする
func TestCaseStatementErrors(t *testing.T) {
	// ログレベルを設定
	logger.SetLogLevel(logger.LevelDebug)
	
	// 関数外でのcase文使用のテスト（エラーになるべき）
	input := `
	case 1 == 1:
		"This should not work" >> 💩
	`
	
	l := lexer.NewLexer(input)
	tokens, _ := l.Tokenize()
	p := NewParser(tokens)
	_, err := p.ParseProgram()
	
	// エラーが発生することを確認
	if err == nil {
		t.Fatalf("Expected parser error for case statement outside function, but got nil")
	}
	
	// エラーメッセージに "関数ブロック内" という文言が含まれることを確認
	errors := p.Errors()
	foundError := false
	for _, errMsg := range errors {
		if strings.Contains(errMsg, "関数ブロック") {
			foundError = true
			break
		}
	}
	
	if !foundError {
		t.Errorf("Expected error message about case statement only allowed inside function body, got: %v", errors)
	}
}

// TestDefaultCaseStatementErrors はdefault文のエラーケースをテストする
func TestDefaultCaseStatementErrors(t *testing.T) {
	// ログレベルを設定
	logger.SetLogLevel(logger.LevelDebug)
	
	// 関数外でのdefault文使用のテスト（エラーになるべき）
	input := `
	default:
		"This should not work" >> 💩
	`
	
	l := lexer.NewLexer(input)
	tokens, _ := l.Tokenize()
	p := NewParser(tokens)
	_, err := p.ParseProgram()
	
	// エラーが発生することを確認
	if err == nil {
		t.Fatalf("Expected parser error for default statement outside function, but got nil")
	}
	
	// エラーメッセージに "関数ブロック内" という文言が含まれることを確認
	errors := p.Errors()
	foundError := false
	for _, errMsg := range errors {
		if strings.Contains(errMsg, "関数ブロック") {
			foundError = true
			break
		}
	}
	
	if !foundError {
		t.Errorf("Expected error message about default statement only allowed inside function body, got: %v", errors)
	}
}
