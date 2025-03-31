package parser

import (
	"fmt"
	
	"github.com/uncode/ast"
	"github.com/uncode/logger"
	"github.com/uncode/token"
)

// parseStatement は文を解析する
func (p *Parser) parseStatement() ast.Statement {
	switch p.curToken.Type {
	case token.GLOBAL:
		return p.parseGlobalStatement()
	case token.CASE:
		// 関数内でのみcase文を許可するチェック
		if !p.insideFunctionBody {
			p.errors = append(p.errors, fmt.Sprintf("%d行目: case文は関数ブロック内でのみ使用できます。関数定義内で使用してください", p.curToken.Line))
			logger.ParserDebug("関数外でのcase文使用を検出: エラー報告 (insideFunctionBody=%v)", p.insideFunctionBody)
			return nil
		}
		
		// ネストしたブロック内のcase文を禁止する追加チェック
		// 直接関数の本体内でないcase文は禁止
		if p.isNestedBlock() {
			p.errors = append(p.errors, fmt.Sprintf("%d行目: case文は関数のルート階層でのみ使用できます。ネストされたブロック内では使用できません", p.curToken.Line))
			logger.ParserDebug("ネストされたブロック内でのcase文使用を検出: エラー報告")
			return nil
		}
		
		logger.ParserDebug("関数内でのcase文使用を検出 (insideFunctionBody=%v)", p.insideFunctionBody)
		return p.parseCaseStatement()
	case token.DEFAULT:
		// 関数内でのみdefault文を許可するチェック
		if !p.insideFunctionBody {
			p.errors = append(p.errors, fmt.Sprintf("%d行目: default文は関数ブロック内でのみ使用できます。関数定義内で使用してください", p.curToken.Line))
			logger.ParserDebug("関数外でのdefault文使用を検出: エラー報告 (insideFunctionBody=%v)", p.insideFunctionBody)
			return nil
		}
		
		// ネストしたブロック内のdefault文を禁止する追加チェック
		if p.isNestedBlock() {
			p.errors = append(p.errors, fmt.Sprintf("%d行目: default文は関数のルート階層でのみ使用できます。ネストされたブロック内では使用できません", p.curToken.Line))
			logger.ParserDebug("ネストされたブロック内でのdefault文使用を検出: エラー報告")
			return nil
		}
		
		logger.ParserDebug("関数内でのdefault文使用を検出 (insideFunctionBody=%v)", p.insideFunctionBody)
		return p.parseDefaultCaseStatement()
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
	
	// 現在のinsideFunctionBodyフラグを保存
	// ブロック解析中は親のコンテキスト（関数内かどうか）を維持する
	prevInsideFunctionBody := p.insideFunctionBody
	logger.ParserDebug("ブロック解析開始 [insideFunctionBody=%v]", p.insideFunctionBody)

	p.nextToken()

	for !p.curTokenIs(token.RBRACE) && !p.curTokenIs(token.EOF) {
		stmt := p.parseStatement()
		if stmt != nil {
			block.Statements = append(block.Statements, stmt)
		}
		p.nextToken()
	}
	
	// ブロック解析が終わったら元の値に戻す（念のため）
	if p.insideFunctionBody != prevInsideFunctionBody {
		logger.ParserDebug("注意: ブロック解析中にinsideFunctionBodyが変更されました (%v -> %v)", 
			prevInsideFunctionBody, p.insideFunctionBody)
		p.insideFunctionBody = prevInsideFunctionBody
	}
	
	logger.ParserDebug("ブロック解析終了 [insideFunctionBody=%v]", p.insideFunctionBody)
	return block
}

// parseCaseStatement はcase文を解析する
func (p *Parser) parseCaseStatement() *ast.CaseStatement {
	stmt := &ast.CaseStatement{Token: p.curToken}
	logger.ParserDebug("case文の解析開始 at %d:%d [insideFunctionBody=%v]", p.curToken.Line, p.curToken.Column, p.insideFunctionBody)

	// caseの次のトークンを取得
	p.nextToken()
	logger.ParserDebug("case文の条件式の解析開始: 現在のトークン = %s", p.curToken.Literal)

	// 条件式を解析
	stmt.Condition = p.parseExpression(LOWEST)
	if stmt.Condition != nil {
		logger.ParserDebug("case文の条件式の解析完了: %s", stmt.Condition.String())
	} else {
		logger.ParserDebug("case文の条件式の解析に失敗")
		return nil
	}

	// コロンを期待
	if !p.expectPeek(token.COLON) {
		logger.ParserDebug("case文の解析エラー: コロンが見つかりませんでした")
		p.errors = append(p.errors, fmt.Sprintf("%d行目: case文の後にコロンが必要です", p.curToken.Line))
		return nil
	}

	// コロンの次のトークンを取得
	p.nextToken()
	logger.ParserDebug("case文のブロック解析開始")

	// 結果ブロックを解析
	// 関数内のcase文では常にConsequenceフィールドを使用する
	blockStmt := p.parseBlockStatement()
	logger.ParserDebug("case文のブロック解析完了")
	
	// 常にConsequenceフィールドを使用
	stmt.Consequence = blockStmt
	logger.ParserDebug("case文として解析完了")
	
	return stmt
}

// parseDefaultCaseStatement はdefault case文を解析する
func (p *Parser) parseDefaultCaseStatement() *ast.DefaultCaseStatement {
	stmt := &ast.DefaultCaseStatement{Token: p.curToken}
	logger.ParserDebug("default文の解析開始 at %d:%d [insideFunctionBody=%v]", p.curToken.Line, p.curToken.Column, p.insideFunctionBody)
	
	// コロンを期待
	if !p.expectPeek(token.COLON) {
		logger.ParserDebug("default文の解析エラー: コロンが見つかりませんでした")
		p.errors = append(p.errors, fmt.Sprintf("%d行目: default文の後にコロンが必要です", p.curToken.Line))
		return nil
	}
	
	// コロンの次のトークンを取得
	p.nextToken()
	logger.ParserDebug("default文のブロック解析開始")
	
	// ブロックを解析
	stmt.Body = p.parseBlockStatement()
	logger.ParserDebug("default文のブロック解析完了")
	
	logger.ParserDebug("default文として解析完了")
	return stmt
}
