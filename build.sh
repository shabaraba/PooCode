#!/bin/bash

# uncodeインタプリタのビルドスクリプト

# ディレクトリ構造の確認
mkdir -p bin

# ビルド
echo "uncodeインタプリタをビルドしています..."
cd src && GO111MODULE=on go build -o ../bin/uncode cmd/uncode/main.go

if [ $? -eq 0 ]; then
    echo "ビルド成功！"
    echo "使用方法: ./bin/uncode <ファイル名>"
else
    echo "ビルドに失敗しました。"
    exit 1
fi
