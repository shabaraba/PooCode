package evaluator

import (
	"github.com/uncode/ast"
	"github.com/uncode/logger"
	"github.com/uncode/object"
)

// evalInfixExpressionWithNode は中置式をノード情報も含めて評価する
// 特に条件式の場合、🍕の参照方法を強化するために使用
func evalInfixExpressionWithNode(node *ast.InfixExpression, env *object.Environment) object.Object {
	// 特殊演算子の処理
	switch node.Operator {
	case "|>":
		// パイプライン演算子
		return evalPipeline(node, env)
	case ">>", ">>=":
		// リダイレクト演算子（代入/追加）
		return evalAssignment(node, env)
	case "=", ":=":
		// 代入演算子
		return evalAssignment(node, env)
	case "+>", "map": // map演算子
		if logger.IsLevelEnabled(mapFilterDebugLevel) {
			logger.Log(mapFilterDebugLevel, "map パイプ演算子 (%s) を検出しました", node.Operator)
		}
		// map関数の処理を実行
		return evalMapOperation(node, env)
	case "?>", "filter": // filter演算子
		if logger.IsLevelEnabled(mapFilterDebugLevel) {
			logger.Log(mapFilterDebugLevel, "filter パイプ演算子 (%s) を検出しました", node.Operator)
		}
		// filter関数の処理を実行
		return evalFilterOperation(node, env)
	}

	// ピザリテラルが含まれる場合のチェック
	// 左辺がピザリテラルの場合
	if _, ok := node.Left.(*ast.PizzaLiteral); ok {
		// currentFunctionの🍕から値を取得
		if currentFunction != nil {
			if pizzaVal := currentFunction.GetPizzaValue(); pizzaVal != nil {
				// 🍕の値を左辺として使用
				logger.Debug("中置式の左辺に🍕を使用: %s", pizzaVal.Inspect())
				left := pizzaVal

				// 右辺を評価
				right := Eval(node.Right, env)
				if right.Type() == object.ERROR_OBJ {
					return right
				}

				// 演算子を適用
				return evalInfixExpression(node.Operator, left, right)
			}
		}

		// それ以外は通常の🍕評価
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

	// 右辺がピザリテラルの場合
	if _, ok := node.Right.(*ast.PizzaLiteral); ok {
		// 左辺を評価
		left := Eval(node.Left, env)
		if left.Type() == object.ERROR_OBJ {
			return left
		}

		// currentFunctionの🍕から値を取得
		if currentFunction != nil {
			if pizzaVal := currentFunction.GetPizzaValue(); pizzaVal != nil {
				// 🍕の値を右辺として使用
				logger.Debug("中置式の右辺に🍕を使用: %s", pizzaVal.Inspect())
				right := pizzaVal

				// 演算子を適用
				return evalInfixExpression(node.Operator, left, right)
			}
		}

		// 環境から🍕を取得（バックアップ）
		if right, ok := env.Get("🍕"); ok {
			return evalInfixExpression(node.Operator, left, right)
		}

		return createError("🍕が定義されていません")
	}

	// 通常の中置式評価
	return evalStandardInfixExpression(node, env)
}
