package parser

import (
	"fmt"

	"github.com/uncode/ast"
	"github.com/uncode/lexer"
	"github.com/uncode/logger"
	"github.com/uncode/token"
)

// 演算子の優先順位
const (
	_ int = iota
	LOWEST
	ASSIGN      // >>
	PIPE        // | |>
	LOGICAL     // && ||
	EQUALS      // == !=
	LESSGREATER // > < >= <=
	SUM         // + -
	PRODUCT     // * / %
	PREFIX      // -X or !X
	CALL        // myFunction(X)
	INDEX       // array[index]
	PROPERTY    // obj.prop or obj's prop
)

// 演算子の優先順位マップ
var precedences = map[token.TokenType]int{
	token.ASSIGN:       ASSIGN,
	token.EQUAL:        ASSIGN, // = も代入演算子として扱う
	token.EQ:           EQUALS,
	token.NOT_EQ:       EQUALS,
	token.LT:           LESSGREATER,
	token.GT:           LESSGREATER,
	token.LE:           LESSGREATER,
	token.GE:           LESSGREATER,
	token.PLUS:         SUM,
	token.MINUS:        SUM,
	token.SLASH:        PRODUCT,
	token.ASTERISK:     PRODUCT,
	token.MODULO:       PRODUCT,
	token.LPAREN:       CALL,
	token.LBRACKET:     INDEX,
	token.DOT:          PROPERTY,
	token.APOSTROPHE_S: PROPERTY,
	token.AND:          LOGICAL,
	token.OR:           LOGICAL,
	token.PIPE:         PIPE,
	token.PIPE_PAR:     PIPE,
	token.MAP_PIPE:     PIPE,
	token.FILTER_PIPE:  PIPE,
	token.DOTDOT:       SUM, // 範囲演算子の優先順位
}

type (
	prefixParseFn func() ast.Expression
	infixParseFn  func(ast.Expression) ast.Expression
)

// Parser はトークン列を解析して抽象構文木を生成する
type Parser struct {
	l         *lexer.Lexer
	tokens    []token.Token
	position  int
	curToken  token.Token
	peekToken token.Token
	errors    []string

	prefixParseFns    map[token.TokenType]prefixParseFn
	infixParseFns     map[token.TokenType]infixParseFn
	insideFunctionBody bool // 関数本体内かどうかのフラグ
}

