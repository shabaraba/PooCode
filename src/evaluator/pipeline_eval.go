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

// パイプライン処理で関数呼び出しを評価する（改善版）
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
	}

	// 通常のケース: 右側がシンプルな関数呼び出し（例: func(arg1, arg2)）
	// 引数を評価（一時環境で評価することで🍕の影響を分離）
	args := evalExpressions(callExpr.Arguments, env)
	for _, arg := range args {
		if arg.Type() == object.ERROR_OBJ {
			return arg
		}
	}

	// デバッグ出力
	logger.Debug("パイプラインの関数名: %s, 左辺値: %s, 引数: %v\n",
		funcName, left.Inspect(), args)

	// 通常の関数呼び出しの場合（例: 左辺 |> func arg1 arg2）
	// 全引数リストを作成（第一引数は左辺の値、第二引数以降は関数呼び出しの引数）
	logger.Debug("通常の関数呼び出し: 引数リストを作成します\n")
	allArgs := []object.Object{left}
	allArgs = append(allArgs, args...)
	args = allArgs

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

