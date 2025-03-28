package evaluator

import (
	"github.com/uncode/ast"
	"github.com/uncode/logger"
	"github.com/uncode/object"
)

// conditionDebugLevel は条件式評価のデバッグレベルを保持します
var conditionDebugLevel = logger.LevelOff

// SetConditionDebugLevel は条件式評価のデバッグレベルを設定します
func SetConditionDebugLevel(level logger.LogLevel) {
	conditionDebugLevel = level
	logger.Debug("条件式評価のデバッグレベルを設定: %s", logger.LevelNames[level])
}

// isConditionDebugEnabled はデバッグが有効かどうかを返します
func isConditionDebugEnabled() bool {
	return conditionDebugLevel > logger.LevelOff
}

// logConditionDebug はデバッグログを出力します
func logConditionDebug(format string, args ...interface{}) {
	if isConditionDebugEnabled() {
		logger.Log(conditionDebugLevel, format, args...)
	}
}

// evalConditionalExpression は条件式を評価します
// 関数オブジェクトの🍕メンバーを優先的に使用するよう修正
func evalConditionalExpression(fn *object.Function, args []object.Object, env *object.Environment) (bool, object.Object) {
	// 条件式が存在しない場合はtrueを返す
	if fn.Condition == nil {
		logConditionDebug("条件式が存在しないため、常にtrueとして評価します")
		return true, &object.Boolean{Value: true}
	}

	// 条件式の詳細な情報をログ出力
	logger.Log(logger.LevelInfo, "条件式評価: 対象値=%s, 条件=%v", args[0].Inspect(), fn.Condition)
	
	// 条件式評価のために独立した環境を作成
	// 重要: 外部環境を継承するようにします - 外部の環境にアクセスできるように
	condEnv := object.NewEnclosedEnvironment(env)
	
	// 🍕メンバーの設定（重要な改善点）
	if len(args) > 0 {
		// 1. 関数オブジェクトに🍕値を直接設定
		logConditionDebug("関数オブジェクトに🍕値を設定: %s (%s)", args[0].Inspect(), args[0].Type())
		fn.SetPizzaValue(args[0])
		
		// 2. 環境にも🍕値を設定（互換性維持のため）
		logConditionDebug("条件評価環境にも🍕値を設定: %s", args[0].Inspect())
		condEnv.Set("🍕", args[0])
	} else {
		logConditionDebug("引数が指定されていないため、🍕値は設定されません")
	}
	
	// 条件式評価前のデバッグ出力
	if isConditionDebugEnabled() {
		// 条件式の詳細表示
		logConditionDebug("-------- 条件式の詳細 --------")
		if infixExpr, ok := fn.Condition.(*ast.InfixExpression); ok {
			logConditionDebug("条件式タイプ: 中置式")
			logConditionDebug("  演算子: %s", infixExpr.Operator)
			logConditionDebug("  左辺: %T - %v", infixExpr.Left, infixExpr.Left)
			logConditionDebug("  右辺: %T - %v", infixExpr.Right, infixExpr.Right)
		} else {
			logConditionDebug("条件式タイプ: %T", fn.Condition)
		}
		
		// 環境内の🍕値の状態表示
		if pizzaVal, ok := condEnv.Get("🍕"); ok {
			logConditionDebug("環境内の🍕変数: タイプ=%s, 値=%s", pizzaVal.Type(), pizzaVal.Inspect())
		} else {
			logConditionDebug("環境内の🍕変数: 未設定")
		}
		
		// 関数オブジェクト内の🍕値の状態表示
		if pizzaVal := fn.GetPizzaValue(); pizzaVal != nil {
			logConditionDebug("関数オブジェクト内の🍕値: タイプ=%s, 値=%s", pizzaVal.Type(), pizzaVal.Inspect())
		} else {
			logConditionDebug("関数オブジェクト内の🍕値: nil")
		}
		logConditionDebug("------------------------------")
	}

	// 条件式評価前に、evalInfixExpression が関数オブジェクトから🍕値を取得できるように
	// 現在の関数コンテキストを設定
	prevFunction := currentFunction
	currentFunction = fn
	
	// 条件式を評価
	condResult := Eval(fn.Condition, condEnv)
	
	// 関数コンテキストを元に戻す
	currentFunction = prevFunction
	
	if condResult.Type() == object.ERROR_OBJ {
		logConditionDebug("条件評価でエラーが発生しました: %s", condResult.Inspect())
		return false, condResult
	}
	
	// 条件式の評価結果を解釈
	var isTrue bool
	
	if condResult.Type() == object.BOOLEAN_OBJ {
		isTrue = condResult.(*object.Boolean).Value
		logConditionDebug("条件式の真偽値（Boolean型）: %v", isTrue)
	} else {
		isTrue = isTruthy(condResult)
		logConditionDebug("条件式の真偽値（非Boolean型）: %v", isTrue)
	}
	
	// 条件式評価直後に評価結果を出力
	logger.Log(logger.LevelInfo, "条件式評価結果: 対象値=%s, 評価結果=%v", args[0].Inspect(), isTrue)
	
	return isTrue, condResult
}
