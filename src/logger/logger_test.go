package logger

import (
	"bytes"
	"strings"
	"testing"
)

// TestLogLevels はログレベルのフィルタリングをテストする
func TestLogLevels(t *testing.T) {
	// テスト用の出力バッファ
	var buf bytes.Buffer
	
	// カラー出力なしのロガーを作成
	logger := NewLogger(
		WithWriter(&buf),
		WithColor(false),
		WithTime(false),
		WithLevel(LevelInfo),
	)
	
	// 各レベルでログを出力
	logger.Error("エラーメッセージ")
	logger.Warn("警告メッセージ")
	logger.Info("情報メッセージ")
	logger.Debug("デバッグメッセージ") // 出力されないはず
	logger.Trace("トレースメッセージ") // 出力されないはず
	
	output := buf.String()
	
	// Error, Warn, Infoは出力されるべき
	if !strings.Contains(output, "エラーメッセージ") {
		t.Errorf("Errorレベルのメッセージが出力されていません")
	}
	if !strings.Contains(output, "警告メッセージ") {
		t.Errorf("Warnレベルのメッセージが出力されていません")
	}
	if !strings.Contains(output, "情報メッセージ") {
		t.Errorf("Infoレベルのメッセージが出力されていません")
	}
	
	// Debug, Traceは出力されないべき
	if strings.Contains(output, "デバッグメッセージ") {
		t.Errorf("Debugレベルのメッセージが出力されています: %s", output)
	}
	if strings.Contains(output, "トレースメッセージ") {
		t.Errorf("Traceレベルのメッセージが出力されています: %s", output)
	}
}

// TestComponentLogging はコンポーネント別ログ出力をテストする
func TestComponentLogging(t *testing.T) {
	var buf bytes.Buffer
	
	logger := NewLogger(
		WithWriter(&buf),
		WithColor(false),
		WithTime(false),
		WithLevel(LevelInfo),
	)
	
	// パーサーコンポーネントのレベルをDebugに設定
	logger.SetComponentLevel(ComponentParser, LevelDebug)
	
	// 各コンポーネントでログを出力
	logger.ComponentInfo(ComponentLexer, "レキサー情報")
	logger.ComponentDebug(ComponentLexer, "レキサーデバッグ") // 出力されないはず
	logger.ComponentInfo(ComponentParser, "パーサー情報")
	logger.ComponentDebug(ComponentParser, "パーサーデバッグ") // 出力されるはず
	
	output := buf.String()
	
	// Infoレベルは両方のコンポーネントで出力されるべき
	if !strings.Contains(output, "レキサー情報") {
		t.Errorf("レキサーのInfoレベルのメッセージが出力されていません")
	}
	if !strings.Contains(output, "パーサー情報") {
		t.Errorf("パーサーのInfoレベルのメッセージが出力されていません")
	}
	
	// Debugレベルはパーサーのみ出力されるべき
	if strings.Contains(output, "レキサーデバッグ") {
		t.Errorf("レキサーのDebugレベルのメッセージが出力されています: %s", output)
	}
	if !strings.Contains(output, "パーサーデバッグ") {
		t.Errorf("パーサーのDebugレベルのメッセージが出力されていません")
	}
}

// TestSpecialLogLevels は特殊ログレベルの制御をテストする
func TestSpecialLogLevels(t *testing.T) {
	var buf bytes.Buffer
	
	logger := NewLogger(
		WithWriter(&buf),
		WithColor(false),
		WithTime(false),
		WithLevel(LevelOff), // 通常のログは全て無効
	)
	
	// 特殊ログレベルを有効化
	logger.SetSpecialLevelEnabled(LevelTypeInfo, true)
	
	// 各レベルでログを出力
	logger.Error("エラーメッセージ") // 出力されないはず
	logger.Info("情報メッセージ") // 出力されないはず
	logger.TypeInfo("型情報メッセージ") // 出力されるはず
	logger.EvalDebug("評価デバッグ") // 出力されないはず
	
	output := buf.String()
	
	// 通常のログは出力されないべき
	if strings.Contains(output, "エラーメッセージ") {
		t.Errorf("Errorレベルのメッセージが出力されています: %s", output)
	}
	if strings.Contains(output, "情報メッセージ") {
		t.Errorf("Infoレベルのメッセージが出力されています: %s", output)
	}
	
	// TypeInfoは有効なので出力されるべき
	if !strings.Contains(output, "型情報メッセージ") {
		t.Errorf("TypeInfoレベルのメッセージが出力されていません")
	}
	
	// EvalDebugは無効なので出力されないべき
	if strings.Contains(output, "評価デバッグ") {
		t.Errorf("EvalDebugレベルのメッセージが出力されています: %s", output)
	}
}

// TestLoggerEnable はロガーの有効/無効化をテストする
func TestLoggerEnable(t *testing.T) {
	var buf bytes.Buffer
	
	logger := NewLogger(
		WithWriter(&buf),
		WithColor(false),
		WithTime(false),
		WithLevel(LevelInfo),
	)
	
	// 有効状態でログを出力
	logger.Info("有効時のメッセージ")
	
	// ロガーを無効化
	logger.Disable()
	logger.Info("無効時のメッセージ")
	
	// ロガーを再度有効化
	logger.Enable()
	logger.Info("再有効時のメッセージ")
	
	output := buf.String()
	
	// 有効時と再有効時のメッセージは出力されるべき、無効時は出力されないべき
	if !strings.Contains(output, "有効時のメッセージ") {
		t.Errorf("有効時のメッセージが出力されていません")
	}
	if strings.Contains(output, "無効時のメッセージ") {
		t.Errorf("無効時のメッセージが出力されています: %s", output)
	}
	if !strings.Contains(output, "再有効時のメッセージ") {
		t.Errorf("再有効時のメッセージが出力されていません")
	}
}

// TestLoggerFormat はログフォーマットをテストする
func TestLoggerFormat(t *testing.T) {
	var buf bytes.Buffer
	
	// タイムスタンプなしのロガー
	logger := NewLogger(
		WithWriter(&buf),
		WithColor(false),
		WithTime(false),
	)
	
	logger.Info("テストメッセージ")
	
	output := buf.String()
	
	// フォーマットは "[INFO] テストメッセージ" のようになるはず
	expected := "[INFO] テストメッセージ"
	if !strings.Contains(output, expected) {
		t.Errorf("ログフォーマットが正しくありません。\n期待値: %s\n実際: %s", expected, output)
	}
}

// TestLoggerParseLevel はログレベル文字列解析をテストする
func TestLoggerParseLevel(t *testing.T) {
	tests := []struct {
		input    string
		expected LogLevel
	}{
		{"OFF", LevelOff},
		{"ERROR", LevelError},
		{"WARN", LevelWarn},
		{"INFO", LevelInfo},
		{"DEBUG", LevelDebug},
		{"TRACE", LevelTrace},
		{"TYPE", LevelTypeInfo},
		{"EVAL", LevelEvalDebug},
		{"UNKNOWN", LevelInfo}, // 不明な場合はInfoになるはず
	}
	
	for _, tt := range tests {
		level := ParseLogLevel(tt.input)
		if level != tt.expected {
			t.Errorf("ParseLogLevel(%q) = %v, 期待値 %v", tt.input, level, tt.expected)
		}
	}
}
