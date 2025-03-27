package evaluator

import (
	"fmt"
	"strings"
	"github.com/uncode/logger"
	"github.com/uncode/object"
)

// log variables for debugging
var builtinLogLevel = logger.LevelTrace

// logIfEnabled logs a message if logging is enabled for the given level
func logIfEnabled(level logger.LogLevel, format string, args ...interface{}) {
	if logger.IsSpecialLevelEnabled(level) || logger.GetComponentLevel(logger.ComponentBuiltin) >= level {
		logger.ComponentDebug(logger.ComponentBuiltin, format, args...)
	}
}

// registerArrayBuiltins は配列関連の組み込み関数を登録する
func registerArrayBuiltins() {
	// 配列を連結して文字列にする関数
	Builtins["join"] = &object.Builtin{
		Name: "join",
		Fn: func(args ...object.Object) object.Object {
			logIfEnabled(builtinLogLevel, "join関数が呼び出されました: 引数数=%d", len(args))
			
			if len(args) \!= 2 {
				logger.ComponentError(logger.ComponentBuiltin, "join関数は2つの引数が必要です: %d個与えられました", len(args))
				return createError("join関数は2つの引数が必要です: %d個与えられました", len(args))
			}
