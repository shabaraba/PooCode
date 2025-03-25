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
	return out.String()
}
