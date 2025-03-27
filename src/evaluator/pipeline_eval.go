package evaluator

import (
	"strconv"

	"github.com/uncode/ast"
	"github.com/uncode/logger"
	"github.com/uncode/object"
)

// maybeConvertToInteger は文字列を整数に変換する試みを行う
// 特に、パイプラインからのprint結果などを数値に変換するのに役立つ
func maybeConvertToInteger(obj object.Object) object.Object {
	if obj.Type() != object.STRING_OBJ {
		return obj // 文字列以外はそのまま返す
	}

	strValue := obj.(*object.String).Value

	// 文字列が数値として解釈可能かを試みる
	if intValue, err := strconv.ParseInt(strValue, 10, 64); err == nil {
		return &object.Integer{Value: intValue}
	}

	// 特定の文字列だけを変換する
	if strValue == "0" {
		return &object.Integer{Value: 0}
	} else if strValue == "1" {
		return &object.Integer{Value: 1}
	}

	// 変換できなければそのまま返す
	return obj
}

// evalPipeline は|>演算子のパイプライン処理を評価する
func evalPipeline(node *ast.InfixExpression, env *object.Environment) object.Object {
	logger.Debug("パイプライン演算子を検出しました")
	
	// 現在の🍕変数の値を保存（もし存在すれば）
	originalPizza, hasPizza := env.Get("🍕")
	if hasPizza {
		logger.Debug("元の🍕変数の値を保存: %s", originalPizza.Inspect())
	}
	
	// |>演算子の場合、左辺の結果を右辺の関数に渡す
	// 左辺には別のパイプライン式が含まれている可能性があります
	left := Eval(node.Left, env)
	if left.Type() == object.ERROR_OBJ {
		return left
	}

	logger.Debug("パイプラインの左辺評価結果: タイプ=%s, 値=%s", left.Type(), left.Inspect())

	// パイプライン処理のための一時環境を作成
	tempEnv := object.NewEnclosedEnvironment(env)
	
	// 明示的に🍕変数に左辺の値を設定（条件式の評価で必要）
	logger.Debug("パイプラインで🍕に値を明示的に設定します: %s\n", left.Inspect())
	// nullを無視（printの結果などがnullの場合に問題が発生）
	if left.Type() != object.NULL_OBJ {
		tempEnv.Set("🍕", left)
	} else {
		logger.Debug("左辺値がnullのため、🍕の設定をスキップします")
	}

	var result object.Object

	// 右辺の式がCallExpressionの場合（関数呼び出し）
	if callExpr, ok := node.Right.(*ast.CallExpression); ok {
		logger.Debug("パイプラインの右辺がCallExpressionです")
		result = evalPipelineWithCallExpression(left, callExpr, tempEnv)
	} else {
		// 右辺が識別子の場合（関数名のみ）
		if ident, ok := node.Right.(*ast.Identifier); ok {
			logger.Debug("識別子としてのパイプライン先: %s\n", ident.Value)
			logger.Debug("パイプラインから 関数を呼び出します (関数名: %s)\n", ident.Value)

			// 環境変数 🍕 を設定して関数呼び出しへ処理を委譲
			// ここで左辺の値を唯一の引数として渡す
			args := []object.Object{left}

			// 組み込み関数を直接取得して呼び出す (特にmapやfilterの場合)
			if builtin, ok := Builtins[ident.Value]; ok {
				logger.Debug("ビルトイン関数 '%s' を実行します\n", ident.Value)
				result = builtin.Fn(args...)
				logger.Debug("ビルトイン関数 '%s' の実行結果: タイプ=%s, 値=%s\n",
					ident.Value, result.Type(), result.Inspect())
			} else {
				// 名前付き関数を適用する
				result = applyNamedFunction(tempEnv, ident.Value, args)
				logger.Debug("パイプライン: 関数 '%s' の実行結果: タイプ=%s, 値=%s\n",
					ident.Value, result.Type(), result.Inspect())
			}
		} else {
			// その他の場合は処理できない
			return createError("パイプラインの右側が関数または識別子ではありません: %T", node.Right)
		}
	}

	// 元の🍕変数を環境に戻す（必要に応じて）
	if hasPizza {
		logger.Debug("元の🍕変数を復元します: %s", originalPizza.Inspect())
		env.Set("🍕", originalPizza)
	}

	logger.Debug("パイプラインの最終結果: タイプ=%s, 値=%s", result.Type(), result.Inspect())
	return result
}

