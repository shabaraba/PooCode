package evaluator

// This is a helper file to resolve initialization cycles
// NullObj is declared here for all packages to reference
// This approach ensures there are no circular dependencies between files

import (
	"fmt"
	
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

// CreateError creates an error object with formatted message
// This is exported to avoid initialization cycle issues
func CreateError(format string, a ...interface{}) *object.Error {
	return &object.Error{Message: fmt.Sprintf(format, a...)}
}
