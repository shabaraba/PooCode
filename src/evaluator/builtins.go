package evaluator

import (
	"fmt"
	"strings"

	"github.com/uncode/logger"
	"github.com/uncode/object"
)

// Forward declarations to avoid initialization cycles
var Builtins map[string]*object.Builtin

// init initializes the builtin functions map
func init() {
	// Initialize the Builtins map after creating function literals
	Builtins = map[string]*object.Builtin{
		"print": &object.Builtin{
			Name: "print",
			Fn: func(args ...object.Object) object.Object {
				// デバッグ情報：受け取った引数の詳細を出力
				logger.Debug("DEBUG: print関数が受け取った引数: %d個", len(args))
				for i, arg := range args {
					logger.Debug("DEBUG: 引数%d - タイプ: %s, 値: %s", i, arg.Type(), arg.Inspect())
					
					// arg.Inspect()ではなく実際の値を表示
					switch arg.Type() {
					case object.INTEGER_OBJ:
						intVal := arg.(*object.Integer).Value
						logger.Debug("DEBUG: 整数値として %d を出力", intVal)
						fmt.Println(intVal)
					case object.STRING_OBJ:
						strVal := arg.(*object.String).Value
						logger.Debug("DEBUG: 文字列として \"%s\" を出力", strVal)
						fmt.Println(strVal)
					case object.BOOLEAN_OBJ:
						boolVal := arg.(*object.Boolean).Value
						logger.Debug("DEBUG: 真偽値として %t を出力", boolVal)
						fmt.Println(boolVal)
					default:
						inspectVal := arg.Inspect()
						logger.Debug("DEBUG: デフォルト - %s を出力", inspectVal)
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
		},
		"add": &object.Builtin{
			Name: "add",
			Fn: func(args ...object.Object) object.Object {
				if len(args) == 0 {
					return createError("add関数は少なくとも1つの引数が必要です")
				}
				
				// 文字列加算の場合
				if str, ok := args[0].(*object.String); ok {
					logger.Debug("add関数: 文字列連結モード")
					
					// 第2引数があれば文字列に変換して連結
					if len(args) > 1 {
						var rightStr string
						switch right := args[1].(type) {
						case *object.String:
							rightStr = right.Value
						case *object.Integer:
							rightStr = fmt.Sprintf("%d", right.Value)
						case *object.Boolean:
							rightStr = fmt.Sprintf("%t", right.Value)
						default:
							rightStr = right.Inspect()
						}
						
						return &object.String{Value: str.Value + rightStr}
					}
					return str
				}
				
				// 整数加算（単一引数の場合は値をそのまま返す）
				left, ok := args[0].(*object.Integer)
				if !ok {
					return createError("add関数の第1引数は整数または文字列である必要があります: %s", args[0].Type())
				}
				
				// 第2引数がない場合は値をそのまま返す
				if len(args) == 1 {
					logger.Debug("add関数: 単一引数 %d をそのまま返します", left.Value)
					return left
				}
				
				// 第2引数があれば加算
				right, ok := args[1].(*object.Integer)
				if !ok {
					return createError("add関数の第2引数は整数である必要があります: %s", args[1].Type())
				}
				
				result := &object.Integer{Value: left.Value + right.Value}
				logger.Debug("add関数: %d + %d = %d", left.Value, right.Value, result.Value)
				return result
			},
			ReturnType: object.ANY_OBJ, // 文字列または整数を返す可能性あり
			ParamTypes: []object.ObjectType{object.ANY_OBJ},
		},
		"sub": &object.Builtin{
			Name: "sub",
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 2 {
					return createError("sub関数は2つの引数が必要です: %d個与えられました", len(args))
				}
				
				// 整数減算
				left, ok := args[0].(*object.Integer)
				if !ok {
					return createError("sub関数の第1引数は整数である必要があります: %s", args[0].Type())
				}
				
				right, ok := args[1].(*object.Integer)
				if !ok {
					return createError("sub関数の第2引数は整数である必要があります: %s", args[1].Type())
				}
				
				return &object.Integer{Value: left.Value - right.Value}
			},
			ReturnType: object.INTEGER_OBJ,
			ParamTypes: []object.ObjectType{object.INTEGER_OBJ, object.INTEGER_OBJ},
		},
		"mul": &object.Builtin{
			Name: "mul",
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 2 {
					return createError("mul関数は2つの引数が必要です: %d個与えられました", len(args))
				}
				
				// 整数乗算
				left, ok := args[0].(*object.Integer)
				if !ok {
					return createError("mul関数の第1引数は整数である必要があります: %s", args[0].Type())
				}
				
				right, ok := args[1].(*object.Integer)
				if !ok {
					return createError("mul関数の第2引数は整数である必要があります: %s", args[1].Type())
				}
				
				return &object.Integer{Value: left.Value * right.Value}
			},
			ReturnType: object.INTEGER_OBJ,
			ParamTypes: []object.ObjectType{object.INTEGER_OBJ, object.INTEGER_OBJ},
		},
		"div": &object.Builtin{
			Name: "div",
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 2 {
					return createError("div関数は2つの引数が必要です: %d個与えられました", len(args))
				}
				
				// 整数除算
				left, ok := args[0].(*object.Integer)
				if !ok {
					return createError("div関数の第1引数は整数である必要があります: %s", args[0].Type())
				}
				
				right, ok := args[1].(*object.Integer)
				if !ok {
					return createError("div関数の第2引数は整数である必要があります: %s", args[1].Type())
				}
				
				// ゼロ除算チェック
				if right.Value == 0 {
					return createError("ゼロによる除算: %d / 0", left.Value)
				}
				
				return &object.Integer{Value: left.Value / right.Value}
			},
			ReturnType: object.INTEGER_OBJ,
			ParamTypes: []object.ObjectType{object.INTEGER_OBJ, object.INTEGER_OBJ},
		},
		"mod": &object.Builtin{
			Name: "mod",
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 2 {
					return createError("mod関数は2つの引数が必要です: %d個与えられました", len(args))
				}
				
				// 整数剰余
				left, ok := args[0].(*object.Integer)
				if !ok {
					return createError("mod関数の第1引数は整数である必要があります: %s", args[0].Type())
				}
				
				right, ok := args[1].(*object.Integer)
				if !ok {
					return createError("mod関数の第2引数は整数である必要があります: %s", args[1].Type())
				}
				
				// ゼロ除算チェック
				if right.Value == 0 {
					return createError("ゼロによるモジュロ: %d %% 0", left.Value)
				}
				
				return &object.Integer{Value: left.Value % right.Value}
			},
			ReturnType: object.INTEGER_OBJ,
			ParamTypes: []object.ObjectType{object.INTEGER_OBJ, object.INTEGER_OBJ},
		},
		"pow": &object.Builtin{
			Name: "pow",
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 2 {
					return createError("pow関数は2つの引数が必要です: %d個与えられました", len(args))
				}
				
				// べき乗
				base, ok := args[0].(*object.Integer)
				if !ok {
					return createError("pow関数の第1引数は整数である必要があります: %s", args[0].Type())
				}
				
				exp, ok := args[1].(*object.Integer)
				if !ok {
					return createError("pow関数の第2引数は整数である必要があります: %s", args[1].Type())
				}
				
				// 負の指数のチェック
				if exp.Value < 0 {
					return createError("pow関数の指数は0以上である必要があります: %d", exp.Value)
				}
				
				result := int64(1)
				for i := int64(0); i < exp.Value; i++ {
					result *= base.Value
				}
				
				return &object.Integer{Value: result}
			},
			ReturnType: object.INTEGER_OBJ,
			ParamTypes: []object.ObjectType{object.INTEGER_OBJ, object.INTEGER_OBJ},
		},
		"to_string": &object.Builtin{
			Name: "to_string",
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return createError("to_string関数は1つの引数が必要です: %d個与えられました", len(args))
				}
				
				switch arg := args[0].(type) {
				case *object.String:
					return arg // 既に文字列
				case *object.Integer:
					return &object.String{Value: fmt.Sprintf("%d", arg.Value)}
				case *object.Boolean:
					return &object.String{Value: fmt.Sprintf("%t", arg.Value)}
				default:
					return &object.String{Value: arg.Inspect()}
				}
			},
			ReturnType: object.STRING_OBJ,
			ParamTypes: []object.ObjectType{object.ANY_OBJ},
		},
		"length": &object.Builtin{
			Name: "length",
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return createError("length関数は1つの引数が必要です: %d個与えられました", len(args))
				}
				
				switch arg := args[0].(type) {
				case *object.String:
					return &object.Integer{Value: int64(len(arg.Value))}
				case *object.Array:
					return &object.Integer{Value: int64(len(arg.Elements))}
				default:
					return createError("length関数は文字列または配列に対してのみ使用できます: %s", args[0].Type())
				}
			},
			ReturnType: object.INTEGER_OBJ,
			ParamTypes: []object.ObjectType{object.ANY_OBJ},
		},
		"eq": &object.Builtin{
			Name: "eq",
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 2 {
					return createError("eq関数は2つの引数が必要です: %d個与えられました", len(args))
				}
				
				switch left := args[0].(type) {
				case *object.Integer:
					if right, ok := args[1].(*object.Integer); ok {
						return &object.Boolean{Value: left.Value == right.Value}
					}
				case *object.String:
					if right, ok := args[1].(*object.String); ok {
						return &object.Boolean{Value: left.Value == right.Value}
					}
				case *object.Boolean:
					if right, ok := args[1].(*object.Boolean); ok {
						return &object.Boolean{Value: left.Value == right.Value}
					}
				}
				
				return &object.Boolean{Value: false}
			},
			ReturnType: object.BOOLEAN_OBJ,
			ParamTypes: []object.ObjectType{object.ANY_OBJ, object.ANY_OBJ},
		},
		"not": &object.Builtin{
			Name: "not",
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return createError("not関数は1つの引数が必要です: %d個与えられました", len(args))
				}
				
				if b, ok := args[0].(*object.Boolean); ok {
					return &object.Boolean{Value: !b.Value}
				}
				
				return &object.Boolean{Value: false} // デフォルトはfalse
			},
			ReturnType: object.BOOLEAN_OBJ,
			ParamTypes: []object.ObjectType{object.BOOLEAN_OBJ},
		},
		"split": &object.Builtin{
			Name: "split",
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 2 {
					return createError("split関数は2つの引数が必要です: %d個与えられました", len(args))
				}
				
				// 第1引数は対象文字列
				if args[0].Type() != object.STRING_OBJ {
					return createError("split関数の第1引数は文字列である必要があります: %s", args[0].Type())
				}
				str, _ := args[0].(*object.String)
				
				// 第2引数は区切り文字
				if args[1].Type() != object.STRING_OBJ {
					return createError("split関数の第2引数は文字列である必要があります: %s", args[1].Type())
				}
				delimiter, _ := args[1].(*object.String)
				
				// 文字列を分割
				parts := strings.Split(str.Value, delimiter.Value)
				
				// 配列を作成
				elements := make([]object.Object, len(parts))
				for i, part := range parts {
					elements[i] = &object.String{Value: part}
				}
				
				return &object.Array{Elements: elements}
			},
			ReturnType: object.ARRAY_OBJ,
			ParamTypes: []object.ObjectType{object.STRING_OBJ, object.STRING_OBJ},
		},
		"join": &object.Builtin{
			Name: "join",
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 2 {
					return createError("join関数は2つの引数が必要です: %d個与えられました", len(args))
				}
				
				// 第1引数は配列
				if args[0].Type() != object.ARRAY_OBJ {
					return createError("join関数の第1引数は配列である必要があります: %s", args[0].Type())
				}
				array, _ := args[0].(*object.Array)
				
				// 第2引数は区切り文字
				if args[1].Type() != object.STRING_OBJ {
					return createError("join関数の第2引数は文字列である必要があります: %s", args[1].Type())
				}
				delimiter, _ := args[1].(*object.String)
				
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
				
				return &object.String{Value: strings.Join(elements, delimiter.Value)}
			},
			ReturnType: object.STRING_OBJ,
			ParamTypes: []object.ObjectType{object.ARRAY_OBJ, object.STRING_OBJ},
		},
		"substring": &object.Builtin{
			Name: "substring",
			Fn: func(args ...object.Object) object.Object {
				// 引数の数をチェック
				if len(args) < 2 || len(args) > 3 {
					return createError("substring関数は2-3個の引数が必要です: %d個与えられました", len(args))
				}
				
				// 第1引数は文字列
				if args[0].Type() != object.STRING_OBJ {
					return createError("substring関数の第1引数は文字列である必要があります: %s", args[0].Type())
				}
				str, _ := args[0].(*object.String)
				
				// 第2引数は開始位置
				if args[1].Type() != object.INTEGER_OBJ {
					return createError("substring関数の第2引数は整数である必要があります: %s", args[1].Type())
				}
				start, _ := args[1].(*object.Integer)
				
				// 文字列の長さを取得
				strLen := int64(len(str.Value))
				
				// 開始位置のバリデーション
				if start.Value < 0 {
					start.Value = 0
				}
				if start.Value >= strLen {
					return &object.String{Value: ""}
				}
				
				// 第3引数がある場合は終了位置
				if len(args) == 3 {
					if args[2].Type() != object.INTEGER_OBJ {
						return createError("substring関数の第3引数は整数である必要があります: %s", args[2].Type())
					}
					end, _ := args[2].(*object.Integer)
					
					// 終了位置のバリデーション
					if end.Value < start.Value {
						return &object.String{Value: ""}
					}
					if end.Value > strLen {
						end.Value = strLen
					}
					
					return &object.String{Value: str.Value[start.Value:end.Value]}
				}
				
				// 第3引数がない場合は文字列の最後まで
				return &object.String{Value: str.Value[start.Value:]}
			},
			ReturnType: object.STRING_OBJ,
			ParamTypes: []object.ObjectType{object.STRING_OBJ, object.INTEGER_OBJ, object.INTEGER_OBJ},
		},
		"to_upper": &object.Builtin{
			Name: "to_upper",
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return createError("to_upper関数は1つの引数が必要です: %d個与えられました", len(args))
				}
				
				if args[0].Type() != object.STRING_OBJ {
					return createError("to_upper関数の引数は文字列である必要があります: %s", args[0].Type())
				}
				str, _ := args[0].(*object.String)
				
				return &object.String{Value: strings.ToUpper(str.Value)}
			},
			ReturnType: object.STRING_OBJ,
			ParamTypes: []object.ObjectType{object.STRING_OBJ},
		},
		"to_lower": &object.Builtin{
			Name: "to_lower",
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return createError("to_lower関数は1つの引数が必要です: %d個与えられました", len(args))
				}
				
				if args[0].Type() != object.STRING_OBJ {
					return createError("to_lower関数の引数は文字列である必要があります: %s", args[0].Type())
				}
				str, _ := args[0].(*object.String)
				
				return &object.String{Value: strings.ToLower(str.Value)}
			},
			ReturnType: object.STRING_OBJ,
			ParamTypes: []object.ObjectType{object.STRING_OBJ},
		},
		"range": &object.Builtin{
			Name: "range",
			Fn: func(args ...object.Object) object.Object {
				// 引数の数をチェック: 1または2つの引数を受け付ける
				if len(args) < 1 || len(args) > 2 {
					return createError("range関数は1-2個の引数が必要です: %d個与えられました", len(args))
				}
				
				var start, end int64
				
				// 1つの引数の場合: 0からその値まで
				if len(args) == 1 {
					if args[0].Type() != object.INTEGER_OBJ {
						return createError("range関数の引数は整数である必要があります: %s", args[0].Type())
					}
					endVal, _ := args[0].(*object.Integer)
					
					start = 0
					end = endVal.Value
				} else {
					// 2つの引数の場合: startからendまで
					if args[0].Type() != object.INTEGER_OBJ || args[1].Type() != object.INTEGER_OBJ {
						return createError("range関数の引数は整数である必要があります")
					}
					
					startVal, _ := args[0].(*object.Integer)
					endVal, _ := args[1].(*object.Integer)
					
					start = startVal.Value
					end = endVal.Value
				}
				
				// 開始位置が終了位置より大きい場合は空の配列を返す
				if start > end {
					return &object.Array{Elements: []object.Object{}}
				}
				
				// 配列を作成
				elements := make([]object.Object, end-start)
				for i := start; i < end; i++ {
					elements[i-start] = &object.Integer{Value: i}
				}
				
				return &object.Array{Elements: elements}
			},
			ReturnType: object.ARRAY_OBJ,
			ParamTypes: []object.ObjectType{object.INTEGER_OBJ, object.INTEGER_OBJ},
		},
		"sum": &object.Builtin{
			Name: "sum",
			Fn: func(args ...object.Object) object.Object {
				// 引数の数をチェック
				if len(args) != 1 {
					return createError("sum関数は1つの引数が必要です: %d個与えられました", len(args))
				}
				
				// 配列かどうかチェック
				if args[0].Type() != object.ARRAY_OBJ {
					return createError("sum関数の引数は配列である必要があります: %s", args[0].Type())
				}
				
				array, _ := args[0].(*object.Array)
				
				// 合計を計算
				sum := int64(0)
				for _, elem := range array.Elements {
					if elem.Type() != object.INTEGER_OBJ {
						return createError("sum関数の配列要素はすべて整数である必要があります: %s", elem.Type())
					}
					
					intVal, _ := elem.(*object.Integer)
					sum += intVal.Value
				}
				
				return &object.Integer{Value: sum}
			},
			ReturnType: object.INTEGER_OBJ,
			ParamTypes: []object.ObjectType{object.ARRAY_OBJ},
		},
		"typeof": &object.Builtin{
			Name: "typeof",
			Fn: func(args ...object.Object) object.Object {
				if len(args) != 1 {
					return createError("typeof関数は1つの引数が必要です: %d個与えられました", len(args))
				}

				// 引数が文字列の場合、組み込み関数名として解釈
				if str, ok := args[0].(*object.String); ok {
					funcName := str.Value
					// This is potentially recursive, but it should be fine for normal use cases
					if builtin, exists := Builtins[funcName]; exists {
						return &object.String{Value: string(builtin.ReturnType)}
					}
					return createError("組み込み関数 '%s' は存在しません", funcName)
				}

				// その他の型はそのまま型情報を返す
				return &object.String{Value: string(args[0].Type())}
			},
			ReturnType: object.STRING_OBJ,
			ParamTypes: []object.ObjectType{object.ANY_OBJ},
		},
	}
}

// 組み込み関数の型情報を取得する関数
// GetBuiltinReturnType は組み込み関数の戻り値の型を返す
func GetBuiltinReturnType(name string) object.ObjectType {
	if builtin, ok := Builtins[name]; ok {
		return builtin.ReturnType
	}
	return object.NULL_OBJ
}

// GetBuiltinParamTypes は組み込み関数のパラメータの型を返す
func GetBuiltinParamTypes(name string) []object.ObjectType {
	if builtin, ok := Builtins[name]; ok {
		return builtin.ParamTypes
	}
	return nil
}
