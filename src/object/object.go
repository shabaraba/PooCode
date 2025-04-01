package object

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
	
	// 特殊な型
	ANY_OBJ          = "ANY"     // どの型でも受け付ける
)

// Object はすべての値のインターフェース
type Object interface {
	Type() ObjectType
	Inspect() string
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

// Statement は文を表す
type Statement interface{}

// Expression は式を表す
type Expression interface{}

// BlockStatement は関数のボディを表す
type BlockStatement struct {
	Statements []Statement
}

// Identifier は識別子を表す
type Identifier struct {
	Value string
}
