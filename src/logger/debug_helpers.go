package logger

// テスト用のヘルパー関数

// IsDebugEnabled は通常デバッグログが有効かどうかを返す
func IsDebugEnabled() bool {
	return GetLogger().GetComponentLevel(ComponentGlobal) >= LevelDebug
}

// IsEvalDebugEnabled は評価器専用デバッグログが有効かどうかを返す
func IsEvalDebugEnabled() bool {
	return GetLogger().IsSpecialLevelEnabled(LevelEvalDebug)
}

// IsParserDebugEnabled はパーサー専用デバッグログが有効かどうかを返す
func IsParserDebugEnabled() bool {
	return GetLogger().IsSpecialLevelEnabled(LevelParserDebug)
}

// EnableDebug は通常デバッグログを有効にする
func EnableDebug() {
	GetLogger().SetComponentLevel(ComponentGlobal, LevelDebug)
}

// DisableDebug は通常デバッグログを無効にする
func DisableDebug() {
	GetLogger().SetComponentLevel(ComponentGlobal, LevelInfo)
}

// EnableEvalDebug は評価器専用デバッグログを有効にする
func EnableEvalDebug() {
	GetLogger().SetSpecialLevelEnabled(LevelEvalDebug, true)
}

// DisableEvalDebug は評価器専用デバッグログを無効にする
func DisableEvalDebug() {
	GetLogger().SetSpecialLevelEnabled(LevelEvalDebug, false)
}

// EnableParserDebug はパーサー専用デバッグログを有効にする
func EnableParserDebug() {
	GetLogger().SetSpecialLevelEnabled(LevelParserDebug, true)
}

// DisableParserDebug はパーサー専用デバッグログを無効にする
func DisableParserDebug() {
	GetLogger().SetSpecialLevelEnabled(LevelParserDebug, false)
}
