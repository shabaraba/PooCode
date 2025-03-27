package evaluator

import (
	"github.com/uncode/ast"
	"github.com/uncode/logger"
	"github.com/uncode/object"
)

// evalRangeExpression evaluates range expressions
func evalRangeExpression(node *ast.RangeExpression, env *object.Environment) object.Object {
	logger.Debug("レンジ式を評価: %v..%v", node.Start, node.End)
	var startObj, endObj object.Object
	
	if node.Start != nil {
		startObj = Eval(node.Start, env)
		if startObj.Type() == object.ERROR_OBJ {
			logger.Debug("レンジ式の開始値評価でエラー: %s", startObj.Inspect())
			return startObj
		}
		logger.Debug("レンジ式の開始値: %s", startObj.Inspect())
	} else {
		// Default start for [..end] is 1
		startObj = &object.Integer{Value: 1}
		logger.Debug("レンジ式の開始値がないため、デフォルト値1を使用")
	}
	
	if node.End != nil {
		endObj = Eval(node.End, env)
		if endObj.Type() == object.ERROR_OBJ {
			logger.Debug("レンジ式の終了値評価でエラー: %s", endObj.Inspect())
			return endObj
		}
		logger.Debug("レンジ式の終了値: %s", endObj.Inspect())
	} else {
		// Default end for [start..] is startObj + 10 (just a convention for this language)
		if startObj.Type() == object.INTEGER_OBJ {
			endObj = &object.Integer{Value: startObj.(*object.Integer).Value + 9}
			logger.Debug("レンジ式の終了値がないため、開始値+9の値 %d を使用", endObj.(*object.Integer).Value)
		} else {
			// If start is not an integer, default to empty array
			logger.Debug("レンジ式の開始値が整数でないため、空の配列を返します")
			return &object.Array{Elements: []object.Object{}}
		}
	}
	
	// Create integer range
	if startObj != nil && endObj != nil {
		if startObj.Type() == object.INTEGER_OBJ && endObj.Type() == object.INTEGER_OBJ {
			start := startObj.(*object.Integer).Value
			end := endObj.(*object.Integer).Value
			logger.Debug("整数レンジを作成: %d..%d", start, end)
			
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
			
			result := &object.Array{Elements: elements}
			logger.Debug("レンジ式の評価結果: %s (要素数: %d)", result.Inspect(), len(elements))
			return result
		}
	}
	
	// Empty array as default
	logger.Debug("レンジ式がinteger型でないため、空の配列を返します")
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
