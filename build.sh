#!/bin/bash

# uncodeインタプリタのビルドスクリプト

# デフォルト設定
DEBUG=false
LOG_LEVEL="INFO"

# コマンドラインオプションの解析
while getopts "dl:" opt; do
  case $opt in
    d) DEBUG=true ;;
    l) LOG_LEVEL="$OPTARG" ;;
    \?) echo "使用方法: $0 [-d] [-l LOG_LEVEL]" >&2
        echo "  -d: デバッグ情報を表示"
        echo "  -l: ログレベル (OFF, ERROR, WARN, INFO, DEBUG, TRACE)"
        exit 1 ;;
  esac
done

# ディレクトリ構造の確認
mkdir -p bin
mkdir -p ai

# ビルド
echo "uncodeインタプリタをビルドしています..." | tee ai/output.log

if [ "$DEBUG" = true ]; then
  echo "デバッグモードでビルドしています..." | tee -a ai/output.log
  cd src && GO111MODULE=on go build -o ../bin/uncode -ldflags "-X main.debugMode=true" ./main.go 2>&1 | tee -a ../ai/output.log
else
  cd src && GO111MODULE=on go build -o ../bin/uncode ./main.go 2>&1 | tee -a ../ai/output.log
fi

if [ ${PIPESTATUS[0]} -eq 0 ]; then
    echo "ビルド成功！" | tee -a ../ai/output.log
    echo "使用方法: ./bin/uncode [オプション] <ファイル名>" | tee -a ../ai/output.log
    echo "詳細なオプションは './bin/uncode --help' で確認できます。" | tee -a ../ai/output.log
    
    # バージョン情報
    BUILD_DATE=$(date "+%Y-%m-%d %H:%M:%S")
    echo "ビルド日時: $BUILD_DATE" | tee -a ../ai/output.log
    
    # コミットメッセージの例（セマンティックバージョニング）
    echo "" | tee -a ../ai/output.log
    echo "コミット例:" | tee -a ../ai/output.log
    echo "git commit -m 'feat: add debug mode to build script'" | tee -a ../ai/output.log
    echo "git commit -m 'fix: resolve missing null object reference'" | tee -a ../ai/output.log
else
    echo "ビルドに失敗しました。" | tee -a ../ai/output.log
    echo "詳細はai/output.logを確認してください。" | tee -a ../ai/output.log
    exit 1
fi
