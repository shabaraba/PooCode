package main

import (
	"fmt"
	"os"

	"github.com/uncode/config"
	"github.com/uncode/evaluator"
	"github.com/uncode/logger"
	"github.com/uncode/runtime"
)

func main() {
	// コマンドラインフラグのパース
	err := config.ParseFlags()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		config.PrintUsage()
		os.Exit(1)
	}

	// ロガーの設定
	err = config.SetupLogger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ロガーの初期化エラー: %s\n", err)
		os.Exit(1)
	}

	// バージョン情報のログ
	logger.Info("PooCode インタプリタ バージョン 0.1.0")
	logger.Debug("デバッグモード: %v", config.GlobalConfig.DebugMode)
	logger.Debug("ログレベル: %s", logger.LevelNames[config.GlobalConfig.LogLevel])
	logger.Debug("ソースファイル: %s", config.GlobalConfig.SourceFile)
	
	// 組み込み関数のログレベルを設定
	if config.GlobalConfig.ShowBuiltinDebug {
		evaluator.SetBuiltinLogLevel(logger.LevelDebug)
	} else {
		evaluator.SetBuiltinLogLevel(logger.LevelInfo)
	}

	// ソースファイルの実行
	result, err := runtime.ExecuteSourceFile(config.GlobalConfig.SourceFile)
	if err != nil {
		// エラーはruntime内でログ出力されるので、ここでは終了コードだけ設定
		os.Exit(result.ExitCode)
	}

	// 正常終了
	os.Exit(0)
}
