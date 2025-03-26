package parser

import (
	"testing"

	"github.com/uncode/lexer"
)

// TestBasicParsing は基本的なパース機能をテストする
func TestBasicParsing(t *testing.T) {
	input := `
def x = 5;
def y = 10;
def sum = x + y;
`

	l := lexer.NewLexer(input)
	tokens, err := l.Tokenize()
	if err != nil {
		t.Fatalf("Tokenization error: %v", err)
	}

	p := NewParser(tokens)
	program, err := p.ParseProgram()
	if err != nil {
		t.Fatalf("Parser error: %v", err)
	}

	if program == nil {
		t.Fatalf("ParseProgram() returned nil")
	}

	if len(program.Statements) != 3 {
		t.Fatalf("program.Statements does not contain 3 statements. got=%d", len(program.Statements))
	}

	// プログラムの検証は具体的すぎるので、ここではステートメントの数だけ確認
}

// TestNextToken はパーサーのnextToken関数をテストする
func TestNextToken(t *testing.T) {
	input := "def x = 5;"
	l := lexer.NewLexer(input)
	tokens, _ := l.Tokenize()
	
	p := NewParser(tokens)
	
	// 初期状態
	if p.curToken.Type != "def" {
		t.Errorf("Initial curToken should be 'def', got %q", p.curToken.Type)
	}
	
	if p.peekToken.Type != "IDENT" {
		t.Errorf("Initial peekToken should be 'IDENT', got %q", p.peekToken.Type)
	}
	
	// 1つ進める
	p.nextToken()
	if p.curToken.Type != "IDENT" || p.curToken.Literal != "x" {
		t.Errorf("After first nextToken, curToken should be 'x', got %q", p.curToken.Literal)
	}
	
	// 最後まで進める
	for i := 0; i < 3; i++ {
		p.nextToken()
	}
	
	// EOF確認
	p.nextToken()
	if p.curToken.Type != "EOF" {
		t.Errorf("After consuming all tokens, curToken should be 'EOF', got %q", p.curToken.Type)
	}
}

// TestTokenValidation はトークン検証メソッドをテストする
func TestTokenValidation(t *testing.T) {
	input := "def x = 5;"
	l := lexer.NewLexer(input)
	tokens, _ := l.Tokenize()
	
	p := NewParser(tokens)
	
	// curTokenIs
	if !p.curTokenIs("def") {
		t.Errorf("curTokenIs should return true for 'def', got false")
	}
	
	if p.curTokenIs("IDENT") {
		t.Errorf("curTokenIs should return false for 'IDENT', got true")
	}
	
	// peekTokenIs
	if !p.peekTokenIs("IDENT") {
		t.Errorf("peekTokenIs should return true for 'IDENT', got false")
	}
	
	if p.peekTokenIs("=") {
		t.Errorf("peekTokenIs should return false for '=', got true")
	}
	
	// expectPeek
	if !p.expectPeek("IDENT") {
		t.Errorf("expectPeek should return true and advance token for 'IDENT'")
	}
	
	if p.curToken.Type != "IDENT" || p.curToken.Literal != "x" {
		t.Errorf("After expectPeek, curToken should be 'x', got %q", p.curToken.Literal)
	}
}

// TestParseErrors はパースエラー処理をテストする
func TestParseErrors(t *testing.T) {
	input := "def x 5;" // '=' が欠けている
	l := lexer.NewLexer(input)
	tokens, _ := l.Tokenize()
	
	p := NewParser(tokens)
	_, err := p.ParseProgram()
	
	if err == nil {
		t.Errorf("ParseProgram should return error for invalid syntax, got nil")
	}
	
	errors := p.Errors()
	if len(errors) == 0 {
		t.Errorf("Parser should report errors for invalid syntax, got none")
	}
}

// TestPrecedence は演算子の優先順位が正しく処理されるかテストする
func TestPrecedence(t *testing.T) {
	tests := []struct {
		input    string
		expected string
	}{
		{
			"a + b * c",
			"(a + (b * c))",
		},
		{
			"a * b + c",
			"((a * b) + c)",
		},
		{
			"a + b + c",
			"((a + b) + c)",
		},
		{
			"a + b - c",
			"((a + b) - c)",
		},
		{
			"a * b * c",
			"((a * b) * c)",
		},
		{
			"a * b / c",
			"((a * b) / c)",
		},
		{
			"a + b / c",
			"(a + (b / c))",
		},
		{
			"a + b * c + d / e - f",
			"(((a + (b * c)) + (d / e)) - f)",
		},
		{
			"3 + 4 * 5 == 3 * 1 + 4 * 5",
			"((3 + (4 * 5)) == ((3 * 1) + (4 * 5)))",
		},
		{
			"5 > 4 == 3 < 4",
			"((5 > 4) == (3 < 4))",
		},
		{
			"5 < 4 != 3 > 4",
			"((5 < 4) != (3 > 4))",
		},
		{
			"3 + 4 * 5 * 6",
			"(3 + ((4 * 5) * 6))",
		},
		{
			"3 * 1 + 4 * 5",
			"((3 * 1) + (4 * 5))",
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
