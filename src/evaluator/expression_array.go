package evaluator

import (
	"fmt"

	"github.com/uncode/ast"
	"github.com/uncode/logger"
	"github.com/uncode/object"
)

// evalRangeExpression は範囲式を評価する
func evalRangeExpression(node *ast.RangeExpression, env *object.Environment) object.Object {
	// デバッグ出力
	logger.ComponentDebug(logger.ComponentEval, "範囲式を評価中: %s", node.String())
	
	// 開始値と終了値を評価
	var startObj, endObj object.Object
	
	if node.Start \!= nil {
		startObj = Eval(node.Start, env)
		if startObj.Type() == object.ERROR_OBJ {
			return startObj
		}
		logger.ComponentTrace(logger.ComponentEval, "範囲式の開始値: %s", startObj.Inspect())
	}
	
	if node.End \!= nil {
		endObj = Eval(node.End, env)
		if endObj.Type() == object.ERROR_OBJ {
			return endObj
		}
		logger.ComponentTrace(logger.ComponentEval, "範囲式の終了値: %s", endObj.Inspect())
	}
	
	// 両方nilの場合（[..]形式）は空の配列を返す
	if node.Start == nil && node.End == nil {
		logger.ComponentDebug(logger.ComponentEval, "範囲式の開始・終了値が両方nilのため、空配列を返します")
		return &object.Array{Elements: []object.Object{}}
	}
	
	// 整数の範囲
	if (startObj == nil || startObj.Type() == object.INTEGER_OBJ) && 
	   (endObj == nil || endObj.Type() == object.INTEGER_OBJ) {
		logger.ComponentDebug(logger.ComponentEval, "整数の範囲を生成します")
		return generateIntegerRange(startObj, endObj)
	}
	
	// 文字列（文字）の範囲
	if (startObj == nil || startObj.Type() == object.STRING_OBJ) && 
	   (endObj == nil || endObj.Type() == object.STRING_OBJ) {
		logger.ComponentDebug(logger.ComponentEval, "文字の範囲を生成します")
		return generateStringRange(startObj, endObj)
	}
	
	// その他の型の範囲はエラー
	logger.ComponentError(logger.ComponentEval, "サポートされていない範囲式の型: %s..%s", 
		getTypeName(startObj), getTypeName(endObj))
	return createError("サポートされていない範囲式の型: %s..%s", 
		getTypeName(startObj), getTypeName(endObj))
}

// generateIntegerRange は整数の範囲を生成する
func generateIntegerRange(startObj, endObj object.Object) object.Object {
	var start, end int64
	
	// 開始値がない場合は0とする
	if startObj == nil {
		start = 0
		logger.ComponentTrace(logger.ComponentEval, "整数範囲の開始値がnilのため、0を使用します")
	} else {
		start = startObj.(*object.Integer).Value
	}
	
	// 終了値がない場合は開始値とする（単一要素の配列）
	if endObj == nil {
		end = start
		logger.ComponentTrace(logger.ComponentEval, "整数範囲の終了値がnilのため、開始値と同じ %d を使用します", start)
	} else {
		end = endObj.(*object.Integer).Value
	}
	
	// 範囲を生成
	var elements []object.Object
	if start <= end {
		// 昇順
		logger.ComponentDebug(logger.ComponentEval, "整数の昇順範囲を生成: %d..%d", start, end)
		for i := start; i <= end; i++ {
			elements = append(elements, &object.Integer{Value: i})
		}
	} else {
		// 降順
		logger.ComponentDebug(logger.ComponentEval, "整数の降順範囲を生成: %d..%d", start, end)
		for i := start; i >= end; i-- {
			elements = append(elements, &object.Integer{Value: i})
		}
	}
	
	return &object.Array{Elements: elements}
}

