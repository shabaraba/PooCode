package evaluator

import (
	"fmt"

	"github.com/uncode/logger"
	"github.com/uncode/object"
)

// createEvalError creates an evaluation error with formatted message and logs it
func createEvalError(format string, a ...interface{}) *object.Error {
	msg := fmt.Sprintf(format, a...)
	logger.Error("評価エラー: %s", msg)
	return createError(msg)
}
