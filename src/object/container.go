package object

import (
	"bytes"
	"fmt"
	"strings"
)

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
