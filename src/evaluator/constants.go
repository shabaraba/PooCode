package evaluator

import (
	"github.com/uncode/object"
)

var (
	TRUE  = &object.Boolean{Value: true}
	FALSE = &object.Boolean{Value: false}
	NULL  = &object.Null{}
	NullObj = &object.Null{}
)
