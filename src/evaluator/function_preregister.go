package evaluator

import (
	"fmt"

	"github.com/uncode/ast"
	"github.com/uncode/logger"
	"github.com/uncode/object"
)

// PreregisterFunctions traverses the AST to find and register all function declarations before execution.
// This ensures that functions are available regardless of their position in the code.
func PreregisterFunctions(program *ast.Program, env *object.Environment) {
	logger.Debug("関数事前登録: プログラムノードを処理中...")
	
	// デバッグレベルが有効な場合はプログラムの概要を表示
	if logger.IsDebugEnabled() {
		logger.Debug("プログラムの概要:")
		logger.Debug("  ステートメント数: %d", len(program.Statements))
	}
	
	// プログラム内のすべてのステートメントを走査
	for i, stmt := range program.Statements {
		logger.Debug("関数事前登録: ステートメント %d を処理中 (%T)", i+1, stmt)
		
		// 式文の場合
		if exprStmt, ok := stmt.(*ast.ExpressionStatement); ok {
			expr := exprStmt.Expression
			logger.Debug("関数事前登録: 式文の中身を処理中 (%T)", expr)
			
			// 関数リテラル式の場合
			if fnLiteral, ok := expr.(*ast.FunctionLiteral); ok && fnLiteral.Name != nil {
				logger.Debug("関数事前登録: 関数リテラル '%s' を登録します", fnLiteral.Name.Value)
				
				// 関数オブジェクトを作成
				params := make([]*object.Identifier, len(fnLiteral.Parameters))
				for i, p := range fnLiteral.Parameters {
					params[i] = &object.Identifier{Value: p.Value}
				}
				
				function := &object.Function{
					Parameters: params,
					ASTBody:    fnLiteral.Body,
					Env:        env,
					InputType:  fnLiteral.InputType,
					ReturnType: fnLiteral.ReturnType,
					Condition:  fnLiteral.Condition,
				}
				
				// 関数名がある場合は環境に登録
				if fnLiteral.Name.Value != "" {
					logger.Debug("関数事前登録: 関数 '%s' を環境に登録します", fnLiteral.Name.Value)
					
					// 条件付き関数の場合、特別な名前で登録（上書きを防ぐため）
					if fnLiteral.Condition != nil {
						// 既存の同名関数の数をカウント
						existingFuncs := env.GetAllFunctionsByName(fnLiteral.Name.Value)
						uniqueName := fmt.Sprintf("%s#%d", fnLiteral.Name.Value, len(existingFuncs))
						
						logger.Debug("関数事前登録: 条件付き関数 '%s' を '%s' として登録します", 
							fnLiteral.Name.Value, uniqueName)
						
						// 特別な名前で登録
						env.Set(uniqueName, function)
						
						// 検索用に元の名前も関連付け
						env.Set(fnLiteral.Name.Value, function)
					} else {
						// 条件なし関数は通常通り登録
						env.Set(fnLiteral.Name.Value, function)
						logger.Debug("関数事前登録: 条件なし関数 '%s' を登録しました", fnLiteral.Name.Value)
					}
				}
			}
		}
	}
	
	logger.Debug("関数事前登録: 完了 - すべての関数が登録されました")
}
