package runtime

import (
	"fmt"
	"os"

	"github.com/uncode/ast"
	"github.com/uncode/config"
	"github.com/uncode/evaluator"
	"github.com/uncode/lexer"
	"github.com/uncode/logger"
	"github.com/uncode/object"
	"github.com/uncode/parser"
	"github.com/uncode/token"
)

// SourceCodeResult は処理結果を表す構造体
type SourceCodeResult struct {
	Tokens   []token.Token
	Program  *ast.Program
	Result   object.Object
	ExitCode int
}

// SetupBuiltins は組み込み関数を環境に設定する
func SetupBuiltins(env *object.Environment) {
	// プリント関数を追加
	env.Set("print", &object.Builtin{
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}
			return evaluator.NULL
		},
	})
	
	// その他の組み込み関数があればここに追加
}

// ExecuteSourceFile はソースファイルを読み込んで実行する
func ExecuteSourceFile(filePath string) (*SourceCodeResult, error) {
	result := &SourceCodeResult{
		ExitCode: 0,
	}

	// ファイル読み込み
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("ファイルを読み込めませんでした: %w", err)
	}

	// ファイル内容をデバッグ出力
	if config.GlobalConfig.ShowLexerDebug {
		logger.Debug("ファイル内容:\n%s\n", string(content))
	}

	// レキサーでトークン化
	l := lexer.NewLexer(string(content))
	tokens, err := l.Tokenize()
	if err != nil {
		logger.Error("レキサーエラー: %s\n", err)
		result.ExitCode = 1
		return result, err
	}
	result.Tokens = tokens

	// トークン列をデバッグ出力
	if config.GlobalConfig.ShowLexerDebug {
		logger.Debug("トークン列:")
		for i, tok := range tokens {
			logger.Debug("%d: %s\n", i, tok.String())
		}
	}

	// パーサーで構文解析
	p := parser.NewParser(tokens)
	program, err := p.ParseProgram()
	if err != nil {
		logger.Error("パーサーエラー: %s\n", err)
		result.ExitCode = 1
		return result, err
	}
	result.Program = program

	// 構文木をデバッグ出力
	if config.GlobalConfig.ShowParserDebug {
		logger.Debug("構文木:")
		logger.Debug(program.String())
	}

	// インタプリタで実行
	env := object.NewEnvironment()
	SetupBuiltins(env)

	// 型情報のデバッグ出力を設定
	if config.GlobalConfig.ShowTypeInfo {
		logger.SetLevel(logger.LevelTypeInfo)
	}

	// 評価フェーズのデバッグ出力
	if config.GlobalConfig.ShowEvalDebug {
		logger.Debug("評価フェーズ開始...")
	}

	evalResult := evaluator.Eval(program, env)
	result.Result = evalResult
	
	if evalResult != nil && evalResult.Type() == object.ERROR_OBJ {
		logger.Error("実行時エラー: %s\n", evalResult.Inspect())
		result.ExitCode = 1
		return result, fmt.Errorf("実行時エラー: %s", evalResult.Inspect())
	}

	// 実行結果を表示
	if evalResult != nil && config.GlobalConfig.ShowEvalDebug {
		logger.Info("実行結果: %s\n", evalResult.Inspect())
	}

	return result, nil
}
