package runtime

import (
	"fmt"

	"github.com/uncode/ast"
	"github.com/uncode/logger"
	"github.com/uncode/object"
)

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