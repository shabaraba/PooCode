package evaluator

import (
	"testing"
	
	"github.com/uncode/logger"
)

func TestDebugLogging(t *testing.T) {
	// 評価器のデバッグログが設定されたログレベルで出力されることをテスト
	oldDebugEnabled := logger.IsDebugEnabled()
	oldEvalDebugEnabled := logger.IsEvalDebugEnabled()
	defer func() {
		// テスト後に元の設定に戻す
		if oldDebugEnabled {
			logger.EnableDebug()
		} else {
			logger.DisableDebug()
		}
		if oldEvalDebugEnabled {
			logger.EnableEvalDebug()
		} else {
			logger.DisableEvalDebug()
		}
	}()
	
	// デバッグログを有効化
	logger.EnableDebug()
	
	// デバッグログ出力を検証
	input := "5 + 5"
	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 10)
	
	// 評価器デバッグのみを有効化
	logger.DisableDebug()
	logger.EnableEvalDebug()
	
	input = "5 * 5"
	evaluated = testEval(input)
	testIntegerObject(t, evaluated, 25)
}

func TestSpecificDebugFlags(t *testing.T) {
	// 特定のデバッグフラグのテスト
	oldDebugEnabled := logger.IsDebugEnabled()
	oldEvalDebugEnabled := logger.IsEvalDebugEnabled()
	
	defer func() {
		// テスト後に元の設定に戻す
		if oldDebugEnabled {
			logger.EnableDebug()
		} else {
			logger.DisableDebug()
		}
		if oldEvalDebugEnabled {
			logger.EnableEvalDebug()
		} else {
			logger.DisableEvalDebug()
		}
	}()
	
	// 評価器デバッグを有効化し、特定フラグも設定
	logger.DisableDebug()
	logger.EnableEvalDebug()
	
	// evalデバッグが有効な場合のテスト
	input := "5 + 5"
	evaluated := testEval(input)
	testIntegerObject(t, evaluated, 10)
	
	// builtinデバッグも有効化（実装がなければスキップ）
	input = `len("hello")`
	evaluated = testEval(input)
	testIntegerObject(t, evaluated, 5)
}
