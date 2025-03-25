package evaluator

import (
	"fmt"
	"strings"

	"github.com/uncode/object"
)

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
	"show": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}
			return NULL
		},
	},
	"add": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("add関数は2つの引数が必要です: %d個与えられました", len(args))
			}
			
			// 文字列の場合は連結
			if args[0].Type() == object.STRING_OBJ {
				str, ok := args[0].(*object.String)
				if !ok {
					return newError("文字列の変換に失敗しました")
				}
				
				// 第2引数を文字列に変換
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
			
			// 整数加算
			left, ok := args[0].(*object.Integer)
			if !ok {
				return newError("add関数の第1引数は整数または文字列である必要があります: %s", args[0].Type())
			}
			
			right, ok := args[1].(*object.Integer)
			if !ok {
				return newError("add関数の第2引数は整数である必要があります: %s", args[1].Type())
			}
			
			return &object.Integer{Value: left.Value + right.Value}
		},
	},
	"sub": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("sub関数は2つの引数が必要です: %d個与えられました", len(args))
			}
			
			// 整数減算
			left, ok := args[0].(*object.Integer)
			if !ok {
				return newError("sub関数の第1引数は整数である必要があります: %s", args[0].Type())
			}
			
			right, ok := args[1].(*object.Integer)
			if !ok {
				return newError("sub関数の第2引数は整数である必要があります: %s", args[1].Type())
			}
			
			return &object.Integer{Value: left.Value - right.Value}
		},
	},
	"mul": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("mul関数は2つの引数が必要です: %d個与えられました", len(args))
			}
			
			// 整数乗算
			left, ok := args[0].(*object.Integer)
			if !ok {
				return newError("mul関数の第1引数は整数である必要があります: %s", args[0].Type())
			}
			
			right, ok := args[1].(*object.Integer)
			if !ok {
				return newError("mul関数の第2引数は整数である必要があります: %s", args[1].Type())
			}
			
			return &object.Integer{Value: left.Value * right.Value}
		},
	},
	"div": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("div関数は2つの引数が必要です: %d個与えられました", len(args))
			}
			
			// 整数除算
			left, ok := args[0].(*object.Integer)
			if !ok {
				return newError("div関数の第1引数は整数である必要があります: %s", args[0].Type())
			}
			
			right, ok := args[1].(*object.Integer)
			if !ok {
				return newError("div関数の第2引数は整数である必要があります: %s", args[1].Type())
			}
			
			// ゼロ除算チェック
			if right.Value == 0 {
				return newError("ゼロによる除算: %d / 0", left.Value)
			}
			
			return &object.Integer{Value: left.Value / right.Value}
		},
	},
	"mod": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("mod関数は2つの引数が必要です: %d個与えられました", len(args))
			}
			
			// 整数剰余
			left, ok := args[0].(*object.Integer)
			if !ok {
				return newError("mod関数の第1引数は整数である必要があります: %s", args[0].Type())
			}
			
			right, ok := args[1].(*object.Integer)
			if !ok {
				return newError("mod関数の第2引数は整数である必要があります: %s", args[1].Type())
			}
			
			// ゼロ除算チェック
			if right.Value == 0 {
				return newError("ゼロによるモジュロ: %d %% 0", left.Value)
			}
			
			return &object.Integer{Value: left.Value % right.Value}
		},
	},
	"pow": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("pow関数は2つの引数が必要です: %d個与えられました", len(args))
			}
			
			// べき乗
			base, ok := args[0].(*object.Integer)
			if !ok {
				return newError("pow関数の第1引数は整数である必要があります: %s", args[0].Type())
			}
			
			exp, ok := args[1].(*object.Integer)
			if !ok {
				return newError("pow関数の第2引数は整数である必要があります: %s", args[1].Type())
			}
			
			// 負の指数のチェック
			if exp.Value < 0 {
				return newError("pow関数の指数は0以上である必要があります: %d", exp.Value)
			}
			
			result := int64(1)
			for i := int64(0); i < exp.Value; i++ {
				result *= base.Value
			}
			
			return &object.Integer{Value: result}
		},
	},
	"to_string": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("to_string関数は1つの引数が必要です: %d個与えられました", len(args))
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
	},
	"length": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("length関数は1つの引数が必要です: %d個与えられました", len(args))
			}
			
			switch arg := args[0].(type) {
			case *object.String:
				return &object.Integer{Value: int64(len(arg.Value))}
			case *object.Array:
				return &object.Integer{Value: int64(len(arg.Elements))}
			default:
				return newError("length関数は文字列または配列に対してのみ使用できます: %s", args[0].Type())
			}
		},
	},
	"eq": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("eq関数は2つの引数が必要です: %d個与えられました", len(args))
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
	},
	"not": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("not関数は1つの引数が必要です: %d個与えられました", len(args))
			}
			
			if b, ok := args[0].(*object.Boolean); ok {
				return &object.Boolean{Value: !b.Value}
			}
			
			return &object.Boolean{Value: false} // デフォルトはfalse
		},
	},
	"split": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("split関数は2つの引数が必要です: %d個与えられました", len(args))
			}
			
			// 第1引数は対象文字列
			if args[0].Type() != object.STRING_OBJ {
				return newError("split関数の第1引数は文字列である必要があります: %s", args[0].Type())
			}
			str, _ := args[0].(*object.String)
			
			// 第2引数は区切り文字
			if args[1].Type() != object.STRING_OBJ {
				return newError("split関数の第2引数は文字列である必要があります: %s", args[1].Type())
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
	},
	"join": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("join関数は2つの引数が必要です: %d個与えられました", len(args))
			}
			
			// 第1引数は配列
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("join関数の第1引数は配列である必要があります: %s", args[0].Type())
			}
			array, _ := args[0].(*object.Array)
			
			// 第2引数は区切り文字
			if args[1].Type() != object.STRING_OBJ {
				return newError("join関数の第2引数は文字列である必要があります: %s", args[1].Type())
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
	},
	"substring": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			// 引数の数をチェック
			if len(args) < 2 || len(args) > 3 {
				return newError("substring関数は2-3個の引数が必要です: %d個与えられました", len(args))
			}
			
			// 第1引数は文字列
			if args[0].Type() != object.STRING_OBJ {
				return newError("substring関数の第1引数は文字列である必要があります: %s", args[0].Type())
			}
			str, _ := args[0].(*object.String)
			
			// 第2引数は開始位置
			if args[1].Type() != object.INTEGER_OBJ {
				return newError("substring関数の第2引数は整数である必要があります: %s", args[1].Type())
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
					return newError("substring関数の第3引数は整数である必要があります: %s", args[2].Type())
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
	},
	"to_upper": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("to_upper関数は1つの引数が必要です: %d個与えられました", len(args))
			}
			
			if args[0].Type() != object.STRING_OBJ {
				return newError("to_upper関数の引数は文字列である必要があります: %s", args[0].Type())
			}
			str, _ := args[0].(*object.String)
			
			return &object.String{Value: strings.ToUpper(str.Value)}
		},
	},
	"to_lower": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 1 {
				return newError("to_lower関数は1つの引数が必要です: %d個与えられました", len(args))
			}
			
			if args[0].Type() != object.STRING_OBJ {
				return newError("to_lower関数の引数は文字列である必要があります: %s", args[0].Type())
			}
			str, _ := args[0].(*object.String)
			
			return &object.String{Value: strings.ToLower(str.Value)}
		},
	},
}
