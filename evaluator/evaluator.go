package evaluator

import (
	"fmt"
	"strings"

	"github.com/uncode/ast"
	"github.com/uncode/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

// Eval は抽象構文木を評価する
func Eval(node interface{}, env *object.Environment) object.Object {
	fmt.Printf("評価中のノード: %T\n", node)
	
	switch node := node.(type) {
	case *ast.Program:
		fmt.Println("プログラムノードを評価")
		return evalProgram(node, env)
		
	case *ast.ExpressionStatement:
		fmt.Println("式文ノードを評価")
		return Eval(node.Expression, env)
		
	case *ast.StringLiteral:
		fmt.Println("文字列リテラルを評価")
		return &object.String{Value: node.Value}
		
	case *ast.IntegerLiteral:
		fmt.Println("整数リテラルを評価")
		return &object.Integer{Value: node.Value}
		
	case *ast.BooleanLiteral:
		fmt.Println("真偽値リテラルを評価")
		if node.Value {
			return TRUE
		}
		return FALSE
		
	case *ast.PizzaLiteral:
		fmt.Println("ピザリテラルを評価")
		// 🍕はパイプラインで渡された値を参照する特別な変数
		if val, ok := env.Get("🍕"); ok {
			return val
		}
		return newError("🍕が定義されていません（関数の外部またはパイプラインを通じて呼び出されていません）")
		
	case *ast.PooLiteral:
		fmt.Println("💩リテラルを評価")
		// 💩は関数の戻り値として扱う特別なリテラル
		return &object.ReturnValue{}
		
	case *ast.PrefixExpression:
		fmt.Println("前置式を評価:", node.Operator)
		right := Eval(node.Right, env)
		if right.Type() == object.ERROR_OBJ {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
		
	case *ast.FunctionLiteral:
		fmt.Println("関数リテラルを評価")
		// ast.Identifierをobject.Identifierに変換
		params := make([]*object.Identifier, len(node.Parameters))
		for i, p := range node.Parameters {
			params[i] = &object.Identifier{Value: p.Value}
		}
		
		// ast.BlockStatementをオブジェクトとして保存
		function := &object.Function{
			Parameters: params,
			ASTBody:    node.Body,
			Env:        env,
			InputType:  node.InputType,
			ReturnType: node.ReturnType,
			Condition:  node.Condition,
		}
		
		// 関数に名前がある場合は環境に登録
		if node.Name != nil {
			fmt.Printf("関数名 %s を環境に登録します\n", node.Name.Value)
			env.Set(node.Name.Value, function)
		}
		
		return function
		
	case *ast.InfixExpression:
		fmt.Println("中置式を評価")
		// パイプライン演算子のチェック
		if node.Operator == "|>" {
			fmt.Println("パイプライン演算子を検出しました")
			// |>演算子の場合、左辺の結果を右辺の関数に渡す
			left := Eval(node.Left, env)
			if left.Type() == object.ERROR_OBJ {
				return left
			}
			
			// 右辺が識別子の場合、関数として評価
			if ident, ok := node.Right.(*ast.Identifier); ok {
				fmt.Printf("識別子としてのパイプライン先: %s\n", ident.Value)
				function := evalIdentifier(ident, env)
				if function.Type() == object.ERROR_OBJ {
					return function
				}
				
				// 専用の環境変数 🍕 に値を設定して関数を呼び出す
				if fn, ok := function.(*object.Function); ok {
					extendedEnv := object.NewEnclosedEnvironment(fn.Env)
					extendedEnv.Set("🍕", left)
					
					// ASTBodyをast.BlockStatementに型アサーション
					astBody, ok := fn.ASTBody.(*ast.BlockStatement)
					if !ok {
						return newError("関数の本体がBlockStatementではありません")
					}
					result := evalBlockStatement(astBody, extendedEnv)
					
					// 💩値を返す（関数の戻り値）
					if obj, ok := result.(*object.ReturnValue); ok {
						return obj.Value
					}
					return result
				} else if builtin, ok := function.(*object.Builtin); ok {
					// 組み込み関数の場合はそのまま引数として渡す
					return builtin.Fn(left)
				}
				
				return newError("関数ではありません: %s", function.Type())
			}
			
			// 右辺が関数呼び出しの場合
			if callExpr, ok := node.Right.(*ast.CallExpression); ok {
				fmt.Println("関数呼び出しとしてのパイプライン先")
				function := Eval(callExpr.Function, env)
				if function.Type() == object.ERROR_OBJ {
					return function
				}
				
				args := evalExpressions(callExpr.Arguments, env)
				
				// 関数オブジェクトの場合、専用の環境変数🍕に左辺の値を設定
				if fn, ok := function.(*object.Function); ok {
					extendedEnv := object.NewEnclosedEnvironment(fn.Env)
					
					// 通常の引数を環境にバインド
					if len(args) != len(fn.Parameters) {
						return newError("引数の数が一致しません: 期待=%d, 実際=%d", len(fn.Parameters), len(args))
					}
					
					for i, param := range fn.Parameters {
						extendedEnv.Set(param.Value, args[i])
					}
					
					// パイプラインからの値を🍕にセット
					extendedEnv.Set("🍕", left)
					
					// 関数本体を評価
					astBody, ok := fn.ASTBody.(*ast.BlockStatement)
					if !ok {
						return newError("関数の本体がBlockStatementではありません")
					}
					result := evalBlockStatement(astBody, extendedEnv)
					
					// 💩値を返す（関数の戻り値）
					if obj, ok := result.(*object.ReturnValue); ok {
						return obj.Value
					}
					return result
				} else if builtin, ok := function.(*object.Builtin); ok {
					// 組み込み関数の場合、leftを第一引数として追加
					args = append([]object.Object{left}, args...)
					return builtin.Fn(args...)
				}
				
				return newError("関数ではありません: %s", function.Type())
			}
			
			return newError("パイプラインの右側が関数または識別子ではありません: %T", node.Right)
		} else if node.Operator == ">>" {
			fmt.Println("代入演算子を検出しました")
			// >>演算子の場合、右辺の変数に左辺の値を代入する
			right := node.Right
			
			// 右辺が識別子の場合は変数に代入
			if ident, ok := right.(*ast.Identifier); ok {
				left := Eval(node.Left, env)
				if left.Type() == object.ERROR_OBJ {
					return left
				}
				
				env.Set(ident.Value, left)
				return left
			}
			
			// 右辺がPooLiteralの場合は戻り値として扱う
			if _, ok := right.(*ast.PooLiteral); ok {
				fmt.Println("💩への代入を検出しました - 戻り値として扱います")
				left := Eval(node.Left, env)
				if left.Type() == object.ERROR_OBJ {
					return left
				}
				return &object.ReturnValue{Value: left}
			}
			
			return newError("代入先が識別子または💩ではありません: %T", right)
		} else {
			// その他の中置演算子
			left := Eval(node.Left, env)
			if left.Type() == object.ERROR_OBJ {
				return left
			}
			
			right := Eval(node.Right, env)
			if right.Type() == object.ERROR_OBJ {
				return right
			}
			
			return evalInfixExpression(node.Operator, left, right)
		}
		
	case *ast.CallExpression:
		fmt.Println("関数呼び出し式を評価")
		function := Eval(node.Function, env)
		if function.Type() == object.ERROR_OBJ {
			return function
		}
		args := evalExpressions(node.Arguments, env)
		
		// 通常の関数呼び出しでは第一引数を🍕として設定しない
		if fn, ok := function.(*object.Function); ok {
			// 引数の数をチェック
			if len(args) != len(fn.Parameters) {
				return newError("引数の数が一致しません: 期待=%d, 実際=%d", len(fn.Parameters), len(args))
			}
			
			// 新しい環境を作成
			extendedEnv := object.NewEnclosedEnvironment(fn.Env)
			
			// 引数を環境にバインド
			for i, param := range fn.Parameters {
				extendedEnv.Set(param.Value, args[i])
			}
			
			// 通常の関数呼び出しでは、🍕を設定しない
			// （修正後の仕様では、🍕はパイプラインで渡された値のみを表す）
			
			// 関数本体を評価
			astBody, ok := fn.ASTBody.(*ast.BlockStatement)
			if !ok {
				return newError("関数の本体がBlockStatementではありません")
			}
			result := evalBlockStatement(astBody, extendedEnv)
			
			// 💩値を返す（関数の戻り値）
			if obj, ok := result.(*object.ReturnValue); ok {
				return obj.Value
			}
			return result
		} else if builtin, ok := function.(*object.Builtin); ok {
			return builtin.Fn(args...)
		}
		
		return newError("関数ではありません: %s", function.Type())
		
	case *ast.Identifier:
		fmt.Println("識別子を評価")
		return evalIdentifier(node, env)
		
	case *ast.AssignStatement:
		fmt.Println("代入文を評価")
		
		// 右辺を評価
		right := Eval(node.Value, env)
		if right.Type() == object.ERROR_OBJ {
			return right
		}
		
		// 左辺が識別子の場合は変数に代入
		if ident, ok := node.Left.(*ast.Identifier); ok {
			fmt.Printf("変数 %s に代入します\n", ident.Value)
			env.Set(ident.Value, right)
			return right
		} else {
			// その他の場合は左辺を評価してから処理
			left := Eval(node.Left, env)
			if left.Type() == object.ERROR_OBJ {
				return left
			}
			
			// 💩リテラルへの代入は特殊な意味を持つ (関数からの戻り値)
			if _, ok := node.Value.(*ast.PooLiteral); ok {
				fmt.Println("💩への代入を検出しました (戻り値)")
				return &object.ReturnValue{Value: left}
			}
		}
		
		return right
		
	// その他のケース
	default:
		fmt.Printf("未実装のノードタイプ: %T\n", node)
		return NULL
	}
}

// evalProgram はプログラムを評価する
func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object = NULL

	for _, statement := range program.Statements {
		result = Eval(statement, env)
	}
	
	return result
}

// evalBlockStatement はブロック文を評価する
func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object = NULL

	for _, statement := range block.Statements {
		result = Eval(statement, env)
		
		// 特殊なケース: >>💩 は関数からの戻り値
		if returnValue, ok := result.(*object.ReturnValue); ok {
			return returnValue
		}
		
		// 代入文の場合、PooLiteralへの代入は特別な意味を持つ
		if assignStmt, ok := statement.(*ast.AssignStatement); ok {
			if _, ok := assignStmt.Value.(*ast.PooLiteral); ok {
				fmt.Println("💩への代入を検出しました - 戻り値として扱います")
				// 右辺の値を取得
				rightVal := Eval(assignStmt.Left, env)
				if rightVal.Type() == object.ERROR_OBJ {
					return rightVal
				}
				return &object.ReturnValue{Value: rightVal}
			}
		}
	}
	
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

// エラー生成用ヘルパー関数
func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

// evalExpressions は複数の式を評価する
func evalExpressions(exps []ast.Expression, env *object.Environment) []object.Object {
	var result []object.Object
	
	for _, e := range exps {
		evaluated := Eval(e, env)
		if evaluated != nil {
			result = append(result, evaluated)
		}
	}
	
	return result
}

// applyFunction は関数を適用する
func applyFunction(fn object.Object, args []object.Object) object.Object {
	switch fn := fn.(type) {
	case *object.Function:
		// 関数呼び出しの実装
		fmt.Println("関数を呼び出します:", fn.Inspect())
		
		// 引数の数をチェック
		if len(args) != len(fn.Parameters) {
			return newError("引数の数が一致しません: 期待=%d, 実際=%d", len(fn.Parameters), len(args))
		}
		
		// 新しい環境を作成
		extendedEnv := object.NewEnclosedEnvironment(fn.Env)
		
		// 引数を環境にバインド
		for i, param := range fn.Parameters {
			extendedEnv.Set(param.Value, args[i])
		}
		
		// 修正後の仕様では、🍕はパイプラインで渡された値のみを表す
		// 通常の関数呼び出しでは🍕は設定しない
		
		// 関数本体を評価（ASTBodyをast.BlockStatementに型アサーション）
		astBody, ok := fn.ASTBody.(*ast.BlockStatement)
		if !ok {
			return newError("関数の本体がBlockStatementではありません")
		}
		result := evalBlockStatement(astBody, extendedEnv)
		
		// 💩値を返す（関数の戻り値）
		if obj, ok := result.(*object.ReturnValue); ok {
			return obj.Value
		}
		return result
		
	case *object.Builtin:
		return fn.Fn(args...)
		
	default:
		return newError("関数ではありません: %s", fn.Type())
	}
}

// evalInfixExpression は中置式を評価する
func evalInfixExpression(operator string, left, right object.Object) object.Object {
	// 整数の演算
	if left.Type() == object.INTEGER_OBJ && right.Type() == object.INTEGER_OBJ {
		return evalIntegerInfixExpression(operator, left, right)
	}
	
	// 文字列の演算
	if left.Type() == object.STRING_OBJ && right.Type() == object.STRING_OBJ {
		return evalStringInfixExpression(operator, left, right)
	}
	
	// 真偽値の演算
	if left.Type() == object.BOOLEAN_OBJ && right.Type() == object.BOOLEAN_OBJ {
		return evalBooleanInfixExpression(operator, left, right)
	}
	
	// 型の不一致
	if left.Type() != right.Type() {
		return newError("型の不一致: %s %s %s", left.Type(), operator, right.Type())
	}
	
	return newError("未知の演算子: %s %s %s", left.Type(), operator, right.Type())
}

// evalIntegerInfixExpression は整数の中置式を評価する
func evalIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Integer).Value
	rightVal := right.(*object.Integer).Value
	
	switch operator {
	case "+":
		return &object.Integer{Value: leftVal + rightVal}
	case "-":
		return &object.Integer{Value: leftVal - rightVal}
	case "*":
		return &object.Integer{Value: leftVal * rightVal}
	case "/":
		// ゼロ除算チェック
		if rightVal == 0 {
			return newError("ゼロによる除算: %d / 0", leftVal)
		}
		return &object.Integer{Value: leftVal / rightVal}
	case "%":
		// ゼロ除算チェック
		if rightVal == 0 {
			return newError("ゼロによるモジュロ: %d %% 0", leftVal)
		}
		return &object.Integer{Value: leftVal % rightVal}
	case "**":
		// べき乗演算子
		result := int64(1)
		for i := int64(0); i < rightVal; i++ {
			result *= leftVal
		}
		return &object.Integer{Value: result}
	case "&":
		// ビット論理積
		return &object.Integer{Value: leftVal & rightVal}
	case "|":
		// ビット論理和（または並列パイプ）
		return &object.Integer{Value: leftVal | rightVal}
	case "^":
		// ビット排他的論理和
		return &object.Integer{Value: leftVal ^ rightVal}
	case "<<":
		// 左シフト
		return &object.Integer{Value: leftVal << uint64(rightVal)}
	case ">>":
		// 右シフト
		return &object.Integer{Value: leftVal >> uint64(rightVal)}
	case "==", "eq":
		return &object.Boolean{Value: leftVal == rightVal}
	case "!=":
		return &object.Boolean{Value: leftVal != rightVal}
	case "<":
		return &object.Boolean{Value: leftVal < rightVal}
	case ">":
		return &object.Boolean{Value: leftVal > rightVal}
	case "<=":
		return &object.Boolean{Value: leftVal <= rightVal}
	case ">=":
		return &object.Boolean{Value: leftVal >= rightVal}
	default:
		return newError("未知の演算子: %s %s %s", left.Type(), operator, right.Type())
	}
}

