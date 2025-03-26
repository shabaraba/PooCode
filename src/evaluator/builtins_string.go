package evaluator

import (
	"fmt"
	"strings"
	"github.com/uncode/object"
)

// registerStringBuiltins は文字列関連の組み込み関数を登録する
func registerStringBuiltins() {
	// 文字列を作成する関数
	Builtins["to_string"] = &object.Builtin{
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
	}

	// 文字列または配列の長さを取得する関数
	Builtins["length"] = &object.Builtin{
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
	}

	// 文字列を分割する関数
	Builtins["split"] = &object.Builtin{
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
	}

	// 部分文字列を取得する関数
	Builtins["substring"] = &object.Builtin{
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
	}

	// 大文字に変換する関数
	Builtins["to_upper"] = &object.Builtin{
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
	}

	// 小文字に変換する関数
	Builtins["to_lower"] = &object.Builtin{
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
	}
}
