package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/uncode/evaluator"
	"github.com/uncode/lexer"
	"github.com/uncode/object"
	"github.com/uncode/parser"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Println("使用方法: uncode <ファイル名>")
		os.Exit(1)
	}

	filename := os.Args[1]
	ext := filepath.Ext(filename)
	if ext != ".poo" && ext != ".💩" {
		fmt.Printf("エラー: サポートされていないファイル拡張子です: %s\n", ext)
		fmt.Println("サポートされている拡張子: .poo, .💩")
		os.Exit(1)
	}

	content, err := os.ReadFile(filename)
	if err != nil {
		fmt.Printf("エラー: ファイルを読み込めませんでした: %s\n", err)
		os.Exit(1)
	}

	// レキサーでトークン化
	l := lexer.NewLexer(string(content))
	tokens, err := l.Tokenize()
	if err != nil {
		fmt.Printf("レキサーエラー: %s\n", err)
		os.Exit(1)
	}

	// パーサーで構文解析
	p := parser.NewParser(tokens)
	program, err := p.ParseProgram()
	if err != nil {
		fmt.Printf("パーサーエラー: %s\n", err)
		os.Exit(1)
	}

	// インタプリタで実行
	env := object.NewEnvironment()
	result := evaluator.Eval(program, env)
	if result != nil && result.Type() == object.ERROR_OBJ {
		fmt.Printf("実行時エラー: %s\n", result.Inspect())
		os.Exit(1)
	}
}
