#!/bin/bash
# 簡易テスト用スクリプト: デバッグ情報なしでfizzbuzzを実行

cd /Users/t002451/my_work/private/PooCode

# コマンドライン引数を取得
FILE=${1:-examples/fizzbuzz.poo}

# ビルドして実行（デバッグフラグなし）
bash build.sh

echo "テストファイル $FILE を実行しています（通常モード）..."
./bin/uncode --color --timestamp --log-level=INFO --log=ai/output.log $FILE | tee ai/output_clean.log
