package evaluator

import (
	"fmt"

	"github.com/uncode/logger"
	"github.com/uncode/object"
)

// NullObj is a singleton instance of the Null object
var NullObj = &object.Null{}

// createError creates a new Error object with the given message
func createError(format string, a ...interface{}) *object.Error {
	msg := fmt.Sprintf(format, a...)
	logger.Error("エラー: %s", msg)
	return &object.Error{Message: msg}
}

// createEvalError creates an evaluation error with formatted message and logs it
func createEvalError(format string, a ...interface{}) *object.Error {
	msg := fmt.Sprintf(format, a...)
	logger.Error("評価エラー: %s", msg)
	return createError(msg)
}
