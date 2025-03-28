package evaluator

import (
	"fmt"
	"strings"
	
	"github.com/uncode/ast"
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
		result := applyFunctionWithPizza(functions[0], args)
		if result != nil {
			return result
		}
		// nilが返された場合は、引数が合わなかった
		logger.Debug("単独関数の引数が合いませんでした")
		return createEvalError("関数 '%s' の引数が合いません", name)
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
		result := applyFunctionWithPizza(matchedCondFunc, args)
		if result != nil {
			return result
		}
		// nilが返された場合は、パラメータとして引数が合わなかった
		logger.Debug("条件付き関数の引数が合いませんでした")
	}

	// 条件付き関数が該当しなかった場合、デフォルト関数を使用
	logger.Debug("デフォルト関数を %d 個見つけました", len(defaultFuncs))
	
	// ステップ1: 明示的に宣言されたデフォルト関数を探す
	if len(defaultFuncs) == 0 {
		// ステップ2: 専用の名前（funcName#default）でデフォルト関数を探してみる
		defaultFuncName := fmt.Sprintf("%s#default", name)
		logger.Debug("デフォルト関数が見つからないので、特殊名 '%s' を探します...", defaultFuncName)
		
		// 環境内のすべての変数をデバッグログに表示
		logger.Debug("=== 環境内の現在の変数（デフォルト関数探索）===")
		for k, v := range env.GetVariables() {
			if funcObj, ok := v.(*object.Function); ok {
				hasCondition := "なし"
				if funcObj.Condition != nil {
					hasCondition = "あり"
				}
				logger.Debug("  変数 '%s': 関数オブジェクト (条件=%s, アドレス=%p)", k, hasCondition, funcObj)
			} else {
				logger.Debug("  変数 '%s': %s", k, v.Type())
			}
		}
		logger.Debug("=================================================")
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
			
			// まず、条件なし関数を探す
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
			
			// それでも見つからなければ、特殊名でもう一度探す
			if len(defaultFuncs) == 0 {
				// 特殊名を使った再検索
				// 非公開のstoreフィールドではなくGetVariables()を使用
				for key, obj := range env.GetVariables() {
					// test#0 などの形式のキーも検索
					if len(key) > len(name) && strings.HasPrefix(key, name+"#") {
						if fn, ok := obj.(*object.Function); ok && fn.Condition == nil {
							logger.Debug("特殊名 '%s' でデフォルト関数が見つかりました", key)
							defaultFuncs = append(defaultFuncs, fn)
							break
						}
					}
				}
			}
			
			// それでも見つからなければ、もう一度試しに最初の関数を使用（最終手段）
			if len(defaultFuncs) == 0 && len(freshFunctions) > 0 {
				logger.Debug("条件なし関数が見つからないため、最初の関数を最後の手段として使用します")
				defaultFuncs = append(defaultFuncs, freshFunctions[0])
			}
		}
	}
	
	// 見つかったデフォルト関数を実行
	if len(defaultFuncs) > 0 {
		logger.Debug("デフォルト関数を使用します: %s", name)
		result := applyFunctionWithPizza(defaultFuncs[0], args)
		if result != nil {
			return result
		}
		// nilが返された場合は、パラメータとして引数が合わなかった
		logger.Debug("デフォルト関数の引数が合いませんでした")
		return createEvalError("関数 '%s' の引数が合いません", name)
	} else {
		// どのような関数も見つからなかった場合、エラーを返す
		logger.Debug("適切なデフォルト関数が見つかりませんでした")
		return createEvalError("条件に一致する関数 '%s' が見つかりません", name)
	}

	// この行は実行されません（上記のif-elseで必ずreturnするため）
}

