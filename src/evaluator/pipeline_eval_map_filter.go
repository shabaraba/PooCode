package evaluator

import (
	"github.com/uncode/ast"
	"github.com/uncode/logger"
	"github.com/uncode/object"
)

// デバッグログレベル設定
var (
	// mapFilterDebugLevel はmap/filter演算子のデバッグレベルを保持します
	mapFilterDebugLevel = logger.LevelDebug
	
	// argumentsDebugLevel は関数引数のバインディングのデバッグレベルを保持します
	argumentsDebugLevel = logger.LevelDebug
	
	// isArgumentsDebugEnabled は関数引数デバッグが有効かどうかを示します
	isArgumentsDebugEnabled = false
)

// SetMapFilterDebugLevel はmap/filter演算子のデバッグレベルを設定します
func SetMapFilterDebugLevel(level logger.LogLevel) {
	mapFilterDebugLevel = level
	logger.Debug("map/filter演算子のデバッグレベルを %d に設定しました", level)
}

// SetArgumentsDebugLevel は関数引数のバインディングのデバッグレベルを設定します
func SetArgumentsDebugLevel(level logger.LogLevel) {
	argumentsDebugLevel = level
	logger.Debug("関数引数バインディングのデバッグレベルを %d に設定しました", level)
}

// EnableArgumentsDebug は関数引数のデバッグを有効にします
func EnableArgumentsDebug() {
	isArgumentsDebugEnabled = true
	logger.Debug("関数引数デバッグを有効にしました")
}

// DisableArgumentsDebug は関数引数のデバッグを無効にします
func DisableArgumentsDebug() {
	isArgumentsDebugEnabled = false
	logger.Debug("関数引数デバッグを無効にしました")
}

// LogArgumentBinding は関数引数のバインディングをログに記録します（デバッグが有効な場合のみ）
func LogArgumentBinding(funcName string, paramName string, value object.Object) {
	if isArgumentsDebugEnabled && logger.IsLevelEnabled(argumentsDebugLevel) {
		logger.Log(argumentsDebugLevel, "関数 '%s': パラメータ '%s' に値 '%s' をバインドしました", 
			funcName, paramName, value.Inspect())
	}
}

// evalInfixExpressionWithNode は中置式を評価する
func evalInfixExpressionWithNode(node *ast.InfixExpression, env *object.Environment) object.Object {
	logger.Debug("中置式を評価します: %s", node.Operator)

	switch node.Operator {
	case "+>": // map演算子
		logger.Debug("map パイプ演算子 (%s) を検出しました", node.Operator)
		// map関数の処理を実行
		return evalMapOperation(node, env)
	case "?>": // filter演算子
		logger.Debug("filter パイプ演算子 (%s) を検出しました", node.Operator)
		// filter関数の処理を実行
		return evalFilterOperation(node, env)
	case "|>": // 標準パイプライン
		logger.Debug("標準パイプライン演算子 (|>) を検出しました")
		return evalPipeline(node, env)
	case "|": // 並列パイプ
		logger.Debug("並列パイプ演算子 (|) を検出しました")
		// 並列パイプの処理は通常評価
		return evalStandardInfixExpression(node, env)
	case ">>": // 代入演算子
		logger.Debug("代入演算子 (>>) を検出しました")
		return evalAssignment(node, env)
	case "=": // 通常の代入演算子
		logger.Debug("通常の代入演算子 (=) を検出しました")
		return evalAssignment(node, env)
	default:
		// その他の演算子は通常の中置式評価
		return evalStandardInfixExpression(node, env)
	}
}

// evalMapOperation はmap演算子(+>)を処理する
func evalMapOperation(node *ast.InfixExpression, env *object.Environment) object.Object {
	logger.Debug("mapパイプライン演算子(+>)の処理を開始")

	// 左辺値の評価（通常は配列）
	left := Eval(node.Left, env)
	if left == nil {
		return createError("mapオペレーション: 左辺の評価結果がnilです")
	}
	if left.Type() == object.ERROR_OBJ {
		return left
	}
	
	// 配列であることを確認
	arr, ok := left.(*object.Array)
	if !ok {
		return createError("map演算子の左辺は配列である必要があります")
	}

	logger.Debug("+> 左辺の評価結果: %s (タイプ: %s)", left.Inspect(), left.Type())

	// 右辺値の評価（関数または関数呼び出し）
	var funcName string
	var funcArgs []object.Object

	switch right := node.Right.(type) {
	case *ast.Identifier:
		// 識別子の場合、関数名として扱う
		logger.Debug("右辺が識別子: %s", right.Value)
		funcName = right.Value
	case *ast.CallExpression:
		logger.Debug("右辺が関数呼び出し式")
		
		// 関数名を取得
		if ident, ok := right.Function.(*ast.Identifier); ok {
			funcName = ident.Value
			logger.Debug("関数名: %s", funcName)
			
			// 追加引数を評価
			funcArgs = evalExpressions(right.Arguments, env)
			if len(funcArgs) > 0 && funcArgs[0] != nil && funcArgs[0].Type() == object.ERROR_OBJ {
				return funcArgs[0]
			}
		} else {
			return createError("関数呼び出し式の関数部分が識別子ではありません: %T", right.Function)
		}
		
		// 別のケース（CallExpressionの処理）は元のコードをそのまま利用
		leftElements := arr.Elements
		// マップ処理の実行
		resultElements := make([]object.Object, 0, len(leftElements))
		for _, leftElement := range leftElements {
			result := evalPipelineWithCallExpression(leftElement, right, env)
			resultElements = append(resultElements, result)
		}
		return &object.Array{Elements: resultElements}
	default:
		return createError("map演算子の右辺が関数または識別子ではありません: %T", node.Right)
	}

	// 直接配列の各要素に対して処理を行う
	resultElements := make([]object.Object, 0, len(arr.Elements))
	
	for _, elem := range arr.Elements {
		// 一時環境を作成し、🍕に要素をセット
		tempEnv := object.NewEnclosedEnvironment(env)
		tempEnv.Set("🍕", elem)
		
		// 現在の要素に対して適切な関数を選択・実行
		// 引数にはelemを含める
		args := []object.Object{elem}
		if funcArgs != nil {
			args = append(args, funcArgs...)
		}
		
		// 関数を取得（環境から検索）
		functions := env.GetAllFunctionsByName(funcName)
		if len(functions) == 0 {
			// 組み込み関数を確認
			if builtin, ok := Builtins[funcName]; ok {
				logger.Debug("ビルトイン関数 '%s' をマップ操作で呼び出します", funcName)
				result := builtin.Fn(args...)
				if result == nil || result.Type() == object.ERROR_OBJ {
					return result
				}
				resultElements = append(resultElements, result)
				continue
			}
			return createError("関数 '%s' が見つかりません", funcName)
		}
		
		// 関数を適用
		logger.Debug("要素 %s に対して関数 %s を適用", elem.Inspect(), funcName)
		result := applyFunctionWithPizza(functions[0], args)
		
		if result == nil || result.Type() == object.ERROR_OBJ {
			logger.Debug("関数 %s の適用中にエラーが発生: %s", funcName, result.Inspect())
			return result
		}
		
		resultElements = append(resultElements, result)
	}
	
	return &object.Array{Elements: resultElements}
}

