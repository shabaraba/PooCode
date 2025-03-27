package evaluator

import (
	"github.com/uncode/object"
)

func registerArrayBuiltins() {
		// map function
	Builtins["map"] = &object.Builtin{
		Name: "map",
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return createError("map function requires 2 arguments")
			}
			
			arr, ok := args[0].(*object.Array)
			if !ok {
				return createError("First argument to map must be an array")
			}
			
			// Check if the second argument is a function (either user-defined or builtin)
			var mapFn func(object.Object) object.Object
			
			switch fn := args[1].(type) {
			case *object.Function:
				// User-defined function
				mapFn = func(elem object.Object) object.Object {
					extendedEnv := object.NewEnclosedEnvironment(fn.Env)
					extendedEnv.Set("🍕", elem)
					
					result := Eval(fn.ASTBody, extendedEnv)
					
					if errObj, ok := result.(*object.Error); ok {
						return errObj
					}
					
					if retVal, ok := result.(*object.ReturnValue); ok {
						result = retVal.Value
					}
					
					return result
				}
			case *object.Builtin:
				// Builtin function
				mapFn = func(elem object.Object) object.Object {
					result := fn.Fn(elem)
					return result
				}
			default:
				return createError("Second argument to map must be a function")
			}
			
			resultElements := make([]object.Object, 0, len(arr.Elements))
			
			for _, elem := range arr.Elements {
				resultElements = append(resultElements, mapFn(elem))
			}
			
			return &object.Array{Elements: resultElements}
		},
		ReturnType: object.ARRAY_OBJ,
		ParamTypes: []object.ObjectType{object.ARRAY_OBJ, object.FUNCTION_OBJ},
	}
	
	// filter function
	Builtins["filter"] = &object.Builtin{
		Name: "filter",
		Fn: func(args ...object.Object) object.Object {
			if len(args) != 2 {
				return createError("filter function requires 2 arguments")
			}
			
			arr, ok := args[0].(*object.Array)
			if !ok {
				return createError("First argument to filter must be an array")
			}
			
			// Check if the second argument is a function (either user-defined or builtin)
			var filterFn func(object.Object) object.Object
			
			switch fn := args[1].(type) {
			case *object.Function:
				// User-defined function
				filterFn = func(elem object.Object) object.Object {
					extendedEnv := object.NewEnclosedEnvironment(fn.Env)
					extendedEnv.Set("🍕", elem)
					
					result := Eval(fn.ASTBody, extendedEnv)
					
					if errObj, ok := result.(*object.Error); ok {
						return errObj
					}
					
					if retVal, ok := result.(*object.ReturnValue); ok {
						result = retVal.Value
					}
					
					return result
				}
			case *object.Builtin:
				// Builtin function
				filterFn = func(elem object.Object) object.Object {
					result := fn.Fn(elem)
					return result
				}
			default:
				return createError("Second argument to filter must be a function")
			}
			
			resultElements := make([]object.Object, 0, len(arr.Elements))
			
			for _, elem := range arr.Elements {
				result := filterFn(elem)
				
				// Check for errors
				if errObj, ok := result.(*object.Error); ok {
					return errObj
				}
				
				if isTruthy(result) {
					resultElements = append(resultElements, elem)
				}
			}
			
			return &object.Array{Elements: resultElements}
		},
		ReturnType: object.ARRAY_OBJ,
		ParamTypes: []object.ObjectType{object.ARRAY_OBJ, object.FUNCTION_OBJ},
	}
}
