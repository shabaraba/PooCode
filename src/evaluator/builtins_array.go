package evaluator

import (
	"fmt"
	"strings"
	"github.com/uncode/logger"
	"github.com/uncode/object"
)

// registerArrayBuiltins は配列関連の組み込み関数を登録する
func registerArrayBuiltins() {
	// map関数 - 配列の各要素に関数を適用する
	Builtins["map"] = &object.Builtin{
		Name: "map",
		Fn: func(args ...object.Object) object.Object {
			// 引数の数チェック
			if len(args) != 2 {
				return createError("map関数は2つの引数が必要です: 配列, 関数")
			}
			
			// 第1引数が配列かチェック
			arr, ok := args[0].(*object.Array)
			if !ok {
				return createError("map関数の第1引数は配列である必要があります: %s", args[0].Type())
			}
			
			// 第2引数が関数かチェック
			fn, ok := args[1].(*object.Function)
			if !ok {
				return createError("map関数の第2引数は関数である必要があります: %s", args[1].Type())
			}
			
			// map関数の引数のパラメータは空である必要がある
			if len(fn.Parameters) > 0 {
				return createError("map関数に渡された関数はパラメーターを取るべきではありません")
			}
			
			// 結果の配列
			resultElements := make([]object.Object, 0, len(arr.Elements))
			
			// 配列の各要素に関数を適用
			for _, elem := range arr.Elements {
				// 関数の環境を拡張して🍕に現在の要素を設定
				extendedEnv := object.NewEnclosedEnvironment(fn.Env)
				extendedEnv.Set("🍕", elem)
				
				// 関数を評価
				result := Eval(fn.ASTBody, extendedEnv)
				
				// エラーチェック
				if errObj, ok := result.(*object.Error); ok {
					return errObj
				}
				
				// ReturnValueをアンラップ
				if retVal, ok := result.(*object.ReturnValue); ok {
					result = retVal.Value
				}
				
				// 結果を配列に追加
				resultElements = append(resultElements, result)
			}
			
			return &object.Array{Elements: resultElements}
		},
		ReturnType: object.ARRAY_OBJ,
		ParamTypes: []object.ObjectType{object.ARRAY_OBJ, object.FUNCTION_OBJ},
	}
	
	// filter関数 - 条件に合致する要素のみを抽出する
	Builtins["filter"] = &object.Builtin{
		Name: "filter",
		Fn: func(args ...object.Object) object.Object {
			// 引数の数チェック
			if len(args) != 2 {
				return createError("filter関数は2つの引数が必要です: 配列, 関数")
			}
			
			// 第1引数が配列かチェック
			arr, ok := args[0].(*object.Array)
			if !ok {
				return createError("filter関数の第1引数は配列である必要があります: %s", args[0].Type())
			}
			
			// 第2引数が関数かチェック
			fn, ok := args[1].(*object.Function)
			if !ok {
				return createError("filter関数の第2引数は関数である必要があります: %s", args[1].Type())
			}
			
			// filter関数の引数のパラメータは空である必要がある
			if len(fn.Parameters) > 0 {
				return createError("filter関数に渡された関数はパラメーターを取るべきではありません")
			}
			
			// 結果の配列
			resultElements := make([]object.Object, 0)
			
			// 配列の各要素に条件関数を適用
			for _, elem := range arr.Elements {
				// 関数の環境を拡張して🍕に現在の要素を設定
				extendedEnv := object.NewEnclosedEnvironment(fn.Env)
				extendedEnv.Set("🍕", elem)
				
				// 条件関数を評価
				result := Eval(fn.ASTBody, extendedEnv)
				
				// エラーチェック
				if errObj, ok := result.(*object.Error); ok {
					return errObj
				}
				
				// ReturnValueをアンラップ
				if retVal, ok := result.(*object.ReturnValue); ok {
					result = retVal.Value
				}
				
				// 結果が真の場合、要素を結果配列に追加
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
				
				// Check for errors
				if errObj, ok := result.(*object.Error); ok {
					return errObj
				}
				
				// Unwrap return value
				if retVal, ok := result.(*object.ReturnValue); ok {
					result = retVal.Value
				}
				
				// Only add element if condition is true
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
