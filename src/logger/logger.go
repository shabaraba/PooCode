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
)

// LevelNames はログレベルと名前のマッピング
var LevelNames = map[LogLevel]string{
	LevelOff:   "OFF",
	LevelError: "ERROR",
	LevelWarn:  "WARN",
	LevelInfo:  "INFO",
	LevelDebug: "DEBUG",
	LevelTrace: "TRACE",
}

var levelColors = map[LogLevel]string{
	LevelError: "\033[31m", // 赤
	LevelWarn:  "\033[33m", // 黄
	LevelInfo:  "\033[32m", // 緑
	LevelDebug: "\033[36m", // シアン
	LevelTrace: "\033[35m", // マゼンタ
}

const (
	colorReset = "\033[0m"
)

// Logger はログを記録するための構造体
type Logger struct {
	level       LogLevel
	writer      io.Writer
	fileWriter  io.Writer
	mu          sync.Mutex
	isEnabled   bool
	useColor    bool
	showTime    bool
}

var (
	defaultLogger *Logger
	once          sync.Once
)

// LoggerOption はロガーの設定用オプション関数型
type LoggerOption func(*Logger)

// WithLevel はログレベルを設定するオプション
func WithLevel(level LogLevel) LoggerOption {
	return func(l *Logger) {
		l.level = level
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
		level:     LevelInfo,
		writer:    os.Stdout,
		fileWriter: nil,
		isEnabled: true,
		useColor:  true,
		showTime:  true,
	}
	
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

// SetLevel はログレベルを設定する
func (l *Logger) SetLevel(level LogLevel) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.level = level
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

	if !l.isEnabled || level > l.level {
		return
	}

	// フォーマット済みのメッセージを生成
	formattedMsg := l.formatLogMessage(level, format, args...)
	
	// 標準出力に書き込み
	fmt.Fprintln(l.writer, formattedMsg)
	
	// ファイルにも書き込み（設定されている場合）
	if l.fileWriter != nil {
		// ファイルには色なしで書き込む
		plainMsg := l.formatLogMessage(level, format, args...)
		fmt.Fprintln(l.fileWriter, plainMsg)
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

// グローバル関数
func SetLevel(level LogLevel) {
	GetLogger().SetLevel(level)
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
	default:
		return LevelInfo // デフォルトはINFO
	}
}
