package evaluator

import (
	"github.com/uncode/ast"
	"github.com/uncode/logger"
	"github.com/uncode/object"
)

// evalAssignment は>>演算子による代入を評価する
func evalAssignment(node *ast.InfixExpression, env *object.Environment) object.Object {
	logger.Debug("代入演算子を検出しました")
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
		logger.Debug("💩への代入を検出しました - 戻り値として扱います")
		left := Eval(node.Left, env)
		if left.Type() == object.ERROR_OBJ {
			return left
		}
		logger.Debug("💩に戻り値として %s を設定します\n", left.Inspect())
		return &object.ReturnValue{Value: left}
	}

	// 右辺がCaseStatementの場合はcase文として扱う
	if caseStmt, ok := right.(*ast.CaseStatement); ok {
		logger.Debug("case文への代入を検出しました")
		return evalCaseStatement(caseStmt, env)
	}

	return createError("代入先が識別子または💩ではありません: %T", right)
}
