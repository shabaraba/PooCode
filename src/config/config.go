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
	SourceFile         string
	DebugMode          bool
	LogLevel           logger.LogLevel
	ComponentLogLevels map[logger.ComponentType]logger.LogLevel
	SpecialLogLevels   map[logger.LogLevel]bool  // 特殊なログレベルの有効/無効
	LogFile            string
	OutputFile         string
	ColorOutput        bool
	ShowTimestamp      bool
	ShowTypeInfo       bool
	ShowLexerDebug     bool
	ShowParserDebug    bool
	ShowEvalDebug      bool
	ShowBuiltinDebug   bool
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
	// マップの初期化
	GlobalConfig.ComponentLogLevels = make(map[logger.ComponentType]logger.LogLevel)
	GlobalConfig.SpecialLogLevels = make(map[logger.LogLevel]bool)
	
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
	logLevelStr := flag.String("log-level", "", "グローバルログレベル (OFF, ERROR, WARN, INFO, DEBUG, TRACE)")
	
	// コンポーネント別ログレベルの設定
	lexerLogLevelStr := flag.String("lexer-log-level", "", "レキサーのログレベル (OFF, ERROR, WARN, INFO, DEBUG, TRACE)")
	parserLogLevelStr := flag.String("parser-log-level", "", "パーサーのログレベル")
	evalLogLevelStr := flag.String("eval-log-level", "", "評価器のログレベル")
	runtimeLogLevelStr := flag.String("runtime-log-level", "", "ランタイムのログレベル")
	builtinLogLevelStr := flag.String("builtin-log-level", "", "組み込み関数のログレベル")

	flag.Parse()

	// グローバルログレベルの設定
	if *logLevelStr != "" {
		GlobalConfig.LogLevel = logger.ParseLogLevel(*logLevelStr)
	} else if GlobalConfig.DebugMode {
		GlobalConfig.LogLevel = logger.LevelDebug
	} else {
		GlobalConfig.LogLevel = logger.LevelInfo
	}
	
	// コンポーネント別ログレベルの設定
	if *lexerLogLevelStr != "" {
		GlobalConfig.ComponentLogLevels[logger.ComponentLexer] = logger.ParseLogLevel(*lexerLogLevelStr)
	}
	
	if *parserLogLevelStr != "" {
		GlobalConfig.ComponentLogLevels[logger.ComponentParser] = logger.ParseLogLevel(*parserLogLevelStr)
	}
	
	if *evalLogLevelStr != "" {
		GlobalConfig.ComponentLogLevels[logger.ComponentEval] = logger.ParseLogLevel(*evalLogLevelStr)
	}
	
	if *runtimeLogLevelStr != "" {
		GlobalConfig.ComponentLogLevels[logger.ComponentRuntime] = logger.ParseLogLevel(*runtimeLogLevelStr)
	}
	
	if *builtinLogLevelStr != "" {
		GlobalConfig.ComponentLogLevels[logger.ComponentBuiltin] = logger.ParseLogLevel(*builtinLogLevelStr)
	}

	// デバッグフラグを設定した場合は自動的に対応するデバッグを有効にする
	if GlobalConfig.DebugMode {
		GlobalConfig.ShowLexerDebug = true
		GlobalConfig.ShowParserDebug = true 
		GlobalConfig.ShowEvalDebug = true
		GlobalConfig.ShowBuiltinDebug = true
		GlobalConfig.ShowConditionDebug = true
		
		// 特殊デバッグログレベルも有効化
		GlobalConfig.SpecialLogLevels[logger.LevelTypeInfo] = GlobalConfig.ShowTypeInfo
		GlobalConfig.SpecialLogLevels[logger.LevelEvalDebug] = GlobalConfig.ShowEvalDebug
	} else {
		// デバッグモードでない場合の設定
		GlobalConfig.SpecialLogLevels[logger.LevelTypeInfo] = GlobalConfig.ShowTypeInfo
		GlobalConfig.SpecialLogLevels[logger.LevelEvalDebug] = GlobalConfig.ShowEvalDebug
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
	// グローバルログレベルの設定を適用
	logger.SetLevel(GlobalConfig.LogLevel)
	
	// コンポーネント別ログレベルの設定を適用
	for component, level := range GlobalConfig.ComponentLogLevels {
		logger.SetComponentLevel(component, level)
	}
	
	// コンポーネント別デバッグフラグからログレベルを設定
	if GlobalConfig.ShowLexerDebug && GlobalConfig.ComponentLogLevels[logger.ComponentLexer] == 0 {
		logger.SetComponentLevel(logger.ComponentLexer, logger.LevelDebug)
	}
	
	if GlobalConfig.ShowParserDebug && GlobalConfig.ComponentLogLevels[logger.ComponentParser] == 0 {
		logger.SetComponentLevel(logger.ComponentParser, logger.LevelDebug)
	}
	
	if GlobalConfig.ShowEvalDebug && GlobalConfig.ComponentLogLevels[logger.ComponentEval] == 0 {
		logger.SetComponentLevel(logger.ComponentEval, logger.LevelDebug)
	}
	
	if GlobalConfig.ShowBuiltinDebug && GlobalConfig.ComponentLogLevels[logger.ComponentBuiltin] == 0 {
		logger.SetComponentLevel(logger.ComponentBuiltin, logger.LevelDebug)
	}
	
	// 特殊ログレベルの設定を適用
	for level, enabled := range GlobalConfig.SpecialLogLevels {
		logger.SetSpecialLevelEnabled(level, enabled)
	}
	
	// 特殊デバッグに関連するコンポーネントのログレベルを設定
	if GlobalConfig.ShowEvalDebug {
		// 評価器デバッグログを有効にする場合、評価デバッグレベルも有効化
		logger.SetSpecialLevelEnabled(logger.LevelEvalDebug, true)
	}
	
	if GlobalConfig.ShowTypeInfo {
		// 型情報表示を有効にする場合、型情報デバッグレベルを有効化
		logger.SetSpecialLevelEnabled(logger.LevelTypeInfo, true)
	}
	
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
	
	// config.goのGlobalConfigを初期化して全てのフラグ定義を呼び出す
	GlobalConfig.ComponentLogLevels = make(map[logger.ComponentType]logger.LogLevel)
	GlobalConfig.SpecialLogLevels = make(map[logger.LogLevel]bool)
	
	// コマンドラインフラグの定義（しかしParseはしない）
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
	
	flag.String("log-level", "", "グローバルログレベル (OFF, ERROR, WARN, INFO, DEBUG, TRACE)")
	flag.String("lexer-log-level", "", "レキサーのログレベル (OFF, ERROR, WARN, INFO, DEBUG, TRACE)")
	flag.String("parser-log-level", "", "パーサーのログレベル")
	flag.String("eval-log-level", "", "評価器のログレベル")
	flag.String("runtime-log-level", "", "ランタイムのログレベル")
	flag.String("builtin-log-level", "", "組み込み関数のログレベル")
	
	flag.PrintDefaults()
	fmt.Println("\nサポートされている拡張子: .poo, .💩")
}
