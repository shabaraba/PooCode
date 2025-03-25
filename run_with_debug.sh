#!/bin/bash
# デバッグ情報ありでfizzbuzzを実行するスクリプト

cd /Users/t002451/my_work/private/PooCode

# ビルドして実行（--debugフラグ付き）
bash build.sh
./bin/uncode --debug examples/fizzbuzz.poo | tee ai/output.log
