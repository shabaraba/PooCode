package ast

import (
	"bytes"
	"strings"

	"github.com/uncode/token"
)

// FunctionLiteral は関数リテラルを表すノード
type FunctionLiteral struct {
	Token      token.Token // 'def' トークン
	Name       *Identifier
	Parameters []*Identifier
	Body       *BlockStatement
	ReturnType string
	InputType  string
	Condition  Expression // 条件付き関数定義の条件部分
	Cases      []*CaseStatement // case文のリスト
}

func (fl *FunctionLiteral) expressionNode()      {}
func (fl *FunctionLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FunctionLiteral) String() string {
	var out bytes.Buffer
	params := []string{}
	for _, p := range fl.Parameters {
		params = append(params, p.String())
	}
	out.WriteString(fl.TokenLiteral() + " ")
	if fl.Name != nil {
		out.WriteString(fl.Name.String())
	}
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(")")
	if fl.InputType != "" || fl.ReturnType != "" {
		out.WriteString(": " + fl.InputType + " -> " + fl.ReturnType)
	}
	if fl.Condition != nil {
		out.WriteString(" if ")
		out.WriteString(fl.Condition.String())
	}
	out.WriteString(" ")
	out.WriteString(fl.Body.String())
	for _, cs := range fl.Cases {
		out.WriteString(cs.String())
	}
	return out.String()
}

// CaseStatement はcase文を表すノード
type CaseStatement struct {
	Token       token.Token // 'case' トークン
	Condition   Expression
	Consequence *BlockStatement
	Body        *BlockStatement // 関数外のcase文で使用
}

func (cs *CaseStatement) statementNode()       {}
func (cs *CaseStatement) expressionNode()      {}
func (cs *CaseStatement) TokenLiteral() string { return cs.Token.Literal }
func (cs *CaseStatement) String() string {
	var out bytes.Buffer
	out.WriteString(cs.TokenLiteral() + " ")
	out.WriteString(cs.Condition.String())
	out.WriteString(": ")
	if cs.Consequence != nil {
		out.WriteString(cs.Consequence.String())
	} else if cs.Body != nil {
		out.WriteString(cs.Body.String())
	}
	return out.String()
}
