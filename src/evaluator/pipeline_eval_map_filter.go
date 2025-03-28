package evaluator

import (
	"github.com/uncode/ast"
	"github.com/uncode/logger"
	"github.com/uncode/object"
)

// ��ð�����-�
var (
	// mapFilterDebugLevel omap/filter�Pn��ð����W~Y
	mapFilterDebugLevel = logger.LevelDebug
	
	// argumentsDebugLevel o�ppnФ�ǣ�n��ð����W~Y
	argumentsDebugLevel = logger.LevelDebug
	
	// isArgumentsDebugEnabled o�pp��ðL	�KiFK�:W~Y
	isArgumentsDebugEnabled = false
)

// SetMapFilterDebugLevel omap/filter�Pn��ð���-�W~Y
func SetMapFilterDebugLevel(level logger.LogLevel) {
	mapFilterDebugLevel = level
	logger.Debug("map/filter�Pn��ð��� %d k-�W~W_", level)
}

// SetArgumentsDebugLevel o�ppnФ�ǣ�n��ð���-�W~Y
func SetArgumentsDebugLevel(level logger.LogLevel) {
	argumentsDebugLevel = level
	logger.Debug("�ppФ�ǣ�n��ð��� %d k-�W~W_", level)
}

// EnableArgumentsDebug o�ppn��ð�	�kW~Y
func EnableArgumentsDebug() {
	isArgumentsDebugEnabled = true
	logger.Debug("�pp��ð�	�kW~W_")
}

// DisableArgumentsDebug o�ppn��ð�!�kW~Y
func DisableArgumentsDebug() {
	isArgumentsDebugEnabled = false
	logger.Debug("�pp��ð�!�kW~W_")
}

// LogArgumentBinding o�ppnФ�ǣ���k2W~Y��ðL	�j4n	
func LogArgumentBinding(funcName string, paramName string, value object.Object) {
	if isArgumentsDebugEnabled && logger.IsLevelEnabled(argumentsDebugLevel) {
		logger.Log(argumentsDebugLevel, "�p '%s': ����� '%s' k$ '%s' �Ф��W~W_", 
			funcName, paramName, value.Inspect())
	}
}

// evalMapOperation omap�P(+>)��Y�
// X $hMn!�k��Y��Fk�c
func evalMapOperation(node *ast.InfixExpression, env *object.Environment) object.Object {
	logger.Debug("mapѤ����P(+>)n����")

	// �$nU�
	left := Eval(node.Left, env)
	if left == nil {
		return createError("map�������: �nU�P�LnilgY")
	}
	if left.Type() == object.ERROR_OBJ {
		return left
	}
	
	// MKX n$K���Wij��LF
	var elements []object.Object
	var isSingleValue bool
	
	if arrayObj, ok := left.(*object.Array); ok {
		// Mn4o]n� �(
		elements = arrayObj.Elements
		isSingleValue = false
		logger.Debug("+> �nU�P�: M %s (���: %s)", left.Inspect(), left.Type())
	} else {
		// X n$n4o� 1dnMhWfqF
		elements = []object.Object{left}
		isSingleValue = true
		logger.Debug("+> �nU�P�: X $ %s (���: %s) �� 1dnMhWfqD~Y", left.Inspect(), left.Type())
	}

	// �$nU��p~_o�p|s�W	
	var funcName string
	var funcArgs []object.Object

	switch right := node.Right.(type) {
	case *ast.Identifier:
		// X%Pn4�phWfqF
		logger.Debug("�LX%P: %s", right.Value)
		funcName = right.Value
	case *ast.CallExpression:
		logger.Debug("�L�p|s�W")
		
		// �p�֗
		if ident, ok := right.Function.(*ast.Identifier); ok {
			funcName = ident.Value
			logger.Debug("�p: %s", funcName)
			
			// ��p�U�
			funcArgs = evalExpressions(right.Arguments, env)
			if len(funcArgs) > 0 && funcArgs[0] != nil && funcArgs[0].Type() == object.ERROR_OBJ {
				return funcArgs[0]
			}
		} else {
			return createError("�p|s�Wn�p�LX%PgoB�~[�: %T", right.Function)
		}
		
		// CallExpressionn4� k�WfevalPipelineWithCallExpression�i(
		resultElements := make([]object.Object, 0, len(elements))
		for _, element := range elements {
			result := evalPipelineWithCallExpression(element, right, env)
			resultElements = append(resultElements, result)
		}
		
		// X $���n4o nP�`Q��Y
		if isSingleValue && len(resultElements) > 0 {
			return resultElements[0]
		}
		return &object.Array{Elements: resultElements}
	default:
		return createError("map�Pn�L�p~_oX%PgoB�~[�: %T", node.Right)
	}

	// ��� k�Wf��LF - SSL�c�@
	resultElements := make([]object.Object, 0, len(elements))
	
	for _, elem := range elements {
		// p���� �, phW��npLB�p��	
		args := []object.Object{elem}
		if funcArgs != nil {
			args = append(args, funcArgs...)
		}
		
		logger.Debug("����\: �  %s k�Wf�p %s ��", elem.Inspect(), funcName)
		
		// D��p���
		if builtin, ok := Builtins[funcName]; ok {
			logger.Debug("��Ȥ�p '%s' �����\g|s�W~Y", funcName)
			result := builtin.Fn(args...)
			if result == nil || result.Type() == object.ERROR_OBJ {
				return result
			}
			resultElements = append(resultElements, result)
			continue
		}
		
		// q U�_���է���(Wf�p�|s�Y
		// applyNamedFunction oa��M�pnijx�hi(�LF
		logger.Debug("�  %s k�WfapplyNamedFunction�|s�W", elem.Inspect())
		result := applyNamedFunction(env, funcName, args)
		
		// �����ï
		if result == nil {
			return createError("�p '%s' ni(P�LnilgY", funcName)
		}
		if result.Type() == object.ERROR_OBJ {
			return result
		}
		
		// P��Mk��
		resultElements = append(resultElements, result)
	}
	
	// X $���n4o nP�`Q��Y
	if isSingleValue && len(resultElements) > 0 {
		return resultElements[0]
	}
	
	return &object.Array{Elements: resultElements}
}

