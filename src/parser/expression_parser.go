package parser

import (
	"github.com/uncode/ast"
	"github.com/uncode/logger"
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

// parseBlockExpression は{...}で囲まれたブロックを式として解析する
func (p *Parser) parseBlockExpression() ast.Expression {
	logger.ParserDebug("ブロック式の解析開始")
	
	// 現在のトークン（{）を含むブロック文を作成
	blockStmt := p.parseBlockStatement()
	if blockStmt == nil {
		logger.ParserDebug("ブロック文の解析に失敗")
		return nil
	}
	
	// ブロック式を作成
	blockExpr := &ast.BlockExpression{
		Token: blockStmt.Token,
		Block: blockStmt,
	}
	
	logger.ParserDebug("ブロック式の解析完了")
	return blockExpr
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
	
	// スライス表記の処理 (array[start..end])
	if p.curTokenIs(token.DOTDOT) {
		// array[..end] の形式
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
			// array[..] の形式（両方なし）
			if !p.expectPeek(token.RBRACKET) {
				return nil
			}
		}
		
		// 範囲式を直接返す
		return rangeExp
	}
	
	// 通常の添字またはスライス表記の開始値
	exp.Index = p.parseExpression(LOWEST)
	
	// array[start..end] または array[start..] の形式の場合
	if p.peekTokenIs(token.DOTDOT) {
		p.nextToken() // ..へ
		rangeExp := &ast.RangeExpression{
			Token: p.curToken,
			Start: exp.Index, // すでに解析した開始値
		}
		
		p.nextToken() // ..の次のトークンへ
		
		// 終了値を解析
		if !p.curTokenIs(token.RBRACKET) {
			rangeExp.End = p.parseExpression(LOWEST)
			if !p.expectPeek(token.RBRACKET) {
				return nil
			}
		} else {
			// array[start..] の形式（終了値なし）
			if !p.expectPeek(token.RBRACKET) {
				return nil
			}
		}
		
		// 範囲式を直接返す
		return rangeExp
	}
	
	// 通常の添字アクセスの場合
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
	
	// デバッグ出力
	if left != nil {
		logger.ParserDebug("パイプライン式の解析開始：演算子=%s, 左辺=%s", pipeToken.Literal, left.String())
	} else {
		logger.ParserDebug("パイプライン式の解析開始：演算子=%s, 左辺=nil", pipeToken.Literal)
		return createErrorExpression(pipeToken, "パイプラインの左辺がnilです")
	}
	
	// パイプ演算子の優先順位を取得
	precedence := p.curPrecedence()
	
	// 次のトークンに進む
	p.nextToken()
	
	// 現在のトークンと次のトークンを記録（デバッグ用）
	logger.ParserDebug("パイプライン右辺の解析中：現在のトークン=%s, 次のトークン=%s", 
		p.curToken.Literal, p.peekToken.Literal)
	
	// パイプの右側の式を解析する
	// パイプの種類に応じた処理
	pipeType := pipeToken.Type
	pipeOp := pipeToken.Literal
	
	// デバッグ情報を追加
	logger.Debug("パイプタイプ解析: Type=%s, Literal=%s", pipeType, pipeOp)
	
	// 通常のパイプライン処理
	// 右辺が関数名やその他の式である場合
	rightExp := p.parseExpression(precedence)
	
	// 解析された右辺の式を記録
	if rightExp != nil {
		logger.ParserDebug("パイプライン右辺の解析結果：タイプ=%T, 値=%s", rightExp, rightExp.String())
	} else {
		logger.ParserDebug("パイプライン右辺の解析結果：nil")
		// nilの場合はエラー処理
		return createErrorExpression(pipeToken, "パイプラインの右辺がnilです")
	}
	
	// パイプタイプに応じて処理を分ける
	if pipeToken.Type == token.PIPE_PAR {
		// 並列パイプ (|) の場合
		return &ast.InfixExpression{
			Token:    pipeToken,
			Operator: pipeToken.Literal,
			Left:     left,
			Right:    rightExp,
		}
	} else {
		// |> 演算子の場合（パイプライン）
		
		// 識別子の後に引数が続く場合の特別処理
		if ident, ok := rightExp.(*ast.Identifier); ok {
			logger.ParserDebug("識別子 '%s' が検出されました。次のトークンをチェック中...", ident.Value)
			
			// 次のトークンが識別子、整数、文字列、配列など、有効な引数となりうるトークンであれば
			// それを関数の引数として処理する
			if !p.peekTokenIs(token.PIPE) && !p.peekTokenIs(token.PIPE_PAR) && 
			   !p.peekTokenIs(token.MAP_PIPE) && !p.peekTokenIs(token.FILTER_PIPE) &&
			   !p.peekTokenIs(token.ASSIGN) && !p.peekTokenIs(token.SEMICOLON) &&
			   !p.peekTokenIs(token.RPAREN) && !p.peekTokenIs(token.RBRACE) &&
			   !p.peekTokenIs(token.RBRACKET) && !p.peekTokenIs(token.COMMA) {
				
				logger.ParserDebug("引数として処理可能なトークンが続きます: %s (%s)", p.peekToken.Literal, p.peekToken.Type)
				
				// 次のトークンを取得
				p.nextToken()
				
				// 引数を収集
				var args []ast.Expression
				
				// 最初の引数を解析
				if p.curToken.Type == token.PIZZA {
					// 🍕トークンが引数の場合、特別処理
					arg := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
					args = append(args, arg)
					logger.ParserDebug("🍕が第1引数として検出されました")
				} else {
					// 通常の引数解析
					arg := p.parseExpression(LOWEST)
					if arg != nil {
						args = append(args, arg)
						logger.ParserDebug("解析された第1引数: %s (タイプ: %T)", arg.String(), arg)
					} else {
						logger.ParserDebug("解析された第1引数: nil")
					}
				}
				
				// さらに引数がある場合
				for p.peekTokenIs(token.IDENT) || p.peekTokenIs(token.INT) || 
					p.peekTokenIs(token.STRING) || p.peekTokenIs(token.BOOLEAN) ||
					p.peekTokenIs(token.PIZZA) || p.peekTokenIs(token.LBRACKET) {
					
					p.nextToken()
					
					if p.curToken.Type == token.PIZZA {
						// 🍕トークンが引数の場合、特別処理
						arg := &ast.Identifier{Token: p.curToken, Value: p.curToken.Literal}
						args = append(args, arg)
						logger.ParserDebug("🍕が追加の引数として検出されました")
					} else {
						// 通常の引数解析
						arg := p.parseExpression(LOWEST)
						if arg != nil {
							args = append(args, arg)
							logger.ParserDebug("解析された追加引数: %s (タイプ: %T)", arg.String(), arg)
						} else {
							logger.ParserDebug("解析された追加引数: nil")
						}
					}
					
					// パイプやセミコロンが来たらループを抜ける
					if p.peekTokenIs(token.PIPE) || p.peekTokenIs(token.PIPE_PAR) || 
					   p.peekTokenIs(token.ASSIGN) || p.peekTokenIs(token.SEMICOLON) {
						break
					}
				}
				
				// CallExpressionを生成
				callExpr := &ast.CallExpression{
					Token:     pipeToken,
					Function:  ident,
					Arguments: args,
				}
				
				// パイプライン式の右辺としてCallExpressionを使用
				logger.ParserDebug("関数呼び出し式を生成: %s(引数: %d個)", ident.Value, len(args))
				rightExp = callExpr
			} else {
				logger.ParserDebug("引数なしの識別子: %s、次のトークン: %s", ident.Value, p.peekToken.Literal)
			}
			
			// 引数がない場合は通常のパイプライン
			return &ast.InfixExpression{
				Token:    pipeToken,
				Operator: pipeToken.Literal,
				Left:     left,
				Right:    rightExp,
			}
		}
		
		// 通常のパイプライン式として処理
		return &ast.InfixExpression{
			Token:    pipeToken,
			Operator: pipeToken.Literal,
			Left:     left,
			Right:    rightExp,
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

// createErrorExpression はエラーメッセージを含む式を作成する
func createErrorExpression(token token.Token, message string) ast.Expression {
	// エラーメッセージをログに出力
	logger.ParserDebug("パースエラー: %s", message)
	
	// エラーメッセージを含む文字列リテラルを作成
	return &ast.StringLiteral{
		Token: token,
		Value: "エラー: " + message,
	}
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

// parseCaseStatement は既に statement_parser.go で定義されているため注釈のみ残しています
// func (p *Parser) parseCaseStatement() *ast.CaseStatement { ... }

// parseFunctionLiteral は既に function_parser.go で定義されているため注釈のみ残しています
// func (p *Parser) parseFunctionLiteral() ast.Expression { ... }
