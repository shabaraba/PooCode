package evaluator

import (
	"fmt"
	
	"github.com/uncode/ast"
	"github.com/uncode/object"
)

// applyFunction は関数を適用する
func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		// 関数呼び出しの実装
		fmt.Println("関数を呼び出します:", fn.Inspect())
		
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
		return fn.Fn(args...)
		
	default:
		return newError("関数ではありません: %s", fn.Type())
	}
}

// applyFunctionWithPizza は関数を適用する（パイプラインの場合同様に🍕も設定）
func applyFunctionWithPizza(fn *object.Function, args []object.Object) object.Object {
	// 関数呼び出しの実装
	fmt.Println("パイプライン対応で関数を呼び出します:", fn.Inspect())
	
	// 新しい環境を作成
	extendedEnv := object.NewEnclosedEnvironment(fn.Env)
	
	// 引数とパラメータのデバッグ出力
	fmt.Printf("関数パラメータ数: %d, 引数数: %d\n", len(fn.Parameters), len(args))
	for i, param := range fn.Parameters {
		fmt.Printf("  パラメータ %d: %s\n", i, param.Value)
	}
	for i, arg := range args {
		fmt.Printf("  引数 %d: %s\n", i, arg.Inspect())
	}
	
	// パイプラインを利用する関数では:
	// - 第1引数は常に🍕として設定される
	// - パラメータがある場合、引数の残りをパラメータにマッピングする
	if len(args) > 0 {
		// 🍕 変数を設定
		extendedEnv.Set("🍕", args[0])
		fmt.Printf("🍕 に値 %s を設定しました\n", args[0].Inspect())
		
		// パラメータを環境にバインド
		if len(fn.Parameters) > 0 {
			// パラメータ名を取得
			paramName := fn.Parameters[0].Value
			
			// numパラメータにどの値を設定するか
			if len(args) > 1 {
				// 複数引数の場合: 第2引数をnumに設定
				extendedEnv.Set(paramName, args[1])
				fmt.Printf("パラメータ '%s' に値 %s を設定しました\n", 
					paramName, args[1].Inspect())
			} else {
				// 単一引数の場合: 🍕と同じ値をnumに設定
				extendedEnv.Set(paramName, args[0])
				fmt.Printf("単一引数: パラメータ '%s' に値 %s を設定しました\n", 
					paramName, args[0].Inspect())
			}
		}
	}
	
	// 関数本体を評価（ASTBodyをast.BlockStatementに型アサーション）
	astBody, ok := fn.ASTBody.(*ast.BlockStatement)
	if !ok {
		return newError("関数の本体がBlockStatementではありません")
	}
	
	fmt.Println("関数本体を評価します...")
	result := evalBlockStatement(astBody, extendedEnv)
	
	// 💩値を返す（関数の戻り値）
	if obj, ok := result.(*object.ReturnValue); ok {
		fmt.Printf("関数から戻り値が見つかりました: %s\n", obj.Value.Inspect())
		return obj.Value
	}
	
	fmt.Printf("関数から戻り値なしで実行完了: %s\n", result.Inspect())
	return result
}

// applyNamedFunction は名前付き関数を検索し、適用する
// 同じ名前で複数の関数が存在する場合は、条件に基づいて適切な関数を選択する
func applyNamedFunction(env *object.Environment, name string, args []object.Object) object.Object {
	fmt.Println("***** applyNamedFunction が呼び出されました *****")
	fmt.Printf("関数名: %s、引数の数: %d\n", name, len(args))
	
	// デバッグ: 環境内のすべての変数を表示
	fmt.Println("現在の環境に登録されている変数:")
	for k, v := range env.GetVariables() {
		fmt.Printf("  %s: %s\n", k, v.Type())
	}
	
	// ビルトイン関数を確認
	if builtin, ok := Builtins[name]; ok {
		if debugMode {
			fmt.Printf("ビルトイン関数 '%s' を呼び出します\n", name)
		}
		return builtin.Fn(args...)
	}
	
	// 環境から同名のすべての関数を取得
	functions := env.GetAllFunctionsByName(name)
	
	if len(functions) == 0 {
		return newError("関数 '%s' が見つかりません", name)
	}
	
	// デバッグ情報
	if debugMode {
		fmt.Printf("関数 '%s' を呼び出します: %d 個の候補が見つかりました\n", name, len(functions))
	}
	
	// 関数が1つだけの場合は直接適用
	if len(functions) == 1 {
		fmt.Println("関数が1つだけ見つかりました")
		return applyFunctionWithPizza(functions[0], args)
	}
	
	fmt.Println("複数の関数が見つかりました:", len(functions))
	
	// 🍕 を設定（もし必要なら）
	if len(args) > 0 {
		fmt.Printf("🍕 に値 %s を設定します\n", args[0].Inspect())
		env.Set("🍕", args[0])
	} else {
		fmt.Println("引数が見つからないため、🍕は設定しません")
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
	fmt.Printf("条件付き関数を %d 個見つけました\n", len(conditionalFuncs))
	for i, fn := range conditionalFuncs {
		fmt.Printf("条件付き関数候補 %d を評価中...\n", i+1)
		fmt.Printf("条件式: %v\n", fn.Condition)
		
		// 条件式を評価するための環境を作成
		condEnv := object.NewEnclosedEnvironment(env)
		if len(args) > 0 {
			condEnv.Set("🍕", args[0])
		}
		
		// 条件式を評価
		condResult := Eval(fn.Condition, condEnv)
		fmt.Printf("条件式の評価結果: %s\n", condResult.Inspect())
		
		// 条件が真なら、この関数を使用
		// Booleanオブジェクトの場合はそのValueを使用、それ以外の場合はisTruthyで評価
		isTrue := false
		if condResult.Type() == object.BOOLEAN_OBJ {
			isTrue = condResult.(*object.Boolean).Value
			fmt.Printf("条件式の真偽値: %v\n", isTrue)
		} else {
			isTrue = isTruthy(condResult)
			fmt.Printf("条件式の評価結果（非Boolean）が %v と判定されました\n", isTrue)
		}
		
		if isTrue {
			fmt.Println("条件が真であるため、この関数を使用します")
			return applyFunctionWithPizza(fn, args)
		} else {
			fmt.Println("条件が偽であるため、この関数をスキップします")
		}
	}
	
	// 条件付き関数が該当しなかった場合、デフォルト関数を使用
	fmt.Printf("デフォルト関数を %d 個見つけました\n", len(defaultFuncs))
	if len(defaultFuncs) > 0 {
		fmt.Println("デフォルト関数を使用します")
		return applyFunctionWithPizza(defaultFuncs[0], args)
	}
	
	// 適用可能な関数が見つからない場合
	return newError("条件に一致する関数 '%s' が見つかりません", name)
}
