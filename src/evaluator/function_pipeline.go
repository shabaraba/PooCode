package evaluator

import (
	"github.com/uncode/ast"
	"github.com/uncode/logger"
	"github.com/uncode/object"
)

// applyPipelineFunction は関数を適用する（パイプラインの場合同様に🍕も設定）
func applyPipelineFunction(fn *object.Function, args []object.Object) object.Object {
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

	// 入力型のチェック
	if len(args) > 0 && fn.InputType != "" {
		logger.Debug("入力型チェック: 関数=%s, 入力型=%s, 実際=%s", 
			fn.Inspect(), fn.InputType, args[0].Type())
		if ok, err := checkInputType(args[0], fn.InputType); !ok {
			return createError("%s", err.Error())
		}
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
		return createError("関数の本体がBlockStatementではありません")
	}

	logger.Debug("関数本体を評価します...")
	result := evalBlockStatement(astBody, extendedEnv)

	// 💩値を返す（関数の戻り値）
	if obj, ok := result.(*object.ReturnValue); ok {
		logger.Debug("関数から戻り値が見つかりました: %s\n", obj.Value.Inspect())
		
		// 戻り値の型チェック
		if fn.ReturnType != "" {
			logger.Debug("戻り値型チェック: 関数=%s, 戻り値型=%s, 実際=%s",
				fn.Inspect(), fn.ReturnType, obj.Value.Type())
			if ok, err := checkReturnType(obj.Value, fn.ReturnType); !ok {
				return createError("%s", err.Error())
			}
		}
		
		return obj.Value
	}

	logger.Debug("関数から戻り値なしで実行完了: %s\n", result.Inspect())
	return result
}
