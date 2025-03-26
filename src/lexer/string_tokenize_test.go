package lexer

import (
	"testing"

	"github.com/uncode/token"
)

// TestStringLiterals は文字列リテラルをテストする
func TestStringLiterals(t *testing.T) {
	input := `"hello world"
"テスト文字列"
"special chars: !@#$%^&*()"
"empty string follows"
""
`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.STRING, "hello world"},
		{token.STRING, "テスト文字列"},
		{token.STRING, "special chars: !@#$%^&*()"},
		{token.STRING, "empty string follows"},
		{token.STRING, ""},
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

// TestTokenize はトークナイザー全体の機能をテストする
func TestTokenize(t *testing.T) {
	input := `def add(x, y) {
  x + y;
}
`

	expectedTokenCount := 14 // 実際のトークン数 + EOF

	tokens, err := NewLexer(input).Tokenize()
	if err != nil {
		t.Fatalf("Error during tokenization: %v", err)
	}

	if len(tokens) != expectedTokenCount {
		t.Fatalf("Expected %d tokens, got %d", expectedTokenCount, len(tokens))
	}

	// 最初と最後のトークンを検証
	if tokens[0].Type != token.FUNCTION {
		t.Fatalf("First token should be FUNCTION, got %q", tokens[0].Type)
	}

	if tokens[len(tokens)-1].Type != token.EOF {
		t.Fatalf("Last token should be EOF, got %q", tokens[len(tokens)-1].Type)
	}
}

// TestSpecialCharacters は特殊文字や絵文字をテストする
func TestSpecialCharacters(t *testing.T) {
	input := `🍕 |> add(5) >> 💩
's // アポストロフィS
// コメント
&& || // 論理演算子
`

	tests := []struct {
		expectedType    token.TokenType
		expectedLiteral string
	}{
		{token.PIZZA, "🍕"},
		{token.PIPE, "|>"},
		{token.IDENT, "add"},
		{token.LPAREN, "("},
		{token.INT, "5"},
		{token.RPAREN, ")"},
		{token.ASSIGN, ">>"},
		{token.POO, "💩"},
		{token.APOSTROPHE_S, "'s"},
		{token.AND, "&&"},
		{token.PIPE_PAR, "|"},
		{token.PIPE_PAR, "|"},
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
