package ast

import (
	"bytes"
	"strings"

	"github.com/uncode/token"
)

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
