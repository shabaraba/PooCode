#!/bin/bash
# パイプ演算子のテストを実行するスクリプト

# ディレクトリの設定
PROJECT_DIR="/Users/t002451/my_work/private/PooCode"
EXAMPLE_DIR="$PROJECT_DIR/examples"
OUTPUT_FILE="$PROJECT_DIR/ai/output.log"

# バイナリのパス
BINARY="$PROJECT_DIR/bin/uncode"

# 色の設定
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
NC='\033[0m' # No Color

echo -e "${YELLOW}新しいパイプ演算子のテストを開始...${NC}"

# コンパイル
echo -e "${YELLOW}コンパイル中...${NC}"
(cd $PROJECT_DIR && go build -o bin/uncode ./src)

if [ $? -ne 0 ]; then
    echo -e "\n${RED}コンパイルに失敗しました。エラーを確認してください。${NC}"
    exit 1
fi

echo -e "${GREEN}コンパイル成功！${NC}\n"

# 新しいパイプ演算子のテスト
echo -e "${YELLOW}新しいパイプ演算子のテスト実行中...${NC}"
$BINARY $EXAMPLE_DIR/test_new_pipes.poo | tee $OUTPUT_FILE

echo -e "\n${GREEN}テスト完了！${NC}"
echo -e "出力ログは $OUTPUT_FILE に保存されました。"
