package object

import (
	"fmt"
	"hash/fnv"
	"strings"
	"bytes"
)

// ObjectType は値の型を表す
type ObjectType string

const (
	INTEGER_OBJ      = "INTEGER"
	FLOAT_OBJ        = "FLOAT"
	BOOLEAN_OBJ      = "BOOLEAN"
	STRING_OBJ       = "STRING"
	NULL_OBJ         = "NULL"
	RETURN_VALUE_OBJ = "RETURN_VALUE"
	ERROR_OBJ        = "ERROR"
	FUNCTION_OBJ     = "FUNCTION"
	BUILTIN_OBJ      = "BUILTIN"
	ARRAY_OBJ        = "ARRAY"
	HASH_OBJ         = "HASH"
	CLASS_OBJ        = "CLASS"
	INSTANCE_OBJ     = "INSTANCE"
)

// Object はすべての値のインターフェース
type Object interface {
	Type() ObjectType
	Inspect() string
	GetPooValue() Object // 💩メンバの値を取得
	SetPooValue(Object)  // 💩メンバの値を設定
}

// Hashable はハッシュキーとして使用可能なオブジェクトのインターフェース
type Hashable interface {
	HashKey() HashKey
}

// HashKey はハッシュマップのキーを表す
type HashKey struct {
	Type  ObjectType
	Value uint64
}

// Integer は整数値を表す
type Integer struct {
	Value int64
	Poo   Object // 💩メンバ
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) GetPooValue() Object {
	if i.Poo == nil {
		i.Poo = i // デフォルトでは自分自身
	}
	return i.Poo
}
func (i *Integer) SetPooValue(val Object) { i.Poo = val }
func (i *Integer) HashKey() HashKey {
	return HashKey{Type: i.Type(), Value: uint64(i.Value)}
}

// Float は浮動小数点値を表す
type Float struct {
	Value float64
	Poo   Object // 💩メンバ
}

func (f *Float) Type() ObjectType { return FLOAT_OBJ }
func (f *Float) Inspect() string  { return fmt.Sprintf("%g", f.Value) }
func (f *Float) GetPooValue() Object {
	if f.Poo == nil {
		f.Poo = f // デフォルトでは自分自身
	}
	return f.Poo
}
func (f *Float) SetPooValue(val Object) { f.Poo = val }
func (f *Float) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(fmt.Sprintf("%g", f.Value)))
	return HashKey{Type: f.Type(), Value: h.Sum64()}
}

// Boolean は真偽値を表す
type Boolean struct {
	Value bool
	Poo   Object // 💩メンバ
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string {
	if b.Value {
		return "true"
	}
	return "false"
}
func (b *Boolean) GetPooValue() Object {
	if b.Poo == nil {
		b.Poo = b // デフォルトでは自分自身
	}
	return b.Poo
}
func (b *Boolean) SetPooValue(val Object) { b.Poo = val }
func (b *Boolean) HashKey() HashKey {
	var value uint64
	if b.Value {
		value = 1
	} else {
		value = 0
	}
	return HashKey{Type: b.Type(), Value: value}
}

// String は文字列値を表す
type String struct {
	Value string
	Poo   Object // 💩メンバ
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return s.Value }
func (s *String) GetPooValue() Object {
	if s.Poo == nil {
		s.Poo = s // デフォルトでは自分自身
	}
	return s.Poo
}
func (s *String) SetPooValue(val Object) { s.Poo = val }
func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))
	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

// Null はnull値を表す
type Null struct {
	Poo Object // 💩メンバ
}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "null" }
func (n *Null) GetPooValue() Object {
	if n.Poo == nil {
		n.Poo = n // デフォルトでは自分自身
	}
	return n.Poo
}
func (n *Null) SetPooValue(val Object) { n.Poo = val }

// ReturnValue はreturn文の値を表す
type ReturnValue struct {
	Value Object
	Poo   Object // 💩メンバ
}

func (rv *ReturnValue) Type() ObjectType { return RETURN_VALUE_OBJ }
func (rv *ReturnValue) Inspect() string  { return rv.Value.Inspect() }
func (rv *ReturnValue) GetPooValue() Object {
	if rv.Poo == nil {
		rv.Poo = rv.Value // デフォルトでは内部の値
	}
	return rv.Poo
}
func (rv *ReturnValue) SetPooValue(val Object) { rv.Poo = val }

// Error はエラー値を表す
type Error struct {
	Message string
	Poo     Object // 💩メンバ
}

func (e *Error) Type() ObjectType { return ERROR_OBJ }
func (e *Error) Inspect() string  { return "ERROR: " + e.Message }
func (e *Error) GetPooValue() Object {
	if e.Poo == nil {
		e.Poo = e // デフォルトでは自分自身
	}
	return e.Poo
}
func (e *Error) SetPooValue(val Object) { e.Poo = val }

// Environment は変数環境を表す
type Environment struct {
	store map[string]Object
	outer *Environment
}

// NewEnvironment は新しい環境を生成する
func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s, outer: nil}
}

// NewEnclosedEnvironment は外部環境を持つ新しい環境を生成する
func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

// Get は環境から変数を取得する
func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

// Set は環境に変数を設定する
func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}

// BlockStatement は関数のボディを表す
type BlockStatement struct {
	Statements []Statement
}

// Statement は文を表す
type Statement interface{}

// Identifier は識別子を表す
type Identifier struct {
	Value string
}

// Function は関数を表す
type Function struct {
	Parameters []*Identifier
	ASTBody    interface{} // ASTのBlockStatementを保持するフィールド
	Env        *Environment
	InputType  string
	ReturnType string
	Condition  interface{} // 条件式
	Poo        Object      // 💩メンバ
}

