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
	left := Eval(node.Left, env)
	if left.Type() == object.ERROR_OBJ {
		return left
	}

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
		return result
	} else {
		// 右辺が識別子の場合（関数名のみ）
		if ident, ok := node.Right.(*ast.Identifier); ok {
			logger.Debug("識別子としてのパイプライン先: %s\n", ident.Value)
			logger.Debug("パイプラインから applyNamedFunction を呼び出します (関数名: %s)\n", ident.Value)

			// 環境変数 🍕 を設定して名前付き関数呼び出しへ処理を委譲
			// ここで左辺の値を唯一の引数として渡す
			args := []object.Object{left}

			// 名前付き関数を適用する
			result = applyNamedFunction(tempEnv, ident.Value, args)
			logger.Debug("パイプライン: 関数 '%s' の実行結果: タイプ=%s, 値=%s\n",
				ident.Value, result.Type(), result.Inspect())
		} else {
			// その他の場合は処理できない
			return createEvalError("パイプラインの右側が関数または識別子ではありません: %T", node.Right)
		}
	}

	// 元の🍕変数を環境に戻す（必要に応じて）
	if hasPizza {
		logger.Debug("元の🍕変数を復元します: %s", originalPizza.Inspect())
		env.Set("🍕", originalPizza)
	}

	return result
}

// パイプライン処理で関数呼び出しを評価する（改善版）
func evalPipelineWithCallExpression(left object.Object, callExpr *ast.CallExpression, env *object.Environment) object.Object {
	// 関数名を取得
	var funcName string
	if ident, ok := callExpr.Function.(*ast.Identifier); ok {
		funcName = ident.Value
	} else {
		return createEvalError("パイプラインの右側の関数名を取得できません: %T", callExpr.Function)
	}

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

	// 全引数リストを作成（第一引数は左辺の値、第二引数以降は関数呼び出しの引数）
	allArgs := []object.Object{left}
	allArgs = append(allArgs, args...)

	// デバッグ: 最終的な引数リストを表示
	logger.Debug("applyNamedFunction に渡す最終引数リスト: %d 個\n", len(allArgs))
	for i, arg := range allArgs {
		logger.Debug("  引数 %d: タイプ=%s, 値=%s\n", i, arg.Type(), arg.Inspect())
	}

	// 組み込み関数を直接取得して呼び出す
	if builtin, ok := Builtins[funcName]; ok {
		logger.Debug("ビルトイン関数 '%s' を実行: 全引数 %d 個\n", funcName, len(allArgs))
		result := builtin.Fn(allArgs...)
		logger.Debug("ビルトイン関数 '%s' の結果: %s\n", funcName, result.Inspect())
		return result
	}

	// 名前付き関数（ユーザー定義関数）を適用する
	result := applyNamedFunction(env, funcName, allArgs)
	logger.Debug("関数 '%s' の適用結果: %s\n", funcName, result.Inspect())
	return result
}