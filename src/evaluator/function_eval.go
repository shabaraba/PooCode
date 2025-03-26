package evaluator

import (
	"github.com/uncode/ast"
	"github.com/uncode/logger"
	"github.com/uncode/object"
)

// applyFunction は関数を適用する
func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		// 関数呼び出しの実装
		logger.Debug("関数を呼び出します:", fn.Inspect())

		// 修正: 引数は1つまでだけ許可（パイプライン以外）
		if len(fn.Parameters) > 1 {
			return newError("関数は最大1つのパラメータしか持てません（パイプライン以外）: %s", fn.Inspect())
		}

		// 引数の数をチェック
		if len(args) != len(fn.Parameters) {
			return newError("引数の数が一致しません: 期待=%d, 実際=%d", len(fn.Parameters), len(args))
		}

		// 新しい環境を作成
		extendedEnv := object.NewEnclosedEnvironment(fn.Env)

		// 引数を環境にバインド
		for i, param := range fn.Parameters {
			extendedEnv.Set(param.Value, args[i])
		}

		// 修正後の仕様では、🍕はパイプラインで渡された値のみを表す
		// 通常の関数呼び出しでは🍕は設定しない

		// 関数本体を評価（ASTBodyをast.BlockStatementに型アサーション）
		astBody, ok := fn.ASTBody.(*ast.BlockStatement)
		if !ok {
			return newError("関数の本体がBlockStatementではありません")
		}
		result := evalBlockStatement(astBody, extendedEnv)

		// 💩値を返す（関数の戻り値）
		if obj, ok := result.(*object.ReturnValue); ok {
			return obj.Value
		}
		return result

	case *object.Builtin:
		// 修正: ビルトイン関数も引数を1つまでに制限（ただし print や数学関数など一部の例外を除く）
		if len(args) > 1 && fn.Name != "print" && fn.Name != "range" && fn.Name != "sum" {
			logger.Debug("ビルトイン関数 %s は引数を1つしか取れません: 実際の引数数=%d\n", fn.Name, len(args))
			// 引数が多い場合は最初の引数だけを使用
			return fn.Fn(args[0])
		}
		return fn.Fn(args...)

	default:
		return newError("関数ではありません: %s", fn.Type())
	}
}

// applyFunctionWithPizza は関数を適用する（パイプラインの場合同様に🍕も設定）
func applyFunctionWithPizza(fn *object.Function, args []object.Object) object.Object {
	// 関数呼び出しの実装
	logger.Debug("パイプライン対応で関数を呼び出します:", fn.Inspect())

	// 引数とパラメータのデバッグ出力
	logger.Debug("関数パラメータ数: %d, 引数数: %d\n", len(fn.Parameters), len(args))
	for i, param := range fn.Parameters {
		logger.Debug("  パラメータ %d: %s\n", i, param.Value)
	}
	for i, arg := range args {
		logger.Debug("  引数 %d: %s\n", i, arg.Inspect())
	}

	// 修正: パイプライン関数の引数チェック
	// - 第1引数は🍕として常に渡される
	// - 追加の引数は最大1つまで
	if len(args) > 2 {
		logger.Debug("警告: パイプラインでは最大1つの追加引数しか使用できません（現在: %d）\n", len(args)-1)
		// 余分な引数は無視して最初の2つだけを使用（🍕 + 1つの引数）
		args = args[:2]
	}

	// 新しい環境を作成
	extendedEnv := object.NewEnclosedEnvironment(fn.Env)

	// パイプラインを利用する関数では:
	// - 第1引数は常に🍕として設定される
	// - パラメータがある場合、引数の残りをパラメータにマッピングする
	if len(args) > 0 {
		// 🍕 変数を設定
		extendedEnv.Set("🍕", args[0])
		logger.Debug("🍕 に値 %s を設定しました\n", args[0].Inspect())

		// パラメータを環境にバインド
		if len(fn.Parameters) > 0 {
			// パラメータ名を取得
			paramName := fn.Parameters[0].Value

			// numパラメータにどの値を設定するか
			if len(args) > 1 {
				// 複数引数の場合: 第2引数をnumに設定
				extendedEnv.Set(paramName, args[1])
				logger.Debug("パラメータ '%s' に値 %s を設定しました\n",
					paramName, args[1].Inspect())
			} else {
				// 単一引数の場合: 🍕と同じ値をnumに設定
				extendedEnv.Set(paramName, args[0])
				logger.Debug("単一引数: パラメータ '%s' に値 %s を設定しました\n",
					paramName, args[0].Inspect())
			}
		}
	}

	// 関数本体を評価（ASTBodyをast.BlockStatementに型アサーション）
	astBody, ok := fn.ASTBody.(*ast.BlockStatement)
	if !ok {
		return newError("関数の本体がBlockStatementではありません")
	}

	logger.Debug("関数本体を評価します...")
	result := evalBlockStatement(astBody, extendedEnv)

	// 💩値を返す（関数の戻り値）
	if obj, ok := result.(*object.ReturnValue); ok {
		logger.Debug("関数から戻り値が見つかりました: %s\n", obj.Value.Inspect())
		return obj.Value
	}

	logger.Debug("関数から戻り値なしで実行完了: %s\n", result.Inspect())
	return result
}

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
		// 余分な引数は無視して最初の引数だけを使用
		args = args[:1]
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

