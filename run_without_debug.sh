#!/bin/bash
# 簡易テスト用スクリプト: デバッグ情報なしでfizzbuzzを実行

cd /Users/t002451/my_work/private/PooCode

# ビルドして実行（デバッグフラグなし）
bash build.sh
./bin/uncode examples/fizzbuzz.poo | tee ai/output.log
