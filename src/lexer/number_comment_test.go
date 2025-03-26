package lexer

import (
	"testing"

	"github.com/uncode/token"
)

// TestNumberLiterals は整数と浮動小数点数のリテラルをテストする
func TestNumberLiterals(t *testing.T) {
	input := `42
3.14
99.99
0.5
`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.INT, "42"},
		{token.FLOAT, "3.14"},
		{token.FLOAT, "99.99"},
		{token.FLOAT, "0.5"},
		{token.EOF, ""},
	}

	l := NewLexer(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}

// TestComments はコメント機能をテストする
func TestComments(t *testing.T) {
	input := `// これはコメントです
42 // 数値の後のコメント
// もう一つのコメント
def // キーワードの後のコメント
`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.INT, "42"},
		{token.FUNCTION, "def"},
		{token.EOF, ""},
	}

	l := NewLexer(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}

		if tok.Literal != tt.expectedLiteral {
			t.Fatalf("tests[%d] - literal wrong. expected=%q, got=%q", i, tt.expectedLiteral, tok.Literal)
		}
	}
}

// TestLineColumn は行番号と列番号の計算が正しいかをテストする
func TestLineColumn(t *testing.T) {
	input := `def add(x, y)
{
  x + y;
}
`

	tests := []struct {
		expectedType   token.TokenType
		expectedLine   int
		expectedColumn int
	}{
		{token.FUNCTION, 1, 1}, // 'def' の位置
		{token.IDENT, 1, 5},    // 'add' の位置
		{token.LPAREN, 1, 8},   // '(' の位置
		{token.IDENT, 1, 9},    // 'x' の位置
		{token.COMMA, 1, 10},   // ',' の位置
		{token.IDENT, 1, 12},   // 'y' の位置
		{token.RPAREN, 1, 13},  // ')' の位置
		{token.LBRACE, 2, 1},   // '{' の位置
		{token.IDENT, 3, 3},    // 'x' の位置
		{token.PLUS, 3, 5},     // '+' の位置
		{token.IDENT, 3, 7},    // 'y' の位置
		{token.SEMICOLON, 3, 8}, // ';' の位置
		{token.RBRACE, 4, 1},   // '}' の位置
	}

	l := NewLexer(input)

	for i, tt := range tests {
		tok := l.NextToken()

		if tok.Type != tt.expectedType {
			t.Fatalf("tests[%d] - tokentype wrong. expected=%q, got=%q", i, tt.expectedType, tok.Type)
		}

		if tok.Line != tt.expectedLine {
			t.Fatalf("tests[%d] - line wrong. expected=%d, got=%d", i, tt.expectedLine, tok.Line)
		}

		if tok.Column != tt.expectedColumn {
			t.Fatalf("tests[%d] - column wrong. expected=%d, got=%d", i, tt.expectedColumn, tok.Column)
		}
	}
}
