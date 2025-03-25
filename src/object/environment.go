package object

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
