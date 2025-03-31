package parser

import (
	"testing"
	
	"github.com/uncode/ast"
	"github.com/uncode/lexer"
	"github.com/uncode/logger"
)

// TestNestedCaseStatement はネストされたブロック内のcase文の解析をテストする
func TestNestedCaseStatement(t *testing.T) {
	// グローバルのログレベルを設定
	logger.SetLevel(logger.LevelDebug)
	// パーサーデバッグログを有効化
	logger.SetSpecialLevelEnabled(logger.LevelParserDebug, true)
	
	// ネストされたブロック内でのcase文使用のテスト
	input := `
	def testFunc() {
		let x = 10
		{
			case x > 5:
				"greater than 5" >> 💩
			case default:
				"less than or equal to 5" >> 💩
		}
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
	
	// 関数名を検証
	if funcLit.Name == nil || funcLit.Name.Value != "testFunc" {
		t.Fatalf("function name is not 'testFunc'. got=%v", funcLit.Name)
	}
	
	// 関数本体内のブロック文を検証
	if len(funcLit.Body.Statements) != 2 {
		t.Fatalf("function body does not contain 2 statements. got=%d", len(funcLit.Body.Statements))
	}
	
	// 2番目のステートメントがブロック文であることを検証
	blockStmt, ok := funcLit.Body.Statements[1].(*ast.BlockStatement)
	if !ok {
		t.Fatalf("function body statement[1] is not ast.BlockStatement. got=%T", funcLit.Body.Statements[1])
	}
	
	// ネストされたブロック内のcase文を検証
	if len(blockStmt.Statements) != 2 {
		t.Fatalf("nested block does not contain 2 statements. got=%d", len(blockStmt.Statements))
	}
	
	// 1番目のステートメントがcase文であることを検証
	caseStmt, ok := blockStmt.Statements[0].(*ast.CaseStatement)
	if !ok {
		t.Fatalf("nested block statement[0] is not ast.CaseStatement. got=%T", blockStmt.Statements[0])
	}
	
	// case文のConditionフィールドを検証
	if caseStmt.Condition == nil {
		t.Fatalf("case statement condition is nil")
	}
	
	// case文のConsequenceフィールドを検証
	if caseStmt.Consequence == nil {
		t.Fatalf("case statement consequence is nil")
	}
	
	// 2番目のステートメントがdefault case文であることを検証
	defaultStmt, ok := blockStmt.Statements[1].(*ast.DefaultCaseStatement)
	if !ok {
		t.Fatalf("nested block statement[1] is not ast.DefaultCaseStatement. got=%T", blockStmt.Statements[1])
	}
	
	// default case文のBodyフィールドを検証
	if defaultStmt.Body == nil {
		t.Fatalf("default case statement body is nil")
	}
}

// TestNestedCaseInBlockStatement はさらに深くネストされたブロック内のcase文の解析をテストする
func TestNestedCaseInBlockStatement(t *testing.T) {
	// グローバルのログレベルを設定
	logger.SetLevel(logger.LevelDebug)
	// パーサーデバッグログを有効化
	logger.SetSpecialLevelEnabled(logger.LevelParserDebug, true)
	
	// より複雑なネスト構造を持つcase文使用のテスト
	input := `
	def complexFunc() {
		let x = 10
		{
			let y = 20
			{
				case x + y > 25:
					"sum greater than 25" >> 💩
				case default:
					"sum less than or equal to 25" >> 💩
			}
		}
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
	
	// 関数名を検証
	if funcLit.Name == nil || funcLit.Name.Value != "complexFunc" {
		t.Fatalf("function name is not 'complexFunc'. got=%v", funcLit.Name)
	}
	
	// 関数本体内のブロック文を検証（1段目のネスト）
	if len(funcLit.Body.Statements) != 2 {
		t.Fatalf("function body does not contain 2 statements. got=%d", len(funcLit.Body.Statements))
	}
	
	// 2番目のステートメントがブロック文であることを検証
	outerBlock, ok := funcLit.Body.Statements[1].(*ast.BlockStatement)
	if !ok {
		t.Fatalf("function body statement[1] is not ast.BlockStatement. got=%T", funcLit.Body.Statements[1])
	}
	
	// 外側ブロック内の文を検証（2段目のネスト）
	if len(outerBlock.Statements) != 2 {
		t.Fatalf("outer block does not contain 2 statements. got=%d", len(outerBlock.Statements))
	}
	
	// 2番目のステートメントがブロック文であることを検証
	innerBlock, ok := outerBlock.Statements[1].(*ast.BlockStatement)
	if !ok {
		t.Fatalf("outer block statement[1] is not ast.BlockStatement. got=%T", outerBlock.Statements[1])
	}
	
	// 内側ブロック内のcase文を検証
	if len(innerBlock.Statements) != 2 {
		t.Fatalf("inner block does not contain 2 statements. got=%d", len(innerBlock.Statements))
	}
	
	// 1番目のステートメントがcase文であることを検証
	caseStmt, ok := innerBlock.Statements[0].(*ast.CaseStatement)
	if !ok {
		t.Fatalf("inner block statement[0] is not ast.CaseStatement. got=%T", innerBlock.Statements[0])
	}
	
	// case文のConditionフィールドを検証
	if caseStmt.Condition == nil {
		t.Fatalf("case statement condition is nil")
	}
	
	// case文のConsequenceフィールドを検証
	if caseStmt.Consequence == nil {
		t.Fatalf("case statement consequence is nil")
	}
	
	// 2番目のステートメントがdefault case文であることを検証
	defaultStmt, ok := innerBlock.Statements[1].(*ast.DefaultCaseStatement)
	if !ok {
		t.Fatalf("inner block statement[1] is not ast.DefaultCaseStatement. got=%T", innerBlock.Statements[1])
	}
	
	// default case文のBodyフィールドを検証
	if defaultStmt.Body == nil {
		t.Fatalf("default case statement body is nil")
	}
}
