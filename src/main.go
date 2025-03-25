package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/uncode/evaluator"
	"github.com/uncode/lexer"
	"github.com/uncode/object"
	"github.com/uncode/parser"
)

var debugMode bool

func main() {
	// コマンドラインフラグのパース
	flag.BoolVar(&debugMode, "debug", false, "デバッグモードを有効にする")
	flag.Parse()

	args := flag.Args()
	if len(args) != 1 {
		fmt.Println("使用方法: uncode [オプション] <ファイル名>")
		fmt.Println("オプション:")
		flag.PrintDefaults()
		os.Exit(1)
	}

	filename := args[0]
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

	// デバッグモードを設定
	evaluator.SetDebugMode(debugMode)
	
	// デバッグモードの場合、ファイル内容を表示
	if debugMode {
		fmt.Printf("ファイル内容:\n%s\n", string(content))
	}

	// レキサーでトークン化
	l := lexer.NewLexer(string(content))
	tokens, err := l.Tokenize()
	if err != nil {
		fmt.Printf("レキサーエラー: %s\n", err)
		os.Exit(1)
	}

	// デバッグモードの場合、トークン列を表示
	if debugMode {
		fmt.Println("トークン列:")
		for i, tok := range tokens {
			fmt.Printf("%d: %s\n", i, tok.String())
		}
	}

	// パーサーで構文解析
	p := parser.NewParser(tokens)
	program, err := p.ParseProgram()
	if err != nil {
		fmt.Printf("パーサーエラー: %s\n", err)
		os.Exit(1)
	}

	// デバッグモードの場合、構文木を表示
	if debugMode {
		fmt.Println("構文木:")
		fmt.Println(program.String())
	}

	// インタプリタで実行
	env := object.NewEnvironment()
	// プリント関数を追加
	env.Set("print", &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}
			return evaluator.NULL
		},
	})
	
	result := evaluator.Eval(program, env)
	if result != nil && result.Type() == object.ERROR_OBJ {
		fmt.Printf("実行時エラー: %s\n", result.Inspect())
		os.Exit(1)
	}

	// デバッグモードの場合、実行結果を表示
	if debugMode && result != nil {
		fmt.Printf("実行結果: %s\n", result.Inspect())
	}
}
