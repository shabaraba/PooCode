package evaluator

import (
	"github.com/uncode/ast"
	"github.com/uncode/logger"
	"github.com/uncode/object"
)

// evalBlockExpression はブロック式を評価する
func evalBlockExpression(node *ast.BlockExpression, env *object.Environment) object.Object {
	logger.Debug("ブロック式の評価開始")
	
	// ブロック式はそのまま内部のBlockStatementを評価する
	return evalBlockStatement(node.Block, env)
}
