package evaluator

import (
	"fmt"

	"github.com/uncode/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

// Eval は抽象構文木を評価する
func Eval(node interface{}, env *object.Environment) object.Object {
	switch node := node.(type) {
	// 文
	case *object.BlockStatement:
		return evalBlockStatement(node, env)
	// その他のケース
	default:
		return NULL
	}
}

// evalBlockStatement はブロック文を評価する
func evalBlockStatement(block *object.BlockStatement, env *object.Environment) object.Object {
	var result object.Object = NULL

	// 実際の実装ではここでステートメントを評価する
	
	return result
}

// 組み込み関数のマップ
var builtins = map[string]*object.Builtin{
	"print": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}
			return NULL
		},
	},
}
