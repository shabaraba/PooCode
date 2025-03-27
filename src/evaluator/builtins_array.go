package evaluator

import (
	"github.com/uncode/ast"
	"github.com/uncode/object"
	"github.com/uncode/logger"
)

// registerArrayBuiltins は配列関連のビルトイン関数を登録する
func registerArrayBuiltins() {
	logger.Debug("配列関連のビルトイン関数を登録します")
	// map function
	Builtins["map"] = &object.Builtin{
		Name: "map",
		Fn: func(args ...object.Object) object.Object {
			// 現在の環境を取得
			env := GetEvalEnv()
			
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
			
			// 関数に渡す追加の固定引数（第3引数以降）
			var funcFixedArgs []object.Object
			if len(args) > 2 {
				funcFixedArgs = args[2:]
				logger.Debug("map関数に追加の引数: %d個", len(funcFixedArgs))
				for i, arg := range funcFixedArgs {
					logger.Debug("  追加引数 %d: %s", i, arg.Inspect())
				}
			}
			
			// 直接関数型で渡された場合
			switch fn := args[1].(type) {
			case *object.Function:
				// ユーザー定義関数
				mapFn = func(elemArgs []object.Object) object.Object {
					extendedEnv := object.NewEnclosedEnvironment(fn.Env)
					
					// 必ず最初の引数（配列要素）を🍕に設定
					if len(elemArgs) > 0 {
						extendedEnv.Set("🍕", elemArgs[0])
					}
					
					// ParamValues フィールドが設定されている場合はそちらから引数を設定
					if len(fn.ParamValues) > 0 {
						// 引数の数をチェック
						paramCount := len(fn.Parameters)
						if len(fn.ParamValues) < paramCount {
							logger.Debug("関数の引数が少なすぎます: 期待=%d, 実際=%d", paramCount, len(fn.ParamValues))
						}
						
						// パラメータ値をバインド
						for i := 0; i < paramCount && i < len(fn.ParamValues); i++ {
							extendedEnv.Set(fn.Parameters[i].Value, fn.ParamValues[i])
						}
					} else if len(fn.Parameters) > 0 && len(funcFixedArgs) > 0 {
						// 後方互換性のために従来の方法もサポート
						// 引数の数を確認
						paramCount := len(fn.Parameters)
						if len(funcFixedArgs) < paramCount {
							logger.Debug("関数の引数が少なすぎます: 期待=%d, 実際=%d", paramCount, len(funcFixedArgs))
						}
						
						// 引数をバインド
						for i := 0; i < paramCount && i < len(funcFixedArgs); i++ {
							extendedEnv.Set(fn.Parameters[i].Value, funcFixedArgs[i])
						}
					}
					
					// 関数本体を評価
					astBody, ok := fn.ASTBody.(*ast.BlockStatement)
					if !ok {
						logger.Error("関数本体がBlockStatementではありません: %T", fn.ASTBody)
						return createError("関数本体がBlockStatementではありません: %T", fn.ASTBody)
					}
					
					result := evalBlockStatement(astBody, extendedEnv)
					
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
					// 配列要素と固定引数を組み合わせる
					allArgs := []object.Object{elemArgs[0]}
					allArgs = append(allArgs, funcFixedArgs...)
					
					return fn.Fn(allArgs...)
				}
			default:
				// 文字列として関数名を取得し、環境から関数を検索
				funcName := args[1].Inspect()
				if funcObj, exists := env.Get(funcName); exists {
					switch fn := funcObj.(type) {
					case *object.Function:
						// ユーザー定義関数
						mapFn = func(elemArgs []object.Object) object.Object {
							extendedEnv := object.NewEnclosedEnvironment(fn.Env)
							
							// 必ず最初の引数（配列要素）を🍕に設定
							if len(elemArgs) > 0 {
								extendedEnv.Set("🍕", elemArgs[0])
							}
							
							// 関数が引数パラメータを持つ場合、funcFixedArgsから引数を設定
							if len(fn.Parameters) > 0 && len(funcFixedArgs) > 0 {
								// 引数の数を確認
								paramCount := len(fn.Parameters)
								if len(funcFixedArgs) < paramCount {
									logger.Debug("関数の引数が少なすぎます: 期待=%d, 実際=%d", paramCount, len(funcFixedArgs))
								}
								
								// 引数をバインド
								for i := 0; i < paramCount && i < len(funcFixedArgs); i++ {
									extendedEnv.Set(fn.Parameters[i].Value, funcFixedArgs[i])
								}
							}
							
							// 関数本体を評価
							astBody, ok := fn.ASTBody.(*ast.BlockStatement)
							if !ok {
								logger.Error("関数本体がBlockStatementではありません: %T", fn.ASTBody)
								return createError("関数本体がBlockStatementではありません: %T", fn.ASTBody)
							}
							
							result := evalBlockStatement(astBody, extendedEnv)
							
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
							// 配列要素と固定引数を組み合わせる
							allArgs := []object.Object{elemArgs[0]}
							allArgs = append(allArgs, funcFixedArgs...)
							
							return fn.Fn(allArgs...)
						}
					default:
						return createError("関数 '%s' は有効な関数ではありません: %T", funcName, funcObj)
					}
				} else {
					return createError("Second argument to map must be a function")
				}
			}
			
			// マップ処理の実行
			resultElements := make([]object.Object, 0, len(arr.Elements))
			
			for _, elem := range arr.Elements {
				// 各要素を単一引数として関数に渡す
				elemArgs := []object.Object{elem}
				
				// 各要素に関数を適用
				result := mapFn(elemArgs)
				
				// デバッグ情報を出力
				logger.Debug("map: 要素 %s に関数を適用した結果: %s", 
					elem.Inspect(), result.Inspect())
				
				resultElements = append(resultElements, result)
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
			// 現在の環境を取得
			env := GetEvalEnv()
			
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
					
					// 関数本体をBlockStatementに変換
					astBody, ok := fn.ASTBody.(*ast.BlockStatement)
					if !ok {
						logger.Error("関数本体がBlockStatementではありません: %T", fn.ASTBody)
						return createError("関数本体がBlockStatementではありません: %T", fn.ASTBody)
					}
					
					result := evalBlockStatement(astBody, extendedEnv)
					
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
				// 文字列として関数名を取得し、環境から関数を検索
				funcName := args[1].Inspect()
				if funcObj, exists := env.Get(funcName); exists {
					switch fn := funcObj.(type) {
					case *object.Function:
						// User-defined function
						filterFn = func(elem object.Object) object.Object {
							extendedEnv := object.NewEnclosedEnvironment(fn.Env)
							extendedEnv.Set("🍕", elem)
							
							// 関数本体をBlockStatementに変換
							astBody, ok := fn.ASTBody.(*ast.BlockStatement)
							if !ok {
								logger.Error("関数本体がBlockStatementではありません: %T", fn.ASTBody)
								return createError("関数本体がBlockStatementではありません: %T", fn.ASTBody)
							}
							
							result := evalBlockStatement(astBody, extendedEnv)
							
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
						return createError("関数 '%s' は有効な関数ではありません: %T", funcName, funcObj)
					}
				} else {
					return createError("Second argument to filter must be a function")
				}
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