// evalFilterOperation ofilter�P(?>)��Y�
// �LX $n4n���Ȃ��
func evalFilterOperation(node *ast.InfixExpression, env *object.Environment) object.Object {
	if logger.IsLevelEnabled(mapFilterDebugLevel) {
		logger.Debug("filter�P(?>)n����")
	}

	// �$nU�
	left := Eval(node.Left, env)
	if left == nil {
		return createError("filter�������: �nU�P�LnilgY")
	}
	if left.Type() == object.ERROR_OBJ {
		return left
	}
	
	// MKX n$K���Wij��LF
	var elements []object.Object
	var isSingleValue bool
	
	if arrayObj, ok := left.(*object.Array); ok {
		// Mn4o]n� �(
		elements = arrayObj.Elements
		isSingleValue = false
		logger.Debug("?> �nU�P�: M %s (���: %s)", left.Inspect(), left.Type())
	} else {
		// X n$n4o� 1dnMhWfqF
		elements = []object.Object{left}
		isSingleValue = true
		logger.Debug("?> �nU�P�: X $ %s (���: %s) �� 1dnMhWfqD~Y", left.Inspect(), left.Type())
	}

	// �$nU��p~_o�p|s�W	
	var funcName string
	var funcArgs []object.Object

	switch right := node.Right.(type) {
	case *ast.Identifier:
		// X%Pn4�phWfqF
		if logger.IsLevelEnabled(mapFilterDebugLevel) {
			logger.Debug("�LX%P: %s", right.Value)
		}
		funcName = right.Value
	case *ast.CallExpression:
		// �p|s�Wn4
		if logger.IsLevelEnabled(mapFilterDebugLevel) {
			logger.Debug("�L�p|s�W")
		}
		if ident, ok := right.Function.(*ast.Identifier); ok {
			// �p�֗
			funcName = ident.Value
			if logger.IsLevelEnabled(mapFilterDebugLevel) {
				logger.Debug("�p: %s", funcName)
			}

			// pnU�
			funcArgs = evalExpressions(right.Arguments, env)
			if len(funcArgs) > 0 && funcArgs[0] != nil && funcArgs[0].Type() == object.ERROR_OBJ {
				return funcArgs[0]
			}
		} else {
			return createError("�p|s�Wn�p�LX%PgoB�~[�: %T", right.Function)
		}
		
		// CallExpressionn4evalPipelineWithCallExpression�(WfU�
		resultElements := make([]object.Object, 0)
		for _, element := range elements {
			// � k�Wf�p�i(
			result := evalPipelineWithCallExpression(element, right, env)
			
			// P�Ltruthyj4nP�k+��
			if isTruthy(result) {
				resultElements = append(resultElements, element)
			}
		}
		
		// X $���n4P�LB�pCn$�jQ�pnull��Y
		if isSingleValue {
			if len(resultElements) > 0 {
				return left // CnX $��Y
			}
			return NULL
		}
		
		return &object.Array{Elements: resultElements}
	default:
		return createError("filter�Pn�L�p~_oX%PgoB�~[�: %T", node.Right)
	}

	// ��Mn� k�Wf��LF - SS��k�c
	resultElements := make([]object.Object, 0)
	
	for _, elem := range elements {
		// p���
		args := []object.Object{elem}
		if funcArgs != nil {
			args = append(args, funcArgs...)
		}
		
		logger.Debug("գ���\: �  %s k�Wf�p %s ��", elem.Inspect(), funcName)
		
		// D��p���
		if builtin, ok := Builtins[funcName]; ok {
			logger.Debug("��Ȥ�p '%s' �գ���\g|s�W~Y", funcName)
			result := builtin.Fn(args...)
			if result == nil || result.Type() == object.ERROR_OBJ {
				return result
			}
			
			// P�Ltruthyj4nP�k+��
			if isTruthy(result) {
				resultElements = append(resultElements, elem)
			}
			continue
		}
		
		// q U�_���է���(Wf�p�|s�Y
		logger.Debug("�  %s k�WfapplyNamedFunction�|s�W", elem.Inspect())
		result := applyNamedFunction(env, funcName, args)
		
		// �����ï
		if result == nil {
			continue // գ�jngP�Lniln4o����
		}
		if result.Type() == object.ERROR_OBJ {
			return result
		}
		
		// P�Ltruthyj4nCn� �P�k+��
		if isTruthy(result) {
			resultElements = append(resultElements, elem)
		}
	}
	
	// X $���n4P�LB�pCn$�jQ�pnull��Y
	if isSingleValue {
		if len(resultElements) > 0 {
			return left // CnX $��Y
		}
		return NULL
	}
	
	return &object.Array{Elements: resultElements}
}
