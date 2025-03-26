#!/bin/bash

# uncodeインタプリタのビルドスクリプト

set -x  # デバッグ出力を有効化

# ディレクトリ構造の確認
mkdir -p bin
mkdir -p ai

# 現在のディレクトリを確認
pwd
ls -la

# ビルド
echo "uncodeインタプリタをビルドしています..."
cd src
pwd
ls -la
echo "実行: GO111MODULE=on go build -o ../bin/uncode ./main.go"
GO111MODULE=on go build -v -o ../bin/uncode ./main.go 2>&1

BUILD_RESULT=$?
if [ $BUILD_RESULT -eq 0 ]; then
    echo "ビルド成功！"
    echo "使用方法: ./bin/uncode [オプション] <ファイル名>"
    echo "詳細なオプションは './bin/uncode --help' で確認できます。"
else
    echo "ビルドに失敗しました。エラーコード: $BUILD_RESULT"
    exit 1
fi
