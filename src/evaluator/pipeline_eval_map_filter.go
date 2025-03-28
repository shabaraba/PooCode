package evaluator

import (
	"github.com/uncode/ast"
	"github.com/uncode/logger"
	"github.com/uncode/object"
)

// デバッグレベル設定
var (
	// mapFilterDebugLevel はmap/filter操作のデバッグレベル
	mapFilterDebugLevel = logger.LevelDebug
	
	// argumentsDebugLevel は引数バインディングのデバッグレベル
	argumentsDebugLevel = logger.LevelDebug
	
	// isArgumentsDebugEnabled は引数デバッグが有効かどうか
	isArgumentsDebugEnabled = false
)

// SetMapFilterDebugLevel はmap/filter操作のデバッグレベルを設定します
func SetMapFilterDebugLevel(level logger.LogLevel) {
	mapFilterDebugLevel = level
	logger.Debug("map/filter操作のデバッグレベルを %d に設定しました", level)
}

// SetArgumentsDebugLevel は引数バインディングのデバッグレベルを設定します
func SetArgumentsDebugLevel(level logger.LogLevel) {
	argumentsDebugLevel = level
	logger.Debug("引数バインディングのデバッグレベルを %d に設定しました", level)
}

// EnableArgumentsDebug は引数のデバッグを有効にします
func EnableArgumentsDebug() {
	isArgumentsDebugEnabled = true
	logger.Debug("引数デバッグを有効にしました")
}

// DisableArgumentsDebug は引数のデバッグを無効にします
func DisableArgumentsDebug() {
	isArgumentsDebugEnabled = false
	logger.Debug("引数デバッグを無効にしました")
}

// LogArgumentBinding は引数のバインディングをログ出力します（デバッグが有効な場合）
func LogArgumentBinding(funcName string, paramName string, value object.Object) {
	if isArgumentsDebugEnabled && logger.IsLevelEnabled(argumentsDebugLevel) {
		logger.Log(argumentsDebugLevel, "関数 '%s': パラメータ '%s' に値 '%s' をバインドしました", 
			funcName, paramName, value.Inspect())
	}
}

// evalMapOperation はmap操作(+>)を評価する
// 各要素に関数を適用して結果を返す
func evalMapOperation(node *ast.InfixExpression, env *object.Environment) object.Object {
	logger.Debug("mapオペレーション(+>)の評価開始")

	// 左側の評価
	left := Eval(node.Left, env)
	if left == nil {
		return createError("map操作エラー: 左の評価結果がnilです")
	}
	if left.Type() == object.ERROR_OBJ {
		return left
	}
	
	// 配列か単一の値かを判断して準備する
	var elements []object.Object
	var isSingleValue bool
	
	if arrayObj, ok := left.(*object.Array); ok {
		// 配列の場合はその要素
		elements = arrayObj.Elements
		isSingleValue = false
		logger.Debug("+> 左の評価結果: 配列 %s (型: %s)", left.Inspect(), left.Type())
	} else {
		// 単一の値の場合は 1個の配列として扱う
		elements = []object.Object{left}
		isSingleValue = true
		logger.Debug("+> 左の評価結果: 単一値 %s (型: %s) を 1個の配列として扱います", left.Inspect(), left.Type())
	}

	// 右側の評価：関数または関数呼び出し
	var funcName string
	var funcArgs []object.Object

	switch right := node.Right.(type) {
	case *ast.Identifier:
		// 識別子の場合、関数として扱う
		logger.Debug("右は識別子: %s", right.Value)
		funcName = right.Value
	case *ast.CallExpression:
		logger.Debug("右は関数呼び出し式")
		
		// 関数名取得
		if ident, ok := right.Function.(*ast.Identifier); ok {
			funcName = ident.Value
			logger.Debug("関数名: %s", funcName)
			
			// 追加引数評価
			funcArgs = evalExpressions(right.Arguments, env)
			if len(funcArgs) > 0 && funcArgs[0] != nil && funcArgs[0].Type() == object.ERROR_OBJ {
				return funcArgs[0]
			}
		} else {
			return createError("関数呼び出し式の関数が識別子ではありません: %T", right.Function)
		}
		
		// CallExpressionの場合、特別に処理するevalPipelineWithCallExpressionを使用
		resultElements := make([]object.Object, 0, len(elements))
		for _, element := range elements {
			result := evalPipelineWithCallExpression(element, right, env)
			resultElements = append(resultElements, result)
		}
		
		// 単一値処理の場合は先頭の結果だけ返す
		if isSingleValue && len(resultElements) > 0 {
			return resultElements[0]
		}
		return &object.Array{Elements: resultElements}
	default:
		return createError("map操作の右が関数または識別子ではありません: %T", node.Right)
	}

	// 各要素に対して処理する - 標準ケース実装
	resultElements := make([]object.Object, 0, len(elements))
	
	for _, elem := range elements {
		// 引数準備：要素自身を第一引数として、追加の引数も設定
		args := []object.Object{elem}
		if funcArgs != nil {
			args = append(args, funcArgs...)
		}
		
		logger.Debug("マップ処理: 要素 %s に対して関数 %s を適用", elem.Inspect(), funcName)
		
		// 組み込み関数処理
		if builtin, ok := Builtins[funcName]; ok {
			logger.Debug("組み込み関数 '%s' をマップ処理で呼び出します", funcName)
			result := builtin.Fn(args...)
			if result == nil || result.Type() == object.ERROR_OBJ {
				return result
			}
			resultElements = append(resultElements, result)
			continue
		}
		
		// 環境変数内に定義された関数を呼び出す
		// applyNamedFunction は任意の名前付き関数を扱える関数
		logger.Debug("要素 %s に対してapplyNamedFunctionを呼び出し", elem.Inspect())
		result := applyNamedFunction(env, funcName, args)
		
		// エラー処理
		if result == nil {
			return createError("関数 '%s' の結果がnilです", funcName)
		}
		if result.Type() == object.ERROR_OBJ {
			return result
		}
		
		// 結果を配列に追加
		resultElements = append(resultElements, result)
	}
	
	// 単一値処理の場合は先頭の結果だけ返す
	if isSingleValue && len(resultElements) > 0 {
		return resultElements[0]
	}
	
	return &object.Array{Elements: resultElements}
}

