package ast

import (
	"bytes"

	"github.com/uncode/token"
)

// Node はすべてのAST（抽象構文木）ノードのインターフェース
type Node interface {
	TokenLiteral() string
	String() string
}

// Statement は文を表すノード
type Statement interface {
	Node
	statementNode()
}

// Expression は式を表すノード
type Expression interface {
	Node
	expressionNode()
}

// Program はプログラム全体を表すルートノード
type Program struct {
	Statements []Statement
}

func (p *Program) TokenLiteral() string {
	if len(p.Statements) > 0 {
		return p.Statements[0].TokenLiteral()
	}
	return ""
}

func (p *Program) String() string {
	var out bytes.Buffer
	for _, s := range p.Statements {
		out.WriteString(s.String())
	}
	return out.String()
}
