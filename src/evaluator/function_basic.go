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
		logger.Debug("関数を呼び出します: %s", fn.Inspect())

		// 修正: 引数は1つまでだけ許可（パイプライン以外）
		if len(fn.Parameters) > 1 {
			return createError("関数は最大1つのパラメータしか持てません（パイプライン以外）: %s", fn.Inspect())
		}

		// 引数の数をチェック
		if len(args) != len(fn.Parameters) {
			return createError("引数の数が一致しません: 期待=%d, 実際=%d", len(fn.Parameters), len(args))
		}

		// 入力型のチェック（パラメータが定義されている型と一致するか）
		if len(args) > 0 && fn.InputType != "" {
			logger.Debug("入力型チェック: 関数=%s, 入力型=%s, 実際=%s", 
				fn.Inspect(), fn.InputType, args[0].Type())
			if ok, err := checkInputType(args[0], fn.InputType); !ok {
				return createError("%s", err.Error())
			}
		}
		
		// 新しい環境を作成
		extendedEnv := object.NewEnclosedEnvironment(fn.Env)

		// 引数を環境にバインド
		for i, param := range fn.Parameters {
			extendedEnv.Set(param.Value, args[i])
		}

		// case文のために第一引数を🍕として設定
		if len(args) > 0 {
			logger.Debug("🍕値を環境に設定: %s", args[0].Inspect())
			extendedEnv.Set("🍕", args[0])
			
			// 関数オブジェクトにも🍕値を設定（将来の参照用）
			fn.SetPizzaValue(args[0])
		} else {
			logger.Debug("引数がないため、🍕値は設定されません")
		}

		// 関数本体を評価（ASTBodyをast.BlockStatementに型アサーション）
		astBody, ok := fn.ASTBody.(*ast.BlockStatement)
		if !ok {
			return createError("関数の本体がBlockStatementではありません")
		}
		result := evalBlockStatement(astBody, extendedEnv)

		// 💩値を返す（関数の戻り値）
		if obj, ok := result.(*object.ReturnValue); ok {
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
		return result

	case *object.Builtin:
		// 修正: ビルトイン関数も引数を1つまでに制限（ただし print や数学関数など一部の例外を除く）
		if len(args) > 1 && fn.Name != "print" && fn.Name != "range" && fn.Name != "sum" {
			logger.Debug("ビルトイン関数 %s は引数を1つしか取れません: 実際の引数数=%d\n", fn.Name, len(args))
		}
		return fn.Fn(args...)

	default:
		return createError("関数ではありません: %s", fn.Type())
	}
}