// Expression は式を表す
type Expression interface{}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	return fmt.Sprintf("function with %d params", len(f.Parameters))
}
func (f *Function) GetPooValue() Object {
	if f.Poo == nil {
		f.Poo = f // デフォルトでは自分自身
	}
	return f.Poo
}
func (f *Function) SetPooValue(val Object) { f.Poo = val }

// BuiltinFunction は組み込み関数の型
type BuiltinFunction func(args ...Object) Object

// Builtin は組み込み関数を表す
type Builtin struct {
	Fn  BuiltinFunction
	Poo Object // 💩メンバ
}

func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string  { return "builtin function" }
func (b *Builtin) GetPooValue() Object {
	if b.Poo == nil {
		b.Poo = b // デフォルトでは自分自身
	}
	return b.Poo
}
func (b *Builtin) SetPooValue(val Object) { b.Poo = val }

// Array は配列を表す
type Array struct {
	Elements []Object
	Poo      Object // 💩メンバ
}

func (a *Array) Type() ObjectType { return ARRAY_OBJ }
func (a *Array) Inspect() string {
	var out bytes.Buffer
	elements := []string{}
	for _, e := range a.Elements {
		elements = append(elements, e.Inspect())
	}
	out.WriteString("[")
	out.WriteString(strings.Join(elements, ", "))
	out.WriteString("]")
	return out.String()
}
func (a *Array) GetPooValue() Object {
	if a.Poo == nil {
		a.Poo = a // デフォルトでは自分自身
	}
	return a.Poo
}
func (a *Array) SetPooValue(val Object) { a.Poo = val }

// HashPair はハッシュマップの値のペアを表す
type HashPair struct {
	Key   Object
	Value Object
}

// Hash はハッシュマップを表す
type Hash struct {
	Pairs map[HashKey]HashPair
	Poo   Object // 💩メンバ
}

func (h *Hash) Type() ObjectType { return HASH_OBJ }
func (h *Hash) Inspect() string {
	var out bytes.Buffer
	pairs := []string{}
	for _, pair := range h.Pairs {
		pairs = append(pairs, fmt.Sprintf("%s: %s",
			pair.Key.Inspect(), pair.Value.Inspect()))
	}
	out.WriteString("{")
	out.WriteString(strings.Join(pairs, ", "))
	out.WriteString("}")
	return out.String()
}
func (h *Hash) GetPooValue() Object {
	if h.Poo == nil {
		h.Poo = h // デフォルトでは自分自身
	}
	return h.Poo
}
func (h *Hash) SetPooValue(val Object) { h.Poo = val }

// Class はクラスを表す
type Class struct {
	Name       string
	Properties map[string]*PropertyDefinition
	Methods    map[string]*Function
	Extends    *Class // 継承元クラス
	Poo        Object // 💩メンバ
}

func (c *Class) Type() ObjectType { return CLASS_OBJ }
func (c *Class) Inspect() string  { return fmt.Sprintf("class %s", c.Name) }
func (c *Class) GetPooValue() Object {
	if c.Poo == nil {
		c.Poo = c // デフォルトでは自分自身
	}
	return c.Poo
}
func (c *Class) SetPooValue(val Object) { c.Poo = val }

// PropertyDefinition はクラスのプロパティ定義を表す
type PropertyDefinition struct {
	Name       string
	Type       string
	Visibility string // "public" または "private"
}

// Instance はクラスのインスタンスを表す
type Instance struct {
	Class      *Class
	Properties map[string]Object
	Poo        Object // 💩メンバ
}

func (i *Instance) Type() ObjectType { return INSTANCE_OBJ }
func (i *Instance) Inspect() string {
	return fmt.Sprintf("instance of %s", i.Class.Name)
}
func (i *Instance) GetPooValue() Object {
	if i.Poo == nil {
		i.Poo = i // デフォルトでは自分自身
	}
	return i.Poo
}
func (i *Instance) SetPooValue(val Object) { i.Poo = val }
func (i *Instance) GetProperty(name string) (Object, bool) {
	// まず自身のプロパティを検索
	if val, ok := i.Properties[name]; ok {
		propDef, exists := i.Class.Properties[name]
		if !exists || propDef.Visibility == "public" {
			return val, true
		}
	}

	// 継承元クラスのプロパティを検索
	class := i.Class
	for class.Extends != nil {
		class = class.Extends
		propDef, exists := class.Properties[name]
		if exists && propDef.Visibility == "public" {
			if val, ok := i.Properties[name]; ok {
				return val, true
			}
		}
	}

	return nil, false
}
func (i *Instance) SetProperty(name string, val Object) bool {
	// プロパティの存在と可視性をチェック
	propDef, exists := i.Class.Properties[name]
	if exists {
		i.Properties[name] = val
		return true
	}

	// 継承元クラスのプロパティをチェック
	class := i.Class
	for class.Extends != nil {
		class = class.Extends
		propDef, exists = class.Properties[name]
		if exists && propDef.Visibility == "public" {
			i.Properties[name] = val
			return true
		}
	}

	return false
}
func (i *Instance) GetMethod(name string) (*Function, bool) {
	// まず自身のクラスのメソッドを検索
	if method, ok := i.Class.Methods[name]; ok {
		return method, true
	}

	// 継承元クラスのメソッドを検索
	class := i.Class
	for class.Extends != nil {
		class = class.Extends
		if method, ok := class.Methods[name]; ok {
			return method, true
		}
	}

	return nil, false
}
