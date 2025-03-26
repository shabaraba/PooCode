package logger

import (
	"fmt"
	"io"
	"os"
	"strings"
	"sync"
	"time"
)

// LogLevel はログレベルを表す型
type LogLevel int

const (
	// ログレベル定義
	LevelOff LogLevel = iota
	LevelError
	LevelWarn
	LevelInfo
	LevelDebug
	LevelTrace
	
	// 特殊デバッグ情報レベル
	LevelTypeInfo  // 型情報のみを表示
	LevelEvalDebug // 評価器専用デバッグ情報
)

// コンポーネント種別を表す型
type ComponentType string

const (
	// コンポーネント定義
	ComponentGlobal  ComponentType = "global"  // グローバル（デフォルト）
	ComponentLexer   ComponentType = "lexer"   // レキサー
	ComponentParser  ComponentType = "parser"  // パーサー
	ComponentEval    ComponentType = "eval"    // 評価器
	ComponentRuntime ComponentType = "runtime" // ランタイム
	ComponentBuiltin ComponentType = "builtin" // 組み込み関数
)

// LevelNames はログレベルと名前のマッピング
var LevelNames = map[LogLevel]string{
	LevelOff:      "OFF",
	LevelError:    "ERROR",
	LevelWarn:     "WARN",
	LevelInfo:     "INFO",
	LevelDebug:    "DEBUG",
	LevelTrace:    "TRACE",
	LevelTypeInfo: "TYPE",
	LevelEvalDebug: "EVAL",
}

var levelColors = map[LogLevel]string{
	LevelError:    "\033[31m", // 赤
	LevelWarn:     "\033[33m", // 黄
	LevelInfo:     "\033[32m", // 緑
	LevelDebug:    "\033[36m", // シアン
	LevelTrace:    "\033[35m", // マゼンタ
	LevelTypeInfo: "\033[34m", // 青
	LevelEvalDebug: "\033[33;1m", // 太字黄色
}

const (
	colorReset = "\033[0m"
)

// Logger はログを記録するための構造体
type Logger struct {
	globalLevel     LogLevel                // グローバルログレベル
	componentLevels map[ComponentType]LogLevel  // コンポーネント別ログレベル
	specialLevels   map[LogLevel]bool      // 特殊なログレベルの有効/無効状態
	writer          io.Writer
	fileWriter      io.Writer
	mu              sync.Mutex
	isEnabled       bool
	useColor        bool
	showTime        bool
}

var (
	defaultLogger *Logger
	once          sync.Once
)

// LoggerOption はロガーの設定用オプション関数型
type LoggerOption func(*Logger)

// WithLevel はグローバルログレベルを設定するオプション
func WithLevel(level LogLevel) LoggerOption {
	return func(l *Logger) {
		l.globalLevel = level
	}
}

// WithComponentLevel はコンポーネント別ログレベルを設定するオプション
func WithComponentLevel(component ComponentType, level LogLevel) LoggerOption {
	return func(l *Logger) {
		l.componentLevels[component] = level
	}
}

// WithWriter は出力先を設定するオプション
func WithWriter(w io.Writer) LoggerOption {
	return func(l *Logger) {
		l.writer = w
	}
}

// WithFileWriter はファイル出力先を設定するオプション
func WithFileWriter(w io.Writer) LoggerOption {
	return func(l *Logger) {
		l.fileWriter = w
	}
}

// WithColor はカラー出力を設定するオプション
func WithColor(useColor bool) LoggerOption {
	return func(l *Logger) {
		l.useColor = useColor
	}
}

// WithTime はタイムスタンプ表示を設定するオプション
func WithTime(showTime bool) LoggerOption {
	return func(l *Logger) {
		l.showTime = showTime
	}
}

// NewLogger は新しいロガーインスタンスを生成する
func NewLogger(options ...LoggerOption) *Logger {
	logger := &Logger{
		globalLevel:     LevelInfo,
		componentLevels: make(map[ComponentType]LogLevel),
		specialLevels:   make(map[LogLevel]bool),
		writer:          os.Stdout,
		fileWriter:      nil,
		isEnabled:       true,
		useColor:        true,
		showTime:        true,
	}
	
	// デフォルトのコンポーネントレベルを設定
	logger.componentLevels[ComponentGlobal] = LevelInfo
	logger.componentLevels[ComponentLexer] = LevelInfo
	logger.componentLevels[ComponentParser] = LevelInfo
	logger.componentLevels[ComponentEval] = LevelInfo
	logger.componentLevels[ComponentRuntime] = LevelInfo
	logger.componentLevels[ComponentBuiltin] = LevelInfo
	
	// デフォルトの特殊ログレベル設定（無効）
	logger.specialLevels[LevelTypeInfo] = false
	logger.specialLevels[LevelEvalDebug] = false
	
	// オプションを適用
	for _, option := range options {
		option(logger)
	}
	
	return logger
}

