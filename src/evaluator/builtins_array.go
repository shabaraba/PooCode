package evaluator

import (
	"fmt"
	"strings"
	"github.com/uncode/object"
)

// registerArrayBuiltins は配列関連の組み込み関数を登録する
func registerArrayBuiltins() {
	// 配列を連結して文字列にする関数
	Builtins["join"] = &object.Builtin{
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
	}

	// 数値シーケンスを作成する関数
	Builtins["range"] = &object.Builtin{
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
	}
}
