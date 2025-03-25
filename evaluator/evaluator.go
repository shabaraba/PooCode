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
				
				return applyFunction(function, []object.Object{left})
			}
			
			// 右辺が関数呼び出しの場合
			if callExpr, ok := node.Right.(*ast.CallExpression); ok {
				fmt.Println("関数呼び出しとしてのパイプライン先")
				function := Eval(callExpr.Function, env)
				if function.Type() == object.ERROR_OBJ {
					return function
				}
				
				args := evalExpressions(callExpr.Arguments, env)
				// 左辺の結果を第1引数に挿入
				args = append([]object.Object{left}, args...)
				
				return applyFunction(function, args)
			}
			
			return newError("パイプラインの右側が関数または識別子ではありません: %T", node.Right)
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
		args := evalExpressions(node.Arguments, env)
		return applyFunction(function, args)
		
	case *ast.Identifier:
		fmt.Println("識別子を評価")
		return evalIdentifier(node, env)
		
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
	"add": &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return newError("add関数は2つの引数が必要です: %d個与えられました", len(args))
			}
			
			// ここでは単純な整数加算のみ実装
			left, ok := args[0].(*object.Integer)
			if !ok {
				return newError("add関数の第1引数は整数である必要があります: %s", args[0].Type())
			}
			
			right, ok := args[1].(*object.Integer)
			if !ok {
				return newError("add関数の第2引数は整数である必要があります: %s", args[1].Type())
			}
			
			return &object.Integer{Value: left.Value + right.Value}
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
		return &object.Integer{Value: leftVal / rightVal}
	case "%":
		return &object.Integer{Value: leftVal % rightVal}
	case "==":
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
	case "|":
		// 並列パイプの場合、最初の値を返す（実際の実装は環境に応じて）
		return left
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
	case "==":
		return &object.Boolean{Value: leftVal == rightVal}
	case "!=":
		return &object.Boolean{Value: leftVal != rightVal}
	default:
		return newError("未知の演算子: %s %s %s", left.Type(), operator, right.Type())
	}
}

// evalBooleanInfixExpression は真偽値の中置式を評価する
func evalBooleanInfixExpression(operator string, left, right object.Object) object.Object {
	leftVal := left.(*object.Boolean).Value
	rightVal := right.(*object.Boolean).Value
	
	switch operator {
	case "==":
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
