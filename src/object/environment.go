package object

import (
	"github.com/uncode/logger"
)

// Environment は変数環境を表す
type Environment struct {
	store map[string]Object
	outer *Environment
}

// NewEnvironment は新しい環境を生成する
func NewEnvironment() *Environment {
	s := make(map[string]Object)
	return &Environment{store: s, outer: nil}
}

// NewEnclosedEnvironment は外部環境を持つ新しい環境を生成する
func NewEnclosedEnvironment(outer *Environment) *Environment {
	env := NewEnvironment()
	env.outer = outer
	return env
}

// Get は環境から変数を取得する
func (e *Environment) Get(name string) (Object, bool) {
	obj, ok := e.store[name]
	if !ok && e.outer != nil {
		obj, ok = e.outer.Get(name)
	}
	return obj, ok
}

// Set は環境に変数を設定する
func (e *Environment) Set(name string, val Object) Object {
	e.store[name] = val
	return val
}

// GetNextFunction は同じ名前の別の関数を取得する
// 特に、条件付き関数が false を返した場合に次の候補を探すために使用
func (e *Environment) GetNextFunction(name string, currentFn *Function) *Function {
	// フラットな関数リストを作成
	var functions []*Function
	var collectFunctions func(*Environment)

	collectFunctions = func(env *Environment) {
		// 現在の環境から同名の関数を探す
		if obj, ok := env.store[name]; ok {
			if fn, ok := obj.(*Function); ok {
				// 現在の関数と異なる関数だけを追加
				if fn != currentFn {
					functions = append(functions, fn)
				}
			}
		}

		// 外部環境も検索
		if env.outer != nil {
			collectFunctions(env.outer)
		}
	}

	collectFunctions(e)

	// 見つかった関数がなければnilを返す
	if len(functions) == 0 {
		return nil
	}

	// 最初に見つかった別の関数を返す
	return functions[0]
}

