package object

import (
	"fmt"
)

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

// GetNextFunction は同じ名前の別の関数を取得する
// 特に、条件付き関数が false を返した場合に次の候補を探すために使用
func (e *Environment) GetNextFunction(name string, currentFn *Function) *Function {
	// フラットな関数リストを作成
	var functions []*Function
	var collectFunctions func(*Environment)
	
	collectFunctions = func(env *Environment) {
		// 現在の環境から同名の関数を探す
		if obj, ok := env.store[name]; ok {
			if fn, ok := obj.(*Function); ok {
				// 現在の関数と異なる関数だけを追加
				if fn != currentFn {
					functions = append(functions, fn)
				}
			}
		}
		
		// 外部環境も検索
		if env.outer != nil {
			collectFunctions(env.outer)
		}
	}
	
	collectFunctions(e)
	
	// 見つかった関数がなければnilを返す
	if len(functions) == 0 {
		return nil
	}
	
	// 最初に見つかった別の関数を返す
	return functions[0]
}

// GetAllFunctionsByName は同名のすべての関数を取得する
// 条件付き関数を正しく呼び出すために使用
func (e *Environment) GetAllFunctionsByName(name string) []*Function {
	var functions []*Function
	var collectFunctions func(*Environment)
	
	// 既に処理済みの関数オブジェクトを記録するマップ
	// ポインタ比較で重複を排除
	processed := make(map[*Function]bool)
	
	collectFunctions = func(env *Environment) {
		// 現在の環境から同名の関数を探す
		if obj, ok := env.store[name]; ok {
			// 単一のオブジェクト（通常の値やシングルトン関数）の場合
			if fn, ok := obj.(*Function); ok {
				// まだ処理済みでなければ追加
				if !processed[fn] {
					processed[fn] = true
					functions = append(functions, fn)
				}
			}
		}
		
		// 条件付き関数は複数定義できるようにするために特殊な名前で保存されている可能性がある
		// 例: "name#1", "name#2" などの形式で同じname向けの複数の関数が定義されている可能性
		// そのようなキーも検索し、関数オブジェクトを取得
		for key, obj := range env.store {
			// 「name#」で始まるエントリを探す
			if len(key) > len(name) && key[:len(name)+1] == name+"#" {
				if fn, ok := obj.(*Function); ok {
					// まだ処理済みでなければ追加
					if !processed[fn] {
						processed[fn] = true
						functions = append(functions, fn)
					}
				}
			}
		}
		
		// 外部環境も検索
		if env.outer != nil {
			collectFunctions(env.outer)
		}
	}
	
	collectFunctions(e)
	
	// デバッグ情報
	fmt.Printf("関数 '%s' の候補を %d 個見つけました\n", name, len(functions))
	for i, fn := range functions {
		// 条件の有無を表示
		hasCondition := "なし"
		if fn.Condition != nil {
			hasCondition = "あり"
		}
		fmt.Printf("  関数候補 %d: 条件=%s\n", i+1, hasCondition)
	}
	
	return functions
}

// GetVariables は環境内のすべての変数を取得する（デバッグ用）
func (e *Environment) GetVariables() map[string]Object {
	// 現在の環境のすべての変数を取得
	vars := make(map[string]Object)
	
	// まず現在の環境の変数をコピー
	for k, v := range e.store {
		vars[k] = v
	}
	
	// 外部環境の変数も取得（ただし、内部環境で上書きされている場合は取得しない）
	if e.outer != nil {
		for k, v := range e.outer.GetVariables() {
			if _, exists := vars[k]; !exists {
				vars[k] = v
			}
		}
	}
	
	return vars
}
