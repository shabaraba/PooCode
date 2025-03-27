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
			
			fn, ok := args[1].(*object.Function)
			if !ok {
				return createError("Second argument to map must be a function")
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
				
				resultElements = append(resultElements, result)
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
			
			fn, ok := args[1].(*object.Function)
			if !ok {
				return createError("Second argument to filter must be a function")
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
