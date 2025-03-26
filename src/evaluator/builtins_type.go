package evaluator

import (
	"github.com/uncode/object"
)

// registerTypeBuiltins は型関連の組み込み関数を登録する
func registerTypeBuiltins() {
	// 等価性判定関数
	Builtins["eq"] = &object.Builtin{
		Name: "eq",
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return createError("eq関数は2つの引数が必要です: %d個与えられました", len(args))
			}
			
			switch left := args[0].(type) {
			case *object.Integer:
				if right, ok := args[1].(*object.Integer); ok {
					return &object.Boolean{Value: left.Value == right.Value}
				}
			case *object.String:
				if right, ok := args[1].(*object.String); ok {
					return &object.Boolean{Value: left.Value == right.Value}
				}
			case *object.Boolean:
				if right, ok := args[1].(*object.Boolean); ok {
					return &object.Boolean{Value: left.Value == right.Value}
				}
			}
			
			return &object.Boolean{Value: false}
		},
		ReturnType: object.BOOLEAN_OBJ,
		ParamTypes: []object.ObjectType{object.ANY_OBJ, object.ANY_OBJ},
	}

	// 論理否定関数
	Builtins["not"] = &object.Builtin{
		Name: "not",
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return createError("not関数は1つの引数が必要です: %d個与えられました", len(args))
			}
			
			if b, ok := args[0].(*object.Boolean); ok {
				return &object.Boolean{Value: !b.Value}
			}
			
			return &object.Boolean{Value: false} // デフォルトはfalse
		},
		ReturnType: object.BOOLEAN_OBJ,
		ParamTypes: []object.ObjectType{object.BOOLEAN_OBJ},
	}

	// 型情報取得関数
	Builtins["typeof"] = &object.Builtin{
		Name: "typeof",
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return createError("typeof関数は1つの引数が必要です: %d個与えられました", len(args))
			}

			// 引数が文字列の場合、組み込み関数名として解釈
			if str, ok := args[0].(*object.String); ok {
				funcName := str.Value
				if builtin, exists := Builtins[funcName]; exists {
					return &object.String{Value: string(builtin.ReturnType)}
				}
				return createError("組み込み関数 '%s' は存在しません", funcName)
			}

			// その他の型はそのまま型情報を返す
			return &object.String{Value: string(args[0].Type())}
		},
		ReturnType: object.STRING_OBJ,
		ParamTypes: []object.ObjectType{object.ANY_OBJ},
	}
}
