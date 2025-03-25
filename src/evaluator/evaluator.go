package evaluator

import (
	"fmt"

	"github.com/uncode/ast"
	"github.com/uncode/object"
)

var (
	NULL  = &object.Null{}
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

// デバッグフラグ
var debugMode = false

// SetDebugMode はデバッグモードを設定する
func SetDebugMode(mode bool) {
	debugMode = mode
}

// Eval は抽象構文木を評価する
func Eval(node interface{}, env *object.Environment) object.Object {
	if debugMode {
		fmt.Printf("評価中のノード: %T\n", node)
	}
	
	switch node := node.(type) {
	case *ast.Program:
		if debugMode {
			fmt.Println("プログラムノードを評価")
		}
		return evalProgram(node, env)
		
	case *ast.ExpressionStatement:
		if debugMode {
			fmt.Println("式文ノードを評価")
		}
		return Eval(node.Expression, env)
		
	case *ast.StringLiteral:
		if debugMode {
			fmt.Println("文字列リテラルを評価")
		}
		return &object.String{Value: node.Value}
		
	case *ast.IntegerLiteral:
		if debugMode {
			fmt.Println("整数リテラルを評価")
		}
		return &object.Integer{Value: node.Value}
		
	case *ast.BooleanLiteral:
		if debugMode {
			fmt.Println("真偽値リテラルを評価")
		}
		if node.Value {
			return TRUE
		}
		return FALSE
		
	case *ast.PizzaLiteral:
		if debugMode {
			fmt.Println("ピザリテラルを評価")
		}
		// 🍕はパイプラインで渡された値を参照する特別な変数
		if val, ok := env.Get("🍕"); ok {
			return val
		}
		return newError("🍕が定義されていません（関数の外部またはパイプラインを通じて呼び出されていません）")
		
	case *ast.PooLiteral:
		if debugMode {
			fmt.Println("💩リテラルを評価")
		}
		// 💩は関数の戻り値として扱う特別なリテラル
		return &object.ReturnValue{}
		
	case *ast.PrefixExpression:
		if debugMode {
			fmt.Println("前置式を評価:", node.Operator)
		}
		right := Eval(node.Right, env)
		if right.Type() == object.ERROR_OBJ {
			return right
		}
		return evalPrefixExpression(node.Operator, right)
		
	case *ast.FunctionLiteral:
		if debugMode {
			fmt.Println("関数リテラルを評価")
		}
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
			if debugMode {
				fmt.Printf("関数名 %s を環境に登録します\n", node.Name.Value)
			}
			env.Set(node.Name.Value, function)
		}
		
		return function
		
	case *ast.InfixExpression:
		if debugMode {
			fmt.Println("中置式を評価")
		}
		// パイプライン演算子のチェック
		if node.Operator == "|>" {
			return evalPipeline(node, env)
		} else if node.Operator == ">>" {
			return evalAssignment(node, env)
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
		if debugMode {
			fmt.Println("関数呼び出し式を評価")
		}
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
		if debugMode {
			fmt.Println("識別子を評価")
		}
		return evalIdentifier(node, env)
		
	case *ast.AssignStatement:
		if debugMode {
			fmt.Println("代入文を評価")
		}
		
		// 右辺を評価
		right := Eval(node.Value, env)
		if right.Type() == object.ERROR_OBJ {
			return right
		}
		
		// 左辺が識別子の場合は変数に代入
		if ident, ok := node.Left.(*ast.Identifier); ok {
			if debugMode {
				fmt.Printf("変数 %s に代入します\n", ident.Value)
			}
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
				if debugMode {
					fmt.Println("💩への代入を検出しました - 戻り値として扱います")
				}
				return &object.ReturnValue{Value: left}
			}
		}
		
		return right
		
	// その他のケース
	default:
		if debugMode {
			fmt.Printf("未実装のノードタイプ: %T\n", node)
		}
		return NULL
	}
}

// エラー生成用ヘルパー関数
func newError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}

// isTruthy は値が真かどうかを判定する
func isTruthy(obj object.Object) bool {
	switch obj {
	case NULL:
		return false
	case TRUE:
		return true
	case FALSE:
		return false
	default:
		// 数値の場合、0以外は真
		if integer, ok := obj.(*object.Integer); ok {
			return integer.Value != 0
		}
		// 文字列の場合、空文字列以外は真
		if str, ok := obj.(*object.String); ok {
			return str.Value != ""
		}
		// それ以外のオブジェクトは真
		return true
	}
}
