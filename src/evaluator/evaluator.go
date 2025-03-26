package evaluator

import (
	"fmt"

	"github.com/uncode/ast"
	"github.com/uncode/logger"
	"github.com/uncode/object"
)

// This file uses variables and functions defined in common.go

// Eval は抽象構文木を評価する
func Eval(node interface{}, env *object.Environment) object.Object {
	logger.Debug("評価中のノード: %T", node)

	switch node := node.(type) {
	case *ast.Program:
		logger.Debug("プログラムノードを評価")
		return evalProgram(node, env)

	case *ast.ExpressionStatement:
		logger.Debug("式文ノードを評価")
		return Eval(node.Expression, env)

	case *ast.StringLiteral:
		logger.Debug("文字列リテラルを評価")
		return &object.String{Value: node.Value}

	case *ast.IntegerLiteral:
		logger.Debug("整数リテラルを評価")
		return &object.Integer{Value: node.Value}

	case *ast.BooleanLiteral:
		logger.Debug("真偽値リテラルを評価")
		if node.Value {
			return TRUE
		}
		return FALSE

	case *ast.PizzaLiteral:
		logger.Debug("ピザリテラルを評価")
		// 🍕はパイプラインで渡された値を参照する特別な変数
		if val, ok := env.Get("🍕"); ok {
			return val
		}
		return createError("🍕が定義されていません（関数の外部またはパイプラインを通じて呼び出されていません）")

	case *ast.PooLiteral:
		logger.Debug("💩リテラルを評価")
		logger.Debug("💩リテラルを検出: 空の戻り値オブジェクトを生成します")

		// Return空のReturnValueオブジェクト
		// 実際の値はpipiline_eval.goのevalAssignment()内で設定される
		return &object.ReturnValue{}

	case *ast.PrefixExpression:
		logger.Debug("前置式を評価: %s", node.Operator)
		right := Eval(node.Right, env)
		if right.Type() == object.ERROR_OBJ {
			return right
		}
		return evalPrefixExpression(node.Operator, right)

	case *ast.FunctionLiteral:
		logger.Debug("関数リテラルを評価")
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
			logger.Debug("関数名 %s を環境に登録します", node.Name.Value)

			// 条件付き関数の場合、特別な名前で登録（上書きを防ぐため）
			if node.Condition != nil {
				// 既存の同名関数の数をカウント
				existingFuncs := env.GetAllFunctionsByName(node.Name.Value)
				uniqueName := fmt.Sprintf("%s#%d", node.Name.Value, len(existingFuncs))

				logger.Debug("条件付き関数 '%s' を '%s' として登録します", node.Name.Value, uniqueName)

				// 特別な名前で登録
				env.Set(uniqueName, function)

				// 検索用に元の名前も関連付け
				env.Set(node.Name.Value, function)
			} else {
				// 条件なし関数は通常通り登録
				env.Set(node.Name.Value, function)
			}
		}

		return function

	case *ast.InfixExpression:
		logger.Debug("中置式を評価")
		// パイプライン演算子と代入演算子のチェック
		if node.Operator == "|>" {
			return evalPipeline(node, env)
		} else if node.Operator == ">>" || node.Operator == "=" {
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
		logger.Debug("関数呼び出し式を評価")
		logger.Trace("関数: %T, 引数の数: %d", node.Function, len(node.Arguments))

		// 関数呼び出しが直接識別子（関数名）の場合、条件付き関数を検索
		if ident, ok := node.Function.(*ast.Identifier); ok {
			// 識別子名で関数を検索
			logger.Debug("識別子 '%s' で関数を検索します", ident.Value)

			// 引数を評価
			args := evalExpressions(node.Arguments, env)
			if len(args) > 0 && args[0].Type() == object.ERROR_OBJ {
				return args[0]
			}

			// デバッグ出力
			logger.Debug("関数 '%s' の引数: %d 個", ident.Value, len(args))
			for i, arg := range args {
				logger.Trace("  引数 %d: %s", i, arg.Inspect())
			}

			// 環境内の同名のすべての関数を検索し、条件に合う関数を適用
			return applyNamedFunction(env, ident.Value, args)
		}

		// 識別子以外（関数リテラルや式の結果など）の場合は従来通り処理
		function := Eval(node.Function, env)
		if function.Type() == object.ERROR_OBJ {
			return function
		}
		args := evalExpressions(node.Arguments, env)
		if len(args) > 0 && args[0].Type() == object.ERROR_OBJ {
			return args[0]
		}

		// 通常の関数呼び出しでは第一引数を🍕として設定しない
		if fn, ok := function.(*object.Function); ok {
			// 引数の数をチェック
			if len(args) != len(fn.Parameters) {
				return createError("引数の数が一致しません: 期待=%d, 実際=%d", len(fn.Parameters), len(args))
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
				return createError("関数の本体がBlockStatementではありません")
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

		return createError("関数ではありません: %s", function.Type())

	case *ast.Identifier:
		logger.Debug("識別子を評価")
		return evalIdentifier(node, env)

	case *ast.AssignStatement:
		logger.Debug("代入文を評価")

		// 右辺を評価
		right := Eval(node.Value, env)
		if right.Type() == object.ERROR_OBJ {
			return right
		}

		// 左辺が識別子の場合は変数に代入
		if ident, ok := node.Left.(*ast.Identifier); ok {
			logger.Debug("変数 %s に代入します", ident.Value)
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
				logger.Debug("💩への代入を検出しました - 戻り値として扱います")
				return &object.ReturnValue{Value: left}
			}
		}

		return right

	// その他のケース
	default:
		logger.Warn("未実装のノードタイプ: %T", node)
		return NullObj
	}
}

// isTruthy は値が真かどうかを判定する
func isTruthy(obj object.Object) bool {
	switch obj {
	case NullObj:
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
