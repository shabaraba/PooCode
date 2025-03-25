package parser

import (
	"github.com/uncode/ast"
	"github.com/uncode/token"
)

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

// parseExpressionStatement は式文を解析する
func (p *Parser) parseExpressionStatement() *ast.ExpressionStatement {
	stmt := &ast.ExpressionStatement{Token: p.curToken}

	stmt.Expression = p.parseExpression(LOWEST)

	if p.peekTokenIs(token.SEMICOLON) {
		p.nextToken()
	}

	return stmt
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
