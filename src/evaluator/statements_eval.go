package evaluator

import (
	"github.com/uncode/ast"
	"github.com/uncode/config"
	"github.com/uncode/logger"
	"github.com/uncode/object"
)

// evalProgram はプログラムを評価する
func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object = NullObj

	// プログラムが空の場合はNULLを返す
	if program == nil || len(program.Statements) == 0 {
		return NullObj
	}
	
	// 事前関数登録を実行（設定が有効な場合のみ）
	if config.GlobalConfig.PreregisterFunctions {
		logger.Debug("プログラム評価前に関数の事前登録を実行します")
		PreregisterFunctions(program, env)
	} else {
		logger.Debug("関数の事前登録はスキップされました（設定が無効です）")
	}
	
	for _, statement := range program.Statements {
		if statement == nil {
			continue
		}
		result = Eval(statement, env)
	}
	
	return result
}

// evalBlockStatement はブロック文を評価する
func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object = NullObj
	
	// デバッグ出力
	logger.Debug("ブロック文の評価を開始します。%d 個のステートメント", len(block.Statements))
	logCaseDebug("ブロック文の評価開始: %d 個のステートメント", len(block.Statements))
	
	// caseステートメントの処理用変数
	var caseEvaluated bool = false      // いずれかのcase文が真となったかを追跡
	var defaultCase *ast.DefaultCaseStatement = nil  // defaultケースの保存用
	
	// 各ステートメントを順番に評価
	for i, statement := range block.Statements {
		if statement == nil {
			logger.Debug("  ステートメント %d は nil です。スキップします", i)
			continue
		}
		
		logger.Debug("  ステートメント %d を評価: %T", i, statement)
		logCaseDebug("ステートメント %d を評価: %T", i, statement)
		
		// case文の処理
		switch stmt := statement.(type) {
		case *ast.CaseStatement:
			// すでにcaseが評価済みなら続行
			if caseEvaluated {
				logger.Debug("  すでにマッチしたcaseがあるため、このcase文をスキップします")
				logCaseDebug("マッチング済みのため case文をスキップ: %s", stmt.Condition.String())
				continue
			}
			
			logger.Debug("  case文を評価します: %s", stmt.Condition.String())
			logCaseDebug("case文の評価: %s", stmt.Condition.String())
			
			// case文の条件を評価
			caseResult := evalCaseStatement(stmt, env)
			
			// エラーチェック
			if isError(caseResult) {
				logger.Debug("  case文の評価でエラーが発生しました: %s", caseResult.Inspect())
				logCaseDebug("case文の評価エラー: %s", caseResult.Inspect())
				return caseResult
			}
			
			// NULLの場合は条件が一致しなかったので続行
			if caseResult == NullObj {
				logCaseDebug("case文の条件が一致しませんでした: %s", stmt.Condition.String())
				continue
			}
			
			// 条件に一致したcase文を見つけた
			logger.Debug("  マッチするcase文を見つけました: %s", stmt.Condition.String())
			logCaseDebug("マッチするcase文を発見: %s - 結果: %s", 
				stmt.Condition.String(), caseResult.Inspect())
			
			result = caseResult
			caseEvaluated = true
			
			// Case文マッチ後のエラーもしくはリターン値の場合は即時リターン
			if result.Type() == object.ERROR_OBJ {
				logCaseDebug("case文の評価結果がエラーのため即時リターン: %s", result.Inspect())
				return result
			}
			
			if returnObj, ok := result.(*object.ReturnValue); ok {
				logCaseDebug("case文の評価結果がリターン値のため即時リターン: %s", returnObj.Inspect())
				return returnObj
			}
			
		case *ast.DefaultCaseStatement:
			// defaultケースを保存（あとで使用）
			defaultCase = stmt
			logger.Debug("  default case文を検出。後で評価します")
			logCaseDebug("default case文を検出。すべてのcaseを確認後に評価します")
			continue
			
		default:
			// 通常の文の評価
			result = Eval(statement, env)
			
			// ReturnValue（関数からの戻り値）が検出された場合は評価を中止して戻る
			if returnValue, ok := result.(*object.ReturnValue); ok {
				logger.Debug("  ReturnValue が検出されました: %s", returnValue.Inspect())
				return returnValue
			}
			
			// ErrorValue が検出された場合も評価を中止して戻る
			if isError(result) {
				logger.Debug("  Error が検出されました: %s", result.Inspect())
				return result
			}
			
			// 代入文の場合、PooLiteralへの代入は特別な意味を持つ（関数からの戻り値）
			if assignStmt, ok := statement.(*ast.AssignStatement); ok {
				if _, ok := assignStmt.Value.(*ast.PooLiteral); ok {
					logger.Debug("  💩への代入を検出しました - 戻り値として扱います")
					// 左辺の値を取得
					leftVal := Eval(assignStmt.Left, env)
					if isError(leftVal) {
						logger.Debug("  💩への代入で左辺の評価エラー: %s", leftVal.Inspect())
						return leftVal
					}
					return &object.ReturnValue{Value: leftVal}
				}
			}
		}
	}
	
	// どのcaseにも一致せず、defaultケースがある場合
	if !caseEvaluated && defaultCase != nil {
		logger.Debug("  マッチするcaseが見つからなかったため、default caseを評価します")
		logCaseDebug("マッチするcaseが見つからず、default caseを評価します")
		
		result = evalDefaultCaseStatement(defaultCase, env)
		
		// default case評価後のエラーまたはリターン値チェック
		if isError(result) {
			logger.Debug("  default caseの評価でエラーが発生しました: %s", result.Inspect())
			logCaseDebug("default caseの評価でエラー: %s", result.Inspect())
			return result
		}
		
		if returnObj, ok := result.(*object.ReturnValue); ok {
			logger.Debug("  default caseからreturn値を検出: %s", returnObj.Inspect())
			logCaseDebug("default caseからreturn値を検出: %s", returnObj.Inspect())
			return returnObj
		}
	}
	
	logger.Debug("ブロック文の評価を完了しました。最終結果: %s", result.Inspect())
	logCaseDebug("ブロック文の評価完了。結果: %s", result.Inspect())
	return result
}

// これらの関数は case_eval.go に移動しました
