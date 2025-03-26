package evaluator

import (
	"fmt"
	"github.com/uncode/logger"
	"github.com/uncode/object"
)

// 組み込み関数のマップ
var Builtins map[string]*object.Builtin

// ビルトイン関数のロギングレベル
var builtinLogLevel = logger.LevelDebug

// デバッグログの有効/無効を設定
func SetBuiltinLogLevel(level logger.LogLevel) {
	builtinLogLevel = level
}

// GetBuiltinLogLevel はビルトイン関数のロギングレベルを取得
func GetBuiltinLogLevel() logger.LogLevel {
	return builtinLogLevel
}

// logIfEnabled はビルトイン関数のデバッグログを出力
func logIfEnabled(level logger.LogLevel, format string, args ...interface{}) {
	componentLevel := logger.GetComponentLevel(logger.ComponentBuiltin)
	if level <= componentLevel {
		// ComponentXXXの関数を使って出力
		switch level {
		case logger.LevelError:
			logger.ComponentError(logger.ComponentBuiltin, format, args...)
		case logger.LevelWarn:
			logger.ComponentWarn(logger.ComponentBuiltin, format, args...)
		case logger.LevelInfo:
			logger.ComponentInfo(logger.ComponentBuiltin, format, args...)
		case logger.LevelDebug:
			logger.ComponentDebug(logger.ComponentBuiltin, format, args...)
		case logger.LevelTrace:
			logger.ComponentTrace(logger.ComponentBuiltin, format, args...)
		default:
			// その他のレベルはそのまま出力
			logger.Debug(format, args...)
		}
	}
}

// init関数でBuiltinsを初期化して循環参照を避ける
func init() {
	Builtins = map[string]*object.Builtin{}

	// 各カテゴリの組み込み関数を登録
	registerMathBuiltins()
	registerStringBuiltins()
	registerArrayBuiltins()
	registerTypeBuiltins()
	registerIOBuiltins()
}

// 組み込み関数の型情報を取得する関数
// GetBuiltinReturnType は組み込み関数の戻り値の型を返す
func GetBuiltinReturnType(name string) object.ObjectType {
	if builtin, ok := Builtins[name]; ok {
		return builtin.ReturnType
	}
	return object.NULL_OBJ
}

// GetBuiltinParamTypes は組み込み関数のパラメータの型を返す
func GetBuiltinParamTypes(name string) []object.ObjectType {
	if builtin, ok := Builtins[name]; ok {
		return builtin.ParamTypes
	}
	return nil
}

// createError はエラーオブジェクトを作成するヘルパー関数
func createError(format string, args ...interface{}) *object.Error {
	errMsg := fmt.Sprintf(format, args...)
	logIfEnabled(logger.LevelError, "エラー: %s", errMsg)
	return &object.Error{Message: errMsg}
}
