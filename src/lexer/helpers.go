package lexer

import (
	"unicode"
	"unicode/utf8"
)

// readChar は次の文字を読み込む
func (l *Lexer) readChar() {
	if l.readPosition >= len(l.input) {
		l.ch = 0 // EOFを表す
	} else {
		// UTF-8文字を正しく読み込む
		r, size := utf8.DecodeRuneInString(l.input[l.readPosition:])
		l.ch = r
		l.position = l.readPosition
		l.readPosition += size
		l.column++
	}

	// 改行文字の処理
	if l.ch == '\n' {
		l.line++
		l.column = 0
	}
}

// peekChar は次の文字を先読みする（位置は進めない）
func (l *Lexer) peekChar() rune {
	if l.readPosition >= len(l.input) {
		return 0
	}
	r, _ := utf8.DecodeRuneInString(l.input[l.readPosition:])
	return r
}

// isLetter は文字が識別子の一部として有効かどうかを判定する
func isLetter(ch rune) bool {
	// 簡単な回避策として、数字以外のすべての文字を許可
	return ch != 0 && !unicode.IsSpace(ch) && !unicode.IsDigit(ch) && 
		ch != '+' && ch != '-' && ch != '*' && ch != '/' && ch != '%' &&
		ch != '=' && ch != '!' && ch != '<' && ch != '>' && ch != '&' &&
		ch != '|' && ch != ',' && ch != ';' && ch != ':' && ch != '(' &&
		ch != ')' && ch != '{' && ch != '}' && ch != '[' && ch != ']' &&
		ch != '.' && ch != '\'' && ch != '"'
}

// isDigit は文字が数字かどうかを判定する
func isDigit(ch rune) bool {
	return unicode.IsDigit(ch)
}
