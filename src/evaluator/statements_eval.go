package evaluator

import (
	"fmt"
	
	"github.com/uncode/ast"
	"github.com/uncode/object"
)

// evalProgram はプログラムを評価する
func evalProgram(program *ast.Program, env *object.Environment) object.Object {
	var result object.Object = NullObj

	for _, statement := range program.Statements {
		result = Eval(statement, env)
	}
	
	return result
}

// evalBlockStatement はブロック文を評価する
func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var result object.Object = NullObj

	for _, statement := range block.Statements {
		result = Eval(statement, env)
		
		// 特殊なケース: >>💩 は関数からの戻り値
		if returnValue, ok := result.(*object.ReturnValue); ok {
			return returnValue
		}
		
		// 代入文の場合、PooLiteralへの代入は特別な意味を持つ
		if assignStmt, ok := statement.(*ast.AssignStatement); ok {
			if _, ok := assignStmt.Value.(*ast.PooLiteral); ok {
				fmt.Println("💩への代入を検出しました - 戻り値として扱います")
				// 右辺の値を取得
				rightVal := Eval(assignStmt.Left, env)
				if rightVal.Type() == object.ERROR_OBJ {
					return rightVal
				}
				return &object.ReturnValue{Value: rightVal}
			}
		}
	}
	
	return result
}
