package parser

import (
	"fmt"
	"strconv"

	"github.com/uncode/ast"
	"github.com/uncode/lexer"
	"github.com/uncode/token"
)

// 演算子の優先順位
const (
	_ int = iota
	LOWEST
	PIPE
	ASSIGN      // >>
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

	prefixParseFns map[token.TokenType]prefixParseFn
	infixParseFns  map[token.TokenType]infixParseFn
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
	p.registerPrefix(token.FUNCTION, p.parseFunctionLiteral)
	p.registerPrefix(token.CLASS, p.parseClassLiteral)
	p.registerPrefix(token.PIZZA, p.parsePizzaLiteral)
	p.registerPrefix(token.POO, p.parsePooLiteral)

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
	p.registerInfix(token.PIPE, p.parsePipeExpression)
	p.registerInfix(token.PIPE_PAR, p.parsePipeExpression)

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

// parseStatement は文を解析する
func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.GLOBAL:
		return p.parseGlobalStatement()
	default:
		return p.parseExpressionStatement()
	}
}

// parseGlobalStatement はグローバル変数宣言を解析する
func (p *Parser) parseGlobalStatement() *ast.GlobalStatement {
	stmt := &ast.GlobalStatement{Token: p.curToken}

	p.nextToken()
	// 型情報があれば解析
	if p.curTokenIs(token.IDENT) {
		// 型情報を取得
		typeStr := p.curToken.Literal
		p.nextToken()
		stmt.Type = typeStr
	}

	// 変数名を解析
	if !p.curTokenIs(token.IDENT) {
		p.peekError(token.IDENT)
		return nil
	}

	stmt.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	return stmt
}

// parseExpressionStatement は式文を解析する
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
}

// parseExpression は式を解析する
func (p *Parser) parseExpression(precedence int) ast.Expression {
	prefix := p.prefixParseFns[p.curToken.Type]
	if prefix == nil {
		p.noPrefixParseFnError(p.curToken.Type)
		return nil
	}
	leftExp := prefix()

	for !p.peekTokenIs(token.SEMICOLON) && precedence < p.peekPrecedence() {
		infix := p.infixParseFns[p.peekToken.Type]
		if infix == nil {
			return leftExp
		}

		p.nextToken()
		leftExp = infix(leftExp)
	}

	return leftExp
}

// parsePipeExpression はパイプライン式を解析する
func (p *Parser) parsePipeExpression(left ast.Expression) ast.Expression {
	// パイプ演算子のトークンを保存
	pipeToken := p.curToken
	
	// パイプ演算子の優先順位を取得
	precedence := p.curPrecedence()
	
	// 次のトークンに進む
	p.nextToken()
	
	// パイプの右側の式を解析する
	// 右側は関数または関数呼び出しである必要がある
	rightExp := p.parseExpression(precedence)
	
	// パイプタイプに応じて処理を分ける
	if pipeToken.Type == token.PIPE {
		// 並列パイプ (|) の場合
		return &ast.InfixExpression{
			Token:    pipeToken,
			Operator: pipeToken.Literal,
			Left:     left,
			Right:    rightExp,
		}
	} else {
		// 通常パイプ (|>) の場合
		// 右辺がCallExpressionかどうかで処理を分ける
		if callExp, ok := rightExp.(*ast.CallExpression); ok {
			// 既存の引数リストの先頭に left を追加
			callExp.Arguments = append([]ast.Expression{left}, callExp.Arguments...)
			return callExp
		} else {
			// 右辺が関数呼び出しでない場合、新しい関数呼び出しとして扱う
			return &ast.CallExpression{
				Token:     pipeToken,
				Function:  rightExp,
				Arguments: []ast.Expression{left},
			}
		}
	}
}

// parseAssignExpression は代入式を解析する
func (p *Parser) parseAssignExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	p.nextToken()
	expression.Right = p.parseExpression(LOWEST)

	return expression
}

// parseIdentifier は識別子を解析する
func (p *Parser) parseIdentifier() ast.Expression {
	return &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
}

