package parser

import (
	"github.com/uncode/ast"
	"github.com/uncode/token"
)

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
