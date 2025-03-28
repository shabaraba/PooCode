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

// evalMapOperation はmap演算子(+>)を処理する
// 単一値と配列の両方に対応するように修正
func evalMapOperation(node *ast.InfixExpression, env *object.Environment) object.Object {
	logger.Debug("mapパイプライン演算子(+>)の処理を開始")

	// 左辺値の評価
	left := Eval(node.Left, env)
	if left == nil {
		return createError("mapオペレーション: 左辺の評価結果がnilです")
	}
	if left.Type() == object.ERROR_OBJ {
		return left
	}
	
	// 配列か単一の値かを確認し、適切な処理を行う
	var elements []object.Object
	var isSingleValue bool
	
	if arrayObj, ok := left.(*object.Array); ok {
		// 配列の場合はその要素を使用
		elements = arrayObj.Elements
		isSingleValue = false
		logger.Debug("+> 左辺の評価結果: 配列 %s (タイプ: %s)", left.Inspect(), left.Type())
	} else {
		// 単一の値の場合は要素1つの配列として扱う
		elements = []object.Object{left}
		isSingleValue = true
		logger.Debug("+> 左辺の評価結果: 単一値 %s (タイプ: %s) を要素1つの配列として扱います", left.Inspect(), left.Type())
	}

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
		
		// CallExpressionの場合、各要素に対してevalPipelineWithCallExpressionを適用
		resultElements := make([]object.Object, 0, len(elements))
		for _, element := range elements {
			result := evalPipelineWithCallExpression(element, right, env)
			resultElements = append(resultElements, result)
		}
		
		// 単一値モードの場合は最初の結果だけを返す
		if isSingleValue && len(resultElements) > 0 {
			return resultElements[0]
		}
		return &object.Array{Elements: resultElements}
	default:
		return createError("map演算子の右辺が関数または識別子ではありません: %T", node.Right)
	}

	// 直接各要素に対して処理を行う
	resultElements := make([]object.Object, 0, len(elements))
	
	for _, elem := range elements {
		// 一時環境を作成し、🍕に要素をセット
		tempEnv := object.NewEnclosedEnvironment(env)
		tempEnv.Set("🍕", elem)
		
		// 現在の要素に対して適切な関数を選択・実行
		// 引数にはelemを含める
		args := []object.Object{elem}
		if funcArgs != nil {
			args = append(args, funcArgs...)
		}
		
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
		
		// 名前付き関数のメカニズムを使用して条件付き関数の正しい選択を行う
		// ここが重要な変更点：条件付き関数を明示的に処理する
		logger.Debug("マップ操作: 要素 %s に対して条件付き関数メカニズムで関数 %s を呼び出し", elem.Inspect(), funcName)
		
		// 環境から全ての関数候補を取得
		functions := env.GetAllFunctionsByName(funcName)
		if len(functions) == 0 {
			logger.Debug("関数 '%s' が見つかりません", funcName)
			return createError("関数 '%s' が見つかりません", funcName)
		}
		
		// 条件付き関数を優先的に探す
		var result object.Object
		var defaultFunc *object.Function
		
		// 現在の要素を処理する一時環境を更新
		tempEnv.Set("🍕", elem)
		
		// まず条件付き関数を試す
		for _, fn := range functions {
			if fn.Condition != nil {
				// 条件式の評価
				isTrue, _ := evalConditionalExpression(fn, args, tempEnv)
				if isTrue {
					// 条件に一致した場合、この関数を使用
					logger.Debug("条件に一致する関数を見つけました: %s", funcName)
					result = applyFunctionWithPizza(fn, args)
					if result != nil {
						// 条件に一致した関数の結果を使用
						resultElements = append(resultElements, result)
						break
					}
				}
			} else {
				// デフォルト関数を保存
				defaultFunc = fn
			}
		}
		
		// 結果がまだ追加されていない場合はデフォルト関数を使う
		if result == nil || result.Type() == object.ERROR_OBJ {
			if defaultFunc != nil {
				logger.Debug("デフォルト関数を使用します: %s", funcName)
				result = applyFunctionWithPizza(defaultFunc, args)
				if result != nil {
					resultElements = append(resultElements, result)
				} else {
					return createError("関数 '%s' の適用に失敗しました", funcName)
				}
			} else {
				return createError("条件に一致する関数が見つかりません: %s", funcName)
			}
		}
	}
	
	// 単一値モードの場合は最初の結果だけを返す
	if isSingleValue && len(resultElements) > 0 {
		return resultElements[0]
	}
	
	return &object.Array{Elements: resultElements}
}

