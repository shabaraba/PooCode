package evaluator

import (
	"fmt"
	
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
		logger.Debug("  %s: %s", k, v.Type())
	}
	logger.Debug("")

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

	// 環境から同名のすべての関数を検索
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
		// case文対応: applyCaseBare を使用して呼び出す
		logCaseDebug("単独関数をcase文対応で実行: %s", functions[0].Inspect())
		return applyCaseBare(functions[0], args)
	}

	logger.Debug("複数の関数が見つかりました: %d", len(functions))

	// 条件付き関数と条件なし関数を正確にグループ化
	var conditionalFuncs []*object.Function
	var defaultFuncs []*object.Function

	// デバッグ情報
	logger.Debug("関数 '%s' を %d 個の候補から分類します", name, len(functions))

	for i, fn := range functions {
		// デバッグ情報: 関数の詳細
		logger.Debug("  関数候補 %d の詳細: Condition=%v, Addr=%p", i+1, fn.Condition, fn)
		
		// 厳密なnilチェックで条件式の有無を判定（重要）
		hasCondition := fn.Condition != nil
		logger.Debug("  条件式判定: %v (nilチェック結果: %v)", fn.Condition, hasCondition)
		
		if hasCondition {
			// 条件付き関数のみを条件付き関数として分類
			conditionalFuncs = append(conditionalFuncs, fn)
			logger.Debug("  関数候補 %d: 条件付き関数として分類（条件式: %v）", i+1, fn.Condition)
			// 追加デバッグ - 関数のすべての属性を表示
			params := ""
			for _, p := range fn.Parameters {
				params += p.Value + ", "
			}
			logger.Debug("    詳細: 入力型=%s, 戻り値型=%s, パラメータ=[%s]", 
				fn.InputType, fn.ReturnType, params)
		} else {
			// 条件式がないものはデフォルト関数として分類
			defaultFuncs = append(defaultFuncs, fn)
			logger.Debug("  関数候補 %d: デフォルト関数として分類（条件式なし）- アドレス: %p", i+1, fn)
			// 追加デバッグ - 関数のすべての属性を表示
			params := ""
			for _, p := range fn.Parameters {
				params += p.Value + ", "
			}
			logger.Debug("    詳細: 入力型=%s, 戻り値型=%s, パラメータ=[%s]", 
				fn.InputType, fn.ReturnType, params)
		}
	}
	
	logger.Debug("分類結果: 条件付き関数=%d個, デフォルト関数=%d個", 
		len(conditionalFuncs), len(defaultFuncs))

	// まず条件付き関数を検索して評価
	logger.Debug("条件付き関数を %d 個見つけました\n", len(conditionalFuncs))
	
	// 条件が真となった関数を格納する変数
	var matchedCondFunc *object.Function
	
	for i, fn := range conditionalFuncs {
		logger.Debug("条件付き関数候補 %d を評価中...\n", i+1)

		// 条件式評価の共通関数を使用
		isTrue, condResult := evalConditionalExpression(fn, args, env)
		
		// エラーが発生した場合、そのエラーを返す
		if condResult != nil && condResult.Type() == object.ERROR_OBJ {
			return condResult
		}

		// 条件が真なら、この関数を使用
		if isTrue {
			logger.Debug("条件が真であるため、この関数を使用します")
			matchedCondFunc = fn
			break // 条件が真となった最初の関数を採用して処理を終了
		} else {
			logger.Debug("条件が偽であるため、この関数をスキップします")
		}
	}
	
	// 条件に一致する関数が見つかった場合、その関数を実行
	if matchedCondFunc != nil {
		logger.Debug("条件に一致する関数を実行します")
		// case文対応: applyCaseBare を使用して呼び出す
		logCaseDebug("条件付き関数をcase文対応で実行: %s", matchedCondFunc.Inspect())
		return applyCaseBare(matchedCondFunc, args)
	}

	// 条件付き関数が該当しなかった場合、デフォルト関数を使用
	logger.Debug("デフォルト関数を %d 個見つけました", len(defaultFuncs))
	
	// ステップ1: 明示的に宣言されたデフォルト関数を探す
	if len(defaultFuncs) == 0 {
		// ステップ2: 専用の名前（funcName#default）でデフォルト関数を探してみる
		defaultFuncName := fmt.Sprintf("%s#default", name)
		logger.Debug("デフォルト関数が見つからないので、'%s' を探します...", defaultFuncName)
		if obj, ok := env.Get(defaultFuncName); ok {
			if function, ok := obj.(*object.Function); ok {
				logger.Debug("専用名でデフォルト関数 '%s' が見つかりました", defaultFuncName)
				defaultFuncs = append(defaultFuncs, function)
			}
		}
		
		// ステップ3: それでも見つからない場合は、一般的な関数を検索
		if len(defaultFuncs) == 0 {
			logger.Debug("一般的な '%s' 関数を検索します...", name)
			// 環境から再度すべての関数を取得（完全リフレッシュ）
			freshFunctions := env.GetAllFunctionsByName(name)
			logger.Debug("見つかった関数: %d 個", len(freshFunctions))
			
			// 条件のない関数を優先して検索
			for _, fn := range freshFunctions {
				// デバッグ出力
				if fn.Condition != nil {
					logger.Debug("  関数: 条件あり - %p", fn)
				} else {
					logger.Debug("  関数: 条件なし - %p", fn)
				}
				
				// 条件のない関数のみを抽出
				if fn.Condition == nil {
					logger.Debug("条件なし関数 '%s' が見つかりました (アドレス: %p)", name, fn)
					defaultFuncs = append(defaultFuncs, fn)
					// 最初の条件なし関数を使用
					break
				}
			}
		}
	}
	
	// 見つかったデフォルト関数を実行
	if len(defaultFuncs) > 0 {
		logger.Debug("デフォルト関数を使用します: %s", name)
		// case文対応: applyCaseBare を使用して呼び出す
		logCaseDebug("デフォルト関数をcase文対応で実行: %s", defaultFuncs[0].Inspect())
		return applyCaseBare(defaultFuncs[0], args)
	} else {
		// どのような関数も見つからなかった場合、エラーを返す
		logger.Debug("適切なデフォルト関数が見つかりませんでした")
		return createEvalError("条件に一致する関数 '%s' が見つかりません", name)
	}

	// この行は実行されません（上記のif-elseで必ずreturnするため）
}

// applyFunctionWithPizza は関数に🍕をセットして実行する
// 注: この関数は後方互換性のために維持しています
// 新しい実装ではapplyCaseBareを使用してください
func applyFunctionWithPizza(fn *object.Function, args []object.Object) object.Object {
	// case文対応: 新しい関数に委譲して実装を一元化
	logCaseDebug("applyFunctionWithPizza は applyCaseBare に委譲します")
	return applyCaseBare(fn, args)
}
