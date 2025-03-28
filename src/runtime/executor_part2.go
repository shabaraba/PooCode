package runtime

import (
	"fmt"

	"github.com/uncode/ast"
	"github.com/uncode/evaluator"
	"github.com/uncode/logger"
	"github.com/uncode/object"
)

// SetupBuiltins は組み込み関数を環境に設定する
func SetupBuiltins(env *object.Environment) {
	// プリント関数を追加
	env.Set("print", &object.Builtin{
		Name: "print",
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}
			return evaluator.NullObj
		},
	})
	
	// 評価器から組み込み関数をすべてインポート
	// evaluator.Builtinsに登録されている関数をすべて環境に追加
	for name, builtin := range evaluator.Builtins {
		logger.Debug("組み込み関数を登録: %s", name)
		env.Set(name, builtin)
	}
}

// convertToObjectIdentifiers は ast.Identifier スライスを object.Identifier スライスに変換する
func convertToObjectIdentifiers(params []*ast.Identifier) []*object.Identifier {
	if params == nil {
		return nil
	}
	
	result := make([]*object.Identifier, len(params))
	for i, param := range params {
		result[i] = &object.Identifier{Value: param.Value}
	}
	return result
}