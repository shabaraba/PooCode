package evaluator

import (
	"fmt"
	
	"github.com/uncode/ast"
	"github.com/uncode/object"
)

// evalPipeline は|>演算子のパイプライン処理を評価する
func evalPipeline(node *ast.InfixExpression, env *object.Environment) object.Object {
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
}

// evalAssignment は>>演算子による代入を評価する
func evalAssignment(node *ast.InfixExpression, env *object.Environment) object.Object {
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
}
