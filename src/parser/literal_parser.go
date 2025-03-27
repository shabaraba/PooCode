package parser

import (
	"fmt"
	"strconv"

	"github.com/uncode/ast"
	"github.com/uncode/token"
)

// parseIdentifier は識別子を解析する
// 修正版: 識別子の後に引数になりうるトークンが続く場合は関数呼び出しとして解析する
func (p *Parser) parseIdentifier() ast.Expression {
	ident := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
	
	// 識別子の後に引数になりうるトークンが続いていて、かつ括弧ではない場合
	// 例: func arg (括弧なしの関数呼び出し)
	if p.peekTokenIs(token.INT) || p.peekTokenIs(token.STRING) || 
	   p.peekTokenIs(token.IDENT) || p.peekTokenIs(token.BOOLEAN) {
		
		// 次のトークンに進む
		p.nextToken()
		
		// 引数を解析
		var arg ast.Expression
		
		// トークンタイプに応じて適切な式を生成
		switch p.curToken.Type {
		case token.INT:
			arg = p.parseIntegerLiteral()
		case token.STRING:
			arg = p.parseStringLiteral()
		case token.BOOLEAN:
			arg = p.parseBooleanLiteral()
		case token.IDENT:
			arg = &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
		}
		
		// 引数が見つかった場合は関数呼び出しとして解析
		if arg != nil {
			return &ast.CallExpression{
				Token:     p.curToken,
				Function:  ident,
				Arguments: []ast.Expression{arg},
			}
		}
	}
	
	// 通常の識別子として扱う
	return ident
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
	
	p.nextToken() // [の次のトークンへ
	
	// [1..10] のような範囲式の場合
	if (p.curTokenIs(token.INT) || p.curTokenIs(token.IDENT) || p.curTokenIs(token.STRING)) && p.peekTokenIs(token.DOTDOT) {
		// 開始値を解析
		startExp := p.parseExpression(LOWEST)
		
		p.nextToken() // ..へ
		rangeExp := &ast.RangeExpression{
			Token: p.curToken,
			Start: startExp,
		}
		
		p.nextToken() // ..の次のトークンへ
		
		// 終了値を解析
		if !p.curTokenIs(token.RBRACKET) {
			rangeExp.End = p.parseExpression(LOWEST)
			if !p.expectPeek(token.RBRACKET) {
				return nil
			}
		} else {
			// [start..] の形式（終了値なし）
			if !p.expectPeek(token.RBRACKET) {
				return nil
			}
		}
		
		return rangeExp
	} else if p.curTokenIs(token.DOTDOT) {
	// [..10] のような開始値なしの範囲式
		// [..end] の形式（開始値なし）
		rangeExp := &ast.RangeExpression{
			Token: p.curToken,
			Start: nil,
		}
		
		p.nextToken() // ..の次のトークンへ
		
		// 終了値を解析
		if !p.curTokenIs(token.RBRACKET) {
			rangeExp.End = p.parseExpression(LOWEST)
			if !p.expectPeek(token.RBRACKET) {
				return nil
			}
		} else {
			// [..] の形式（両方なし）
			if !p.expectPeek(token.RBRACKET) {
				return nil
			}
		}
		
		return rangeExp
	}
	
	// 通常の配列の場合 [1, 2, 3]
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

// parseRangeExpression は範囲式 [start..end] を解析する
func (p *Parser) parseRangeExpression() ast.Expression {
	rangeExp := &ast.RangeExpression{
		Token: p.curToken,
		Start: nil,
	}
	
	p.nextToken() // ..の次のトークンへ
	
	// 終了値を解析
	if !p.curTokenIs(token.RBRACKET) {
		rangeExp.End = p.parseExpression(LOWEST)
		if !p.expectPeek(token.RBRACKET) {
			return nil
		}
	} else {
		// [..] の形式（両方なし）
		if !p.expectPeek(token.RBRACKET) {
			return nil
		}
	}
	
	return rangeExp
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
