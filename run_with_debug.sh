#!/bin/bash
# デバッグ情報ありでファイルを実行するスクリプト

cd /Users/t002451/my_work/private/PooCode

# コマンドライン引数を取得
FILE=${1:-examples/fizzbuzz.poo}

# ビルドして実行（--debugフラグ付き）
bash build.sh

echo "テストファイル $FILE を実行しています（デバッグモード）..."
./bin/uncode --debug --log-level=DEBUG --log=ai/output.log --color --timestamp $FILE | tee ai/output_debug.log
