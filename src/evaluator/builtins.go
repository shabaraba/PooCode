package evaluator

import (
	"fmt"
	"strings"

	"github.com/uncode/object"
)

// 組み込み関数のマップ
var Builtins = map[string]*object.Builtin{
	"print": &object.Builtin{
		Name: "print",
		Fn: func(args ...object.Object) object.Object {
			// デバッグ情報：受け取った引数の詳細を出力
			fmt.Printf("DEBUG: print関数が受け取った引数: %d個\n", len(args))
			for i, arg := range args {
				fmt.Printf("DEBUG: 引数%d - タイプ: %s, 値: %s\n", i, arg.Type(), arg.Inspect())
				
				// arg.Inspect()ではなく実際の値を表示
				switch arg.Type() {
				case object.INTEGER_OBJ:
					intVal := arg.(*object.Integer).Value
					fmt.Printf("DEBUG: 整数値として %d を出力\n", intVal)
					fmt.Println(intVal)
				case object.STRING_OBJ:
					strVal := arg.(*object.String).Value
					fmt.Printf("DEBUG: 文字列として \"%s\" を出力\n", strVal)
					fmt.Println(strVal)
				case object.BOOLEAN_OBJ:
					boolVal := arg.(*object.Boolean).Value
					fmt.Printf("DEBUG: 真偽値として %t を出力\n", boolVal)
					fmt.Println(boolVal)
				default:
					inspectVal := arg.Inspect()
					fmt.Printf("DEBUG: デフォルト - %s を出力\n", inspectVal)
					fmt.Println(inspectVal)
				}
			}
			// 第一引数を返すように変更（パイプラインの連鎖を維持するため）
			if len(args) > 0 {
				return args[0]
			}
			return NULL
		},
	},
	"show": &object.Builtin{
		Name: "show",
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				// arg.Inspect()ではなく実際の値を表示
				switch arg.Type() {
				case object.INTEGER_OBJ:
					fmt.Println(arg.(*object.Integer).Value)
				case object.STRING_OBJ:
					fmt.Println(arg.(*object.String).Value)
				case object.BOOLEAN_OBJ:
					fmt.Println(arg.(*object.Boolean).Value)
				default:
					fmt.Println(arg.Inspect())
				}
			}
			return NULL
		},
	},
	"add": &object.Builtin{
		Name: "add",
		Fn: func(args ...object.Object) object.Object {
			fmt.Printf("add関数が呼び出されました: 引数=%d個\n", len(args))
			// デバッグ: すべての引数を出力
			for i, arg := range args {
				fmt.Printf("  引数 %d: %s (型: %s)\n", i, arg.Inspect(), arg.Type())
			}
			
			if len(args) < 1 {
				return newError("add関数は少なくとも1つの引数が必要です")
			}
			
			// 文字列の場合は連結
			if args[0].Type() == object.STRING_OBJ {
				str, ok := args[0].(*object.String)
				if !ok {
					return newError("文字列の変換に失敗しました")
				}
				
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
				return newError("add関数の第1引数は整数または文字列である必要があります: %s", args[0].Type())
			}
			
			// 第2引数がない場合は値をそのまま返す
			if len(args) == 1 {
				fmt.Printf("add関数: 単一引数 %d をそのまま返します\n", left.Value)
				return left
			}
			
			// 第2引数があれば加算
			right, ok := args[1].(*object.Integer)
			if !ok {
				return newError("add関数の第2引数は整数である必要があります: %s", args[1].Type())
			}
			
			result := &object.Integer{Value: left.Value + right.Value}
			fmt.Printf("add関数: %d + %d = %d\n", left.Value, right.Value, result.Value)
			return result
		},
	},
	"sub": &object.Builtin{
		Name: "sub",
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
		Name: "mul",
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
		Name: "div",
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
		Name: "mod",
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
		Name: "pow",
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
		Name: "to_string",
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
		Name: "length",
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
		Name: "eq",
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
		Name: "not",
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
		Name: "split",
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
		Name: "join",
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
		Name: "substring",
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
		Name: "to_upper",
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
		Name: "to_lower",
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
	"range": &object.Builtin{
		Name: "range",
		Fn: func(args ...object.Object) object.Object {
			// 引数の数をチェック: 1または2つの引数を受け付ける
			if len(args) < 1 || len(args) > 2 {
				return newError("range関数は1-2個の引数が必要です: %d個与えられました", len(args))
			}
			
			var start, end int64
			
			// 1つの引数の場合: 0からその値まで
			if len(args) == 1 {
				if args[0].Type() != object.INTEGER_OBJ {
					return newError("range関数の引数は整数である必要があります: %s", args[0].Type())
				}
				endVal, _ := args[0].(*object.Integer)
				
				start = 0
				end = endVal.Value
			} else {
				// 2つの引数の場合: startからendまで
				if args[0].Type() != object.INTEGER_OBJ || args[1].Type() != object.INTEGER_OBJ {
					return newError("range関数の引数は整数である必要があります")
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
	},
	"sum": &object.Builtin{
		Name: "sum",
		Fn: func(args ...object.Object) object.Object {
			// 引数の数をチェック
			if len(args) != 1 {
				return newError("sum関数は1つの引数が必要です: %d個与えられました", len(args))
			}
			
			// 配列かどうかチェック
			if args[0].Type() != object.ARRAY_OBJ {
				return newError("sum関数の引数は配列である必要があります: %s", args[0].Type())
			}
			
			array, _ := args[0].(*object.Array)
			
			// 合計を計算
			sum := int64(0)
			for _, elem := range array.Elements {
				if elem.Type() != object.INTEGER_OBJ {
					return newError("sum関数の配列要素はすべて整数である必要があります: %s", elem.Type())
				}
				
				intVal, _ := elem.(*object.Integer)
				sum += intVal.Value
			}
			
			return &object.Integer{Value: sum}
		},
	},
}