// GetAllFunctionsByName は同名のすべての関数を取得する
// 条件付き関数を正しく呼び出すために使用
func (e *Environment) GetAllFunctionsByName(name string) []*Function {
	var conditionFunctions []*Function  // 条件付き関数用
	var defaultFunction *Function       // デフォルト関数用
	var collectFunctions func(*Environment)

	logger.Debug("GetAllFunctionsByName: 関数 '%s' の検索開始", name)
	
	// すべての変数を表示
	logger.Debug("GetAllFunctionsByName: 現在の環境にある変数一覧:")
	for k, v := range e.GetVariables() {
		if funcObj, ok := v.(*Function); ok {
			hasCondition := "なし"
			if funcObj.Condition != nil {
				hasCondition = "あり"
			}
			logger.Debug("  変数 '%s': 関数オブジェクト (条件=%s)", k, hasCondition)
		} else {
			logger.Debug("  変数 '%s': %s", k, v.Type())
		}
	}
	
	// 既に処理済みの関数オブジェクトを記録するマップ
	processed := make(map[*Function]bool)

	collectFunctions = func(env *Environment) {
		// ステップ1: まず専用のデフォルト関数キーを優先的に検索
		defaultKey := name + "#default"
		logger.Debug("GetAllFunctionsByName: デフォルト関数キー '%s' を検索中", defaultKey)
		
		if obj, ok := env.store[defaultKey]; ok {
			logger.Debug("GetAllFunctionsByName: キー '%s' が環境に存在します", defaultKey)
			if fn, ok := obj.(*Function); ok {
				logger.Debug("GetAllFunctionsByName: '%s' は関数オブジェクトです", defaultKey)
				if !processed[fn] {
					defaultFunction = fn
					processed[fn] = true
					logger.Debug("専用デフォルト関数 '%s' を見つけました (%p)", defaultKey, fn)
				}
			} else {
				logger.Debug("GetAllFunctionsByName: 警告 - '%s' は関数ではありません: %T", defaultKey, obj)
			}
		} else {
			logger.Debug("GetAllFunctionsByName: キー '%s' が環境に存在しません", defaultKey)
		}

		// ステップ2: 全関数を検索して分類
		for key, obj := range env.store {
			// 基本名と一致するか、「name#数字」のパターンで一致するキーを検索
			// #default は既に処理したので除外
			if key == defaultKey {
				continue // 専用デフォルト関数は既に処理済み
			}
			
			logger.Debug("GetAllFunctionsByName: キー '%s' を評価中", key)
			
			if (key == name) || 
			   (len(key) > len(name) && key[:len(name)+1] == name+"#") {
				logger.Debug("GetAllFunctionsByName: キー '%s' がパターンに一致", key)
				
				if fn, ok := obj.(*Function); ok {
					if !processed[fn] {
						processed[fn] = true
						logger.Debug("GetAllFunctionsByName: キー '%s' の関数オブジェクト (%p) を処理", key, fn)

						// 条件の有無でデフォルト関数と条件付き関数に分類
						if fn.Condition == nil {
							// 条件なし関数は特別扱い
							logger.Debug("GetAllFunctionsByName: キー '%s' は条件なし関数です", key)
							
							// デフォルト関数がまだない場合のみ設定
							if defaultFunction == nil {
								defaultFunction = fn
								logger.Debug("通常名の条件なし関数 '%s' をデフォルト関数として登録 (%p)", key, fn)
							} else {
								// デフォルト関数が既にある場合は一般条件付き関数として扱う
								logger.Debug("別の条件なし関数 '%s' を条件付き関数として登録（デフォルト関数は既に存在）", key)
								conditionFunctions = append(conditionFunctions, fn)
							}
						} else {
							// 条件あり関数は条件付き関数として扱う
							logger.Debug("条件付き関数 '%s' を条件付き関数リストに追加 (%p)", key, fn)
							conditionFunctions = append(conditionFunctions, fn)
						}
					} else {
						logger.Debug("GetAllFunctionsByName: キー '%s' の関数オブジェクトは既に処理済み", key)
					}
				} else {
					logger.Debug("GetAllFunctionsByName: キー '%s' は関数オブジェクトではありません: %T", key, obj)
				}
			}
		}

		// 外部環境も検索（ただし、デフォルト関数がすでに見つかっている場合は条件付き関数のみ検索）
		if env.outer != nil {
			collectFunctions(env.outer)
		}
	}

	collectFunctions(e)

	// 結果を構築
	// ステップ1: まず条件付き関数を追加
	var result []*Function
	result = append(result, conditionFunctions...)
	
	// ステップ2: デフォルト関数があれば最後に追加（フォールバックとして）
	if defaultFunction != nil {
		result = append(result, defaultFunction)
		logger.Debug("デフォルト関数を結果に追加: %p", defaultFunction)
	}

	// デバッグ情報（詳細）
	logger.Debug("関数 '%s' の候補を %d 個見つけました（条件付き：%d, デフォルト：%v）", 
		name, len(result), len(conditionFunctions), defaultFunction != nil)
	
	// 詳細なデバッグ情報を出力
	for i, fn := range result {
		condStatus := "なし"
		if fn.Condition != nil {
			condStatus = "あり"
		}
		logger.Debug("  関数候補 %d の詳細: Condition=%v, Addr=%p, 条件=%s", i+1, fn.Condition, fn, condStatus)
		logger.Debug("  条件式判定: %v (nilチェック結果: %v)", fn.Condition, fn.Condition != nil)
		
		if fn.Condition != nil {
			logger.Debug("  関数候補 %d: 条件付き関数として分類（条件式: %v）", i+1, fn.Condition)
		} else {
			logger.Debug("  関数候補 %d: デフォルト関数として分類（条件式なし）", i+1)
		}
		
		// 関数のパラメータ情報
		params := ""
		for _, p := range fn.Parameters {
			if params != "" {
				params += ", "
			}
			params += p.Value
		}
		logger.Debug("    詳細: 入力型=%s, 戻り値型=%s, パラメータ=[%s]", 
			fn.InputType, fn.ReturnType, params)
	}

	return result
}

// GetVariables は環境内のすべての変数を取得する（デバッグ用）
func (e *Environment) GetVariables() map[string]Object {
	// 現在の環境のすべての変数を取得
	vars := make(map[string]Object)

	// まず現在の環境の変数をコピー
	for k, v := range e.store {
		vars[k] = v
	}

	// 外部環境の変数も取得（ただし、内部環境で上書きされている場合は取得しない）
	if e.outer != nil {
		for k, v := range e.outer.GetVariables() {
			if _, exists := vars[k]; !exists {
				vars[k] = v
			}
		}
	}

	return vars
}
