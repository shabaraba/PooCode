package main

import (
	"fmt"
	"os"

	"github.com/uncode/config"
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
	fmt.Print("logger")
	err = config.SetupLogger()
	if err != nil {
		fmt.Fprintf(os.Stderr, "ロガーの初期化エラー: %s\n", err)
		os.Exit(1)
	}

	// バージョン情報のログ
	logger.Info("PooCode インタプリタ バージョン 0.1.0")
	logger.Debug("デバッグモード: %v", config.GlobalConfig.DebugMode)
	logger.Debug("ログレベル: %s", logger.LevelNames[config.GlobalConfig.LogLevel])

	// ソースファイルの実行
	result, err := runtime.ExecuteSourceFile(config.GlobalConfig.SourceFile)
	if err != nil {
		// エラーはruntime内でログ出力されるので、ここでは終了コードだけ設定
		os.Exit(result.ExitCode)
	}

	// 正常終了
	os.Exit(0)
}