// evalFilterOperation はfilter演算子(?>)を処理する
// 左辺が単一値の場合のサポートも追加
func evalFilterOperation(node *ast.InfixExpression, env *object.Environment) object.Object {
	if logger.IsLevelEnabled(mapFilterDebugLevel) {
		logger.Debug("filter演算子(?>)の処理を開始")
	}

	// 左辺値の評価
	left := Eval(node.Left, env)
	if left == nil {
		return createError("filterオペレーション: 左辺の評価結果がnilです")
	}
	if left.Type() == object.ERROR_OBJ {
		return left
	}
	
	// 配列か単一の値かを確認し、適切な処理を行う
	var elements []object.Object
	var isSingleValue bool
	
	if arrayObj, ok := left.(*object.Array); ok {
		// 配列の場合はその要素を使用
		elements = arrayObj.Elements
		isSingleValue = false
		logger.Debug("?> 左辺の評価結果: 配列 %s (タイプ: %s)", left.Inspect(), left.Type())
	} else {
		// 単一の値の場合は要素1つの配列として扱う
		elements = []object.Object{left}
		isSingleValue = true
		logger.Debug("?> 左辺の評価結果: 単一値 %s (タイプ: %s) を要素1つの配列として扱います", left.Inspect(), left.Type())
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
		resultElements := make([]object.Object, 0)
		for _, element := range elements {
			// 各要素に対して関数を適用
			result := evalPipelineWithCallExpression(element, right, env)
			
			// 結果がtruthyな場合のみ結果に含める
			if isTruthy(result) {
				resultElements = append(resultElements, element)
			}
		}
		
		// 単一値モードの場合、結果があれば元の値を、なければnullを返す
		if isSingleValue {
			if len(resultElements) > 0 {
				return left // 元の単一値を返す
			}
			return NULL
		}
		
		return &object.Array{Elements: resultElements}
	default:
		return createError("filter演算子の右辺が関数または識別子ではありません: %T", node.Right)
	}

	// 直接配列の各要素に対して処理を行う
	resultElements := make([]object.Object, 0)
	
	for _, elem := range elements {
		// 一時環境を作成し、🍕に要素をセット
		tempEnv := object.NewEnclosedEnvironment(env)
		tempEnv.Set("🍕", elem)
		
		// 現在の要素に対して適切な関数を選択・実行
		// 引数にはelemを含める
		args := []object.Object{elem}
		if funcArgs != nil {
			args = append(args, funcArgs...)
		}
		
		// 組み込み関数を確認
		if builtin, ok := Builtins[funcName]; ok {
			logger.Debug("ビルトイン関数 '%s' をフィルター操作で呼び出します", funcName)
			result := builtin.Fn(args...)
			if result == nil || result.Type() == object.ERROR_OBJ {
				return result
			}
			
			// 結果がtruthyな場合のみ結果に含める
			if isTruthy(result) {
				resultElements = append(resultElements, elem)
			}
			continue
		}
		
		// 名前付き関数のメカニズムを使用して条件付き関数の正しい選択を行う
		// ここが重要な変更点：条件付き関数を明示的に処理する
		logger.Debug("フィルター操作: 要素 %s に対して条件付き関数メカニズムで関数 %s を呼び出し", elem.Inspect(), funcName)
		
		// 環境から全ての関数候補を取得
		functions := env.GetAllFunctionsByName(funcName)
		if len(functions) == 0 {
			logger.Debug("関数 '%s' が見つかりません", funcName)
			return createError("関数 '%s' が見つかりません", funcName)
		}
		
		// 条件付き関数を優先的に探す
		var result object.Object
		var defaultFunc *object.Function
		
		// 現在の要素を処理する一時環境を更新
		tempEnv.Set("🍕", elem)
		
		// まず条件付き関数を試す
		for _, fn := range functions {
			if fn.Condition != nil {
				// 条件式の評価
				isTrue, _ := evalConditionalExpression(fn, args, tempEnv)
				if isTrue {
					// 条件に一致した場合、この関数を使用
					logger.Debug("条件に一致する関数を見つけました: %s", funcName)
					result = applyFunctionWithPizza(fn, args)
					if result != nil && isTruthy(result) {
						// 条件に一致した関数の結果がtrue的な場合、元の要素を含める
						resultElements = append(resultElements, elem)
						break
					}
				}
			} else {
				// デフォルト関数を保存
				defaultFunc = fn
			}
		}
		
		// まだ条件付き関数で処理できていない場合はデフォルト関数を使う
		if result == nil && defaultFunc != nil {
			logger.Debug("デフォルト関数を使用します: %s", funcName)
			result = applyFunctionWithPizza(defaultFunc, args)
			if result != nil && isTruthy(result) {
				resultElements = append(resultElements, elem)
			}
		}
	}
	
	// 単一値モードの場合、結果があれば元の値を、なければnullを返す
	if isSingleValue {
		if len(resultElements) > 0 {
			return left // 元の単一値を返す
		}
		return NULL
	}
	
	return &object.Array{Elements: resultElements}
}