// evalStringInfixExpression は文字列の中置式を評価する
func evalStringInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.String).Value
	rightVal := right.(*object.String).Value
	
	switch operator {
	case "+":
		return &object.String{Value: leftVal + rightVal}
	case "==", "eq":
		return &object.Boolean{Value: leftVal == rightVal}
	case "!=":
		return &object.Boolean{Value: leftVal != rightVal}
	case "<":
		return &object.Boolean{Value: leftVal < rightVal}
	case ">":
		return &object.Boolean{Value: leftVal > rightVal}
	case "<=":
		return &object.Boolean{Value: leftVal <= rightVal}
	case ">=":
		return &object.Boolean{Value: leftVal >= rightVal}
	case "contains":
		return &object.Boolean{Value: strings.Contains(leftVal, rightVal)}
	case "starts_with":
		return &object.Boolean{Value: strings.HasPrefix(leftVal, rightVal)}
	case "ends_with":
		return &object.Boolean{Value: strings.HasSuffix(leftVal, rightVal)}
	default:
		return newError("未知の演算子: %s %s %s", left.Type(), operator, right.Type())
	}
}

// evalBooleanInfixExpression は真偽値の中置式を評価する
func evalBooleanInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Boolean).Value
	rightVal := right.(*object.Boolean).Value
	
	switch operator {
	case "==", "eq":
		return &object.Boolean{Value: leftVal == rightVal}
	case "!=":
		return &object.Boolean{Value: leftVal != rightVal}
	case "&&":
		return &object.Boolean{Value: leftVal && rightVal}
	case "||":
		return &object.Boolean{Value: leftVal || rightVal}
	case "|":
		// 並列パイプの場合、最初の真の値を返す
		if leftVal {
			return left
		}
		return right
	default:
		return newError("未知の演算子: %s %s %s", left.Type(), operator, right.Type())
	}
}

