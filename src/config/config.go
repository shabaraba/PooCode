package config

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/uncode/logger"
)

// Config はアプリケーション全体の設定を保持する構造体
type Config struct {
	SourceFile      string
	DebugMode       bool
	LogLevel        logger.LogLevel
	LogFile         string
	OutputFile      string
	ColorOutput     bool
	ShowTimestamp   bool
	ShowTypeInfo    bool
	ShowLexerDebug  bool
	ShowParserDebug bool
	ShowEvalDebug   bool
	ShowBuiltinDebug bool
	ShowConditionDebug bool // 条件式の評価デバッグ表示
}

// GlobalConfig はアプリケーション全体で使用される設定
var GlobalConfig Config

// カスタムエラー型
type InvalidArgsError struct {
	Message string
}

func (e *InvalidArgsError) Error() string {
	return fmt.Sprintf("引数エラー: %s", e.Message)
}

type UnsupportedExtensionError struct {
	Extension string
}

func (e *UnsupportedExtensionError) Error() string {
	return fmt.Sprintf("エラー: サポートされていないファイル拡張子です: %s", e.Extension)
}

// ParseFlags はコマンドライン引数をパースし、設定を行う
func ParseFlags() error {
	// コマンドラインフラグのパース
	flag.BoolVar(&GlobalConfig.DebugMode, "debug", false, "デバッグモードを有効にする")
	flag.StringVar(&GlobalConfig.LogFile, "log", "", "ログファイルのパス (指定がなければ標準出力のみ)")
	flag.StringVar(&GlobalConfig.OutputFile, "output", "", "出力ファイルのパス (tee で出力を記録)")
	flag.BoolVar(&GlobalConfig.ColorOutput, "color", true, "カラー出力を有効にする")
	flag.BoolVar(&GlobalConfig.ShowTimestamp, "timestamp", true, "タイムスタンプを表示する")
	flag.BoolVar(&GlobalConfig.ShowTypeInfo, "show-types", false, "型情報を表示する")
	flag.BoolVar(&GlobalConfig.ShowLexerDebug, "show-lexer", false, "レキサーのデバッグ情報を表示する")
	flag.BoolVar(&GlobalConfig.ShowParserDebug, "show-parser", false, "パーサーのデバッグ情報を表示する")
	flag.BoolVar(&GlobalConfig.ShowEvalDebug, "show-eval", false, "評価時のデバッグ情報を表示する")
	flag.BoolVar(&GlobalConfig.ShowBuiltinDebug, "show-builtin", false, "組み込み関数のデバッグ情報を表示する")
	flag.BoolVar(&GlobalConfig.ShowConditionDebug, "show-condition", false, "条件式評価のデバッグ情報を表示する")

	// ログレベルをフラグで指定できるようにする
	logLevelStr := flag.String("log-level", "", "ログレベル (OFF, ERROR, WARN, INFO, DEBUG, TRACE)")

	flag.Parse()

	// ログレベルの設定
	if *logLevelStr != "" {
		GlobalConfig.LogLevel = logger.ParseLogLevel(*logLevelStr)
	} else if GlobalConfig.DebugMode {
		GlobalConfig.LogLevel = logger.LevelDebug
	} else {
		GlobalConfig.LogLevel = logger.LevelInfo
	}

	// デバッグフラグを設定した場合は自動的に対応するデバッグを有効にする
	if GlobalConfig.DebugMode {
		GlobalConfig.ShowLexerDebug = true
		GlobalConfig.ShowParserDebug = true 
		GlobalConfig.ShowEvalDebug = true
		GlobalConfig.ShowBuiltinDebug = true
		GlobalConfig.ShowConditionDebug = true
	}

	// ソースファイルのパス取得
	args := flag.Args()
	if len(args) != 1 {
		return &InvalidArgsError{
			Message: "ソースファイルが指定されていません",
		}
	}

	GlobalConfig.SourceFile = args[0]

	// ファイル拡張子のチェック
	ext := filepath.Ext(GlobalConfig.SourceFile)
	if ext != ".poo" && ext != ".💩" {
		return &UnsupportedExtensionError{
			Extension: ext,
		}
	}

	return nil
}

// SetupLogger はロガーの設定を行う
func SetupLogger() error {
	// ロガーの設定を適用
	logger.SetLevel(GlobalConfig.LogLevel)
	
	// カラー出力の設定
	if GlobalConfig.ColorOutput {
		logger.EnableColor()
	} else {
		logger.DisableColor()
	}
	
	// タイムスタンプの設定
	if GlobalConfig.ShowTimestamp {
		logger.EnableTimestamp()
	} else {
		logger.DisableTimestamp()
	}
	
	// ログファイルの設定
	if GlobalConfig.LogFile != "" {
		f, err := os.OpenFile(GlobalConfig.LogFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			return fmt.Errorf("ログファイルを開けませんでした: %w", err)
		}
		logger.SetFileOutput(f)
	}
	
	return nil
}

// PrintUsage はコマンドの使用方法を表示する
func PrintUsage() {
	fmt.Println("使用方法: uncode [オプション] <ファイル名>")
	fmt.Println("オプション:")
	flag.PrintDefaults()
	fmt.Println("\nサポートされている拡張子: .poo, .💩")
}
