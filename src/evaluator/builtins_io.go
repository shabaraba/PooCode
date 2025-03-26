package evaluator

import (
	"fmt"
	"github.com/uncode/logger"
	"github.com/uncode/object"
)

// registerIOBuiltins は入出力関連の組み込み関数を登録する
func registerIOBuiltins() {
	// 標準出力に出力する関数
	Builtins["print"] = &object.Builtin{
		Name: "print",
		Fn: func(args ...object.Object) object.Object {
			// デバッグ情報：受け取った引数の詳細を出力
			logIfEnabled(logger.LevelDebug, "print関数が受け取った引数: %d個", len(args))
			for i, arg := range args {
				logIfEnabled(logger.LevelDebug, "引数%d - タイプ: %s, 値: %s", i, arg.Type(), arg.Inspect())
				
				// arg.Inspect()ではなく実際の値を表示
				switch arg.Type() {
				case object.INTEGER_OBJ:
					intVal := arg.(*object.Integer).Value
					logIfEnabled(logger.LevelDebug, "整数値として %d を出力", intVal)
					fmt.Println(intVal)
				case object.STRING_OBJ:
					strVal := arg.(*object.String).Value
					logIfEnabled(logger.LevelDebug, "文字列として \"%s\" を出力", strVal)
					fmt.Println(strVal)
				case object.BOOLEAN_OBJ:
					boolVal := arg.(*object.Boolean).Value
					logIfEnabled(logger.LevelDebug, "真偽値として %t を出力", boolVal)
					fmt.Println(boolVal)
				default:
					inspectVal := arg.Inspect()
					logIfEnabled(logger.LevelDebug, "デフォルト - %s を出力", inspectVal)
					fmt.Println(inspectVal)
				}
			}
			// 第一引数を返すように変更（パイプラインの連鎖を維持するため）
			if len(args) > 0 {
				return args[0]
			}
			return NullObj
		},
		ReturnType: object.ANY_OBJ,
		ParamTypes: []object.ObjectType{object.ANY_OBJ},
	}
}
