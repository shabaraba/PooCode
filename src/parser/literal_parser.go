package parser

import (
	"fmt"
	"strconv"

	"github.com/uncode/ast"
	"github.com/uncode/token"
)

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

// parseArrayLiteral は配列リテラルを解析する
func (p *Parser) parseArrayLiteral() ast.Expression {
	array := &ast.ArrayLiteral{Token: p.curToken}
	array.Elements = p.parseExpressionList(token.RBRACKET)
	return array
}

// parsePizzaLiteral は🍕リテラルを解析する
func (p *Parser) parsePizzaLiteral() ast.Expression {
	return &ast.PizzaLiteral{Token: p.curToken}
}

// parsePooLiteral は💩リテラルを解析する
func (p *Parser) parsePooLiteral() ast.Expression {
	return &ast.PooLiteral{Token: p.curToken}
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