// NewParser は新しいパーサーを生成する
func NewParser(tokens []token.Token) *Parser {
	p := &Parser{
		tokens:         tokens,
		position:       0,
		errors:         []string{},
		prefixParseFns: make(map[token.TokenType]prefixParseFn),
		infixParseFns:  make(map[token.TokenType]infixParseFn),
	}

	// 前置演算子の解析関数を登録
	p.registerPrefix(token.IDENT, p.parseIdentifier)
	p.registerPrefix(token.INT, p.parseIntegerLiteral)
	p.registerPrefix(token.FLOAT, p.parseFloatLiteral)
	p.registerPrefix(token.STRING, p.parseStringLiteral)
	p.registerPrefix(token.BOOLEAN, p.parseBooleanLiteral)
	p.registerPrefix(token.MINUS, p.parsePrefixExpression)
	p.registerPrefix(token.NOT, p.parsePrefixExpression)
	p.registerPrefix(token.BANG, p.parsePrefixExpression)
	p.registerPrefix(token.LPAREN, p.parseGroupedExpression)
	p.registerPrefix(token.LBRACKET, p.parseArrayLiteral)
	p.registerPrefix(token.LBRACE, p.parseBlockExpression)  // ブロック式の解析関数を追加
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)
	p.registerPrefix(token.CLASS, p.parseClassLiteral)
	p.registerPrefix(token.PIZZA, p.parsePizzaLiteral)
	p.registerPrefix(token.POO, p.parsePooLiteral)
	p.registerPrefix(token.DOTDOT, p.parseRangeExpression)
	// EOFトークンに対するダミー解析関数を登録
	p.registerPrefix(token.EOF, p.parseEOF)

	// 中置演算子の解析関数を登録
	p.registerInfix(token.PLUS, p.parseInfixExpression)
	p.registerInfix(token.MINUS, p.parseInfixExpression)
	p.registerInfix(token.SLASH, p.parseInfixExpression)
	p.registerInfix(token.ASTERISK, p.parseInfixExpression)
	p.registerInfix(token.MODULO, p.parseInfixExpression)
	p.registerInfix(token.EQ, p.parseInfixExpression)
	p.registerInfix(token.NOT_EQ, p.parseInfixExpression)
	p.registerInfix(token.LT, p.parseInfixExpression)
	p.registerInfix(token.GT, p.parseInfixExpression)
	p.registerInfix(token.LE, p.parseInfixExpression)
	p.registerInfix(token.GE, p.parseInfixExpression)
	p.registerInfix(token.AND, p.parseInfixExpression)
	p.registerInfix(token.OR, p.parseInfixExpression)
	p.registerInfix(token.LPAREN, p.parseCallExpression)
	p.registerInfix(token.LBRACKET, p.parseIndexExpression)
	p.registerInfix(token.DOT, p.parsePropertyExpression)
	p.registerInfix(token.APOSTROPHE_S, p.parsePropertyExpression)
	p.registerInfix(token.ASSIGN, p.parseAssignExpression)
	p.registerInfix(token.EQUAL, p.parseAssignExpression)  // = も代入演算子として扱う
	p.registerInfix(token.PIPE, p.parsePipeExpression)
	p.registerInfix(token.PIPE_PAR, p.parsePipeExpression)
	p.registerInfix(token.MAP_PIPE, p.parsePipeExpression)
	p.registerInfix(token.FILTER_PIPE, p.parsePipeExpression)

	// 最初の2つのトークンを読み込む
	if len(tokens) > 0 {
		p.curToken = tokens[0]
		if len(tokens) > 1 {
			p.peekToken = tokens[1]
		}
	}

	return p
}

// Errors はパース中に発生したエラーを返す
func (p *Parser) Errors() []string {
	return p.errors
}

// isNestedBlock はパーサーが現在ネストされたブロック内にいるかどうかを判定する
func (p *Parser) isNestedBlock() bool {
	// ブロック式のネスト状態を確認するためのヘルパー
	// パーサーの状態から現在ネストされたブロック内かどうかを判定
	
	// このメソッドは次の条件が両方満たされた場合にtrueを返す:
	// 1. 関数本体内にいる (insideFunctionBodyがtrue)
	// 2. 現在ネストされたブロックを解析中
	
	// 関数内でない場合、ネストしたブロックかどうかは関係ない
	if !p.insideFunctionBody {
		return false
	}
	
	// ネストレベルを確認
	return p.isParsingNestedBlock()
}

// isParsingNestedBlock は現在解析しているブロックがネストされているかを確認する
func (p *Parser) isParsingNestedBlock() bool {
	// このメソッドは現在のパーサーの状態からネストされたブロックを解析中かを判定
	
	// 現在のトークン列と位置を検査して、ネストレベルを判定
	// 簡易実装: 現在のブロックが関数のルート直下かどうかを判定
	
	// 現在の位置から遡って、関数宣言と最初のブレースを探す
	braceCount := 0
	
	// 現在のトークンの前までのトークンを検査
	for i := p.position; i >= 0; i-- {
		if i >= len(p.tokens) {
			continue
		}
		
		token := p.tokens[i]
		
		// ブレースのカウント
		if token.Type == "{" {
			braceCount++
		} else if token.Type == "}" {
			braceCount--
		}
		
		// 関数宣言を見つけた場合
		if token.Type == "def" {
			// 関数宣言後に見つかったブレースが2つ以上なら、ネストされたブロック
			// (1つ目は関数ブロック自体、2つ目以降はネストされたブロック)
			logger.ParserDebug("関数宣言を検出: ブレースカウント=%d", braceCount)
			return braceCount >= 2
		}
	}
	
	// デバッグ情報を記録
	logger.ParserDebug("ネストレベルの判定: braceCount=%d", braceCount)
	
	// 関数宣言が見つからない場合、またはネストレベルが判定できない場合は
	// シンプルな判定として、ブレースが2つ以上あればネストされたブロックとする
	return braceCount >= 2
}

