package evaluator

import (
	"github.com/uncode/ast"
	"github.com/uncode/object"
)

// evalRangeExpression evaluates range expressions
func evalRangeExpression(node *ast.RangeExpression, env *object.Environment) object.Object {
	var startObj, endObj object.Object
	
	if node.Start \!= nil {
		startObj = Eval(node.Start, env)
		if startObj.Type() == object.ERROR_OBJ {
			return startObj
		}
	}
	
	if node.End \!= nil {
		endObj = Eval(node.End, env)
		if endObj.Type() == object.ERROR_OBJ {
			return endObj
		}
	}
	
	// Empty range [..]
	if node.Start == nil && node.End == nil {
		return &object.Array{Elements: []object.Object{}}
	}
	
	// Integer range
	if (startObj == nil || startObj.Type() == object.INTEGER_OBJ) && 
	   (endObj == nil || endObj.Type() == object.INTEGER_OBJ) {
		return generateIntegerRange(startObj, endObj)
	}
	
	// String range
	if (startObj == nil || startObj.Type() == object.STRING_OBJ) && 
	   (endObj == nil || endObj.Type() == object.STRING_OBJ) {
		return generateStringRange(startObj, endObj)
	}
	
	// Error for unsupported types
	return createError("Unsupported range expression type: %s..%s", 
		getTypeName(startObj), getTypeName(endObj))
}

// generateIntegerRange generates integer ranges
func generateIntegerRange(startObj, endObj object.Object) object.Object {
	var start, end int64
	
	if startObj == nil {
		start = 0
	} else {
		start = startObj.(*object.Integer).Value
	}
	
	if endObj == nil {
		end = start
	} else {
		end = endObj.(*object.Integer).Value
	}
	
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

// generateStringRange generates character ranges
func generateStringRange(startObj, endObj object.Object) object.Object {
	var startChar, endChar rune
	
	if startObj == nil {
		startChar = 'a'
	} else {
		startStr := startObj.(*object.String).Value
		if len(startStr) \!= 1 {
			return createError("Start value of character range must be a single character: %s", startStr)
		}
		startChar = []rune(startStr)[0]
	}
	
	if endObj == nil {
		endChar = startChar
	} else {
		endStr := endObj.(*object.String).Value
		if len(endStr) \!= 1 {
			return createError("End value of character range must be a single character: %s", endStr)
		}
		endChar = []rune(endStr)[0]
	}
	
	var elements []object.Object
	if startChar <= endChar {
		for c := startChar; c <= endChar; c++ {
			elements = append(elements, &object.String{Value: string(c)})
		}
	} else {
		for c := startChar; c >= endChar; c-- {
			elements = append(elements, &object.String{Value: string(c)})
		}
	}
	
	return &object.Array{Elements: elements}
}

// evalIndexExpression evaluates index expressions
func evalIndexExpression(left, index object.Object, env *object.Environment) object.Object {
	switch {
	case left.Type() == object.ARRAY_OBJ:
		return evalArrayIndexExpression(left, index)
	case left.Type() == object.STRING_OBJ:
		return evalStringIndexExpression(left, index)
	default:
		return createError("Index operator not supported for: %s", left.Type())
	}
}

// evalArrayIndexExpression evaluates array index expressions
func evalArrayIndexExpression(array, index object.Object) object.Object {
	arrayObj := array.(*object.Array)
	
	// Single index access
	if idx, ok := index.(*object.Integer); ok {
		return evalArraySingleIndex(arrayObj, idx.Value)
	}
	
	return createError("Array index must be an integer, got: %s", index.Type())
}

// evalArraySingleIndex evaluates single array index access
func evalArraySingleIndex(array *object.Array, index int64) object.Object {
	length := int64(len(array.Elements))
	
	// Support negative indices
	if index < 0 {
		index = length + index
	}
	
	// Check index bounds
	if index < 0 || index >= length {
		return createError("Index out of bounds: %d (array length: %d)", index, length)
	}
	
	return array.Elements[index]
}

// evalStringIndexExpression evaluates string index expressions
func evalStringIndexExpression(str, index object.Object) object.Object {
	strValue := str.(*object.String).Value
	strRunes := []rune(strValue)
	length := int64(len(strRunes))
	
	idx, ok := index.(*object.Integer)
	if \!ok {
		return createError("String index must be an integer, got: %s", index.Type())
	}
	
	i := idx.Value
	
	// Support negative indices
	if i < 0 {
		i = length + i
	}
	
	// Check index bounds
	if i < 0 || i >= length {
		return createError("Index out of bounds: %d (string length: %d)", i, length)
	}
	
	return &object.String{Value: string(strRunes[i])}
}

// getTypeName returns the type name of an object
func getTypeName(obj object.Object) string {
	if obj == nil {
		return "undefined"
	}
	return string(obj.Type())
}
