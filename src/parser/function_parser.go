package parser

import (
	"fmt"

	"github.com/uncode/ast"
	"github.com/uncode/logger"
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

	// 修正: 括弧を使った関数定義と括弧なしの関数定義の両方をサポート
	// 次のトークンがパラメータ（IDENT）かどうかをチェック
	if p.peekTokenIs(token.IDENT) {
		// 括弧なしのパラメータ定義: def func param { ... }
		p.nextToken()
		// 引数は1つだけサポート
		param := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		lit.Parameters = []*ast.Identifier{param}
	} else if p.peekTokenIs(token.LPAREN) {
		// 括弧ありのパラメータ定義: def func(param) { ... }
		p.nextToken() // (
		lit.Parameters = p.parseFunctionParameters()
	}

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
// 修正版: 括弧を使った呼び出し「func(arg)」形式をサポート
func (p *Parser) parseCallExpression(function ast.Expression) ast.Expression {
	exp := &ast.CallExpression{Token: p.curToken, Function: function}
	
	// 引数リストを解析
	args := p.parseExpressionList(token.RPAREN)
	
	logger.Debug("関数呼び出し解析中: 関数=%s, 引数数=%d", function.String(), len(args))
	
	// スタンドアロンな関数呼び出しの場合、引数は最大1つまで
	if len(args) > 1 {
		p.errors = append(p.errors, fmt.Sprintf("%d行目: 関数 %s は最大1つの引数しか取れません（パイプラインを除く）", 
			p.curToken.Line, function.String()))
		// エラーの場合でも、最初の引数だけを使用して解析を続行
		exp.Arguments = []ast.Expression{args[0]}
	} else {
		exp.Arguments = args
	}
	
	return exp
}
