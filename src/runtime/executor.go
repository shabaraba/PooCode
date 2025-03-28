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
		Name: "print",
		Fn: func(args ...object.Object) object.Object {
			for _, arg := range args {
				fmt.Println(arg.Inspect())
			}
			return evaluator.NullObj
		},
	})
	
	// 評価器から組み込み関数をすべてインポート
	// evaluator.Builtinsに登録されている関数をすべて環境に追加
	for name, builtin := range evaluator.Builtins {
		logger.Debug("組み込み関数を登録: %s", name)
		env.Set(name, builtin)
	}
}

// convertToObjectIdentifiers は ast.Identifier スライスを object.Identifier スライスに変換する
func convertToObjectIdentifiers(params []*ast.Identifier) []*object.Identifier {
	if params == nil {
		return nil
	}
	
	result := make([]*object.Identifier, len(params))
	for i, param := range params {
		result[i] = &object.Identifier{Value: param.Value}
	}
	return result
}

// findAndRegisterFunctionsInExpression は式の中から関数定義を再帰的に探索する
func findAndRegisterFunctionsInExpression(expr ast.Expression, env *object.Environment, count *int) {
	if expr == nil {
		return
	}

	switch e := expr.(type) {
	case *ast.FunctionLiteral:
		// 関数リテラルを見つけた場合
		if e.Name != nil {
			// 関数名があれば登録
			function := &object.Function{
				Parameters: convertToObjectIdentifiers(e.Parameters),
				ASTBody:    e.Body,
				Env:        env,
				InputType:  e.InputType,
				ReturnType: e.ReturnType,
				Condition:  e.Condition,
			}

			funcName := e.Name.Value
			logger.Debug("埋め込み関数定義: 関数 '%s' の定義を見つけました", funcName)

			// 条件付き関数の場合の特別な処理
			if e.Condition != nil {
				existingFuncs := env.GetAllFunctionsByName(funcName)
				uniqueName := fmt.Sprintf("%s#%d", funcName, len(existingFuncs))
				logger.Debug("埋め込み条件付き関数 '%s' を '%s' として事前登録します", funcName, uniqueName)
				env.Set(uniqueName, function)
			}

			// 通常の名前でも登録
			env.Set(funcName, function)
			*count++
			logger.Debug("埋め込み関数 '%s' を事前登録しました", funcName)
		}

		// 関数本体内のステートメントも探索
		if e.Body != nil {
			for _, stmt := range e.Body.Statements {
				if exprStmt, ok := stmt.(*ast.ExpressionStatement); ok {
					findAndRegisterFunctionsInExpression(exprStmt.Expression, env, count)
				}

				if assignStmt, ok := stmt.(*ast.AssignStatement); ok {
					findAndRegisterFunctionsInExpression(assignStmt.Value, env, count)
				}
			}
		}

	case *ast.InfixExpression:
		// 中置式の場合、左右の式を探索
		findAndRegisterFunctionsInExpression(e.Left, env, count)
		findAndRegisterFunctionsInExpression(e.Right, env, count)

	case *ast.PrefixExpression:
		// 前置式の場合、右の式を探索
		findAndRegisterFunctionsInExpression(e.Right, env, count)

	case *ast.CallExpression:
		// 関数呼び出しの場合、関数と引数を探索
		findAndRegisterFunctionsInExpression(e.Function, env, count)
		for _, arg := range e.Arguments {
			findAndRegisterFunctionsInExpression(arg, env, count)
		}

	case *ast.IndexExpression:
		// 添字式の場合、配列と添字を探索
		findAndRegisterFunctionsInExpression(e.Left, env, count)
		findAndRegisterFunctionsInExpression(e.Index, env, count)
	}
}

// preRegisterFunctions は ASTをトラバースして関数定義を事前に環境に登録する
func preRegisterFunctions(program *ast.Program, env *object.Environment) {
	if program == nil || len(program.Statements) == 0 {
		return
	}

	logger.Debug("関数の事前登録を開始します...")
	registeredCount := 0

	// プログラム内の全てのステートメントを走査
	for _, stmt := range program.Statements {
		// ExpressionStatement 内の FunctionLiteral を検出
		// FUNCTION (def) で始まる式文を見つける
		if exprStmt, ok := stmt.(*ast.ExpressionStatement); ok {
			// 関数リテラルの場合
			if funcLit, ok := exprStmt.Expression.(*ast.FunctionLiteral); ok {
				if funcLit.Name != nil {
					// 関数を環境に登録
					function := &object.Function{
						Parameters: convertToObjectIdentifiers(funcLit.Parameters),
						ASTBody:    funcLit.Body,
						Env:        env,
						InputType:  funcLit.InputType,
						ReturnType: funcLit.ReturnType,
						Condition:  funcLit.Condition,
					}

					// 関数名を取得
					funcName := funcLit.Name.Value
					logger.Debug("FunctionLiteral: 関数 '%s' の定義を見つけました", funcName)

					// 条件付き関数の場合の特別な処理
					if funcLit.Condition != nil {
						// 既存の同名関数の数をカウント
						existingFuncs := env.GetAllFunctionsByName(funcName)
						uniqueName := fmt.Sprintf("%s#%d", funcName, len(existingFuncs))

						logger.Debug("条件付き関数 '%s' を '%s' として事前登録します", funcName, uniqueName)

						// 特別な名前で登録
						env.Set(uniqueName, function)
					}

					// 通常の名前でも登録
					env.Set(funcName, function)
					registeredCount++
					logger.Debug("関数 '%s' を事前登録しました", funcName)
				}
			}
		}

		// AssignStatement 内の FunctionLiteral を検出（関数を変数に代入するケース）
		if assignStmt, ok := stmt.(*ast.AssignStatement); ok {
			if funcLit, ok := assignStmt.Value.(*ast.FunctionLiteral); ok {
				if ident, ok := assignStmt.Left.(*ast.Identifier); ok {
					logger.Debug("AssignStatement: 関数を変数 '%s' に代入する定義を見つけました", ident.Value)
					
					// 関数を環境に登録
					function := &object.Function{
						Parameters: convertToObjectIdentifiers(funcLit.Parameters),
						ASTBody:    funcLit.Body,
						Env:        env,
						InputType:  funcLit.InputType,
						ReturnType: funcLit.ReturnType,
						Condition:  funcLit.Condition,
					}

					// 代入先の変数名を関数名として使用
					funcName := ident.Value
					env.Set(funcName, function)
					registeredCount++
					logger.Debug("代入式の関数 '%s' を事前登録しました", funcName)
				}
			}
		}
		logger.Debug("")
	}

	// 第二パス: すべてのステートメントを再度走査して、埋もれた関数定義を見つける
	// 特に、コメントの後に出現する可能性がある関数定義を見つけるため
	for _, stmt := range program.Statements {
		// トップレベルの式を探索
		if exprStmt, ok := stmt.(*ast.ExpressionStatement); ok {
			// より複雑な式の中にある関数定義を掘り下げる
			findAndRegisterFunctionsInExpression(exprStmt.Expression, env, &registeredCount)
		}
	}

	logger.Debug("関数の事前登録が完了しました。%d 個の関数を登録しました", registeredCount)
}

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
	
	// 関数の事前登録を実行（設定が有効な場合のみ）
	if config.GlobalConfig.PreregisterFunctions {
		logger.Debug("関数の事前登録機能が有効です")
		// preRegisterFunctions(program, env)
		evaluator.PreregisterFunctions(program, env)
	}

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
