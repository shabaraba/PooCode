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
