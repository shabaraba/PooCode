package evaluator

import (
	"github.com/uncode/object"
)

// registerArrayBuiltins registers array-related builtin functions
func registerArrayBuiltins() {
	// map function - applies a function to each element of an array
	Builtins["map"] = &object.Builtin{
		Name: "map",
		Fn: func(args ...object.Object) object.Object {
			// Check number of arguments
			if len(args) \!= 2 {
				return createError("map function requires 2 arguments: array, function")
			}
			
			// Check first argument is an array
			arr, ok := args[0].(*object.Array)
			if \!ok {
				return createError("First argument to map must be an array, got: %s", args[0].Type())
			}
			
			// Check second argument is a function
			fn, ok := args[1].(*object.Function)
			if \!ok {
				return createError("Second argument to map must be a function, got: %s", args[1].Type())
			}
			
			// Function should not have parameters
			if len(fn.Parameters) > 0 {
				return createError("Function passed to map should not take parameters")
			}
			
			// Map each element
			resultElements := make([]object.Object, 0, len(arr.Elements))
			
			for _, elem := range arr.Elements {
				// Create extended environment with 🍕 set to current element
				extendedEnv := object.NewEnclosedEnvironment(fn.Env)
				extendedEnv.Set("🍕", elem)
				
				// Evaluate function with element
				result := Eval(fn.ASTBody, extendedEnv)
				
				// Check for errors
				if errObj, ok := result.(*object.Error); ok {
					return errObj
				}
				
				// Unwrap return value
				if retVal, ok := result.(*object.ReturnValue); ok {
					result = retVal.Value
				}
				
				// Add result to array
				resultElements = append(resultElements, result)
			}
			
			return &object.Array{Elements: resultElements}
		},
		ReturnType: object.ARRAY_OBJ,
		ParamTypes: []object.ObjectType{object.ARRAY_OBJ, object.FUNCTION_OBJ},
	}
	
	// filter function - filters elements based on a condition function
	Builtins["filter"] = &object.Builtin{
		Name: "filter",
		Fn: func(args ...object.Object) object.Object {
			// Check number of arguments
			if len(args) \!= 2 {
				return createError("filter function requires 2 arguments: array, function")
			}
			
			// Check first argument is an array
			arr, ok := args[0].(*object.Array)
			if \!ok {
				return createError("First argument to filter must be an array, got: %s", args[0].Type())
			}
			
			// Check second argument is a function
			fn, ok := args[1].(*object.Function)
			if \!ok {
				return createError("Second argument to filter must be a function, got: %s", args[1].Type())
			}
			
			// Function should not have parameters
			if len(fn.Parameters) > 0 {
				return createError("Function passed to filter should not take parameters")
			}
			
			// Filter elements
			resultElements := make([]object.Object, 0)
			
			for _, elem := range arr.Elements {
				// Create extended environment with 🍕 set to current element
				extendedEnv := object.NewEnclosedEnvironment(fn.Env)
				extendedEnv.Set("🍕", elem)
				
				// Evaluate condition function with element
				result := Eval(fn.ASTBody, extendedEnv)
				
				// Check for errors
				if errObj, ok := result.(*object.Error); ok {
					return errObj
				}
				
				// Unwrap return value
				if retVal, ok := result.(*object.ReturnValue); ok {
					result = retVal.Value
				}
				
				// Only add element if condition is true
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
			
			resultElements := make([]object.Object, 0, len(arr.Elements))
			
			for _, elem := range arr.Elements {
				extendedEnv := object.NewEnclosedEnvironment(fn.Env)
				extendedEnv.Set("🍕", elem)
				
				result := Eval(fn.ASTBody, extendedEnv)
				
				if errObj, ok := result.(*object.Error); ok {
					return errObj
				}
				
				if retVal, ok := result.(*object.ReturnValue); ok {
					result = retVal.Value
				}
				
				resultElements = append(resultElements, result)
			}
			
			return &object.Array{Elements: resultElements}
		},
		ReturnType: object.ARRAY_OBJ,
		ParamTypes: []object.ObjectType{object.ARRAY_OBJ, object.FUNCTION_OBJ},
	}
	
	Builtins["filter"] = &object.Builtin{
		Name: "filter",
		Fn: func(args ...object.Object) object.Object {
			if len(args) \!= 2 {
				return createError("filter function requires 2 arguments")
			}
			
			arr, ok := args[0].(*object.Array)
			if \!ok {
				return createError("First argument to filter must be an array")
			}
			
			fn, ok := args[1].(*object.Function)
			if \!ok {
				return createError("Second argument to filter must be a function")
			}
			
			if len(fn.Parameters) > 0 {
				return createError("Function passed to filter should not take parameters")
			}
			
			resultElements := make([]object.Object, 0)
			
			for _, elem := range arr.Elements {
				extendedEnv := object.NewEnclosedEnvironment(fn.Env)
				extendedEnv.Set("🍕", elem)
				
				result := Eval(fn.ASTBody, extendedEnv)
				
				if errObj, ok := result.(*object.Error); ok {
					return errObj
				}
				
				if retVal, ok := result.(*object.ReturnValue); ok {
					result = retVal.Value
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
