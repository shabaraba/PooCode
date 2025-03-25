package object

import (
	"fmt"
	"hash/fnv"
)

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
