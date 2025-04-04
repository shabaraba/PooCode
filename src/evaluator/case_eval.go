package evaluator

import (
	"github.com/uncode/ast"
	"github.com/uncode/logger"
	"github.com/uncode/object"
)

// case文評価のデバッグレベル
var caseDebugLevel = logger.LevelOff

// SetCaseDebugLevel はcase文評価のデバッグレベルを設定
func SetCaseDebugLevel(level logger.LogLevel) {
	caseDebugLevel = level
	logger.Debug("case文評価のデバッグレベルを設定: %s", logger.LevelNames[level])
}

// isCaseDebugEnabled はデバッグが有効かを返す
func isCaseDebugEnabled() bool {
	return caseDebugLevel > logger.LevelOff
}

// logCaseDebug はデバッグログを出力
func logCaseDebug(format string, args ...interface{}) {
	if isCaseDebugEnabled() {
		logger.Log(caseDebugLevel, format, args...)
	}
}

// evalCaseStatement はcase文を評価
func evalCaseStatement(node *ast.CaseStatement, env *object.Environment) object.Object {
	logCaseDebug("case文の評価を開始: %s", node.Condition.String())
	
	// 🍕変数の存在確認と取得
	pizzaVal, ok := getPizzaValueFromEnv(env)
	if !ok {
		logCaseDebug("case文の評価中: 🍕変数が見つかりません")
		return createError("case文の評価中に🍕変数が見つかりませんでした")
	}
	
	logCaseDebug("case文の評価: 条件=%s, 🍕値=%s", 
		node.Condition.String(), pizzaVal.Inspect())
	
	// 条件式評価中のフラグを設定
	if currentFunction != nil {
		currentFunction.Condition = node.Condition
		// 評価後に元に戻すように後始末
		defer func() {
			currentFunction.Condition = nil
		}()
	}
	
	// 条件式を評価
	condition := Eval(node.Condition, env)
	if isError(condition) {
		logCaseDebug("case文の条件評価でエラー: %s", condition.Inspect())
		// エラーを返して、evalBlockStatementでハンドリングする
		return condition
	}
	
	// 条件式の結果を詳細にログ
	logCaseDebug("条件評価結果: タイプ=%s, 値=%s, isTruthy=%v", 
		condition.Type(), condition.Inspect(), isTruthy(condition))
	
	// 条件が真の場合、ブロックを実行
	if isTruthy(condition) {
		logCaseDebug("条件が真: ブロックを実行")
		if node.Body != nil {
			result := evalBlockStatement(node.Body, env)
			logCaseDebug("case文のブロック評価結果: %s", result.Inspect())
			return result
		} else if node.Consequence != nil {
			result := evalBlockStatement(node.Consequence, env)
			logCaseDebug("case文の結果ブロック評価結果: %s", result.Inspect())
			return result
		}
		logCaseDebug("警告: case文に実行可能なブロックがありません")
		return NullObj
	}
	
	// 条件が偽の場合
	logCaseDebug("条件が偽: 次のcaseへ")
	return NullObj
}

// evalDefaultCaseStatement はdefault文を評価
func evalDefaultCaseStatement(node *ast.DefaultCaseStatement, env *object.Environment) object.Object {
	logCaseDebug("default文の評価を開始")
	// 条件チェックなし、常にブロックを実行
	result := evalBlockStatement(node.Body, env)
	logCaseDebug("default文の評価結果: %s", result.Inspect())
	return result
}

// 🍕変数の取得補助関数
func getPizzaValueFromEnv(env *object.Environment) (object.Object, bool) {
	if obj, ok := env.Get("🍕"); ok {
		logCaseDebug("環境から🍕値を取得: %s", obj.Inspect())
		return obj, true
	}
	
	// 現在の関数からの取得を試みる
	if currentFunction != nil {
		if pizzaVal := currentFunction.GetPizzaValue(); pizzaVal != nil {
			logCaseDebug("現在の関数から🍕値を取得: %s", pizzaVal.Inspect())
			return pizzaVal, true
		}
	}
	
	logCaseDebug("環境に🍕値が設定されていません")
	return nil, false
}

// エッジケース対応のヘルパー関数
func checkCaseConditionSafety(condition object.Object) (bool, object.Object) {
	// NULL値のチェック
	if condition == NullObj {
		logCaseDebug("条件がNULL: 偽として評価")
		return false, NullObj
	}
	
	// エラー値のチェック
	if condition.Type() == object.ERROR_OBJ {
		logCaseDebug("条件がエラー: 評価中止")
		return false, condition
	}
	
	return true, nil
}
