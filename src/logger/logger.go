package logger

import (
	"fmt"
	"io"
	"os"
	"sync"
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

var levelNames = map[LogLevel]string{
	LevelOff:   "OFF",
	LevelError: "ERROR",
	LevelWarn:  "WARN",
	LevelInfo:  "INFO",
	LevelDebug: "DEBUG",
	LevelTrace: "TRACE",
}

// Logger はログを記録するための構造体
type Logger struct {
	level       LogLevel
	writer      io.Writer
	teeWriter   io.Writer  // 複数の出力先に同時に書き込むためのライター
	mu          sync.Mutex
	isEnabled   bool
}

var (
	defaultLogger *Logger
	once          sync.Once
)

// GetLogger はシングルトンロガーインスタンスを返す
func GetLogger() *Logger {
	once.Do(func() {
		defaultLogger = &Logger{
			level:     LevelInfo,
			writer:    os.Stdout,
			teeWriter: nil,
			isEnabled: true,
		}
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
	l.updateWriter()
}

// SetTeeOutput はログを同時に別のライターにも出力するよう設定する
func (l *Logger) SetTeeOutput(w io.Writer) {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.teeWriter = w
	l.updateWriter()
}

// updateWriter は現在の writer と teeWriter の状態に基づいて、
// MultiWriterを作成するか、単一のwriterを使用するか決定する
func (l *Logger) updateWriter() {
	// 現在の状態によってwriter設定を更新
	if l.teeWriter != nil {
		// io.MultiWriterを作成して両方に出力
		l.writer = io.MultiWriter(l.writer, l.teeWriter)
		// teeWriterはnilに戻してMultiWriterが二重に作成されないようにする
		l.teeWriter = nil
	}
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

// log はメッセージを指定されたレベルでログに記録する
func (l *Logger) log(level LogLevel, format string, args ...interface{}) {
	l.mu.Lock()
	defer l.mu.Unlock()

	if !l.isEnabled || level > l.level {
		return
	}

	levelName := levelNames[level]
	msg := fmt.Sprintf(format, args...)
	fmt.Fprintf(l.writer, "[%s] %s\n", levelName, msg)
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
