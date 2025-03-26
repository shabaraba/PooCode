package lexer

import (
	"testing"
)

// TestIsLetter はisLetter関数の挙動をテストする
func TestIsLetter(t *testing.T) {
	tests := []struct {
		input    rune
		expected bool
	}{
		{'a', true},
		{'Z', true},
		{'_', true},
		{'あ', true},  // 日本語
		{'漢', true},  // 漢字
		{'1', false}, // 数字
		{'+', false}, // 演算子
		{' ', false}, // 空白
		{0, false},   // NULL
	}

	for i, tt := range tests {
		result := isLetter(tt.input)
		if result != tt.expected {
			t.Errorf("tests[%d] - isLetter(%q) = %v, want %v", 
				i, string(tt.input), result, tt.expected)
		}
	}
}

// TestIsDigit はisDigit関数の挙動をテストする
func TestIsDigit(t *testing.T) {
	tests := []struct {
		input    rune
		expected bool
	}{
		{'0', true},
		{'9', true},
		{'a', false},
		{'+', false},
		{'.', false},
		{' ', false},
		{0, false},
	}

	for i, tt := range tests {
		result := isDigit(tt.input)
		if result != tt.expected {
			t.Errorf("tests[%d] - isDigit(%q) = %v, want %v", 
				i, string(tt.input), result, tt.expected)
		}
	}
}

// TestReadChar はreadChar関数の基本的な挙動をテストする
func TestReadChar(t *testing.T) {
	input := "abc"
	l := NewLexer(input)

	// 初期状態の確認（NewLexerでreadCharが一度呼ばれる）
	if l.ch != 'a' {
		t.Errorf("Initial character should be 'a', got %q", string(l.ch))
	}

	// 次の文字を読み込む
	l.readChar()
	if l.ch != 'b' {
		t.Errorf("Second character should be 'b', got %q", string(l.ch))
	}

	// さらに次の文字を読み込む
	l.readChar()
	if l.ch != 'c' {
		t.Errorf("Third character should be 'c', got %q", string(l.ch))
	}

	// 最後の文字を超えて読み込む（EOFを表す0が返るべき）
	l.readChar()
	if l.ch != 0 {
		t.Errorf("After last character should be EOF(0), got %q", string(l.ch))
	}
}

// TestPeekChar はpeekChar関数の挙動をテストする
func TestPeekChar(t *testing.T) {
	input := "abc"
	l := NewLexer(input)

	// 初期位置から先読み
	if l.peekChar() != 'b' {
		t.Errorf("Peek should return 'b', got %q", string(l.peekChar()))
	}

	// 位置が変わっていないことを確認
	if l.ch != 'a' {
		t.Errorf("Current character should still be 'a', got %q", string(l.ch))
	}

	// 1文字進めて先読み
	l.readChar()
	if l.peekChar() != 'c' {
		t.Errorf("Peek should return 'c', got %q", string(l.peekChar()))
	}

	// もう1文字進めて先読み（最後の文字まで来た場合）
	l.readChar()
	if l.peekChar() != 0 {
		t.Errorf("Peek at last char should return EOF(0), got %q", string(l.peekChar()))
	}
}