// applyFunctionWithPizza は関数に🍕をセットして実行する
func applyFunctionWithPizza(fn *object.Function, args []object.Object) object.Object {
	// 関数呼び出し用の新しい環境を作成
	extendedEnv := object.NewEnclosedEnvironment(fn.Env)
	funcName, _ := fn.Name()
	
	// デバッグ情報
	if isArgumentsDebugEnabled {
		logger.Debug("関数呼び出し: %s", funcName)
		for i, arg := range args {
			logger.Debug("  引数%d: %s (%s)", i, arg.Inspect(), arg.Type())
		}
	}

	// 引数をバインド
	if len(args) > 0 {
		// 第1引数は🍕に設定
		// 🍕値を環境変数として設定（後方互換性のため）
		logger.Debug("第1引数を🍕にセット: %s", args[0].Inspect())
		extendedEnv.Set("🍕", args[0])
		LogArgumentBinding(funcName, "🍕", args[0])
		
		// 新しい実装: 🍕値を関数オブジェクト自体に設定
		logger.Debug("関数オブジェクトに🍕値を直接設定: %s", args[0].Inspect())
		fn.SetPizzaValue(args[0])
		
		// パラメータがある場合、パラメータに引数をバインド
		// これには二つのケースがある:
		// 1. 引数が1つだけの場合（パイプラインの基本的な動作）: 🍕と最初のパラメータに同じ値をバインド
		// 2. 引数が複数ある場合: 2番目以降の引数を順番にパラメータにバインド
		if len(fn.Parameters) > 0 {
			if len(args) == 1 {
				// 引数が1つの場合、最初のパラメータにも同じ値をバインド（利便性のため）
				paramName := fn.Parameters[0].Value
				extendedEnv.Set(paramName, args[0])
				LogArgumentBinding(funcName, paramName, args[0])
			} else {
				// 引数が複数の場合、2番目以降の引数を順番にパラメータにバインド
				for i := 0; i < len(fn.Parameters) && i+1 < len(args); i++ {
					paramName := fn.Parameters[i].Value
					extendedEnv.Set(paramName, args[i+1])
					LogArgumentBinding(funcName, paramName, args[i+1])
				}
			}
		}

		// デバッグ詳細情報
		if isArgumentsDebugEnabled {
			// パラメータの詳細をログに出力
			for i, param := range fn.Parameters {
				logger.Debug("パラメータ%d: %s", i, param.Value)
			}
			
			// 環境内の全変数をデバッグ出力
			logger.Debug("関数環境内の全変数:")
			for k, v := range extendedEnv.GetVariables() {
				logger.Debug("  %s = %s", k, v.Inspect())
			}
		}
	} else if len(fn.Parameters) > 0 {
		// 引数が必要なのに渡されていない場合はnilを返す
		logger.Debug("引数がまったくありませんが、関数は引数を必要としています")
		return nil
	}

	// 関数本体を評価
	astBody, ok := fn.ASTBody.(*ast.BlockStatement)
	if !ok {
		return createEvalError("関数の本体がBlockStatementではありません")
	}

	// 現在の関数コンテキストを保存
	prevFunction := currentFunction
	
	// 現在の関数コンテキストを設定
	logger.Debug("現在の関数コンテキストを設定: %s", funcName)
	currentFunction = fn
	
	logger.Debug("関数 '%s' の本体を評価中...", funcName)
	evaluated := evalBlockStatement(astBody, extendedEnv)
	logger.Debug("関数 '%s' の評価結果: %s (%T)", funcName, evaluated.Inspect(), evaluated)
	
	// 元の関数コンテキストを復元
	logger.Debug("元の関数コンテキストを復元")
	currentFunction = prevFunction

	// ReturnValue の場合は Value を抽出
	if returnValue, ok := evaluated.(*object.ReturnValue); ok {
		logger.Debug("関数 '%s' から戻り値を受け取りました: %s", funcName, returnValue.Inspect())
		// Valueフィールドがnilの場合は空のオブジェクトを返す
		if returnValue.Value == nil {
			logger.Debug("戻り値が nil です、NULL を返します")
			return NullObj
		}
		return returnValue.Value
	}

	logger.Debug("通常の評価結果を返します: %s", evaluated.Inspect())
	return evaluated
}