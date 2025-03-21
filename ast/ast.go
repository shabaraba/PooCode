package ast

import (
	"bytes"
	"strings"

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
	Token     token.Token // '|>' または '|' トークン
	Left      Expression
	Right     Expression
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

// Identifier は識別子を表すノード
type Identifier struct {
	Token token.Token // IDENT トークン
	Value string
}

func (i *Identifier) expressionNode()      {}
func (i *Identifier) TokenLiteral() string { return i.Token.Literal }
func (i *Identifier) String() string       { return i.Value }

// IntegerLiteral は整数リテラルを表すノード
type IntegerLiteral struct {
	Token token.Token // INT トークン
	Value int64
}

func (il *IntegerLiteral) expressionNode()      {}
func (il *IntegerLiteral) TokenLiteral() string { return il.Token.Literal }
func (il *IntegerLiteral) String() string       { return il.Token.Literal }

// FloatLiteral は浮動小数点リテラルを表すノード
type FloatLiteral struct {
	Token token.Token // FLOAT トークン
	Value float64
}

func (fl *FloatLiteral) expressionNode()      {}
func (fl *FloatLiteral) TokenLiteral() string { return fl.Token.Literal }
func (fl *FloatLiteral) String() string       { return fl.Token.Literal }

// StringLiteral は文字列リテラルを表すノード
type StringLiteral struct {
	Token token.Token // STRING トークン
	Value string
}

func (sl *StringLiteral) expressionNode()      {}
func (sl *StringLiteral) TokenLiteral() string { return sl.Token.Literal }
func (sl *StringLiteral) String() string       { return "\"" + sl.Value + "\"" }

// BooleanLiteral は真偽値リテラルを表すノード
type BooleanLiteral struct {
	Token token.Token // true または false トークン
	Value bool
}

func (bl *BooleanLiteral) expressionNode()      {}
func (bl *BooleanLiteral) TokenLiteral() string { return bl.Token.Literal }
func (bl *BooleanLiteral) String() string {
	if bl.Value {
		return "true"
	}
	return "false"
}

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

// ArrayLiteral は配列リテラルを表すノード
type ArrayLiteral struct {
	Token    token.Token // '[' トークン
	Elements []Expression
}

func (al *ArrayLiteral) expressionNode()      {}
func (al *ArrayLiteral) TokenLiteral() string { return al.Token.Literal }
func (al *ArrayLiteral) String() string {
	var out bytes.Buffer
	elements := []string{}
	for _, el := range al.Elements {
		elements = append(elements, el.String())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
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

// ClassLiteral はクラス定義を表すノード
type ClassLiteral struct {
	Token      token.Token // 'class' トークン
	Name       *Identifier
	Properties []*PropertyDefinition
	Methods    []*FunctionLiteral
	Extends    *Identifier // 継承元クラス
}

func (cl *ClassLiteral) expressionNode()      {}
func (cl *ClassLiteral) TokenLiteral() string { return cl.Token.Literal }
func (cl *ClassLiteral) String() string {
	var out bytes.Buffer
	out.WriteString(cl.TokenLiteral() + " ")
	out.WriteString(cl.Name.String())
	if cl.Extends != nil {
		out.WriteString(" extends " + cl.Extends.String())
	}
	out.WriteString(" {\n")
	for _, p := range cl.Properties {
		out.WriteString("  " + p.String() + "\n")
	}
	for _, m := range cl.Methods {
		out.WriteString("  " + m.String() + "\n")
	}
	out.WriteString("}")
	return out.String()
}

// PropertyDefinition はクラスのプロパティ定義を表すノード
type PropertyDefinition struct {
	Token      token.Token // 'public' または 'private' トークン
	Name       *Identifier
	Type       string
	Visibility string // "public" または "private"
}

func (pd *PropertyDefinition) expressionNode()      {}
func (pd *PropertyDefinition) TokenLiteral() string { return pd.Token.Literal }
func (pd *PropertyDefinition) String() string {
	var out bytes.Buffer
	out.WriteString(pd.Visibility + " ")
	if pd.Type != "" {
		out.WriteString(pd.Type + " ")
	}
	out.WriteString(pd.Name.String())
	return out.String()
}

// EnumLiteral は列挙型定義を表すノード
type EnumLiteral struct {
	Token  token.Token // 'enum' トークン
	Name   *Identifier
	Values []*Identifier
}

func (el *EnumLiteral) expressionNode()      {}
func (el *EnumLiteral) TokenLiteral() string { return el.Token.Literal }
func (el *EnumLiteral) String() string {
	var out bytes.Buffer
	values := []string{}
	for _, v := range el.Values {
		values = append(values, v.String())
	}
	out.WriteString(el.TokenLiteral() + " ")
	out.WriteString(el.Name.String())
	out.WriteString(" {\n")
	out.WriteString("  " + strings.Join(values, "\n  "))
	out.WriteString("\n}")
	return out.String()
}

// PizzaLiteral は🍕リテラルを表すノード
type PizzaLiteral struct {
	Token token.Token // '🍕' トークン
}

func (pl *PizzaLiteral) expressionNode()      {}
func (pl *PizzaLiteral) TokenLiteral() string { return pl.Token.Literal }
func (pl *PizzaLiteral) String() string       { return "🍕" }

// PooLiteral は💩リテラルを表すノード
type PooLiteral struct {
	Token token.Token // '💩' トークン
}

func (pl *PooLiteral) expressionNode()      {}
func (pl *PooLiteral) TokenLiteral() string { return pl.Token.Literal }
func (pl *PooLiteral) String() string       { return "💩" }

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