// nextToken は次のトークンに進む
func (p *Parser) nextToken() {
	p.position++
	if p.position >= len(p.tokens) {
		p.curToken = token.Token{Type: token.EOF, Literal: ""}
		p.peekToken = token.Token{Type: token.EOF, Literal: ""}
	} else {
		p.curToken = p.peekToken
		if p.position+1 < len(p.tokens) {
			p.peekToken = p.tokens[p.position+1]
		} else {
			p.peekToken = token.Token{Type: token.EOF, Literal: ""}
		}
	}
}

// parseEOF はEOFトークンに対する前置解析関数
// EOFに到達したときに呼び出されるだけなので、実質的には何もしない
func (p *Parser) parseEOF() ast.Expression {
	logger.Debug("EOFトークンを検出しました（プログラムの終わり）")
	return nil
}

// ParseProgram はプログラム全体を解析する
func (p *Parser) ParseProgram() (*ast.Program, error) {
	program := &ast.Program{
		Statements: []ast.Statement{},
	}

	for p.curToken.Type != token.EOF {
		stmt := p.parseStatement()
		if stmt != nil {
			program.Statements = append(program.Statements, stmt)
		}
		p.nextToken()
	}

	if len(p.errors) > 0 {
		return nil, fmt.Errorf("パース中にエラーが発生しました: %v", p.errors)
	}

	return program, nil
}

// curTokenIs は現在のトークンが指定した型かどうかを判定する
func (p *Parser) curTokenIs(t token.TokenType) bool {
	return p.curToken.Type == t
}

// peekTokenIs は次のトークンが指定した型かどうかを判定する
func (p *Parser) peekTokenIs(t token.TokenType) bool {
	return p.peekToken.Type == t
}

// expectPeek は次のトークンが指定した型であれば次に進む
func (p *Parser) expectPeek(t token.TokenType) bool {
	if p.peekTokenIs(t) {
		p.nextToken()
		return true
	}
	p.peekError(t)
	return false
}

// peekError は次のトークンが期待と異なる場合にエラーを追加する
func (p *Parser) peekError(t token.TokenType) {
	msg := fmt.Sprintf("%d行目: 次のトークンは %s であることが期待されていますが、実際は %s です",
		p.peekToken.Line, t, p.peekToken.Type)
	p.errors = append(p.errors, msg)
}

// noPrefixParseFnError は前置解析関数がない場合にエラーを追加する
func (p *Parser) noPrefixParseFnError(t token.TokenType) {
	msg := fmt.Sprintf("%d行目: トークン %s に対する前置解析関数がありません",
		p.curToken.Line, t)
	p.errors = append(p.errors, msg)
}

// registerPrefix は前置演算子の解析関数を登録する
func (p *Parser) registerPrefix(tokenType token.TokenType, fn prefixParseFn) {
	p.prefixParseFns[tokenType] = fn
}

// registerInfix は中置演算子の解析関数を登録する
func (p *Parser) registerInfix(tokenType token.TokenType, fn infixParseFn) {
	p.infixParseFns[tokenType] = fn
}

// peekPrecedence は次のトークンの優先順位を返す
func (p *Parser) peekPrecedence() int {
	if p, ok := precedences[p.peekToken.Type]; ok {
		return p
	}
	return LOWEST
}

// curPrecedence は現在のトークンの優先順位を返す
func (p *Parser) curPrecedence() int {
	if p, ok := precedences[p.curToken.Type]; ok {
		return p
	}
	return LOWEST
}