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

	// 関数名のマップを作成して二重登録を防止
	registeredFunctions := make(map[string]bool)
	// 第一パス: すべてのトップレベル関数定義を処理
	for i, stmt := range program.Statements {
		logger.Debug("関数事前登録: ステートメント %d を処理中 (%T)", i+1, stmt)
		registerFunctionsInStatement(stmt, env, registeredFunctions)
	}
	// 第二パス: すべてのステートメント内のネストされた関数定義を再帰的に処理
	logger.Debug("関数事前登録: 第二パス - ネストされた関数定義を検索")
	// すべてのステートメントを走査
	for i, stmt := range program.Statements {
		logger.Debug("関数事前登録: ネスト走査 - ステートメント %d (%T)", i+1, stmt)
		findNestedFunctions(stmt, env, registeredFunctions)
	}
	
	logger.Debug("関数事前登録: 完了 - すべての関数が登録されました")
	
	// 登録された関数の数を表示
	totalFunctions := 0
	for range registeredFunctions {
		totalFunctions++
	}
	logger.Debug("関数事前登録: 合計 %d 個のユニークな関数名が登録されました", totalFunctions)
}

// registerFunctionsInStatement はステートメント内の関数定義を処理して登録する
func registerFunctionsInStatement(stmt ast.Statement, env *object.Environment, registered map[string]bool) {
	if stmt == nil {
		return
	}

	// デバッグ出力
	logger.Debug("関数事前登録: ステートメント %T を処理しています", stmt)

	switch s := stmt.(type) {
	case *ast.ExpressionStatement:
		if s.Expression == nil {
			return
		}
		
		// 式文の中身が関数リテラルの場合
		if fn, ok := s.Expression.(*ast.FunctionLiteral); ok {
			registerFunction(fn, env, registered)
		}
		
	case *ast.AssignStatement:
		// 代入文の右辺が関数リテラルの場合
		if fn, ok := s.Value.(*ast.FunctionLiteral); ok {
			// 関数名が未設定の場合、左辺の識別子を関数名として使用
			if fn.Name == nil {
				if ident, ok := s.Left.(*ast.Identifier); ok {
					fn.Name = ident
					logger.Debug("関数事前登録: 代入文から関数名 '%s' を設定しました", ident.Value)
				}
			}
			registerFunction(fn, env, registered)
		}
	}
}

// registerFunction は関数リテラルを環境に登録する
func registerFunction(fn *ast.FunctionLiteral, env *object.Environment, registered map[string]bool) {
	if fn == nil || fn.Name == nil || fn.Name.Value == "" {
		return
	}
	
	// 関数名を取得
	funcName := fn.Name.Value
	
	// すでに登録済みかチェック
	if _, exists := registered[funcName]; exists {
		logger.Debug("関数 '%s' は既に登録されています", funcName)
		// 条件付き関数の場合は名前を変えて追加登録する
		if fn.Condition != nil {
			existingFuncs := env.GetAllFunctionsByName(funcName)
			uniqueName := fmt.Sprintf("%s#%d", funcName, len(existingFuncs))
			
			// 関数オブジェクトを作成して登録
			createAndRegisterFunction(fn, uniqueName, env)
			
			// 元の名前としても登録（上書きではなく関連付け）
			if obj, ok := env.Get(uniqueName); ok {
				env.Set(funcName, obj)
			}
			
			// 登録済みマップに記録
			uniqueKey := fmt.Sprintf("%s#%d", funcName, len(existingFuncs))
			registered[uniqueKey] = true
			
			logger.Debug("関数事前登録: 条件付き関数 '%s' を '%s' として追加登録しました", 
				funcName, uniqueName)
		}
		return
	}
	
	// 関数オブジェクトを作成
	createAndRegisterFunction(fn, funcName, env)
	
	// 登録済みマップに記録
	registered[funcName] = true
	
	// 条件付き関数の場合、個別の名前でも登録
	if fn.Condition != nil {
		uniqueName := fmt.Sprintf("%s#0", funcName)
		if obj, ok := env.Get(funcName); ok {
			env.Set(uniqueName, obj)
			logger.Debug("関数事前登録: 条件付き関数 '%s' を '%s' としても登録しました", 
				funcName, uniqueName)
		}
	}
	
	logger.Debug("関数事前登録: 関数 '%s' を登録しました", funcName)
}

