#!/bin/bash
# PooCode デバッグ用ユーティリティスクリプト
# 特定のデバッグモードを簡単に実行するためのスクリプト

cd /Users/t002451/my_work/private/PooCode

# 引数の確認
if [ "$#" -lt 1 ]; then
    echo "使用方法: $0 <モード> [ファイル]"
    echo "モード:"
    echo "  all      全てのデバッグ情報を表示"
    echo "  lexer    レキサーのデバッグ情報のみ表示"
    echo "  parser   パーサーのデバッグ情報のみ表示"
    echo "  eval     評価のデバッグ情報のみ表示"
    echo "  type     型情報のみ表示"
    echo "  none     デバッグ情報なしで実行（通常実行）"
    echo "  help     ヘルプメッセージを表示"
    echo ""
    echo "ファイル: 実行するPooCodeファイル（省略時はfizzbuzz.poo）"
    exit 1
fi

# モードとファイル名の取得
MODE="$1"
FILE="${2:-examples/fizzbuzz.poo}"

# ログファイル名
LOG_FILE="ai/output.log"

# ビルド
bash build.sh

# モードによって異なるコマンド実行
case "$MODE" in
    all)
        echo "全てのデバッグ情報を表示してファイル $FILE を実行します..."
        ./bin/uncode --debug --color --timestamp --log-level=DEBUG --log=$LOG_FILE --show-lexer --show-parser --show-eval --show-types $FILE | tee ai/output_debug_all.log
        ;;
    lexer)
        echo "レキサーのデバッグ情報を表示してファイル $FILE を実行します..."
        ./bin/uncode --color --timestamp --log-level=DEBUG --log=$LOG_FILE --show-lexer $FILE | tee ai/output_debug_lexer.log
        ;;
    parser)
        echo "パーサーのデバッグ情報を表示してファイル $FILE を実行します..."
        ./bin/uncode --color --timestamp --log-level=DEBUG --log=$LOG_FILE --show-parser $FILE | tee ai/output_debug_parser.log
        ;;
    eval)
        echo "評価のデバッグ情報を表示してファイル $FILE を実行します..."
        ./bin/uncode --color --timestamp --log-level=DEBUG --log=$LOG_FILE --show-eval $FILE | tee ai/output_debug_eval.log
        ;;
    type)
        echo "型情報を表示してファイル $FILE を実行します..."
        ./bin/uncode --color --timestamp --log-level=INFO --log=$LOG_FILE --show-types $FILE | tee ai/output_debug_type.log
        ;;
    none)
        echo "デバッグ情報なしでファイル $FILE を実行します..."
        ./bin/uncode --color --timestamp --log-level=INFO --log=$LOG_FILE $FILE | tee ai/output_clean.log
        ;;
    help)
        echo "PooCode デバッグユーティリティ:"
        echo "  このスクリプトは様々なデバッグモードでPooCodeを実行するのに役立ちます。"
        echo ""
        echo "使用方法: $0 <モード> [ファイル]"
        echo "モード:"
        echo "  all      全てのデバッグ情報を表示"
        echo "  lexer    レキサーのデバッグ情報のみ表示"
        echo "  parser   パーサーのデバッグ情報のみ表示"
        echo "  eval     評価のデバッグ情報のみ表示"
        echo "  type     型情報のみ表示"
        echo "  none     デバッグ情報なしで実行（通常実行）"
        echo "  help     このヘルプメッセージを表示"
        echo ""
        echo "例:"
        echo "  $0 all examples/fizzbuzz.poo    - 全てのデバッグ情報を表示してfizzbuzzを実行"
        echo "  $0 parser examples/hello.poo    - パーサーのデバッグ情報のみ表示してhelloを実行"
        ;;
    *)
        echo "不明なモード: $MODE"
        echo "使用できるモード: all, lexer, parser, eval, type, none, help"
        exit 1
        ;;
esac

if [ "$MODE" != "help" ]; then
    echo "実行完了。出力は ai/ ディレクトリに保存されています。"
fi
