package object

import (
	"bytes"
	"fmt"
	"strings"
)

// ObjectType represents the different types of objects in the language
type ObjectType string

// Object is the base interface for all object types in the language
type Object interface {
	Type() ObjectType
	Inspect() string
}

// Hashable is an interface for objects that can be used as hash keys
type Hashable interface {
	HashKey() HashKey
}

// HashKey represents a unique key for hashable objects
type HashKey struct {
	Type  ObjectType
	Value uint64
}

// Stringer is an interface for objects that can be converted to a string
type Stringer interface {
	String() string
}
