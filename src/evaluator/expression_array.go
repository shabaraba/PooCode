package evaluator

import (
	"github.com/uncode/ast"
	"github.com/uncode/object"
)

// evalRangeExpression evaluates range expressions
func evalRangeExpression(node *ast.RangeExpression, env *object.Environment) object.Object {
	var startObj, endObj object.Object
	
	if node.Start != nil {
		startObj = Eval(node.Start, env)
		if startObj.Type() == object.ERROR_OBJ {
			return startObj
		}
	}
	
	if node.End != nil {
		endObj = Eval(node.End, env)
		if endObj.Type() == object.ERROR_OBJ {
			return endObj
		}
	}
	
	// Create integer range
	if startObj != nil && endObj != nil {
		if startObj.Type() == object.INTEGER_OBJ && endObj.Type() == object.INTEGER_OBJ {
			start := startObj.(*object.Integer).Value
			end := endObj.(*object.Integer).Value
			
			var elements []object.Object
			if start <= end {
				for i := start; i <= end; i++ {
					elements = append(elements, &object.Integer{Value: i})
				}
			} else {
				for i := start; i >= end; i-- {
					elements = append(elements, &object.Integer{Value: i})
				}
			}
			
			return &object.Array{Elements: elements}
		}
	}
	
	// Empty array as default
	return &object.Array{Elements: []object.Object{}}
}

// evalIndexExpression evaluates index expressions
func evalIndexExpression(left, index object.Object, env *object.Environment) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ:
		return evalArrayIndexExpression(left, index)
	case left.Type() == object.STRING_OBJ:
		return evalStringIndexExpression(left, index)
	default:
		return createError("Index operator not supported")
	}
}

// evalArrayIndexExpression evaluates array index expressions
func evalArrayIndexExpression(array, index object.Object) object.Object {
	arrayObj := array.(*object.Array)
	
	if idx, ok := index.(*object.Integer); ok {
		return evalArraySingleIndex(arrayObj, idx.Value)
	}
	
	return createError("Array index must be an integer")
}

// evalArraySingleIndex evaluates single array index access
func evalArraySingleIndex(array *object.Array, index int64) object.Object {
	length := int64(len(array.Elements))
	
	if index < 0 {
		index = length + index
	}
	
	if index < 0 || index >= length {
		return createError("Index out of bounds")
	}
	
	return array.Elements[index]
}

// evalStringIndexExpression evaluates string index expressions
func evalStringIndexExpression(str, index object.Object) object.Object {
	strValue := str.(*object.String).Value
	strRunes := []rune(strValue)
	length := int64(len(strRunes))
	
	idx, ok := index.(*object.Integer)
	if !ok {
		return createError("String index must be an integer")
	}
	
	i := idx.Value
	
	if i < 0 {
		i = length + i
	}
	
	if i < 0 || i >= length {
		return createError("Index out of bounds")
	}
	
	return &object.String{Value: string(strRunes[i])}
}