// GetLogger はシングルトンロガーインスタンスを返す
func GetLogger() *Logger {
	once.Do(func() {
		defaultLogger = NewLogger()
	})
	return defaultLogger
}

// SetLevel はグローバルログレベルを設定する
func (l *Logger) SetLevel(level LogLevel) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.globalLevel = level
	// グローバルレベルの変更を全コンポーネントに反映（ただし明示的に設定されたものは除く）
	l.componentLevels[ComponentGlobal] = level
}

// SetComponentLevel はコンポーネント別ログレベルを設定する
func (l *Logger) SetComponentLevel(component ComponentType, level LogLevel) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.componentLevels[component] = level
}

// GetComponentLevel はコンポーネント別ログレベルを取得する
func (l *Logger) GetComponentLevel(component ComponentType) LogLevel {
	l.mu.Lock()
	defer l.mu.Unlock()
	level, exists := l.componentLevels[component]
	if !exists {
		return l.componentLevels[ComponentGlobal] // デフォルトはグローバル設定
	}
	return level
}

// SetSpecialLevelEnabled は特殊ログレベルの有効/無効を設定する
func (l *Logger) SetSpecialLevelEnabled(level LogLevel, enabled bool) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.specialLevels[level] = enabled
}

// IsSpecialLevelEnabled は特殊ログレベルが有効かどうかを返す
func (l *Logger) IsSpecialLevelEnabled(level LogLevel) bool {
	l.mu.Lock()
	defer l.mu.Unlock()
	enabled, exists := l.specialLevels[level]
	if !exists {
		return false // デフォルトは無効
	}
	return enabled
}

// SetOutput はログの出力先を設定する
func (l *Logger) SetOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.writer = w
}

// SetFileOutput はファイルへの出力を設定する
func (l *Logger) SetFileOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.fileWriter = w
}

// EnableColor はカラー出力を有効にする
func (l *Logger) EnableColor() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.useColor = true
}

// DisableColor はカラー出力を無効にする
func (l *Logger) DisableColor() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.useColor = false
}

// EnableTimestamp はタイムスタンプ表示を有効にする
func (l *Logger) EnableTimestamp() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.showTime = true
}

// DisableTimestamp はタイムスタンプ表示を無効にする
func (l *Logger) DisableTimestamp() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.showTime = false
}

// Enable はロガーを有効にする
func (l *Logger) Enable() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.isEnabled = true
}

// Disable はロガーを無効にする
func (l *Logger) Disable() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.isEnabled = false
}

// formatLogMessage はログメッセージをフォーマットする
func (l *Logger) formatLogMessage(level LogLevel, format string, args ...interface{}) string {
	levelName := LevelNames[level]
	msg := fmt.Sprintf(format, args...)
	
	var builder strings.Builder
	
	// タイムスタンプを追加
	if l.showTime {
		timestamp := time.Now().Format("2006-01-02 15:04:05.000")
		builder.WriteString(fmt.Sprintf("[%s] ", timestamp))
	}
	
	// レベル名を追加（カラーあり/なし）
	if l.useColor {
		colorCode, hasColor := levelColors[level]
		if hasColor {
			builder.WriteString(fmt.Sprintf("%s[%s]%s ", colorCode, levelName, colorReset))
		} else {
			builder.WriteString(fmt.Sprintf("[%s] ", levelName))
		}
	} else {
		builder.WriteString(fmt.Sprintf("[%s] ", levelName))
	}
	
	// メッセージを追加
	builder.WriteString(msg)
	
	return builder.String()
}

