package evaluator

import (
	"github.com/uncode/ast"
	"github.com/uncode/config"
	"github.com/uncode/logger"
	"github.com/uncode/object"
)

// applyNamedFunction は名前付き関数を検索し、適用する
// 同じ名前で複数の関数が存在する場合は、条件に基づいて適切な関数を選択する
func applyNamedFunction(env *object.Environment, name string, args []object.Object) object.Object {
	logger.Debug("***** applyNamedFunction が呼び出されました *****")
	logger.Debug("関数名: %s、引数の数: %d\n", name, len(args))

	// デバッグ: 環境内のすべての変数を表示
	logger.Debug("現在の環境に登録されている変数:")
	for k, v := range env.GetVariables() {
		logger.Debug("  %s: %s\n", k, v.Type())
	}

	// 修正: 引数の数を制限（パイプライン以外）
	// パイプラインではない通常の呼び出しの場合、引数は1つだけ
	if len(args) > 1 {
		logger.Debug("警告: 関数 '%s' は通常の呼び出しでは1つの引数しか取れません（現在: %d）\n",
			name, len(args))
	}

	// ビルトイン関数を確認
	if builtin, ok := Builtins[name]; ok {
		logger.Debug("ビルトイン関数 '%s' を呼び出します\n", name)
		return builtin.Fn(args...)
	}

	// 環境から同名のすべての関数を取得
	functions := env.GetAllFunctionsByName(name)

	if len(functions) == 0 {
		return createEvalError("関数 '%s' が見つかりません", name)
	}

	// デバッグ情報
	logger.Debug("関数 '%s' を呼び出します: %d 個の候補が見つかりました\n", name, len(functions))
	for i, fn := range functions {
		if fn.Condition != nil {
			logger.Debug("  関数候補 %d: 条件=あり\n", i+1)
		} else {
			logger.Debug("  関数候補 %d: 条件=なし\n", i+1)
		}
	}

	// 修正：関数適用のための独立した環境を作成
	// これにより元の環境の🍕変数が上書きされるのを防ぐ
	funcEnv := object.NewEnclosedEnvironment(env)

	// 🍕 を設定（もし引数があれば）
	if len(args) > 0 {
		logger.Debug("関数適用の環境で🍕に値 %s を設定します\n", args[0].Inspect())
		logger.Debug("🍕の値のタイプ: %s\n", args[0].Type())
		funcEnv.Set("🍕", args[0])
	} else {
		logger.Debug("引数が見つからないため、🍕は設定しません")
	}

	// 関数が1つだけの場合は直接適用
	if len(functions) == 1 {
		logger.Debug("関数が1つだけ見つかりました")
		return applyFunctionWithPizza(functions[0], args)
	}

	logger.Debug("複数の関数が見つかりました: %d", len(functions))

	// 条件付き関数と条件なし関数をグループ化
	var conditionalFuncs []*object.Function
	var defaultFuncs []*object.Function

	for _, fn := range functions {
		if fn.Condition != nil {
			conditionalFuncs = append(conditionalFuncs, fn)
		} else {
			defaultFuncs = append(defaultFuncs, fn)
		}
	}

	// まず条件付き関数を検索して評価
	logger.Debug("条件付き関数を %d 個見つけました\n", len(conditionalFuncs))
	for i, fn := range conditionalFuncs {
		logger.Debug("条件付き関数候補 %d を評価中...\n", i+1)

		// 条件式の詳細を表示（ShowConditionDebugがtrueの場合のみ）
		if config.GlobalConfig.ShowConditionDebug {
			logger.Debug("-------- 条件式の詳細評価 --------")
			logger.Debug("条件式: %v", fn.Condition)
			// AST構造をより詳細に表示
			if infixExpr, ok := fn.Condition.(*ast.InfixExpression); ok {
				logger.Debug("条件式タイプ: 中置式")
				logger.Debug("  演算子: %s", infixExpr.Operator)
				logger.Debug("  左辺: %T - %v", infixExpr.Left, infixExpr.Left)
				logger.Debug("  右辺: %T - %v", infixExpr.Right, infixExpr.Right)
			} else {
				logger.Debug("条件式タイプ: %T", fn.Condition)
			}
			logger.Debug("----------------------------------")
		}

		// 条件式を評価するための環境を作成
		condEnv := object.NewEnclosedEnvironment(funcEnv)

		// 条件式を評価
		condResult := Eval(fn.Condition, condEnv)

		if config.GlobalConfig.ShowConditionDebug {
			logger.Debug("条件式の評価結果: %s", condResult.Inspect())
			logger.Debug("条件式の評価結果のタイプ: %s", condResult.Type())
		}

		// エラーが発生した場合、詳細を出力
		if condResult.Type() == object.ERROR_OBJ {
			logger.Debug("条件評価でエラーが発生しました: %s", condResult.Inspect())
			return condResult
		}

		// 条件が真なら、この関数を使用
		// Booleanオブジェクトの場合はそのValueを使用、それ以外の場合はisTruthyで評価
		isTrue := false
		if condResult.Type() == object.BOOLEAN_OBJ {
			isTrue = condResult.(*object.Boolean).Value
			logger.Debug("条件式の真偽値: %v", isTrue)
		} else {
			isTrue = isTruthy(condResult)
			logger.Debug("条件式の評価結果（非Boolean）が %v と判定されました", isTrue)
		}

		if isTrue {
			logger.Debug("条件が真であるため、この関数を使用します")
			return applyFunctionWithPizza(fn, args)
		} else {
			logger.Debug("条件が偽であるため、この関数をスキップします")
		}
	}

	// 条件付き関数が該当しなかった場合、デフォルト関数を使用
	logger.Debug("デフォルト関数を %d 個見つけました", len(defaultFuncs))
	if len(defaultFuncs) > 0 {
		logger.Debug("デフォルト関数を使用します")
		return applyFunctionWithPizza(defaultFuncs[0], args)
	}

	// 適用可能な関数が見つからない場合
	logger.Debug("条件に一致する関数が見つかりません")
	return createEvalError("条件に一致する関数 '%s' が見つかりません", name)
}

// applyFunctionWithPizza は関数に🍕をセットして実行する
func applyFunctionWithPizza(fn *object.Function, args []object.Object) object.Object {
	// 関数呼び出し用の新しい環境を作成
	extendedEnv := object.NewEnclosedEnvironment(fn.Env)
	funcName, _ := fn.Name()
	logger.Info("function: %s", funcName)
	for _, arg := range args {
		logger.Info("%#v", arg)
	}

	// 引数をバインド
	if len(args) > 0 {
		// 第1引数は特別に🍕にもセット
		extendedEnv.Set("🍕", args[0])

		// 通常の引数セット
		for i, param := range fn.Parameters {
			if i < len(args) {
				extendedEnv.Set(param.Value, args[i+1])
			} else {
				return createEvalError("引数の数が足りません: 期待=%d, 実際=%d", len(fn.Parameters), len(args))
			}
		}
	} else if len(fn.Parameters) > 0 {
		return createEvalError("引数が足りません: 期待=%d, 実際=0", len(fn.Parameters))
	}

	// 関数本体を評価
	astBody, ok := fn.ASTBody.(*ast.BlockStatement)
	if !ok {
		return createEvalError("関数の本体がBlockStatementではありません")
	}

	evaluated := evalBlockStatement(astBody, extendedEnv)

	// ReturnValue の場合は Value を抽出
	if returnValue, ok := evaluated.(*object.ReturnValue); ok {
		return returnValue.Value
	}

	return evaluated
}
