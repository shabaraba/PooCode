package ast

import (
	"bytes"

	"github.com/uncode/token"
)

// ExpressionStatement は式文を表すノード
type ExpressionStatement struct {
	Token      token.Token // 式の最初のトークン
	Expression Expression
}

func (es *ExpressionStatement) statementNode()       {}
func (es *ExpressionStatement) TokenLiteral() string { return es.Token.Literal }
func (es *ExpressionStatement) String() string {
	if es.Expression != nil {
		return es.Expression.String()
	}
	return ""
}

// AssignStatement は代入文を表すノード
type AssignStatement struct {
	Token token.Token // '>>' トークン
	Left  Expression
	Value Expression
}

func (as *AssignStatement) statementNode()       {}
func (as *AssignStatement) TokenLiteral() string { return as.Token.Literal }
func (as *AssignStatement) String() string {
	var out bytes.Buffer
	out.WriteString(as.Left.String())
	out.WriteString(" >> ")
	out.WriteString(as.Value.String())
	return out.String()
}

// PipeStatement はパイプライン文を表すノード
type PipeStatement struct {
	Token      token.Token // '|>' または '|' トークン
	Left       Expression
	Right      Expression
	IsParallel bool
}

func (ps *PipeStatement) statementNode()       {}
func (ps *PipeStatement) TokenLiteral() string { return ps.Token.Literal }
func (ps *PipeStatement) String() string {
	var out bytes.Buffer
	out.WriteString(ps.Left.String())
	if ps.IsParallel {
		out.WriteString(" | ")
	} else {
		out.WriteString(" |> ")
	}
	out.WriteString(ps.Right.String())
	return out.String()
}

// BlockStatement はブロック文を表すノード
type BlockStatement struct {
	Token      token.Token // '{' トークン
	Statements []Statement
}

func (bs *BlockStatement) statementNode()       {}
func (bs *BlockStatement) TokenLiteral() string { return bs.Token.Literal }
func (bs *BlockStatement) String() string {
	var out bytes.Buffer
	out.WriteString("{ ")
	for _, s := range bs.Statements {
		out.WriteString(s.String())
	}
	out.WriteString(" }")
	return out.String()
}

// GlobalStatement はグローバル変数宣言を表すノード
type GlobalStatement struct {
	Token token.Token // 'global' トークン
	Name  *Identifier
	Type  string
}

func (gs *GlobalStatement) statementNode()       {}
func (gs *GlobalStatement) TokenLiteral() string { return gs.Token.Literal }
func (gs *GlobalStatement) String() string {
	var out bytes.Buffer
	out.WriteString(gs.TokenLiteral() + " ")
	if gs.Type != "" {
		out.WriteString(gs.Type + " ")
	}
	out.WriteString(gs.Name.String())
	return out.String()
}
