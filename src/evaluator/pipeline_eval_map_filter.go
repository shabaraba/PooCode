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
	if logger.IsLevelEnabled(mapFilterDebugLevel) {
		logger.Log(mapFilterDebugLevel, "mapオペレーション(+>)の評価開始")
	}

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
		if logger.IsLevelEnabled(mapFilterDebugLevel) {
			logger.Log(mapFilterDebugLevel, "+> 左の評価結果: 配列 %s (型: %s)", left.Inspect(), left.Type())
		}
	} else {
		// 単一の値の場合は 1個の配列として扱う
		elements = []object.Object{left}
		isSingleValue = true
		if logger.IsLevelEnabled(mapFilterDebugLevel) {
			logger.Log(mapFilterDebugLevel, "+> 左の評価結果: 単一値 %s (型: %s) を 1個の配列として扱います", left.Inspect(), left.Type())
		}
	}

	// 右側の評価：関数または関数呼び出し
	var funcName string
	var funcArgs []object.Object

	switch right := node.Right.(type) {
	case *ast.Identifier:
		// 識別子の場合、関数として扱う
		if logger.IsLevelEnabled(mapFilterDebugLevel) {
			logger.Log(mapFilterDebugLevel, "右は識別子: %s", right.Value)
		}
		funcName = right.Value
	case *ast.CallExpression:
		if logger.IsLevelEnabled(mapFilterDebugLevel) {
			logger.Log(mapFilterDebugLevel, "右は関数呼び出し式")
		}
		
		// 関数名取得
		if ident, ok := right.Function.(*ast.Identifier); ok {
			funcName = ident.Value
			if logger.IsLevelEnabled(mapFilterDebugLevel) {
				logger.Log(mapFilterDebugLevel, "関数名: %s", funcName)
			}
			
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
	
	// 環境内のすべての変数をデバッグログに出力
	if logger.IsLevelEnabled(mapFilterDebugLevel) {
		logger.Log(mapFilterDebugLevel, "マップ演算子: 環境内の変数一覧:")
		for k, v := range env.GetVariables() {
			if fn, ok := v.(*object.Function); ok {
				hasCondition := "なし"
				if fn.Condition != nil {
					hasCondition = "あり"
				}
				logger.Log(mapFilterDebugLevel, "  変数 '%s': 関数オブジェクト (条件=%s, アドレス=%p)", k, hasCondition, fn)
			} else {
				logger.Log(mapFilterDebugLevel, "  変数 '%s': %s", k, v.Type())
			}
		}
	}
	
	for _, elem := range elements {
		// 引数準備：要素自身を第一引数として、追加の引数も設定
		args := []object.Object{elem}
		if funcArgs != nil {
			args = append(args, funcArgs...)
		}
		
		if logger.IsLevelEnabled(mapFilterDebugLevel) {
			logger.Log(mapFilterDebugLevel, "マップ処理: 要素 %s に対して関数 '%s' を適用", elem.Inspect(), funcName)
		}
		
		// 組み込み関数処理
		if builtin, ok := Builtins[funcName]; ok {
			if logger.IsLevelEnabled(mapFilterDebugLevel) {
				logger.Log(mapFilterDebugLevel, "組み込み関数 '%s' をマップ処理で呼び出します", funcName)
			}
			result := builtin.Fn(args...)
			if result == nil || result.Type() == object.ERROR_OBJ {
				return result
			}
			resultElements = append(resultElements, result)
			continue
		}
		
		// 環境変数内に定義された関数を呼び出す
		// lookupFunctionInEnvironment は関数検索を統一化したヘルパー関数
		if logger.IsLevelEnabled(mapFilterDebugLevel) {
			logger.Log(mapFilterDebugLevel, "要素 %s に対してlookupFunctionInEnvironmentを呼び出し", elem.Inspect())
		}
		
		// 関数検索ロジックを強化
		// パイプライン演算子と同じ方法で関数を検索する
		function, err := lookupFunctionInEnvironment(env, funcName)
		if err != nil {
			if logger.IsLevelEnabled(mapFilterDebugLevel) {
				logger.Log(mapFilterDebugLevel, "関数 '%s' が見つかりませんでした: %s", funcName, err.Error())
			}
			return createError("関数 '%s' が見つかりません: %s", funcName, err.Error())
		}
		
		if logger.IsLevelEnabled(mapFilterDebugLevel) {
			logger.Log(mapFilterDebugLevel, "関数 '%s' が見つかりました (アドレス: %p)", funcName, function)
		}
		
		// 関数にピザ値をセットして実行
		result := applyFunctionWithPizza(function, args)
		
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
		logger.Log(mapFilterDebugLevel, "filter操作(?>)の評価開始")
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
		if logger.IsLevelEnabled(mapFilterDebugLevel) {
			logger.Log(mapFilterDebugLevel, "?> 左の評価結果: 配列 %s (型: %s)", left.Inspect(), left.Type())
		}
	} else {
		// 単一の値の場合は 1個の配列として扱う
		elements = []object.Object{left}
		isSingleValue = true
		if logger.IsLevelEnabled(mapFilterDebugLevel) {
			logger.Log(mapFilterDebugLevel, "?> 左の評価結果: 単一値 %s (型: %s) を 1個の配列として扱います", left.Inspect(), left.Type())
		}
	}

	// 右側の評価：関数または関数呼び出し
	var funcName string
	var funcArgs []object.Object

	switch right := node.Right.(type) {
	case *ast.Identifier:
		// 識別子の場合、関数として扱う
		if logger.IsLevelEnabled(mapFilterDebugLevel) {
			logger.Log(mapFilterDebugLevel, "右は識別子: %s", right.Value)
		}
		funcName = right.Value
	case *ast.CallExpression:
		// 関数呼び出しの場合
		if logger.IsLevelEnabled(mapFilterDebugLevel) {
			logger.Log(mapFilterDebugLevel, "右は関数呼び出し式")
		}
		if ident, ok := right.Function.(*ast.Identifier); ok {
			// 関数名取得
			funcName = ident.Value
			if logger.IsLevelEnabled(mapFilterDebugLevel) {
				logger.Log(mapFilterDebugLevel, "関数名: %s", funcName)
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
	
	// 環境内のすべての変数をデバッグログに出力
	if logger.IsLevelEnabled(mapFilterDebugLevel) {
		logger.Log(mapFilterDebugLevel, "フィルター演算子: 環境内の変数一覧:")
		for k, v := range env.GetVariables() {
			if fn, ok := v.(*object.Function); ok {
				hasCondition := "なし"
				if fn.Condition != nil {
					hasCondition = "あり"
				}
				logger.Log(mapFilterDebugLevel, "  変数 '%s': 関数オブジェクト (条件=%s, アドレス=%p)", k, hasCondition, fn)
			} else {
				logger.Log(mapFilterDebugLevel, "  変数 '%s': %s", k, v.Type())
			}
		}
	}
	
	for _, elem := range elements {
		// 引数準備
		args := []object.Object{elem}
		if funcArgs != nil {
			args = append(args, funcArgs...)
		}
		
		if logger.IsLevelEnabled(mapFilterDebugLevel) {
			logger.Log(mapFilterDebugLevel, "フィルタ処理: 要素 %s に対して関数 %s を適用", elem.Inspect(), funcName)
		}
		
		// 組み込み関数処理
		if builtin, ok := Builtins[funcName]; ok {
			if logger.IsLevelEnabled(mapFilterDebugLevel) {
				logger.Log(mapFilterDebugLevel, "組み込み関数 '%s' をフィルタ処理で呼び出します", funcName)
			}
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
		if logger.IsLevelEnabled(mapFilterDebugLevel) {
			logger.Log(mapFilterDebugLevel, "要素 %s に対してlookupFunctionInEnvironmentを呼び出し", elem.Inspect())
		}
		
		// 関数検索ロジックを強化
		// パイプライン演算子と同じ方法で関数を検索する
		function, err := lookupFunctionInEnvironment(env, funcName)
		if err != nil {
			if logger.IsLevelEnabled(mapFilterDebugLevel) {
				logger.Log(mapFilterDebugLevel, "関数 '%s' が見つかりませんでした: %s", funcName, err.Error())
			}
			return createError("関数 '%s' が見つかりません: %s", funcName, err.Error())
		}
		
		if logger.IsLevelEnabled(mapFilterDebugLevel) {
			logger.Log(mapFilterDebugLevel, "関数 '%s' が見つかりました (アドレス: %p)", funcName, function)
		}
		
		// 関数にピザ値をセットして実行
		result := applyFunctionWithPizza(function, args)
		
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

// lookupFunctionInEnvironment は環境内から関数を検索するヘルパー関数
// パイプラインとマップ演算子で関数検索ロジックを統一
func lookupFunctionInEnvironment(env *object.Environment, name string) (*object.Function, error) {
	// デバッグログ出力
	if logger.IsLevelEnabled(mapFilterDebugLevel) {
		logger.Log(mapFilterDebugLevel, "環境内から関数 '%s' を検索しています", name)
	}
	
	// ビルトイン関数をチェック
	if _, ok := Builtins[name]; ok {
		return nil, createEvalError("組み込み関数 '%s' は直接検索されません", name)
	}
	
	// 環境から同名のすべての関数を取得
	functions := env.GetAllFunctionsByName(name)
	
	if len(functions) == 0 {
		// 関数が見つからない場合のデバッグ情報
		if logger.IsLevelEnabled(mapFilterDebugLevel) {
			logger.Log(mapFilterDebugLevel, "関数 '%s' は環境内に見つかりませんでした", name)
			
			// 環境内のすべての関数名を表示
			logger.Log(mapFilterDebugLevel, "環境内の登録済み関数:")
			for k, v := range env.GetVariables() {
				if _, ok := v.(*object.Function); ok {
					logger.Log(mapFilterDebugLevel, "  - %s", k)
				}
			}
		}
		return nil, createEvalError("関数 '%s' が見つかりません", name)
	}
	
	// 条件付き関数とデフォルト関数を分類
	var defaultFunction *object.Function
	var conditionalFunctions []*object.Function
	
	for _, fn := range functions {
		if fn.Condition == nil {
			defaultFunction = fn
		} else {
			conditionalFunctions = append(conditionalFunctions, fn)
		}
	}
	
	// デフォルト関数があればそれを返す
	if defaultFunction != nil {
		if logger.IsLevelEnabled(mapFilterDebugLevel) {
			logger.Log(mapFilterDebugLevel, "関数 '%s' のデフォルト関数を返します", name)
		}
		return defaultFunction, nil
	}
	
	// デフォルト関数がなければ最初の条件付き関数を返す
	if len(conditionalFunctions) > 0 {
		if logger.IsLevelEnabled(mapFilterDebugLevel) {
			logger.Log(mapFilterDebugLevel, "関数 '%s' の条件付き関数を返します (デフォルト関数がないため)", name)
		}
		return conditionalFunctions[0], nil
	}
	
	// ここには到達しないはず（上記で関数が見つからない場合はエラーを返すため）
	return nil, createEvalError("関数 '%s' の処理中に予期しないエラーが発生しました", name)
}
