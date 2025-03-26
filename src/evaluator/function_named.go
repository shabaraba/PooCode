package evaluator

import (
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
		return newError("関数 '%s' が見つかりません", name)
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

	// 関数が1つだけの場合は直接適用
	if len(functions) == 1 {
		logger.Debug("関数が1つだけ見つかりました")
		return applyFunctionWithPizza(functions[0], args)
	}

	logger.Debug("複数の関数が見つかりました:", len(functions))

	// 🍕 を設定（もし必要なら）
	if len(args) > 0 {
		logger.Debug("🍕 に値 %s を設定します\n", args[0].Inspect())
		logger.Debug("🍕の値のタイプ: %s\n", args[0].Type())
		env.Set("🍕", args[0])
	} else {
		logger.Debug("引数が見つからないため、🍕は設定しません")
	}

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
		logger.Debug("条件式: %v\n", fn.Condition)

		// 条件式を評価するための環境を作成
		condEnv := object.NewEnclosedEnvironment(env)
		if len(args) > 0 {
			condEnv.Set("🍕", args[0])
			logger.Debug("条件評価のために 🍕 に値 %s を設定しました\n", args[0].Inspect())
		}

		// 条件式を評価
		condResult := Eval(fn.Condition, condEnv)
		logger.Debug("条件式の評価結果: %s\n", condResult.Inspect())
		logger.Debug("条件式の評価結果のタイプ: %s\n", condResult.Type())

		// エラーが発生した場合、詳細を出力
		if condResult.Type() == object.ERROR_OBJ {
			logger.Debug("条件評価でエラーが発生しました: %s\n", condResult.Inspect())
			return condResult
		}

		// 条件が真なら、この関数を使用
		// Booleanオブジェクトの場合はそのValueを使用、それ以外の場合はisTruthyで評価
		isTrue := false
		if condResult.Type() == object.BOOLEAN_OBJ {
			isTrue = condResult.(*object.Boolean).Value
			logger.Debug("条件式の真偽値: %v\n", isTrue)
		} else {
			isTrue = isTruthy(condResult)
			logger.Debug("条件式の評価結果（非Boolean）が %v と判定されました\n", isTrue)
		}

		if isTrue {
			logger.Debug("条件が真であるため、この関数を使用します")
			return applyFunctionWithPizza(fn, args)
		} else {
			logger.Debug("条件が偽であるため、この関数をスキップします")
		}
	}

	// 条件付き関数が該当しなかった場合、デフォルト関数を使用
	logger.Debug("デフォルト関数を %d 個見つけました\n", len(defaultFuncs))
	if len(defaultFuncs) > 0 {
		logger.Debug("デフォルト関数を使用します")
		return applyFunctionWithPizza(defaultFuncs[0], args)
	}

	// 適用可能な関数が見つからない場合
	return newError("条件に一致する関数 '%s' が見つかりません", name)
}