// log はメッセージを指定されたレベルでログに記録する
func (l *Logger) log(level LogLevel, format string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	// 特殊ログレベルの判定（LevelTypeInfo, LevelEvalDebugなど）
	isSpecialLevel := level >= LevelTypeInfo
	
	// 特殊ログレベルの場合は専用の有効/無効設定を使用
	if isSpecialLevel {
		enabled, exists := l.specialLevels[level]
		if !exists || !enabled {
			// 特殊レベルが無効なら出力しない
			return
		}
	} else if !l.isEnabled || level > l.globalLevel {
		// 通常レベルの場合はグローバルログレベルに基づく
		return
	}

	// フォーマット済みのメッセージを生成
	formattedMsg := l.formatLogMessage(level, format, args...)
	
	// 標準出力に書き込み
	fmt.Fprintln(l.writer, formattedMsg)
	
	// ファイルにも書き込み（設定されている場合）
	if l.fileWriter != nil {
		fmt.Fprintln(l.fileWriter, formattedMsg)
	}
}

// logWithComponent はコンポーネント指定付きでログを記録する
func (l *Logger) logWithComponent(component ComponentType, level LogLevel, format string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	// 特殊ログレベルの判定（LevelTypeInfo, LevelEvalDebugなど）
	isSpecialLevel := level >= LevelTypeInfo
	
	if isSpecialLevel {
		// 特殊ログレベルの場合は専用の有効/無効設定を使用
		enabled, exists := l.specialLevels[level]
		if !exists || !enabled {
			// 特殊レベルが無効なら出力しない
			return
		}
	} else {
		// 通常レベルの場合はコンポーネントレベルに基づく
		// コンポーネントのログレベルを確認
		componentLevel, exists := l.componentLevels[component]
		if !exists {
			componentLevel = l.componentLevels[ComponentGlobal]
		}

		// コンポーネントのログレベルに基づきフィルタリング
		if !l.isEnabled || level > componentLevel {
			return
		}
	}

	// コンポーネント名を含むフォーマット済みのメッセージを生成
	msg := fmt.Sprintf(format, args...)
	levelName := LevelNames[level]
	
	var builder strings.Builder
	
	// タイムスタンプを追加
	if l.showTime {
		timestamp := time.Now().Format("2006-01-02 15:04:05.000")
		builder.WriteString(fmt.Sprintf("[%s] ", timestamp))
	}
	
	// レベル名とコンポーネント名を追加（カラーあり/なし）
	if l.useColor {
		colorCode, hasColor := levelColors[level]
		if hasColor {
			builder.WriteString(fmt.Sprintf("%s[%s]%s ", colorCode, levelName, colorReset))
			builder.WriteString(fmt.Sprintf("%s[%s]%s ", colorCode, string(component), colorReset))
		} else {
			builder.WriteString(fmt.Sprintf("[%s] ", levelName))
			builder.WriteString(fmt.Sprintf("[%s] ", string(component)))
		}
	} else {
		builder.WriteString(fmt.Sprintf("[%s] ", levelName))
		builder.WriteString(fmt.Sprintf("[%s] ", string(component)))
	}
	
	// メッセージを追加
	builder.WriteString(msg)
	
	formattedMsg := builder.String()
	
	// 標準出力に書き込み
	fmt.Fprintln(l.writer, formattedMsg)
	
	// ファイルにも書き込み（設定されている場合）
	if l.fileWriter != nil {
		fmt.Fprintln(l.fileWriter, formattedMsg)
	}
}

// Error はエラーレベルのログを記録する
func (l *Logger) Error(format string, args ...interface{}) {
	l.log(LevelError, format, args...)
}

// Warn は警告レベルのログを記録する
func (l *Logger) Warn(format string, args ...interface{}) {
	l.log(LevelWarn, format, args...)
}

// Info は情報レベルのログを記録する
func (l *Logger) Info(format string, args ...interface{}) {
	l.log(LevelInfo, format, args...)
}

// Debug はデバッグレベルのログを記録する
func (l *Logger) Debug(format string, args ...interface{}) {
	l.log(LevelDebug, format, args...)
}

// Trace はトレースレベルのログを記録する
func (l *Logger) Trace(format string, args ...interface{}) {
	l.log(LevelTrace, format, args...)
}

// TypeInfo は型情報レベルのログを記録する
func (l *Logger) TypeInfo(format string, args ...interface{}) {
	l.log(LevelTypeInfo, format, args...)
}

// EvalDebug は評価器専用デバッグ情報を記録する
func (l *Logger) EvalDebug(format string, args ...interface{}) {
	l.log(LevelEvalDebug, format, args...)
}