// evalFilterOperation はfilter演算子(?>)を処理する
func evalFilterOperation(node *ast.InfixExpression, env *object.Environment) object.Object {
	if logger.IsLevelEnabled(mapFilterDebugLevel) {
		logger.Debug("filter演算子(?>)の処理を開始")
	}

	// 左辺値の評価（通常は配列）
	left := Eval(node.Left, env)
	if left == nil {
		return createError("filterオペレーション: 左辺の評価結果がnilです")
	}
	if left.Type() == object.ERROR_OBJ {
		return left
	}
	
	// 配列であることを確認
	arr, ok := left.(*object.Array)
	if !ok {
		return createError("filter演算子の左辺は配列である必要があります")
	}

	if logger.IsLevelEnabled(mapFilterDebugLevel) {
		logger.Debug("?> 左辺の評価結果: %s (タイプ: %s)", left.Inspect(), left.Type())
	}

	// 右辺値の評価（関数または関数呼び出し）
	var funcName string
	var funcArgs []object.Object

	switch right := node.Right.(type) {
	case *ast.Identifier:
		// 識別子の場合、関数名として扱う
		if logger.IsLevelEnabled(mapFilterDebugLevel) {
			logger.Debug("右辺が識別子: %s", right.Value)
		}
		funcName = right.Value
	case *ast.CallExpression:
		// 関数呼び出しの場合
		if logger.IsLevelEnabled(mapFilterDebugLevel) {
			logger.Debug("右辺が関数呼び出し式")
		}
		if ident, ok := right.Function.(*ast.Identifier); ok {
			// 関数名を取得
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
			return createError("関数呼び出し式の関数部分が識別子ではありません: %T", right.Function)
		}
		
		// CallExpressionの場合、evalPipelineWithCallExpressionを使用して評価
		leftElements := arr.Elements
		// フィルター処理の実行
		resultElements := make([]object.Object, 0)
		for _, leftElement := range leftElements {
			// 各要素に対して関数を適用
			result := evalPipelineWithCallExpression(leftElement, right, env)
			
			// 結果がtruthyな場合のみ結果に含める
			if isTruthy(result) {
				resultElements = append(resultElements, leftElement)
			}
		}
		return &object.Array{Elements: resultElements}
	default:
		return createError("filter演算子の右辺が関数または識別子ではありません: %T", node.Right)
	}

	// 直接配列の各要素に対して処理を行う
	resultElements := make([]object.Object, 0)
	
	for _, elem := range arr.Elements {
		// 一時環境を作成し、🍕に要素をセット
		tempEnv := object.NewEnclosedEnvironment(env)
		tempEnv.Set("🍕", elem)
		
		// 現在の要素に対して適切な関数を選択・実行
		// 引数にはelemを含める
		args := []object.Object{elem}
		if funcArgs != nil {
			args = append(args, funcArgs...)
		}
		
		// 関数を取得（環境から検索）
		functions := env.GetAllFunctionsByName(funcName)
		if len(functions) == 0 {
			// 組み込み関数を確認
			if builtin, ok := Builtins[funcName]; ok {
				logger.Debug("ビルトイン関数 '%s' をフィルター操作で呼び出します", funcName)
				result := builtin.Fn(args...)
				if result == nil || result.Type() == object.ERROR_OBJ {
					return result
				}
				if isTruthy(result) {
					resultElements = append(resultElements, elem)
				}
				continue
			}
			return createError("関数 '%s' が見つかりません", funcName)
		}
		
		// 関数を適用
		logger.Debug("要素 %s に対して関数 %s を適用", elem.Inspect(), funcName)
		result := applyFunctionWithPizza(functions[0], args)
		
		if result == nil || result.Type() == object.ERROR_OBJ {
			logger.Debug("関数 %s の適用中にエラーが発生: %s", funcName, result.Inspect())
			return result
		}
		
		// 結果がtruthyな場合のみ結果に含める
		if isTruthy(result) {
			resultElements = append(resultElements, elem)
		}
	}
	
	return &object.Array{Elements: resultElements}
}
