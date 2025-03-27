#!/bin/bash

# run_pipe_debug.sh - パイプライン処理のデバッグログを有効にして実行するスクリプト

# カレントディレクトリをスクリプトの場所に変更
cd "$(dirname "$0")"

# パイプライン処理のデバッグログを有効にする環境変数
export POO_PIPE_DEBUG=1

# 第1引数がない場合はテストファイルを指定
if [ -z "$1" ]; then
  FILE="examples/test_pipe_operators.poo"
else
  FILE="$1"
fi

echo "実行ファイル: $FILE"
echo "パイプライン処理のデバッグログを有効にして実行します"

# 実行して出力をログファイルにも保存
./poo $FILE | tee ai/output.log
