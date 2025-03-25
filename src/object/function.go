package object

import (
	"fmt"
)

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

// Name は関数の名前を取得する
// 環境内で定義されている関数名を特定する必要がある場合に使用
func (f *Function) Name() (string, bool) {
	// 環境をスキャンして関数オブジェクトに対応する名前を見つける
	for name, obj := range f.Env.store {
		if obj == f {
			return name, true
		}
	}
	
	// 外部環境も検索
	if f.Env.outer != nil {
		env := f.Env.outer
		for env != nil {
			for name, obj := range env.store {
				if obj == f {
					return name, true
				}
			}
			env = env.outer
		}
	}
	
	return "", false
}

// BuiltinFunction は組み込み関数の型
type BuiltinFunction func(args ...Object) Object

// Builtin は組み込み関数を表す
type Builtin struct {
	Name string           // 関数名
	Fn   BuiltinFunction
	Poo  Object // 💩メンバ
}

func (b *Builtin) Type() ObjectType { return BUILTIN_OBJ }
func (b *Builtin) Inspect() string {
	if b.Name != "" {
		return fmt.Sprintf("builtin function: %s", b.Name)
	}
	return "builtin function"
}
func (b *Builtin) GetPooValue() Object {
	if b.Poo == nil {
		b.Poo = b // デフォルトでは自分自身
	}
	return b.Poo
}
func (b *Builtin) SetPooValue(val Object) { b.Poo = val }

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
