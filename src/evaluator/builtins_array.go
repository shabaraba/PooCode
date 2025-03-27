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
			
			// 第1引数は配列
			if args[0].Type() \!= object.ARRAY_OBJ {
				logger.ComponentError(logger.ComponentBuiltin, "join関数の第1引数は配列である必要があります: %s", args[0].Type())
				return createError("join関数の第1引数は配列である必要があります: %s", args[0].Type())
			}
			array, _ := args[0].(*object.Array)
			
			// 第2引数は区切り文字
			if args[1].Type() \!= object.STRING_OBJ {
				logger.ComponentError(logger.ComponentBuiltin, "join関数の第2引数は文字列である必要があります: %s", args[1].Type())
				return createError("join関数の第2引数は文字列である必要があります: %s", args[1].Type())
			}
			delimiter, _ := args[1].(*object.String)
			
			logIfEnabled(builtinLogLevel, "join関数: 配列要素数=%d, 区切り文字='%s'", len(array.Elements), delimiter.Value)
			
			// 配列の各要素を文字列に変換
			elements := make([]string, len(array.Elements))
			for i, elem := range array.Elements {
				switch e := elem.(type) {
				case *object.String:
					elements[i] = e.Value
				case *object.Integer:
					elements[i] = fmt.Sprintf("%d", e.Value)
				case *object.Boolean:
					elements[i] = fmt.Sprintf("%t", e.Value)
				default:
					elements[i] = e.Inspect()
				}
			}
			
			result := strings.Join(elements, delimiter.Value)
			logIfEnabled(builtinLogLevel, "join関数: 結果='%s'", result)
			return &object.String{Value: result}
		},
		ReturnType: object.STRING_OBJ,
		ParamTypes: []object.ObjectType{object.ARRAY_OBJ, object.STRING_OBJ},
	}

	// 数値シーケンスを作成する関数
	Builtins["range"] = &object.Builtin{
		Name: "range",
		Fn: func(args ...object.Object) object.Object {
			logIfEnabled(builtinLogLevel, "range関数が呼び出されました: 引数数=%d", len(args))
			
			// 引数の数をチェック: 1または2つの引数を受け付ける
			if len(args) < 1 || len(args) > 2 {
				logger.ComponentError(logger.ComponentBuiltin, "range関数は1-2個の引数が必要です: %d個与えられました", len(args))
				return createError("range関数は1-2個の引数が必要です: %d個与えられました", len(args))
			}
			
			var start, end int64
			
			// 1つの引数の場合: 0からその値まで
			if len(args) == 1 {
				if args[0].Type() \!= object.INTEGER_OBJ {
					logger.ComponentError(logger.ComponentBuiltin, "range関数の引数は整数である必要があります: %s", args[0].Type())
					return createError("range関数の引数は整数である必要があります: %s", args[0].Type())
				}
				endVal, _ := args[0].(*object.Integer)
				
				start = 0
				end = endVal.Value
				logIfEnabled(builtinLogLevel, "range関数: 0から%dまでの範囲を生成", end)
			} else {
				// 2つの引数の場合: startからendまで
				if args[0].Type() \!= object.INTEGER_OBJ || args[1].Type() \!= object.INTEGER_OBJ {
					logger.ComponentError(logger.ComponentBuiltin, "range関数の引数は整数である必要があります")
					return createError("range関数の引数は整数である必要があります")
				}
				
				startVal, _ := args[0].(*object.Integer)
				endVal, _ := args[1].(*object.Integer)
				
				start = startVal.Value
				end = endVal.Value
				logIfEnabled(builtinLogLevel, "range関数: %dから%dまでの範囲を生成", start, end)
			}
			
			// 開始位置が終了位置より大きい場合は空の配列を返す
			if start > end {
				logger.ComponentWarn(logger.ComponentBuiltin, "range関数: 開始値 %d が終了値 %d より大きいため、空配列を返します", start, end)
				return &object.Array{Elements: []object.Object{}}
			}
			
			// 配列を作成
			elements := make([]object.Object, end-start)
			for i := start; i < end; i++ {
				elements[i-start] = &object.Integer{Value: i}
			}
			
			logIfEnabled(builtinLogLevel, "range関数: 結果の配列要素数=%d", len(elements))
			return &object.Array{Elements: elements}
		},
		ReturnType: object.ARRAY_OBJ,
		ParamTypes: []object.ObjectType{object.INTEGER_OBJ, object.INTEGER_OBJ},
	}
	
	// map関数 - 配列の各要素に関数を適用する
	Builtins["map"] = &object.Builtin{
		Name: "map",
		Fn: func(args ...object.Object) object.Object {
			logIfEnabled(builtinLogLevel, "map関数が呼び出されました: 引数数=%d", len(args))
			
			// 引数の数チェック
			if len(args) \!= 2 {
				logger.ComponentError(logger.ComponentBuiltin, "map関数は2つの引数が必要です: 配列, 関数")
				return createError("map関数は2つの引数が必要です: 配列, 関数")
			}
			
			// 第1引数が配列かチェック
			arr, ok := args[0].(*object.Array)
			if \!ok {
				logger.ComponentError(logger.ComponentBuiltin, "map関数の第1引数は配列である必要があります: %s", args[0].Type())
				return createError("map関数の第1引数は配列である必要があります: %s", args[0].Type())
			}
			
			// 第2引数が関数かチェック
			fn, ok := args[1].(*object.Function)
			if \!ok {
				logger.ComponentError(logger.ComponentBuiltin, "map関数の第2引数は関数である必要があります: %s", args[1].Type())
				return createError("map関数の第2引数は関数である必要があります: %s", args[1].Type())
			}
			
			// map関数の引数のパラメータは空である必要がある
			if len(fn.Parameters) > 0 {
				logger.ComponentError(logger.ComponentBuiltin, "map関数に渡された関数はパラメーターを取るべきではありません")
				return createError("map関数に渡された関数はパラメーターを取るべきではありません")
			}
			
			logIfEnabled(builtinLogLevel, "map関数: 配列要素数=%d", len(arr.Elements))
			
			// 結果の配列
			resultElements := make([]object.Object, 0, len(arr.Elements))
			
			// 配列の各要素に関数を適用
			for i, elem := range arr.Elements {
				logIfEnabled(builtinLogLevel, "map関数: 要素 %d を処理中: %s", i, elem.Inspect())
				
				// 関数の環境を拡張して🍕に現在の要素を設定
				extendedEnv := object.NewEnclosedEnvironment(fn.Env)
				extendedEnv.Set("🍕", elem)
				
				// 関数を評価
				result := Eval(fn.ASTBody, extendedEnv)
				
				// エラーチェック
				if errObj, ok := result.(*object.Error); ok {
					logger.ComponentError(logger.ComponentBuiltin, "map関数の処理中にエラーが発生: %s", errObj.Message)
					return errObj
				}
				
				// ReturnValueをアンラップ
				if retVal, ok := result.(*object.ReturnValue); ok {
					result = retVal.Value
				}
				
				logIfEnabled(builtinLogLevel, "map関数: 要素 %d の処理結果: %s", i, result.Inspect())
				
				// 結果を配列に追加
				resultElements = append(resultElements, result)
			}
			
			logIfEnabled(builtinLogLevel, "map関数: 処理完了, 結果の配列要素数=%d", len(resultElements))
			return &object.Array{Elements: resultElements}
		},
		ReturnType: object.ARRAY_OBJ,
		ParamTypes: []object.ObjectType{object.ARRAY_OBJ, object.FUNCTION_OBJ},
	}
	
	// filter関数 - 条件に合致する要素のみを抽出する
	Builtins["filter"] = &object.Builtin{
		Name: "filter",
		Fn: func(args ...object.Object) object.Object {
			logIfEnabled(builtinLogLevel, "filter関数が呼び出されました: 引数数=%d", len(args))
			
			// 引数の数チェック
			if len(args) \!= 2 {
				logger.ComponentError(logger.ComponentBuiltin, "filter関数は2つの引数が必要です: 配列, 関数")
				return createError("filter関数は2つの引数が必要です: 配列, 関数")
			}
			
			// 第1引数が配列かチェック
			arr, ok := args[0].(*object.Array)
			if \!ok {
				logger.ComponentError(logger.ComponentBuiltin, "filter関数の第1引数は配列である必要があります: %s", args[0].Type())
				return createError("filter関数の第1引数は配列である必要があります: %s", args[0].Type())
			}
			
			// 第2引数が関数かチェック
			fn, ok := args[1].(*object.Function)
			if \!ok {
				logger.ComponentError(logger.ComponentBuiltin, "filter関数の第2引数は関数である必要があります: %s", args[1].Type())
				return createError("filter関数の第2引数は関数である必要があります: %s", args[1].Type())
			}
			
			// filter関数の引数のパラメータは空である必要がある
			if len(fn.Parameters) > 0 {
				logger.ComponentError(logger.ComponentBuiltin, "filter関数に渡された関数はパラメーターを取るべきではありません")
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
