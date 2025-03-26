package evaluator

import (
	"fmt"
	"github.com/uncode/logger"
	"github.com/uncode/object"
)

// registerMathBuiltins は数学関連の組み込み関数を登録する
func registerMathBuiltins() {
	Builtins["add"] = &object.Builtin{
		Name: "add",
		Fn: func(args ...object.Object) object.Object {
			if len(args) == 0 {
				return createError("add関数は少なくとも1つの引数が必要です")
			}
			
			// 文字列加算の場合
			if str, ok := args[0].(*object.String); ok {
				logIfEnabled(logger.LevelDebug, "add関数: 文字列連結モード")
				
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
				logIfEnabled(logger.LevelDebug, "add関数: 単一引数 %d をそのまま返します", left.Value)
				return left
			}
			
			// 第2引数があれば加算
			right, ok := args[1].(*object.Integer)
			if !ok {
				return createError("add関数の第2引数は整数である必要があります: %s", args[1].Type())
			}
			
			result := &object.Integer{Value: left.Value + right.Value}
			logIfEnabled(logger.LevelDebug, "add関数: %d + %d = %d", left.Value, right.Value, result.Value)
			return result
		},
		ReturnType: object.ANY_OBJ, // 文字列または整数を返す可能性あり
		ParamTypes: []object.ObjectType{object.ANY_OBJ},
	}

	Builtins["sub"] = &object.Builtin{
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
	}

	Builtins["mul"] = &object.Builtin{
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
	}

	Builtins["div"] = &object.Builtin{
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
	}

	Builtins["mod"] = &object.Builtin{
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
	}

	Builtins["pow"] = &object.Builtin{
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
	}

	// 整数配列の合計を計算する関数
	Builtins["sum"] = &object.Builtin{
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
	}
}
