package object

import (
	"fmt"
	"hash/fnv"
	"math"
	"strconv"
)

const (
	INTEGER_OBJ  = "INTEGER"
	BOOLEAN_OBJ  = "BOOLEAN"
	NULL_OBJ     = "NULL"
	STRING_OBJ   = "STRING"
	FLOAT_OBJ    = "FLOAT"
)

// Integer represents an integer object
type Integer struct {
	Value int64
}

func (i *Integer) Type() ObjectType { return INTEGER_OBJ }
func (i *Integer) Inspect() string  { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(strconv.FormatInt(i.Value, 10)))
	return HashKey{Type: i.Type(), Value: h.Sum64()}
}

// Boolean represents a boolean object
type Boolean struct {
	Value bool
}

func (b *Boolean) Type() ObjectType { return BOOLEAN_OBJ }
func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }
func (b *Boolean) HashKey() HashKey {
	var value uint64
	if b.Value {
		value = 1
	} else {
		value = 0
	}
	return HashKey{Type: b.Type(), Value: value}
}

// Null represents a null object
type Null struct{}

func (n *Null) Type() ObjectType { return NULL_OBJ }
func (n *Null) Inspect() string  { return "null" }

// String represents a string object
type String struct {
	Value string
}

func (s *String) Type() ObjectType { return STRING_OBJ }
func (s *String) Inspect() string  { return s.Value }
func (s *String) HashKey() HashKey {
	h := fnv.New64a()
	h.Write([]byte(s.Value))
	return HashKey{Type: s.Type(), Value: h.Sum64()}
}

// Float represents a float object
type Float struct {
	Value float64
}

func (f *Float) Type() ObjectType { return FLOAT_OBJ }
func (f *Float) Inspect() string {
	// Format float with maximum precision, trimming unnecessary zeros
	return strconv.FormatFloat(f.Value, 'g', -1, 64)
}
func (f *Float) HashKey() HashKey {
	// Use the bit representation of the float for hashing
	bits := math.Float64bits(f.Value)
	return HashKey{Type: f.Type(), Value: bits}
}