// evalPrefixExpression は前置式を評価する
func evalPrefixExpression(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	case "not":
		// 言語仕様で "not" は ! と同様に扱う
		return evalBangOperatorExpression(right)
	default:
		return newError("未知の前置演算子: %s%s", operator, right.Type())
	}
}

// evalBangOperatorExpression は ! 演算子を評価する
func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NULL:
		return TRUE
	default:
		// 真偽値以外の値に対しては false を返す
		if right.Type() == object.BOOLEAN_OBJ {
			if right.(*object.Boolean).Value {
				return FALSE
			}
			return TRUE
		}
		return FALSE
	}
}

// evalMinusPrefixOperatorExpression は - 演算子を評価する
func evalMinusPrefixOperatorExpression(right object.Object) object.Object {
	if right.Type() != object.INTEGER_OBJ {
		return newError("-演算子は整数に対してのみ使用できます: %s", right.Type())
	}
	
	value := right.(*object.Integer).Value
	return &object.Integer{Value: -value}
}

// evalIdentifier は識別子を評価する
func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	// 環境から変数を探す
	if val, ok := env.Get(node.Value); ok {
		return val
	}
	
	// 組み込み関数を探す
	if builtin, ok := builtins[node.Value]; ok {
		return builtin
	}
	
	return newError("識別子が見つかりません: " + node.Value)
}
