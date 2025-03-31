package lexer

import (
	"github.com/uncode/token"
)

// Lexer は入力文字列を解析してトークンを生成する
type Lexer struct {
	input        string // 入力文字列
	position     int    // 現在の位置
	readPosition int    // 次の読み込み位置
	ch           rune   // 現在の文字
	line         int    // 現在の行番号
	column       int    // 現在の列番号
}

// NewLexer は新しいLexerを生成する
func NewLexer(input string) *Lexer {
	l := &Lexer{
		input: input,
		line:  1,
	}
	l.readChar() // 最初の文字を読み込む
	return l
}

// Tokenize は入力文字列を全てトークン化する
func (l *Lexer) Tokenize() ([]token.Token, error) {
	var tokens []token.Token
	for {
		tok := l.NextToken()
		tokens = append(tokens, tok)
		if tok.Type == token.EOF {
			break
		}
	}
	return tokens, nil
}

// NextToken は次のトークンを取得する
func (l *Lexer) NextToken() token.Token {
	l.skipWhitespace()

	var tok token.Token
	tok.Line = l.line
	tok.Column = l.column

	switch l.ch {
	case '-':
		tok = l.newToken(token.MINUS, string(l.ch))
	case '*':
		tok = l.newToken(token.ASTERISK, string(l.ch))
	case '/':
		// コメントのチェック: // が見つかったら行末までスキップ
		if l.peekChar() == '/' {
			l.skipComment()
			return l.NextToken() // コメントをスキップした後で次のトークンを取得
		}
		tok = l.newToken(token.SLASH, string(l.ch))
	case '%':
		tok = l.newToken(token.MODULO, string(l.ch))
	case '=':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = l.newToken(token.EQ, string(ch)+string(l.ch))
		} else {
			tok = l.newToken(token.EQUAL, string(l.ch))
		}
	case '!':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = l.newToken(token.NOT_EQ, string(ch)+string(l.ch))
		} else {
			tok = l.newToken(token.BANG, string(l.ch))
		}
	case '<':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = l.newToken(token.LE, string(ch)+string(l.ch))
		} else {
			tok = l.newToken(token.LT, string(l.ch))
		}
	case '>':
		if l.peekChar() == '=' {
			ch := l.ch
			l.readChar()
			tok = l.newToken(token.GE, string(ch)+string(l.ch))
		} else if l.peekChar() == '>' {
			ch := l.ch
			l.readChar()
			tok = l.newToken(token.ASSIGN, string(ch)+string(l.ch))
		} else {
			tok = l.newToken(token.GT, string(l.ch))
		}
	case '&':
		if l.peekChar() == '&' {
			ch := l.ch
			l.readChar()
			tok = l.newToken(token.AND, string(ch)+string(l.ch))
		} else {
			tok = l.newToken(token.ILLEGAL, string(l.ch))
		}
	case '|':
		if l.peekChar() == '>' {
			ch := l.ch
			l.readChar()
			tok = l.newToken(token.PIPE, string(ch)+string(l.ch))
		} else {
			tok = l.newToken(token.PIPE_PAR, string(l.ch))
		}
	case '+':
		if l.peekChar() == '>' {
			ch := l.ch
			l.readChar()
			// 実際の演算子をリテラルとして使用
			tok = l.newToken(token.MAP_PIPE, string(ch)+string(l.ch))
		} else {
			tok = l.newToken(token.PLUS, string(l.ch))
		}
	case '?':
		if l.peekChar() == '>' {
			ch := l.ch
			l.readChar()
			// 実際の演算子をリテラルとして使用
			tok = l.newToken(token.FILTER_PIPE, string(ch)+string(l.ch))
		} else {
			tok = l.newToken(token.ILLEGAL, string(l.ch))
		}
	case ',':
		tok = l.newToken(token.COMMA, string(l.ch))
	case ';':
		tok = l.newToken(token.SEMICOLON, string(l.ch))
	case ':':
		tok = l.newToken(token.COLON, string(l.ch))
	case '(':
		tok = l.newToken(token.LPAREN, string(l.ch))
	case ')':
		tok = l.newToken(token.RPAREN, string(l.ch))
	case '{':
		tok = l.newToken(token.LBRACE, string(l.ch))
	case '}':
		tok = l.newToken(token.RBRACE, string(l.ch))
	case '[':
		tok = l.newToken(token.LBRACKET, string(l.ch))
	case ']':
		tok = l.newToken(token.RBRACKET, string(l.ch))
	case '.':
		// Check for '..' (range operator)
		if l.peekChar() == '.' {
			ch := l.ch
			l.readChar()
			tok = l.newToken(token.DOTDOT, string(ch)+string(l.ch))
		} else {
			// Check if it's a floating point number
			if isDigit(l.peekChar()) {
				// This is part of a float, go back and let readNumber handle it
				l.position--
				l.readPosition--
				l.column--
				return l.readNumber()
			}
			// It's just a regular period
			tok = l.newToken(token.DOT, string(l.ch))
		}
	case '\'':
		if l.peekChar() == 's' {
			ch := l.ch
			l.readChar()
			tok = l.newToken(token.APOSTROPHE_S, string(ch)+string(l.ch))
		} else {
			tok = l.newToken(token.ILLEGAL, string(l.ch))
		}
	case '"':
		// 文字列リテラルの開始位置を記録
		startLine := l.line
		startColumn := l.column
		
		// 文字列を読み込む
		literal := l.readString()
		
		// トークンを生成
		tok = token.Token{
			Type:    token.STRING,
			Literal: literal,
			Line:    startLine,
			Column:  startColumn,
		}
		return tok
	case '🍕':
		tok = l.newToken(token.PIZZA, string(l.ch))
	case '💩':
		tok = l.newToken(token.POO, string(l.ch))
	case 0:
		tok.Type = token.EOF
		tok.Literal = ""
	default:
		if isLetter(l.ch) {
			// 識別子の開始位置を記録
			startLine := l.line
			startColumn := l.column
			
			// 識別子を読み込む
			literal := l.readIdentifier()
			
			// トークンを生成
			tok = token.Token{
				Type:    token.LookupIdent(literal),
				Literal: literal,
				Line:    startLine,
				Column:  startColumn,
			}
			return tok
		} else if isDigit(l.ch) {
			return l.readNumber()
		} else {
			tok = l.newToken(token.ILLEGAL, string(l.ch))
		}
	}

	l.readChar()
	return tok
}

// newToken は新しいトークンを生成する
func (l *Lexer) newToken(tokenType token.TokenType, literal string) token.Token {
	return token.Token{
		Type:    tokenType,
		Literal: literal,
		Line:    l.line,
		Column:  l.column,
	}
}