// evalPipelineWithCallExpression は関数呼び出しを含むパイプライン処理を評価する
func evalPipelineWithCallExpression(left object.Object, callExpr *ast.CallExpression, env *object.Environment) object.Object {
	// 関数名を取得
	var funcName string
	
	// 特殊ケース: 右側が関数呼び出し式（例: add_num(100)）のケース
	if ident, ok := callExpr.Function.(*ast.Identifier); ok {
		// 関数名を取得
		funcName = ident.Value
		logger.Debug("関数呼び出し式の関数名: %s\n", funcName)
	} else {
		logger.Debug("関数呼び出し式が識別子ではありません: %T\n", callExpr.Function)
		
		// map add_num(100) のような特殊なケースを処理
		if nestedCallExpr, ok := callExpr.Function.(*ast.CallExpression); ok {
			logger.Debug("入れ子の関数呼び出しを検出しました: %T\n", nestedCallExpr)

			// まず、内側の関数名を取得
			if innerIdent, ok := nestedCallExpr.Function.(*ast.Identifier); ok {
				// 内側の関数名（例: add_num）
				innerFuncName := innerIdent.Value
				logger.Debug("内側の関数名: %s\n", innerFuncName)
				
				// 内側の関数に対応する関数オブジェクトを環境から取得
				funcObj, exists := env.Get(innerFuncName)
				if !exists {
					return createError("関数 '%s' が見つかりません", innerFuncName)
				}
				
				// 引数を評価
				args := evalExpressions(nestedCallExpr.Arguments, env)
				if len(args) > 0 && args[0].Type() == object.ERROR_OBJ {
					return args[0]
				}
				
				logger.Debug("内側の関数の引数: %d 個\n", len(args))
				for i, arg := range args {
					logger.Debug("  引数 %d: %s\n", i, arg.Inspect())
				}
				
				// 外側の関数名（例: map）を取得
				// 現在の文脈では通常「map」
				outerFuncName := "map"
				logger.Debug("外側の関数名: %s\n", outerFuncName)
				
				// map ビルトイン関数を取得
				builtin, ok := Builtins[outerFuncName]
				if !ok {
					return createError("ビルトイン関数 '%s' が見つかりません", outerFuncName)
				}
				
				// 内側の関数と引数をまとめて配列に渡す
				switch fn := funcObj.(type) {
				case *object.Function:
					// ユーザー定義関数の場合は引数を設定した新しい関数を作成
					logger.Debug("ユーザー定義関数に引数をセット: %s\n", innerFuncName)
					
					newFunc := &object.Function{
						Parameters:  fn.Parameters,
						ParamValues: args,  // 重要: 引数を保存
						ASTBody:     fn.ASTBody,
						Env:         fn.Env,
						InputType:   fn.InputType,
						ReturnType:  fn.ReturnType,
					}
					
					// 配列と関数を引数リストにして map 関数を呼び出す
					mapArgs := []object.Object{left, newFunc}
					return builtin.Fn(mapArgs...)
					
				case *object.Builtin:
					// ビルトイン関数の場合
					logger.Debug("ビルトイン関数として処理: %s\n", innerFuncName)
					
					// ビルトイン関数と引数を一緒に渡す
					mapArgs := []object.Object{left, fn}
					mapArgs = append(mapArgs, args...)
					return builtin.Fn(mapArgs...)
					
				default:
					return createError("'%s' は有効な関数ではありません: %T", innerFuncName, funcObj)
				}
			} else {
				return createError("入れ子の関数呼び出しの関数名が識別子ではありません: %T", nestedCallExpr.Function)
			}
		} else {
			return createError("パイプラインの右側の関数名を取得できません: %T", callExpr.Function)
		}
	}

	// 通常のケース: 右側がシンプルな関数呼び出し（例: func(arg1, arg2)）
	// 引数を評価（一時環境で評価することで🍕の影響を分離）
	args := evalExpressions(callExpr.Arguments, env)
	if len(args) > 0 && args[0].Type() == object.ERROR_OBJ {
		return args[0]
	}

	// デバッグ出力
	logger.Debug("パイプラインの関数名: %s, 左辺値: %s, 引数: %v\n",
		funcName, left.Inspect(), args)
	
	// 特殊ケース: map(add_num(100))のようなケースを処理
	if funcName == "map" && len(args) == 1 {
		if fn, ok := args[0].(*object.Function); ok {
			if len(fn.Parameters) > 0 && len(callExpr.Arguments) > 1 {
				// map(add_num(100))のようなケース
				logger.Debug("特殊なmap呼び出し検出: map(func(arg))\n")
				
				// 第1引数は配列、第2引数は関数（すでに引数付きで評価済み）
				specialArgs := []object.Object{left, args[0]}
				
				// map関数を呼び出し
				if builtin, ok := Builtins[funcName]; ok {
					return builtin.Fn(specialArgs...)
				}
			}
		}
	}

	// mapやfilterなどのビルトイン関数の特別処理
	if funcName == "map" || funcName == "filter" {
		// mapやfilterのケースでは、第一引数は配列、第二引数は関数
		logger.Debug("map/filter関数のための特別処理を行います\n")
		
		// 左辺の値が配列かどうかを確認
		_, isArray := left.(*object.Array)
		if !isArray {
			logger.Warn("map/filter関数には配列が必要ですが、受け取ったのは %s です\n", left.Type())
		}
		
		if len(args) == 0 {
			// 第一引数は左辺の値（配列）
			logger.Debug("map/filter: 引数がないため、左辺の値のみを使用します\n")
			args = []object.Object{left}
		} else {
			// 関数名を取得できた場合（ユーザー定義関数名など）
			if args[0].Type() == object.STRING_OBJ {
				logger.Debug("map/filter: 第1引数が文字列 '%s' です - 関数名として扱います\n", args[0].Inspect())
				
				// 環境から関数を探す
				funcNameStr := args[0].(*object.String).Value
				if fn, ok := env.Get(funcNameStr); ok {
					logger.Debug("環境から関数 '%s' を見つけました\n", funcNameStr)
					
					// 関数を第2引数として設定し直す
					args[0] = fn
				}
			}
			
			// 第二引数以降は変更なし（例: map add_num 100 -> [left, add_num, 100]）
			logger.Debug("map/filter: 引数リストを作成します: 左辺の値 + %d 個の引数\n", len(args))
			allArgs := []object.Object{left}
			allArgs = append(allArgs, args...)
			args = allArgs
		}
	} else {
		// 通常の関数呼び出しの場合（例: 左辺 |> func arg1 arg2）
		// 全引数リストを作成（第一引数は左辺の値、第二引数以降は関数呼び出しの引数）
		logger.Debug("通常の関数呼び出し: 引数リストを作成します\n")
		allArgs := []object.Object{left}
		allArgs = append(allArgs, args...)
		args = allArgs
	}

	// デバッグ: 最終的な引数リストを表示
	logger.Debug("関数呼び出しに渡す最終引数リスト: %d 個\n", len(args))
	for i, arg := range args {
		logger.Debug("  引数 %d: タイプ=%s, 値=%s\n", i, arg.Type(), arg.Inspect())
	}

	var result object.Object

	// 組み込み関数を直接取得して呼び出す
	if builtin, ok := Builtins[funcName]; ok {
		logger.Debug("ビルトイン関数 '%s' を実行: 全引数 %d 個\n", funcName, len(args))
		result = builtin.Fn(args...)
		logger.Debug("ビルトイン関数 '%s' の結果: タイプ=%s, 値=%s\n", 
			funcName, result.Type(), result.Inspect())
	} else {
		// 名前付き関数（ユーザー定義関数）を適用する
		result = applyNamedFunction(env, funcName, args)
		logger.Debug("関数 '%s' の適用結果: タイプ=%s, 値=%s\n", 
			funcName, result.Type(), result.Inspect())
	}

	return result
}

// この部分は pipeline_call_eval.go と assignment_eval.go に移動しました

// pipeDebugLevel はパイプラインのデバッグレベルを保持します
var pipeDebugLevel = logger.LevelDebug

// SetPipeDebugLevel はパイプラインのデバッグレベルを設定します
func SetPipeDebugLevel(level logger.LogLevel) {
	pipeDebugLevel = level
}