// parseIntegerLiteral は整数リテラルを解析する
func (p *Parser) parseIntegerLiteral() ast.Expression {
	lit := &ast.IntegerLiteral{Token: p.curToken}

	value, err := strconv.ParseInt(p.curToken.Literal, 0, 64)
	if err != nil {
		msg := fmt.Sprintf("整数 '%s' を解析できませんでした", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value
	return lit
}

// parseFloatLiteral は浮動小数点リテラルを解析する
func (p *Parser) parseFloatLiteral() ast.Expression {
	lit := &ast.FloatLiteral{Token: p.curToken}

	value, err := strconv.ParseFloat(p.curToken.Literal, 64)
	if err != nil {
		msg := fmt.Sprintf("浮動小数点数 '%s' を解析できませんでした", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}

	lit.Value = value
	return lit
}

// parseStringLiteral は文字列リテラルを解析する
func (p *Parser) parseStringLiteral() ast.Expression {
	return &ast.StringLiteral{Token: p.curToken, Value: p.curToken.Literal}
}

// parseBooleanLiteral は真偽値リテラルを解析する
func (p *Parser) parseBooleanLiteral() ast.Expression {
	value, err := strconv.ParseBool(p.curToken.Literal)
	if err != nil {
		msg := fmt.Sprintf("真偽値 '%s' を解析できませんでした", p.curToken.Literal)
		p.errors = append(p.errors, msg)
		return nil
	}
	return &ast.BooleanLiteral{Token: p.curToken, Value: value}
}

// parsePrefixExpression は前置式を解析する
func (p *Parser) parsePrefixExpression() ast.Expression {
	expression := &ast.PrefixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
	}

	p.nextToken()
	expression.Right = p.parseExpression(PREFIX)

	return expression
}

// parseInfixExpression は中置式を解析する
func (p *Parser) parseInfixExpression(left ast.Expression) ast.Expression {
	expression := &ast.InfixExpression{
		Token:    p.curToken,
		Operator: p.curToken.Literal,
		Left:     left,
	}

	precedence := p.curPrecedence()
	p.nextToken()
	expression.Right = p.parseExpression(precedence)

	return expression
}

// parseGroupedExpression は括弧で囲まれた式を解析する
func (p *Parser) parseGroupedExpression() ast.Expression {
	p.nextToken()

	exp := p.parseExpression(LOWEST)

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return exp
}

// parseArrayLiteral は配列リテラルを解析する
func (p *Parser) parseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{Token: p.curToken}
	array.Elements = p.parseExpressionList(token.RBRACKET)
	return array
}

// parseExpressionList は式のリストを解析する
func (p *Parser) parseExpressionList(end token.TokenType) []ast.Expression {
	list := []ast.Expression{}

	if p.peekTokenIs(end) {
		p.nextToken()
		return list
	}

	p.nextToken()
	list = append(list, p.parseExpression(LOWEST))

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		list = append(list, p.parseExpression(LOWEST))
	}

	if !p.expectPeek(end) {
		return nil
	}

	return list
}

// parseIndexExpression は添字式を解析する
func (p *Parser) parseIndexExpression(left ast.Expression) ast.Expression {
	exp := &ast.IndexExpression{Token: p.curToken, Left: left}

	p.nextToken()
	exp.Index = p.parseExpression(LOWEST)

	if !p.expectPeek(token.RBRACKET) {
		return nil
	}

	return exp
}

// parsePropertyExpression はプロパティアクセス式を解析する
func (p *Parser) parsePropertyExpression(left ast.Expression) ast.Expression {
	exp := &ast.PropertyAccessExpression{
		Token:  p.curToken,
		Object: left,
	}

	p.nextToken()
	exp.Property = p.parseExpression(PROPERTY)

	return exp
}