// generateStringRange は文字（1文字の文字列）の範囲を生成する
func generateStringRange(startObj, endObj object.Object) object.Object {
	var startChar, endChar rune
	
	// 開始値がない場合は'a'とする
	if startObj == nil {
		startChar = 'a'
		logger.ComponentTrace(logger.ComponentEval, "文字範囲の開始値がnilのため、'a'を使用します")
	} else {
		startStr := startObj.(*object.String).Value
		if len(startStr) \!= 1 {
			logger.ComponentError(logger.ComponentEval, "文字範囲の開始値は1文字の文字列である必要があります: %s", startStr)
			return createError("文字範囲の開始値は1文字の文字列である必要があります: %s", startStr)
		}
		startChar = []rune(startStr)[0]
	}
	
	// 終了値がない場合は開始値とする（単一要素の配列）
	if endObj == nil {
		endChar = startChar
		logger.ComponentTrace(logger.ComponentEval, "文字範囲の終了値がnilのため、開始値と同じ '%c' を使用します", startChar)
	} else {
		endStr := endObj.(*object.String).Value
		if len(endStr) \!= 1 {
			logger.ComponentError(logger.ComponentEval, "文字範囲の終了値は1文字の文字列である必要があります: %s", endStr)
			return createError("文字範囲の終了値は1文字の文字列である必要があります: %s", endStr)
		}
		endChar = []rune(endStr)[0]
	}
	
	// 範囲を生成
	var elements []object.Object
	if startChar <= endChar {
		// 昇順
		logger.ComponentDebug(logger.ComponentEval, "文字の昇順範囲を生成: '%c'...'%c'", startChar, endChar)
		for c := startChar; c <= endChar; c++ {
			elements = append(elements, &object.String{Value: string(c)})
		}
	} else {
		// 降順
		logger.ComponentDebug(logger.ComponentEval, "文字の降順範囲を生成: '%c'...'%c'", startChar, endChar)
		for c := startChar; c >= endChar; c-- {
			elements = append(elements, &object.String{Value: string(c)})
		}
	}
	
	return &object.Array{Elements: elements}
}

// evalIndexExpression はインデックス式を評価する
func evalIndexExpression(left, index object.Object, env *object.Environment) object.Object {
	logger.ComponentDebug(logger.ComponentEval, "インデックス式を評価中: %s[%s]", left.Type(), index.Type())
	
	switch {
	case left.Type() == object.ARRAY_OBJ:
		// 配列のインデックスアクセス
		return evalArrayIndexExpression(left, index)
	case left.Type() == object.STRING_OBJ:
		// 文字列のインデックスアクセス
		return evalStringIndexExpression(left, index)
	default:
		logger.ComponentError(logger.ComponentEval, "インデックス演算子はサポートされていません: %s", left.Type())
		return createError("インデックス演算子はサポートされていません: %s", left.Type())
	}
}

// evalArrayIndexExpression は配列のインデックス式を評価する
func evalArrayIndexExpression(array, index object.Object) object.Object {
	arrayObj := array.(*object.Array)
	
	// インデックスが範囲式の場合はスライスを返す
	if rangeExp, ok := index.(*ast.RangeExpression); ok {
		logger.ComponentDebug(logger.ComponentEval, "配列のスライス式を評価中")
		return evalArraySliceExpression(arrayObj, rangeExp)
	}
	
	// 通常のインデックスアクセス
	if idx, ok := index.(*object.Integer); ok {
		logger.ComponentDebug(logger.ComponentEval, "配列の単一インデックスアクセスを評価中: インデックス=%d", idx.Value)
		return evalArraySingleIndex(arrayObj, idx.Value)
	}
	
	logger.ComponentError(logger.ComponentEval, "配列のインデックスは整数である必要があります: %s", index.Type())
	return createError("配列のインデックスは整数である必要があります: %s", index.Type())
}

// evalArraySingleIndex は配列の単一インデックスアクセスを評価する
func evalArraySingleIndex(array *object.Array, index int64) object.Object {
	// 配列の長さ
	length := int64(len(array.Elements))
	
	// 負のインデックスは末尾からのアクセス
	if index < 0 {
		logger.ComponentDebug(logger.ComponentEval, "負のインデックス %d を配列長 %d に対して調整します", index, length)
		index = length + index
	}
	
	// インデックス範囲チェック
	if index < 0 || index >= length {
		logger.ComponentError(logger.ComponentEval, "インデックスが範囲外です: %d (配列長: %d)", index, length)
		return createError("インデックスが範囲外です: %d (配列長: %d)", index, length)
	}
	
	logger.ComponentTrace(logger.ComponentEval, "配列インデックス %d の要素を返します: %s", index, array.Elements[index].Inspect())
	return array.Elements[index]
}

