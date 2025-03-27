package lexer

import (
	"unicode"
	"github.com/uncode/token"
)

// readString は文字列リテラルを読み込む
func (l *Lexer) readString() string {
	position := l.position + 1
	for {
		l.readChar()
		if l.ch == '"' || l.ch == 0 {
			break
		}
	}
	return l.input[position:l.position]
}

// readIdentifier は識別子を読み込む
func (l *Lexer) readIdentifier() string {
	position := l.position
	
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	
	return l.input[position:l.position]
}

// readNumber は数値を読み込む
func (l *Lexer) readNumber() token.Token {
	position := l.position
	isFloat := false

	for isDigit(l.ch) {
		l.readChar()
		// Check for float point
		if l.ch == '.' {
			// If the next char is also '.', then this is not a float
			// but an integer followed by a range operator'..'
			if l.peekChar() == '.' {
				break // Exit the loop to return INT 
			}
			
			// It's a decimal point
			isFloat = true
			l.readChar()
		}
	}

	literal := l.input[position:l.position]
	if isFloat {
		return token.Token{
			Type:    token.FLOAT,
			Literal: literal,
			Line:    l.line,
			Column:  l.column,
		}
	}
	return token.Token{
		Type:    token.INT,
		Literal: literal,
		Line:    l.line,
		Column:  l.column,
	}
}

// skipWhitespace は空白文字をスキップする
func (l *Lexer) skipWhitespace() {
	for unicode.IsSpace(l.ch) {
		l.readChar()
	}
}

// skipComment はコメントをスキップする
// '//' から行末までをスキップする
func (l *Lexer) skipComment() {
	// 最初の '/' は既に読み込み済み、次の '/' もスキップ
	l.readChar()
	
	// 改行文字または終端に到達するまでスキップ
	for l.ch != '\n' && l.ch != 0 {
		l.readChar()
	}
}
