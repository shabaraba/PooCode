package evaluator

// This is a helper file to resolve initialization cycles
// NullObj is declared here for all packages to reference
// This approach ensures there are no circular dependencies between files

import (
	"fmt"
	
	"github.com/uncode/logger"
	"github.com/uncode/object"
)

// Common values for reuse across the package
var (
	// NullObj is the null object representation
	NullObj = &object.Null{}
	// Common boolean values
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
)

// createError creates an error object with formatted message
// This is exported to avoid initialization cycle issues
func createError(format string, a ...interface{}) *object.Error {
	msg := fmt.Sprintf(format, a...)
	logger.Error("実行時エラー: %s", msg)
	return &object.Error{Message: msg}
}
