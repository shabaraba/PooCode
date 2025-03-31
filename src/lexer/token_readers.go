package lexer

import (
	"unicode"
	"github.com/uncode/token"
)

// readString は文字列リテラルを読み込む
// エスケープシーケンスもサポート
func (l *Lexer) readString() string {
	var result []rune
	l.readChar() // 最初の " をスキップ

	for {
		// 文字列終端または入力終端に達した場合
		if l.ch == '"' || l.ch == 0 {
			break
		}
		
		// エスケープシーケンスの処理
		if l.ch == '\\' {
			l.readChar() // バックスラッシュをスキップして次の文字を読む
			
			switch l.ch {
			case 'n':
				result = append(result, '\n') // 改行
			case 't':
				result = append(result, '\t') // タブ
			case 'r':
				result = append(result, '\r') // キャリッジリターン
			case '\\':
				result = append(result, '\\') // バックスラッシュ
			case '"':
				result = append(result, '"') // 二重引用符
			case '\'':
				result = append(result, '\'') // 一重引用符
			case '0':
				result = append(result, '\x00') // NULL文字
			case '+', '|', '>', ' ', '\n':
				// よく間違えられる文字に対する特別なケース
				// 例えば "string" + var のようなケースで " + var" を別の文字列として解釈しないように
				result = append(result, l.ch)
			default:
				// 未知のエスケープシーケンスの場合はそのまま両方の文字を追加
				result = append(result, '\\')
				result = append(result, l.ch)
			}
		} else {
			// 通常の文字はそのまま追加
			result = append(result, l.ch)
		}
		
		l.readChar()
	}
	
	// 文字列の終端（閉じ二重引用符）がまだ残っていればスキップ
	if l.ch == '"' {
		l.readChar()
	}
	
	return string(result)
}

// readIdentifier は識別子を読み込む
func (l *Lexer) readIdentifier() string {
	position := l.position
	
	// 最初の文字が既に識別子として有効であることは呼び出し側でチェック済み
	for isLetter(l.ch) || isDigit(l.ch) {
		l.readChar()
	}
	
	return l.input[position:l.position]
}

// readNumber は数値を読み込む
func (l *Lexer) readNumber() token.Token {
	position := l.position
	startLine := l.line
	startColumn := l.column
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
			Line:    startLine,
			Column:  startColumn,
		}
	}
	return token.Token{
		Type:    token.INT,
		Literal: literal,
		Line:    startLine,
		Column:  startColumn,
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
