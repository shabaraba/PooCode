package evaluator

import (
	"fmt"
	
	"github.com/uncode/ast"
	"github.com/uncode/object"
)

// evalPipeline は|>演算子のパイプライン処理を評価する
func evalPipeline(node *ast.InfixExpression, env *object.Environment) object.Object {
	if debugMode {
		fmt.Println("パイプライン演算子を検出しました")
	}
	// |>演算子の場合、左辺の結果を右辺の関数に渡す
	left := Eval(node.Left, env)
	if left.Type() == object.ERROR_OBJ {
		return left
	}
	
	// 右辺が識別子の場合、関数として評価
	if ident, ok := node.Right.(*ast.Identifier); ok {
		if debugMode {
			fmt.Printf("識別子としてのパイプライン先: %s\n", ident.Value)
		}
		
		fmt.Printf("パイプラインから applyNamedFunction を呼び出します (関数名: %s)\n", ident.Value)
		
		// 環境変数 🍕 を設定して名前付き関数呼び出しへ処理を委譲
		// ここで左辺の値を唯一の引数として渡す
		args := []object.Object{left}
		
		// 名前付き関数を適用する（条件付き関数の処理も行う）
		// 戻り値を変数に格納して、何が返されるか確認する
		result := applyNamedFunction(env, ident.Value, args)
		fmt.Printf("パイプライン: 関数 '%s' の実行結果: タイプ=%s, 値=%s\n", 
			ident.Value, result.Type(), result.Inspect())
		return result
	}
	
	// 右辺が関数呼び出しの場合
	if callExpr, ok := node.Right.(*ast.CallExpression); ok {
		if debugMode {
			fmt.Println("関数呼び出しとしてのパイプライン先")
		}
		
		// 関数名を識別子から直接取得（デバッグ用）
		if ident, ok := callExpr.Function.(*ast.Identifier); ok {
			fmt.Printf("パイプラインでビルトイン関数 '%s' を呼び出します\n", ident.Value)
			
			// 引数を評価
			args := evalExpressions(callExpr.Arguments, env)
			
			// ビルトイン関数を直接取得して呼び出す
			if builtin, ok := Builtins[ident.Value]; ok {
				// leftを第一引数、その他の引数は後続
				allArgs := []object.Object{left}
				allArgs = append(allArgs, args...)
				
				fmt.Printf("ビルトイン関数 '%s' を実行: 全引数 %d 個\n", ident.Value, len(allArgs))
				for i, arg := range allArgs {
					fmt.Printf("  引数 %d: %s\n", i, arg.Inspect())
				}
				
				result := builtin.Fn(allArgs...)
				fmt.Printf("ビルトイン関数 '%s' の結果: %s\n", ident.Value, result.Inspect())
				return result
			}
			
			// ビルトイン関数でない場合は名前付き関数として呼び出し
			// leftを第一引数、その他の引数は後続
			allArgs := []object.Object{left}
			allArgs = append(allArgs, args...)
			
			// 名前付き関数を適用する（条件付き関数の処理も行う）
			result := applyNamedFunction(env, ident.Value, allArgs)
			fmt.Printf("パイプライン(callExpr): 関数 '%s' の実行結果: タイプ=%s, 値=%s\n", 
				ident.Value, result.Type(), result.Inspect())
			return result
		}
		
		// 識別子以外の関数式を評価
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
				result = obj.Value
				fmt.Printf("関数結果(ReturnValue): %s\n", result.Inspect())
				return result
			}
			fmt.Printf("関数結果(直接): %s\n", result.Inspect())
			return result
		} else if builtin, ok := function.(*object.Builtin); ok {
			// 組み込み関数の場合、leftを第一引数として追加
			fmt.Printf("ビルトイン関数 '%s' を実行\n", callExpr.Function.(*ast.Identifier).Value)
			args = append([]object.Object{left}, args...)
			
			fmt.Printf("引数: %d個\n", len(args))
			for i, arg := range args {
				fmt.Printf("  引数 %d: %s\n", i, arg.Inspect())
			}
			
			result := builtin.Fn(args...)
			fmt.Printf("ビルトイン関数の結果: %s\n", result.Inspect())
			return result
		}
		
		return newError("関数ではありません: %s", function.Type())
	}
	
	// その他の場合は処理できない
	return newError("パイプラインの右側が関数または識別子ではありません: %T", node.Right)
}

// evalAssignment は>>演算子による代入を評価する
func evalAssignment(node *ast.InfixExpression, env *object.Environment) object.Object {
	if debugMode {
		fmt.Println("代入演算子を検出しました")
	}
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
		if debugMode {
			fmt.Println("💩への代入を検出しました - 戻り値として扱います")
		}
		left := Eval(node.Left, env)
		if left.Type() == object.ERROR_OBJ {
			return left
		}
		fmt.Printf("💩に戻り値として %s を設定します\n", left.Inspect())
		return &object.ReturnValue{Value: left}
	}
	
	return newError("代入先が識別子または💩ではありません: %T", right)
}
