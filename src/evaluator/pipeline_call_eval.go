package evaluator

import (
	"github.com/uncode/ast"
	"github.com/uncode/object"
)

// evalCallExpressionForPipeline はパイプライン用に関数呼び出し式を評価する特別な関数
func evalCallExpressionForPipeline(callExpr *ast.CallExpression, env *object.Environment) object.Object {
	// 関数名を取得
	var funcName string
	if ident, ok := callExpr.Function.(*ast.Identifier); ok {
		funcName = ident.Value
	} else {
		return createError("関数名を取得できません: %T", callExpr.Function)
	}
	
	// 引数を評価
	args := evalExpressions(callExpr.Arguments, env)
	if len(args) > 0 && args[0].Type() == object.ERROR_OBJ {
		return args[0]
	}
	
	// 環境から関数を検索
	fn, exists := env.Get(funcName)
	if !exists {
		// ビルトイン関数を確認
		if builtin, ok := Builtins[funcName]; ok {
			return builtin
		}
		return createError("関数 '%s' が見つかりません", funcName)
	}
	
	// 関数オブジェクトの場合
	if function, ok := fn.(*object.Function); ok {
		// 引数付き関数を作成して返す
		// 🍕については後で設定するので、ここでは引数だけを持った関数として返す
		LogPipe("関数 '%s' に引数 %d 個を設定\n", funcName, len(args))
		
		// 新しい関数オブジェクトを作成（パラメータと引数を持つ）
		newFunction := &object.Function{
			Parameters: function.Parameters,
			ASTBody:    function.ASTBody,
			Env:        function.Env,
			InputType:  function.InputType,
			ReturnType: function.ReturnType,
			// 重要: 引数を保存
			ParamValues: args,
		}
		
		return newFunction
	}
	
	// その他のケース（ビルトイン関数など）
	return fn
}
