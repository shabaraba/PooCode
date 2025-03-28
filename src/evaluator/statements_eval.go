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
	
	for i, statement := range block.Statements {
		logger.Debug("  ステートメント %d を評価: %T", i, statement)
		result = Eval(statement, env)
		
		// ReturnValue（関数からの戻り値）が検出された場合は評価を中止して戻る
		if returnValue, ok := result.(*object.ReturnValue); ok {
			logger.Debug("  ReturnValue が検出されました: %s", returnValue.Inspect())
			return returnValue
		}
		
		// ErrorValue が検出された場合も評価を中止して戻る
		if result.Type() == object.ERROR_OBJ {
			logger.Debug("  Error が検出されました: %s", result.Inspect())
			return result
		}
		
		// 代入文の場合、PooLiteralへの代入は特別な意味を持つ
		if assignStmt, ok := statement.(*ast.AssignStatement); ok {
			if _, ok := assignStmt.Value.(*ast.PooLiteral); ok {
				logger.Debug("  💩への代入を検出しました - 戻り値として扱います")
				// 左辺の値を取得
				leftVal := Eval(assignStmt.Left, env)
				if leftVal.Type() == object.ERROR_OBJ {
					return leftVal
				}
				return &object.ReturnValue{Value: leftVal}
			}
		}
	}
	
	logger.Debug("ブロック文の評価を完了しました。最終結果: %s", result.Inspect())
	return result
}
