#!/bin/bash
# デバッグ情報ありでファイルを実行するスクリプト

cd /Users/t002451/my_work/private/PooCode

# コマンドライン引数を取得
FILE=${1:-examples/fizzbuzz.poo}

# コンポーネント別のログレベル設定
LEXER_LOG_LEVEL=${2:-DEBUG}
PARSER_LOG_LEVEL=${3:-DEBUG}
EVAL_LOG_LEVEL=${4:-DEBUG}
BUILTIN_LOG_LEVEL=${5:-DEBUG}
RUNTIME_LOG_LEVEL=${6:-DEBUG}

# ビルドして実行（--debugフラグ付き）
bash build.sh

echo "テストファイル $FILE を実行しています（デバッグモード）..."
./bin/uncode --debug --log-level=DEBUG --log=ai/output.log \
  --lexer-log-level=$LEXER_LOG_LEVEL \
  --parser-log-level=$PARSER_LOG_LEVEL \
  --eval-log-level=$EVAL_LOG_LEVEL \
  --builtin-log-level=$BUILTIN_LOG_LEVEL \
  --runtime-log-level=$RUNTIME_LOG_LEVEL \
  --color --timestamp --show-parser $FILE | tee ai/output_debug.log
