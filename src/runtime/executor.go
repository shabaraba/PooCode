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