// evalArraySliceExpression は配列のスライス式を評価する
func evalArraySliceExpression(array *object.Array, rangeExp *ast.RangeExpression) object.Object {
	// 配列の長さ
	length := int64(len(array.Elements))
	logger.ComponentDebug(logger.ComponentEval, "配列スライス式を評価中: 配列長=%d", length)
	
	// 開始インデックスと終了インデックスを決定
	var start, end int64
	
	// 開始インデックスの評価
	if rangeExp.Start == nil {
		// [..end] の形式
		start = 0
		logger.ComponentTrace(logger.ComponentEval, "スライスの開始インデックスがnilのため、0を使用します")
	} else if startInt, ok := Eval(rangeExp.Start, nil).(*object.Integer); ok {
		start = startInt.Value
		// 負のインデックスは末尾からの相対位置
		if start < 0 {
			logger.ComponentTrace(logger.ComponentEval, "負のスライス開始インデックス %d を配列長 %d に対して調整します", start, length)
			start = length + start
		}
		// 範囲チェック
		if start < 0 {
			logger.ComponentWarn(logger.ComponentEval, "スライスの開始インデックスが負の値になるため、0に調整します")
			start = 0
		}
	} else {
		logger.ComponentError(logger.ComponentEval, "配列スライスの開始インデックスは整数である必要があります")
		return createError("配列スライスの開始インデックスは整数である必要があります")
	}
	
	// 終了インデックスの評価
	if rangeExp.End == nil {
		// [start..] の形式
		end = length
		logger.ComponentTrace(logger.ComponentEval, "スライスの終了インデックスがnilのため、配列長 %d を使用します", length)
	} else if endInt, ok := Eval(rangeExp.End, nil).(*object.Integer); ok {
		end = endInt.Value
		// 負のインデックスは末尾からの相対位置
		if end < 0 {
			logger.ComponentTrace(logger.ComponentEval, "負のスライス終了インデックス %d を配列長 %d に対して調整します", end, length)
			end = length + end
		}
		// 範囲チェック
		if end > length {
			logger.ComponentWarn(logger.ComponentEval, "スライスの終了インデックスが配列長を超えるため、配列長 %d に調整します", length)
			end = length
		}
	} else {
		logger.ComponentError(logger.ComponentEval, "配列スライスの終了インデックスは整数である必要があります")
		return createError("配列スライスの終了インデックスは整数である必要があります")
	}
	
	// 終了インデックスは含まない（スライス表記に合わせる）
	if start > end {
		logger.ComponentWarn(logger.ComponentEval, "スライスの開始インデックス %d が終了インデックス %d より大きいため、空配列を返します", start, end)
		return &object.Array{Elements: []object.Object{}}
	}
	
	logger.ComponentDebug(logger.ComponentEval, "配列スライス: インデックス %d から %d まで", start, end)
	
	// 新しい配列を作成
	elements := make([]object.Object, end-start)
	copy(elements, array.Elements[start:end])
	
	return &object.Array{Elements: elements}
}

// evalStringIndexExpression は文字列のインデックス式を評価する
func evalStringIndexExpression(str, index object.Object) object.Object {
	// 文字列データを取得
	strValue := str.(*object.String).Value
	strRunes := []rune(strValue)
	length := int64(len(strRunes))
	
	logger.ComponentDebug(logger.ComponentEval, "文字列インデックス式を評価中: 文字列='%s', 長さ=%d", strValue, length)
	
	// インデックスが整数かチェック
	idx, ok := index.(*object.Integer)
	if \!ok {
		logger.ComponentError(logger.ComponentEval, "文字列のインデックスは整数である必要があります: %s", index.Type())
		return createError("文字列のインデックスは整数である必要があります: %s", index.Type())
	}
	
	// インデックス値を取得
	i := idx.Value
	
	// 負のインデックスは末尾からのアクセス
	if i < 0 {
		logger.ComponentTrace(logger.ComponentEval, "負のインデックス %d を文字列長 %d に対して調整します", i, length)
		i = length + i
	}
	
	// インデックス範囲チェック
	if i < 0 || i >= length {
		logger.ComponentError(logger.ComponentEval, "インデックスが範囲外です: %d (文字列長: %d)", i, length)
		return createError("インデックスが範囲外です: %d (文字列長: %d)", i, length)
	}
	
	logger.ComponentTrace(logger.ComponentEval, "文字列インデックス %d の文字を返します: '%c'", i, strRunes[i])
	// 文字列の1文字を返す
	return &object.String{Value: string(strRunes[i])}
}

// getTypeName はオブジェクトの型名を返す（nilの場合は「未指定」）
func getTypeName(obj object.Object) string {
	if obj == nil {
		return "未指定"
	}
	return string(obj.Type())
}
