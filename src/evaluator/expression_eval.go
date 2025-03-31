package evaluator

import (
	"strconv"
	"strings"

	"github.com/uncode/ast"
	"github.com/uncode/logger"
	"github.com/uncode/object"
)

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

// evalStandardInfixExpression は標準的な中置式を評価する
func evalStandardInfixExpression(node *ast.InfixExpression, env *object.Environment) object.Object {
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
	
	// 文字列と整数の比較
	if (left.Type() == object.STRING_OBJ && right.Type() == object.INTEGER_OBJ) ||
	   (left.Type() == object.INTEGER_OBJ && right.Type() == object.STRING_OBJ) {
		return evalStringIntegerInfixExpression(operator, left, right)
	}
	
	// 真偽値の演算
	if left.Type() == object.BOOLEAN_OBJ && right.Type() == object.BOOLEAN_OBJ {
		return evalBooleanInfixExpression(operator, left, right)
	}
	
	// 型の不一致
	if left.Type() != right.Type() {
		return createError("型の不一致: %s %s %s", left.Type(), operator, right.Type())
	}
	
	return createError("未知の演算子: %s %s %s", left.Type(), operator, right.Type())
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
			return createError("ゼロによる除算: %d / 0", leftVal)
		}
		return &object.Integer{Value: leftVal / rightVal}
	case "%":
		// ゼロ除算チェック
		if rightVal == 0 {
			return createError("ゼロによるモジュロ: %d %% 0", leftVal)
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
		return createError("未知の演算子: %s %s %s", left.Type(), operator, right.Type())
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
		return createError("未知の演算子: %s %s %s", left.Type(), operator, right.Type())
	}
}

// evalStringIntegerInfixExpression は文字列と整数の比較を評価する
func evalStringIntegerInfixExpression(operator string, left, right object.Object) object.Object {
	// 数値の自動変換を試みる
	var leftNum, rightNum int64
	var leftStr, rightStr string
	var err error

	// 左辺と右辺の型に基づいて変数を設定
	if left.Type() == object.STRING_OBJ {
		leftStr = left.(*object.String).Value
		rightNum = right.(*object.Integer).Value
		// 文字列を数値に変換
		// 文字列が数値として解釈できるなら変換して比較、そうでなければ文字列と数値の特別な比較ルールを適用
		if leftNum, err = strconv.ParseInt(leftStr, 10, 64); err != nil {
			// 文字列が数値ではない場合、ほとんどの演算で常にfalseを返す
			// ただし !=, == の場合は特別処理
			switch operator {
			case "==", "eq":
				return FALSE // 文字列と数値は常に等しくない
			case "!=":
				return TRUE // 文字列と数値は常に異なる
			default:
				return createError("型の不一致による比較: %s %s %s", left.Type(), operator, right.Type())
			}
		}
		// 文字列が数値として解釈できる場合、数値比較として扱う
		return evalIntegerInfixExpression(operator, &object.Integer{Value: leftNum}, right)
	} else {
		// 整数が左辺、文字列が右辺の場合
		leftNum = left.(*object.Integer).Value
		rightStr = right.(*object.String).Value
		// 文字列を数値に変換
		if rightNum, err = strconv.ParseInt(rightStr, 10, 64); err != nil {
			// 文字列が数値ではない場合、ほとんどの演算で常にfalseを返す
			switch operator {
			case "==", "eq":
				return FALSE // 数値と文字列は常に等しくない
			case "!=":
				return TRUE // 数値と文字列は常に異なる
			default:
				return createError("型の不一致による比較: %s %s %s", left.Type(), operator, right.Type())
			}
		}
		// 文字列が数値として解釈できる場合、数値比較として扱う
		return evalIntegerInfixExpression(operator, left, &object.Integer{Value: rightNum})
	}
}
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
		return createError("未知の演算子: %s %s %s", left.Type(), operator, right.Type())
	}
}

// evalPrefixExpression は前置式を評価する
func evalPrefixExpression(operator string, right object.Object) object.Object {
	logger.EvalDebug("<<<評価器デバッグ専用ログ>>> 前置式を評価します: operator=%s, right=%s", operator, right.Inspect())
	
	switch operator {
	case "!":
		return evalBangOperatorExpression(right)
	case "-":
		return evalMinusPrefixOperatorExpression(right)
	case "not":
		// 言語仕様で "not" は ! と同様に扱う
		return evalBangOperatorExpression(right)
	default:
		return createError("未知の前置演算子: %s%s", operator, right.Type())
	}
}

// evalBangOperatorExpression は ! 演算子を評価する
func evalBangOperatorExpression(right object.Object) object.Object {
	switch right {
	case TRUE:
		return FALSE
	case FALSE:
		return TRUE
	case NullObj:
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
		return createError("-演算子は整数に対してのみ使用できます: %s", right.Type())
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
	if builtin, ok := Builtins[node.Value]; ok {
		return builtin
	}
	
	return createError("識別子が見つかりません: " + node.Value)
}