// createAndRegisterFunction は関数オブジェクトを生成して環境に登録する
func createAndRegisterFunction(fn *ast.FunctionLiteral, name string, env *object.Environment) {
	// パラメータをオブジェクト形式に変換
	params := make([]*object.Identifier, len(fn.Parameters))
	for i, p := range fn.Parameters {
		params[i] = &object.Identifier{Value: p.Value}
	}
	
	// 関数オブジェクトを作成
	function := &object.Function{
		Parameters: params,
		ASTBody:    fn.Body,
		Env:        env,
		InputType:  fn.InputType,
		ReturnType: fn.ReturnType,
		Condition:  fn.Condition,
	}
	
	// 環境に登録
	env.Set(name, function)
}

// findNestedFunctions はステートメント内にネストされた関数定義を再帰的に検索する
func findNestedFunctions(node interface{}, env *object.Environment, registered map[string]bool) {
	if node == nil {
		return
	}
	
	switch n := node.(type) {
	case *ast.Program:
		// プログラム内のすべてのステートメントを処理
		for _, stmt := range n.Statements {
			findNestedFunctions(stmt, env, registered)
		}
		
	case *ast.BlockStatement:
		// ブロック内のすべてのステートメントを処理
		for _, stmt := range n.Statements {
			findNestedFunctions(stmt, env, registered)
		}
		
	case *ast.ExpressionStatement:
		// 式文の式を処理
		findNestedFunctions(n.Expression, env, registered)
		
	case *ast.AssignStatement:
		// 代入文の左辺と右辺を処理
		findNestedFunctions(n.Left, env, registered)
		findNestedFunctions(n.Value, env, registered)
		
		// 右辺が関数リテラルの場合
		if fn, ok := n.Value.(*ast.FunctionLiteral); ok {
			// 関数名が未設定の場合、左辺の識別子を関数名として使用
			if fn.Name == nil {
				if ident, ok := n.Left.(*ast.Identifier); ok {
					fn.Name = ident
					logger.Debug("関数事前登録(ネスト): 代入文から関数名 '%s' を設定しました", ident.Value)
				}
			}
			registerFunction(fn, env, registered)
		}
		
	case *ast.FunctionLiteral:
		// 関数リテラルを登録
		if n.Name != nil && n.Name.Value != "" {
			logger.Debug("関数事前登録(ネスト): 関数リテラル '%s' を発見", n.Name.Value)
			registerFunction(n, env, registered)
		}
		
		// 関数本体内もネストされた関数を検索
		if n.Body != nil {
			findNestedFunctions(n.Body, env, registered)
		}
		
	case *ast.PrefixExpression:
		// 前置式の右辺を処理
		findNestedFunctions(n.Right, env, registered)
		
	case *ast.InfixExpression:
		// 中置式の左辺と右辺を処理
		findNestedFunctions(n.Left, env, registered)
		findNestedFunctions(n.Right, env, registered)
		
	// ifを使ったケースを削除または修正
		
	case *ast.CallExpression:
		// 関数式と引数を処理
		findNestedFunctions(n.Function, env, registered)
		for _, arg := range n.Arguments {
			findNestedFunctions(arg, env, registered)
		}
		
	case *ast.IndexExpression:
		// 配列とインデックス式を処理
		findNestedFunctions(n.Left, env, registered)
		findNestedFunctions(n.Index, env, registered)
		
	case *ast.ArrayLiteral:
		// 配列の各要素を処理
		for _, elem := range n.Elements {
			findNestedFunctions(elem, env, registered)
		}
	}
}
