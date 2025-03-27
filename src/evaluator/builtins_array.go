package evaluator

import (
	"github.com/uncode/object"
	"github.com/uncode/logger"
)

func registerArrayBuiltins() {
	// map function
	Builtins["map"] = &object.Builtin{
		Name: "map",
		Fn: func(args ...object.Object) object.Object {
			// 最低2つの引数が必要（配列と関数）
			if len(args) < 2 {
				return createError("map function requires at least 2 arguments (array and function)")
			}
			
			// 第1引数は配列でないといけない
			arr, ok := args[0].(*object.Array)
			if !ok {
				return createError("First argument to map must be an array")
			}
			
			// 第2引数は関数（ユーザー定義関数またはビルトイン関数）
			var mapFn func([]object.Object) object.Object
			
			// 関数の引数（mapへの追加引数があれば保存）
			extraArgs := args[2:]
			
			switch fn := args[1].(type) {
			case *object.Function:
				// ユーザー定義関数
				mapFn = func(elemArgs []object.Object) object.Object {
					extendedEnv := object.NewEnclosedEnvironment(fn.Env)
					
					// 必ず最初の引数を🍕に設定
					if len(elemArgs) > 0 {
						extendedEnv.Set("🍕", elemArgs[0])
					}
					
					// 関数が引数を持つ場合、引数を設定
					if len(fn.Parameters) > 0 && len(elemArgs) > 1 {
						// 引数の数を確認
						paramCount := len(fn.Parameters)
						if len(elemArgs)-1 < paramCount {
							logger.Debug("関数の引数が少なすぎます: 期待=%d, 実際=%d", paramCount, len(elemArgs)-1)
						}
						
						// 引数をバインド（🍕の次の引数から）
						for i := 0; i < paramCount && i+1 < len(elemArgs); i++ {
							extendedEnv.Set(fn.Parameters[i].Value, elemArgs[i+1])
						}
					}
					
					// 関数本体を評価
					result := Eval(fn.ASTBody, extendedEnv)
					
					// エラー処理
					if errObj, ok := result.(*object.Error); ok {
						return errObj
					}
					
					// 戻り値の処理
					if retVal, ok := result.(*object.ReturnValue); ok {
						result = retVal.Value
					}
					
					return result
				}
			case *object.Builtin:
				// ビルトイン関数
				mapFn = func(elemArgs []object.Object) object.Object {
					return fn.Fn(elemArgs...)
				}
			default:
				return createError("Second argument to map must be a function")
			}
			
			// マップ処理の実行
			resultElements := make([]object.Object, 0, len(arr.Elements))
			
			for _, elem := range arr.Elements {
				// 現在の要素と追加引数を組み合わせる
				elemArgs := []object.Object{elem}
				elemArgs = append(elemArgs, extraArgs...)
				
				// 各要素に関数を適用
				resultElements = append(resultElements, mapFn(elemArgs))
			}
			
			return &object.Array{Elements: resultElements}
		},
		ReturnType: object.ARRAY_OBJ,
		ParamTypes: []object.ObjectType{object.ARRAY_OBJ, object.FUNCTION_OBJ},
	}
	
	// filter function
	Builtins["filter"] = &object.Builtin{
		Name: "filter",
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return createError("filter function requires 2 arguments")
			}
			
			arr, ok := args[0].(*object.Array)
			if !ok {
				return createError("First argument to filter must be an array")
			}
			
			// Check if the second argument is a function (either user-defined or builtin)
			var filterFn func(object.Object) object.Object
			
			switch fn := args[1].(type) {
			case *object.Function:
				// User-defined function
				filterFn = func(elem object.Object) object.Object {
					extendedEnv := object.NewEnclosedEnvironment(fn.Env)
					extendedEnv.Set("🍕", elem)
					
					result := Eval(fn.ASTBody, extendedEnv)
					
					if errObj, ok := result.(*object.Error); ok {
						return errObj
					}
					
					if retVal, ok := result.(*object.ReturnValue); ok {
						result = retVal.Value
					}
					
					return result
				}
			case *object.Builtin:
				// Builtin function
				filterFn = func(elem object.Object) object.Object {
					result := fn.Fn(elem)
					return result
				}
			default:
				return createError("Second argument to filter must be a function")
			}
			
			resultElements := make([]object.Object, 0, len(arr.Elements))
			
			for _, elem := range arr.Elements {
				result := filterFn(elem)
				
				// Check for errors
				if errObj, ok := result.(*object.Error); ok {
					return errObj
				}
				
				if isTruthy(result) {
					resultElements = append(resultElements, elem)
				}
			}
			
			return &object.Array{Elements: resultElements}
		},
		ReturnType: object.ARRAY_OBJ,
		ParamTypes: []object.ObjectType{object.ARRAY_OBJ, object.FUNCTION_OBJ},
	}
}
