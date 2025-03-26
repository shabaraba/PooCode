#!/bin/bash
# PooCode 実行用ユーティリティスクリプト

cd /Users/t002451/my_work/private/PooCode

# 引数の確認
if [ "$#" -lt 1 ]; then
    echo "使用方法: $0 <ファイル> [オプション]"
    echo "オプション:"
    echo "  --debug       デバッグモードで実行"
    echo "  --log-level   ログレベル (OFF, ERROR, WARN, INFO, DEBUG, TRACE, TYPE)"
    echo "  --no-color    カラー出力を無効化"
    echo "  --no-time     タイムスタンプを非表示"
    echo "  --show-lexer  レキサーデバッグ情報を表示"
    echo "  --show-parser パーサーデバッグ情報を表示"
    echo "  --show-eval   評価時デバッグ情報を表示"
    echo "  --show-types  型情報を表示"
    exit 1
fi

# ファイル名を取得
FILE="$1"
shift

# 引数をパース
DEBUG=""
LOG_LEVEL="INFO"
COLOR="--color"
TIMESTAMP="--timestamp"
SHOW_LEXER=""
SHOW_PARSER=""
SHOW_EVAL=""
SHOW_TYPES=""
LOG_FILE="ai/output.log"
OUTPUT_FILE="ai/output.log"

while [ "$#" -gt 0 ]; do
    case "$1" in
        --debug)
            DEBUG="--debug"
            LOG_LEVEL="DEBUG"
            ;;
        --log-level=*)
            LOG_LEVEL="${1#*=}"
            ;;
        --no-color)
            COLOR=""
            ;;
        --no-time)
            TIMESTAMP=""
            ;;
        --show-lexer)
            SHOW_LEXER="--show-lexer"
            ;;
        --show-parser)
            SHOW_PARSER="--show-parser"
            ;;
        --show-eval)
            SHOW_EVAL="--show-eval"
            ;;
        --show-types)
            SHOW_TYPES="--show-types"
            ;;
        *)
            echo "不明なオプション: $1"
            exit 1
            ;;
    esac
    shift
done

# ビルド
bash build.sh

# 実行コマンドの組み立て
CMD="./bin/uncode $DEBUG $COLOR $TIMESTAMP --log-level=$LOG_LEVEL --log=$LOG_FILE $SHOW_LEXER $SHOW_PARSER $SHOW_EVAL $SHOW_TYPES $FILE"

# 実行
echo "実行コマンド: $CMD"
eval "$CMD | tee $OUTPUT_FILE"

# 終了メッセージ
echo "実行完了。出力は $OUTPUT_FILE に保存されています。"
