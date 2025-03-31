package ast

import (
	"bytes"
	"strings"

	"github.com/uncode/token"
)

// Identifier は識別子を表すノード
type Identifier struct {
	Token token.Token // IDENT トークン
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

// PrefixExpression は前置式を表すノード
type PrefixExpression struct {
	Token    token.Token // 前置演算子トークン
	Operator string
	Right    Expression
}

func (pe *PrefixExpression) expressionNode()      {}
func (pe *PrefixExpression) TokenLiteral() string { return pe.Token.Literal }
func (pe *PrefixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(pe.Operator)
	out.WriteString(pe.Right.String())
	out.WriteString(")")
	return out.String()
}

// InfixExpression は中置式を表すノード
type InfixExpression struct {
	Token    token.Token // 演算子トークン
	Left     Expression
	Operator string
	Right    Expression
}

func (ie *InfixExpression) expressionNode()      {}
func (ie *InfixExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *InfixExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString(" " + ie.Operator + " ")
	out.WriteString(ie.Right.String())
	out.WriteString(")")
	return out.String()
}

// CallExpression は関数呼び出しを表すノード
type CallExpression struct {
	Token     token.Token // '(' トークン
	Function  Expression  // 関数名または関数リテラル
	Arguments []Expression
}

func (ce *CallExpression) expressionNode()      {}
func (ce *CallExpression) TokenLiteral() string { return ce.Token.Literal }
func (ce *CallExpression) String() string {
	var out bytes.Buffer
	args := []string{}
	for _, a := range ce.Arguments {
		args = append(args, a.String())
	}
	out.WriteString(ce.Function.String())
	out.WriteString("(")
	out.WriteString(strings.Join(args, ", "))
	out.WriteString(")")
	return out.String()
}

// PropertyAccessExpression はプロパティアクセスを表すノード
type PropertyAccessExpression struct {
	Token    token.Token // '.' または "'s" トークン
	Object   Expression
	Property Expression
}

func (pa *PropertyAccessExpression) expressionNode()      {}
func (pa *PropertyAccessExpression) TokenLiteral() string { return pa.Token.Literal }
func (pa *PropertyAccessExpression) String() string {
	var out bytes.Buffer
	out.WriteString(pa.Object.String())
	if pa.Token.Type == token.APOSTROPHE_S {
		out.WriteString("'s ")
	} else {
		out.WriteString(".")
	}
	out.WriteString(pa.Property.String())
	return out.String()
}

// IndexExpression は配列の添字アクセスを表すノード
type IndexExpression struct {
	Token token.Token // '[' トークン
	Left  Expression
	Index Expression
}

func (ie *IndexExpression) expressionNode()      {}
func (ie *IndexExpression) TokenLiteral() string { return ie.Token.Literal }
func (ie *IndexExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	out.WriteString(ie.Left.String())
	out.WriteString("[")
	out.WriteString(ie.Index.String())
	out.WriteString("])")
	return out.String()
}

// RangeExpression は範囲式を表すノード
type RangeExpression struct {
	Token token.Token // '..' トークン
	Start Expression  // 開始値
	End   Expression  // 終了値
}

func (re *RangeExpression) expressionNode()      {}
func (re *RangeExpression) TokenLiteral() string { return re.Token.Literal }
func (re *RangeExpression) String() string {
	var out bytes.Buffer
	out.WriteString("(")
	if re.Start != nil {
		out.WriteString(re.Start.String())
	}
	out.WriteString("..")
	if re.End != nil {
		out.WriteString(re.End.String())
	}
	out.WriteString(")")
	return out.String()
}

// BlockExpression はブロック式を表すノード
type BlockExpression struct {
	Token token.Token     // '{' トークン
	Block *BlockStatement // ブロック内のステートメント
}

func (be *BlockExpression) expressionNode()      {}
func (be *BlockExpression) TokenLiteral() string { return be.Token.Literal }
func (be *BlockExpression) String() string {
	var out bytes.Buffer
	out.WriteString("{ ")
	out.WriteString(be.Block.String())
	out.WriteString(" }")
	return out.String()
}
