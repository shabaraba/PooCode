package parser

import (
	"github.com/uncode/ast"
	"github.com/uncode/token"
)

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

	// 条件付き関数定義の条件部分を解析
	if p.peekTokenIs(token.IF) {
		p.nextToken() // if
		p.nextToken()
		lit.Condition = p.parseExpression(LOWEST)
	}

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
