#!/bin/bash

# uncodeインタプリタのビルドスクリプト

# ディレクトリ構造の確認
mkdir -p bin
mkdir -p ai

# ビルド
echo "uncodeインタプリタをビルドしています..."
cd src && GO111MODULE=on go build -o ../bin/uncode ./main.go

if [ $? -eq 0 ]; then
    echo "ビルド成功！"
    echo "使用方法: ./bin/uncode [オプション] <ファイル名>"
    echo "詳細なオプションは './bin/uncode --help' で確認できます。"
else
    echo "ビルドに失敗しました。"
    exit 1
fi