// コンポーネント指定付きログ関数群
// ComponentError はコンポーネント指定付きでエラーレベルのログを記録する
func (l *Logger) ComponentError(component ComponentType, format string, args ...interface{}) {
	l.logWithComponent(component, LevelError, format, args...)
}

// ComponentWarn はコンポーネント指定付きで警告レベルのログを記録する
func (l *Logger) ComponentWarn(component ComponentType, format string, args ...interface{}) {
	l.logWithComponent(component, LevelWarn, format, args...)
}

// ComponentInfo はコンポーネント指定付きで情報レベルのログを記録する
func (l *Logger) ComponentInfo(component ComponentType, format string, args ...interface{}) {
	l.logWithComponent(component, LevelInfo, format, args...)
}

// ComponentDebug はコンポーネント指定付きでデバッグレベルのログを記録する
func (l *Logger) ComponentDebug(component ComponentType, format string, args ...interface{}) {
	l.logWithComponent(component, LevelDebug, format, args...)
}

// ComponentTrace はコンポーネント指定付きでトレースレベルのログを記録する
func (l *Logger) ComponentTrace(component ComponentType, format string, args ...interface{}) {
	l.logWithComponent(component, LevelTrace, format, args...)
}

// グローバル関数
func SetLevel(level LogLevel) {
	GetLogger().SetLevel(level)
}

func SetComponentLevel(component ComponentType, level LogLevel) {
	GetLogger().SetComponentLevel(component, level)
}

func GetComponentLevel(component ComponentType) LogLevel {
	return GetLogger().GetComponentLevel(component)
}

func SetSpecialLevelEnabled(level LogLevel, enabled bool) {
	GetLogger().SetSpecialLevelEnabled(level, enabled)
}

func IsSpecialLevelEnabled(level LogLevel) bool {
	return GetLogger().IsSpecialLevelEnabled(level)
}

func SetOutput(w io.Writer) {
	GetLogger().SetOutput(w)
}

func SetFileOutput(w io.Writer) {
	GetLogger().SetFileOutput(w)
}

func EnableColor() {
	GetLogger().EnableColor()
}

func DisableColor() {
	GetLogger().DisableColor()
}

func EnableTimestamp() {
	GetLogger().EnableTimestamp()
}

func DisableTimestamp() {
	GetLogger().DisableTimestamp()
}

func Enable() {
	GetLogger().Enable()
}

func Disable() {
	GetLogger().Disable()
}

// コンポーネント指定付きグローバルログ関数
func ComponentError(component ComponentType, format string, args ...interface{}) {
	GetLogger().ComponentError(component, format, args...)
}

func ComponentWarn(component ComponentType, format string, args ...interface{}) {
	GetLogger().ComponentWarn(component, format, args...)
}

func ComponentInfo(component ComponentType, format string, args ...interface{}) {
	GetLogger().ComponentInfo(component, format, args...)
}

func ComponentDebug(component ComponentType, format string, args ...interface{}) {
	GetLogger().ComponentDebug(component, format, args...)
}

func ComponentTrace(component ComponentType, format string, args ...interface{}) {
	GetLogger().ComponentTrace(component, format, args...)
}

func Error(format string, args ...interface{}) {
	GetLogger().Error(format, args...)
}

func Warn(format string, args ...interface{}) {
	GetLogger().Warn(format, args...)
}

func Info(format string, args ...interface{}) {
	GetLogger().Info(format, args...)
}

func Debug(format string, args ...interface{}) {
	GetLogger().Debug(format, args...)
}

func Trace(format string, args ...interface{}) {
	GetLogger().Trace(format, args...)
}

func TypeInfo(format string, args ...interface{}) {
	GetLogger().TypeInfo(format, args...)
}

func EvalDebug(format string, args ...interface{}) {
	GetLogger().EvalDebug(format, args...)
}

// ParseLogLevel は文字列からログレベルを解析する
func ParseLogLevel(levelStr string) LogLevel {
	switch levelStr {
	case "OFF":
		return LevelOff
	case "ERROR":
		return LevelError
	case "WARN":
		return LevelWarn
	case "INFO":
		return LevelInfo
	case "DEBUG":
		return LevelDebug
	case "TRACE":
		return LevelTrace
	case "TYPE":
		return LevelTypeInfo
	case "EVAL":
		return LevelEvalDebug
	default:
		return LevelInfo // デフォルトはINFO
	}
}
