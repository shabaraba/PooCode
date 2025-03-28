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
	
	// 条件なし関数と条件付き関数で処理を分ける
	if fn.Condition == nil {
		// ===== デフォルト関数（条件なし）の処理 =====
		
		// 特殊名を作成
		defaultName := fmt.Sprintf("%s#default", funcName)
		
		// 詳細なデバッグ情報
		logger.Debug("関数事前登録: デフォルト関数 '%s' の登録を開始します（条件なし）", funcName)
		logger.Debug("関数事前登録: 関数情報 - InputType='%s', ReturnType='%s', Parameters=%d個",
			fn.InputType, fn.ReturnType, len(fn.Parameters))
		
		// デフォルト関数として専用の名前で登録（これは必ず行う）
		createAndRegisterFunction(fn, defaultName, env)
		logger.Debug("関数事前登録: デフォルト関数 '%s' を特殊名 '%s' で登録しました", 
			funcName, defaultName)
		
		// 登録済みマップに記録（特殊名）
		registered[defaultName] = true
		
		// 環境を検査して登録状況を確認
		if obj, ok := env.Get(defaultName); ok {
			logger.Debug("関数事前登録: 登録確認 - '%s' が環境に存在します: %T", defaultName, obj)
		} else {
			logger.Debug("関数事前登録: 警告 - '%s' が環境に存在しません", defaultName)
		}
		
		// 元の名前でも登録する（これはどんな場合でも実行）
		// デフォルト関数が優先的に呼ばれるようにするため
		if obj, ok := env.Get(defaultName); ok {
			// 元の名前がすでに登録されている場合でも上書きする
			// これにより、条件なし関数が優先される
			env.Set(funcName, obj)
			logger.Debug("関数事前登録: デフォルト関数 '%s' を元の名前でも登録/上書きしました", funcName)
			
			// 登録済みマップに記録（元の名前）
			registered[funcName] = true
			
			// 登録状況の確認
			if regObj, regOk := env.Get(funcName); regOk {
				logger.Debug("関数事前登録: 登録確認 - '%s' が環境に存在します: %T", funcName, regObj)
			} else {
				logger.Debug("関数事前登録: 警告 - '%s' が環境に存在しません", funcName)
			}
		} else {
			logger.Debug("関数事前登録: 警告 - デフォルト関数の元の名前登録に失敗しました: %s", funcName)
		}
		
	} else {
		// ===== 条件付き関数の処理 =====
		
		// すでに登録済みかチェック
		if _, exists := registered[funcName]; exists {
			logger.Debug("関数 '%s' は既に登録されています", funcName)
			
			// 同名の既存関数を取得して、条件付き関数の数をカウント
			existingFuncs := env.GetAllFunctionsByName(funcName)
			
			// 条件付き関数のカウント
			condFuncCount := 0
			for _, f := range existingFuncs {
				if f.Condition != nil {
					condFuncCount++
				}
			}
			
			// ユニークな名前を生成
			uniqueName := fmt.Sprintf("%s#%d", funcName, condFuncCount)
			
			// 特殊名で関数オブジェクトを作成して登録
			createAndRegisterFunction(fn, uniqueName, env)
			logger.Debug("関数事前登録: 条件付き関数 '%s' を '%s' として追加登録しました", 
				funcName, uniqueName)
				
			// 元の名前としても関連付け（#default 名がない場合のみ）
			// デフォルト関数の方が優先度が高い場合は上書きしない
			defaultKey := fmt.Sprintf("%s#default", funcName)
			if _, hasDefault := env.Get(defaultKey); !hasDefault {
				if obj, ok := env.Get(uniqueName); ok {
					env.Set(funcName, obj)
					logger.Debug("関数事前登録: 条件付き関数 '%s' を元の名前でも関連付けました", funcName)
				}
			}
			
			// 登録済みマップに記録
			registered[uniqueName] = true
			
		} else {
			// 初めての条件付き関数を登録
			
			// 元の名前で登録
			createAndRegisterFunction(fn, funcName, env)
			logger.Debug("関数事前登録: 初めての条件付き関数 '%s' を登録しました", funcName)
			
			// 特殊名でも登録
			uniqueName := fmt.Sprintf("%s#0", funcName)
			if obj, ok := env.Get(funcName); ok {
				env.Set(uniqueName, obj)
				logger.Debug("関数事前登録: 条件付き関数 '%s' を '%s' としても登録しました", 
					funcName, uniqueName)
			}
			
			// 登録済みマップに記録
			registered[funcName] = true
			registered[uniqueName] = true
		}
	}
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
