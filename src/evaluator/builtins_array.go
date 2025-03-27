package evaluator

import (
	"fmt"
	"strings"
	"github.com/uncode/logger"
	"github.com/uncode/object"
)

// log variables for debugging
var builtinLogLevel = logger.LevelTrace

// logIfEnabled logs a message if logging is enabled for the given level
func logIfEnabled(level logger.LogLevel, format string, args ...interface{}) {
	if logger.IsSpecialLevelEnabled(level) || logger.GetComponentLevel(logger.ComponentBuiltin) >= level {
		logger.ComponentDebug(logger.ComponentBuiltin, format, args...)
	}
}

// registerArrayBuiltins は配列関連の組み込み関数を登録する
func registerArrayBuiltins() {
	// 配列を連結して文字列にする関数
	Builtins["join"] = &object.Builtin{
		Name: "join",
		Fn: func(args ...object.Object) object.Object {
			logIfEnabled(builtinLogLevel, "join関数が呼び出されました: 引数数=%d", len(args))
			
			if len(args) \!= 2 {
				logger.ComponentError(logger.ComponentBuiltin, "join関数は2つの引数が必要です: %d個与えられました", len(args))
				return createError("join関数は2つの引数が必要です: %d個与えられました", len(args))
			}
				return createError("filter関数に渡された関数はパラメーターを取るべきではありません")
			}
			
			logIfEnabled(builtinLogLevel, "filter関数: 配列要素数=%d", len(arr.Elements))
			
			// 結果の配列
			resultElements := make([]object.Object, 0)
			
			// 配列の各要素に条件関数を適用
			for i, elem := range arr.Elements {
				logIfEnabled(builtinLogLevel, "filter関数: 要素 %d を処理中: %s", i, elem.Inspect())
				
				// 関数の環境を拡張して🍕に現在の要素を設定
				extendedEnv := object.NewEnclosedEnvironment(fn.Env)
				extendedEnv.Set("🍕", elem)
				
				// 条件関数を評価
				result := Eval(fn.ASTBody, extendedEnv)
				
				// エラーチェック
				if errObj, ok := result.(*object.Error); ok {
					logger.ComponentError(logger.ComponentBuiltin, "filter関数の処理中にエラーが発生: %s", errObj.Message)
					return errObj
				}
				
				// ReturnValueをアンラップ
				if retVal, ok := result.(*object.ReturnValue); ok {
					result = retVal.Value
				}
				
				// 結果が真の場合、要素を結果配列に追加
				if isTruthy(result) {
					logIfEnabled(builtinLogLevel, "filter関数: 要素 %d は条件を満たします", i)
					resultElements = append(resultElements, elem)
				} else {
					logIfEnabled(builtinLogLevel, "filter関数: 要素 %d は条件を満たしません", i)
				}
			}
			
			logIfEnabled(builtinLogLevel, "filter関数: 処理完了, 結果の配列要素数=%d", len(resultElements))
			return &object.Array{Elements: resultElements}
		},
		ReturnType: object.ARRAY_OBJ,
		ParamTypes: []object.ObjectType{object.ARRAY_OBJ, object.FUNCTION_OBJ},
	}
}