// evalFilterOperation はfilter操作(?>)を評価する
// 条件を満たす要素のみを返す
func evalFilterOperation(node *ast.InfixExpression, env *object.Environment) object.Object {
	if logger.IsLevelEnabled(mapFilterDebugLevel) {
		logger.Debug("filter操作(?>)の評価開始")
	}

	// 左側の評価
	left := Eval(node.Left, env)
	if left == nil {
		return createError("filter操作エラー: 左の評価結果がnilです")
	}
	if left.Type() == object.ERROR_OBJ {
		return left
	}
	
	// 配列か単一の値かを判断して準備する
	var elements []object.Object
	var isSingleValue bool
	
	if arrayObj, ok := left.(*object.Array); ok {
		// 配列の場合はその要素
		elements = arrayObj.Elements
		isSingleValue = false
		logger.Debug("?> 左の評価結果: 配列 %s (型: %s)", left.Inspect(), left.Type())
	} else {
		// 単一の値の場合は 1個の配列として扱う
		elements = []object.Object{left}
		isSingleValue = true
		logger.Debug("?> 左の評価結果: 単一値 %s (型: %s) を 1個の配列として扱います", left.Inspect(), left.Type())
	}

	// 右側の評価：関数または関数呼び出し
	var funcName string
	var funcArgs []object.Object

	switch right := node.Right.(type) {
	case *ast.Identifier:
		// 識別子の場合、関数として扱う
		if logger.IsLevelEnabled(mapFilterDebugLevel) {
			logger.Debug("右は識別子: %s", right.Value)
		}
		funcName = right.Value
	case *ast.CallExpression:
		// 関数呼び出しの場合
		if logger.IsLevelEnabled(mapFilterDebugLevel) {
			logger.Debug("右は関数呼び出し式")
		}
		if ident, ok := right.Function.(*ast.Identifier); ok {
			// 関数名取得
			funcName = ident.Value
			if logger.IsLevelEnabled(mapFilterDebugLevel) {
				logger.Debug("関数名: %s", funcName)
			}

			// 引数の評価
			funcArgs = evalExpressions(right.Arguments, env)
			if len(funcArgs) > 0 && funcArgs[0] != nil && funcArgs[0].Type() == object.ERROR_OBJ {
				return funcArgs[0]
			}
		} else {
			return createError("関数呼び出し式の関数が識別子ではありません: %T", right.Function)
		}
		
		// CallExpressionの場合、evalPipelineWithCallExpressionを使って評価
		resultElements := make([]object.Object, 0)
		for _, element := range elements {
			// 要素に対して関数を評価
			result := evalPipelineWithCallExpression(element, right, env)
			
			// 結果がtruthyならば元の要素を保持
			if isTruthy(result) {
				resultElements = append(resultElements, element)
			}
		}
		
		// 単一値処理の場合、結果があれば元の値、なければnullを返す
		if isSingleValue {
			if len(resultElements) > 0 {
				return left // 元の単一値を返す
			}
			return NULL
		}
		
		return &object.Array{Elements: resultElements}
	default:
		return createError("filter操作の右が関数または識別子ではありません: %T", node.Right)
	}

	// 各配列の要素に対して処理する - 標準ケース実装
	resultElements := make([]object.Object, 0)
	
	for _, elem := range elements {
		// 引数準備
		args := []object.Object{elem}
		if funcArgs != nil {
			args = append(args, funcArgs...)
		}
		
		logger.Debug("フィルタ処理: 要素 %s に対して関数 %s を適用", elem.Inspect(), funcName)
		
		// 組み込み関数処理
		if builtin, ok := Builtins[funcName]; ok {
			logger.Debug("組み込み関数 '%s' をフィルタ処理で呼び出します", funcName)
			result := builtin.Fn(args...)
			if result == nil || result.Type() == object.ERROR_OBJ {
				return result
			}
			
			// 結果がtruthyならば元の要素を保持
			if isTruthy(result) {
				resultElements = append(resultElements, elem)
			}
			continue
		}
		
		// 環境変数内に定義された関数を呼び出す
		logger.Debug("要素 %s に対してapplyNamedFunctionを呼び出し", elem.Inspect())
		result := applyNamedFunction(env, funcName, args)
		
		// エラー処理
		if result == nil {
			continue // フィルタです場合、結果がnilの場合は無視
		}
		if result.Type() == object.ERROR_OBJ {
			return result
		}
		
		// 結果がtruthyならば元の要素を元の配列に保持
		if isTruthy(result) {
			resultElements = append(resultElements, elem)
		}
	}
	
	// 単一値処理の場合、結果があれば元の値、なければnullを返す
	if isSingleValue {
		if len(resultElements) > 0 {
			return left // 元の単一値を返す
		}
		return NULL
	}
	
	return &object.Array{Elements: resultElements}
}
