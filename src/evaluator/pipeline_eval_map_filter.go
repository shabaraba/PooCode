package evaluator

import (
	"github.com/uncode/ast"
	"github.com/uncode/logger"
	"github.com/uncode/object"
)

// evalInfixExpressionWithNode は中置式を評価する
func evalInfixExpressionWithNode(node *ast.InfixExpression, env *object.Environment) object.Object {
	logger.Debug("中置式を評価します: %s", node.Operator)

	switch node.Operator {
	case "+>": // map演算子
		logger.Debug("map パイプ演算子 (+>) を検出しました")
		// map関数の処理を実行
		return evalMapOperation(node, env)
	case "?>": // filter演算子
		logger.Debug("filter パイプ演算子 (?>) を検出しました")
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
	logger.Debug("map演算子(+>)の処理を開始")

	// 左辺値の評価（通常は配列）
	left := Eval(node.Left, env)
	if left.Type() == object.ERROR_OBJ {
		return left
	}

	logger.Debug("+> 左辺の評価結果: %s (タイプ: %s)", left.Inspect(), left.Type())

	// 右辺値の評価（関数または関数呼び出し）
	var funcObj object.Object
	var funcArgs []object.Object

	switch right := node.Right.(type) {
	case *ast.Identifier:
		// 識別子の場合、関数名として扱う
		logger.Debug("右辺が識別子: %s", right.Value)
		funcNameObj, exists := env.Get(right.Value)
		if !exists {
			return createError("関数 '%s' が見つかりません", right.Value)
		}
		funcObj = funcNameObj
	case *ast.CallExpression:
		// 関数呼び出しの場合
		logger.Debug("右辺が関数呼び出し式")
		if ident, ok := right.Function.(*ast.Identifier); ok {
			// 関数名を取得
			logger.Debug("関数名: %s", ident.Value)
			funcNameObj, exists := env.Get(ident.Value)
			if !exists {
				return createError("関数 '%s' が見つかりません", ident.Value)
			}
			funcObj = funcNameObj

			// 引数の評価
			funcArgs = evalExpressions(right.Arguments, env)
			if len(funcArgs) > 0 && funcArgs[0].Type() == object.ERROR_OBJ {
				return funcArgs[0]
			}
		} else {
			return createError("関数呼び出し式の関数部分が識別子ではありません: %T", right.Function)
		}
	default:
		return createError("map演算子の右辺が関数または識別子ではありません: %T", node.Right)
	}

	// map関数の呼び出し
	logger.Debug("map関数をビルトイン関数として呼び出し")
	mapBuiltin, ok := Builtins["map"]
	if !ok {
		return createError("map関数がビルトイン関数として見つかりません")
	}

	// 引数リストの構築: [配列, 関数, 追加引数...]
	var mapArgs []object.Object
	mapArgs = append(mapArgs, left)       // 第1引数: 配列
	mapArgs = append(mapArgs, funcObj)    // 第2引数: 関数
	mapArgs = append(mapArgs, funcArgs...) // 追加引数

	// map関数の実行
	logger.Debug("map関数実行: 引数数=%d", len(mapArgs))
	return mapBuiltin.Fn(mapArgs...)
}

// evalFilterOperation はfilter演算子(?>)を処理する
func evalFilterOperation(node *ast.InfixExpression, env *object.Environment) object.Object {
	logger.Debug("filter演算子(?>)の処理を開始")

	// 左辺値の評価（通常は配列）
	left := Eval(node.Left, env)
	if left.Type() == object.ERROR_OBJ {
		return left
	}

	logger.Debug("?> 左辺の評価結果: %s (タイプ: %s)", left.Inspect(), left.Type())

	// 右辺値の評価（関数または関数呼び出し）
	var funcObj object.Object
	var funcArgs []object.Object

	switch right := node.Right.(type) {
	case *ast.Identifier:
		// 識別子の場合、関数名として扱う
		logger.Debug("右辺が識別子: %s", right.Value)
		funcNameObj, exists := env.Get(right.Value)
		if !exists {
			return createError("関数 '%s' が見つかりません", right.Value)
		}
		funcObj = funcNameObj
	case *ast.CallExpression:
		// 関数呼び出しの場合
		logger.Debug("右辺が関数呼び出し式")
		if ident, ok := right.Function.(*ast.Identifier); ok {
			// 関数名を取得
			logger.Debug("関数名: %s", ident.Value)
			funcNameObj, exists := env.Get(ident.Value)
			if !exists {
				return createError("関数 '%s' が見つかりません", ident.Value)
			}
			funcObj = funcNameObj

			// 引数の評価
			funcArgs = evalExpressions(right.Arguments, env)
			if len(funcArgs) > 0 && funcArgs[0].Type() == object.ERROR_OBJ {
				return funcArgs[0]
			}
		} else {
			return createError("関数呼び出し式の関数部分が識別子ではありません: %T", right.Function)
		}
	default:
		return createError("filter演算子の右辺が関数または識別子ではありません: %T", node.Right)
	}

	// filter関数の呼び出し
	logger.Debug("filter関数をビルトイン関数として呼び出し")
	filterBuiltin, ok := Builtins["filter"]
	if !ok {
		return createError("filter関数がビルトイン関数として見つかりません")
	}

	// 引数リストの構築: [配列, 関数, 追加引数...]
	var filterArgs []object.Object
	filterArgs = append(filterArgs, left)       // 第1引数: 配列
	filterArgs = append(filterArgs, funcObj)    // 第2引数: 関数
	filterArgs = append(filterArgs, funcArgs...) // 追加引数

	// filter関数の実行
	logger.Debug("filter関数実行: 引数数=%d", len(filterArgs))
	return filterBuiltin.Fn(filterArgs...)
}