// parseFunctionLiteral は関数リテラルを解析する
func (p *Parser) parseFunctionLiteral() ast.Expression {
	lit := &ast.FunctionLiteral{Token: p.curToken}

	// 関数名があれば解析
	if p.peekTokenIs(token.IDENT) {
		p.nextToken()
		lit.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	}

	// パラメータリストを解析
	if !p.expectPeek(token.LPAREN) {
		return nil
	}
	lit.Parameters = p.parseFunctionParameters()

	// 型注釈があれば解析
	if p.peekTokenIs(token.COLON) {
		p.nextToken() // :
		p.nextToken() // 入力型
		lit.InputType = p.curToken.Literal

		if p.peekTokenIs(token.MINUS) {
			p.nextToken() // -
			if p.peekTokenIs(token.GT) {
				p.nextToken() // >
				p.nextToken() // 出力型
				lit.ReturnType = p.curToken.Literal
			}
		}
	}

	// 条件付き関数定義の条件部分を解析
	if p.peekTokenIs(token.IF) {
		p.nextToken() // if
		p.nextToken()
		lit.Condition = p.parseExpression(LOWEST)
	}

	// 関数本体を解析
	if !p.expectPeek(token.LBRACE) {
		return nil
	}
	lit.Body = p.parseBlockStatement()

	return lit
}

// parseFunctionParameters は関数のパラメータリストを解析する
func (p *Parser) parseFunctionParameters() []*ast.Identifier {
	identifiers := []*ast.Identifier{}

	if p.peekTokenIs(token.RPAREN) {
		p.nextToken()
		return identifiers
	}

	p.nextToken()

	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	identifiers = append(identifiers, ident)

	for p.peekTokenIs(token.COMMA) {
		p.nextToken()
		p.nextToken()
		ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		identifiers = append(identifiers, ident)
	}

	if !p.expectPeek(token.RPAREN) {
		return nil
	}

	return identifiers
}

// parseCallExpression は関数呼び出し式を解析する
func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.curToken, Function: function}
	exp.Arguments = p.parseExpressionList(token.RPAREN)
	return exp
}

// parseBlockStatement はブロック文を解析する
func (p *Parser) parseBlockStatement() *ast.BlockStatement {
	block := &ast.BlockStatement{Token: p.curToken}
	block.Statements = []ast.Statement{}

	p.nextToken()

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}

	return block
}

// parseClassLiteral はクラス定義を解析する
func (p *Parser) parseClassLiteral() ast.Expression {
	lit := &ast.ClassLiteral{Token: p.curToken}

	// クラス名を解析
	if !p.expectPeek(token.IDENT) {
		return nil
	}
	lit.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}

	// 継承があれば解析
	if p.peekTokenIs(token.EXTENDS) {
		p.nextToken() // extends
		if !p.expectPeek(token.IDENT) {
			return nil
		}
		lit.Extends = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	}

	// クラス本体を解析
	if !p.expectPeek(token.LBRACE) {
		return nil
	}

	p.nextToken()

	// プロパティとメソッドを解析
	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		if p.curTokenIs(token.PUBLIC) || p.curTokenIs(token.PRIVATE) {
			// プロパティ定義
			prop := &ast.PropertyDefinition{
				Token:      p.curToken,
				Visibility: p.curToken.Literal,
			}

			p.nextToken()
			// 型情報があれば解析
			if p.peekTokenIs(token.IDENT) && !p.peekTokenIs(token.IDENT) {
				prop.Type = p.curToken.Literal
				p.nextToken()
			}

			// プロパティ名を解析
			if !p.curTokenIs(token.IDENT) {
				p.peekError(token.IDENT)
				return nil
			}
			prop.Name = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
			lit.Properties = append(lit.Properties, prop)
		} else if p.curTokenIs(token.FUNCTION) {
			// メソッド定義
			method := p.parseFunctionLiteral().(*ast.FunctionLiteral)
			lit.Methods = append(lit.Methods, method)
		} else {
			p.errors = append(p.errors, fmt.Sprintf("クラス定義内で予期しないトークンです: %s", p.curToken.Literal))
		}
		p.nextToken()
	}

	return lit
}

// parsePizzaLiteral は🍕リテラルを解析する
func (p *Parser) parsePizzaLiteral() ast.Expression {
	return &ast.PizzaLiteral{Token: p.curToken}
}

// parsePooLiteral は💩リテラルを解析する
func (p *Parser) parsePooLiteral() ast.Expression {
	return &ast.PooLiteral{Token: p.curToken}
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
