package evaluator

import (
	"fmt"

	"github.com/uncode/object"
)

// 型名と対応するObjectTypeのマッピング
var typeMapping = map[string]object.ObjectType{
	"int":    object.INTEGER_OBJ,
	"float":  object.FLOAT_OBJ,
	"bool":   object.BOOLEAN_OBJ,
	"str":    object.STRING_OBJ,
	"null":   object.NULL_OBJ,
	"array":  object.ARRAY_OBJ,
	"hash":   object.HASH_OBJ,
	"class":  object.CLASS_OBJ,
	"object": "", // 任意の型を許可
}

// checkInputType は入力値（🍕）の型をチェックする
func checkInputType(input object.Object, expectedType string) (bool, error) {
	if expectedType == "" || expectedType == "object" {
		// 型指定がないか、任意の型を許可する場合はチェックしない
		return true, nil
	}

	// 期待される型のObjectTypeを取得
	expectedObjType, ok := typeMapping[expectedType]
	if !ok {
		return false, fmt.Errorf("未知の型定義: %s", expectedType)
	}

	// 実際の型をチェック
	actualType := input.Type()
	if actualType != expectedObjType {
		return false, fmt.Errorf("🍕の型が不正です: 期待=%s, 実際=%s", expectedType, mapObjectTypeToName(actualType))
	}

	return true, nil
}

// checkReturnType は戻り値（💩）の型をチェックする
func checkReturnType(result object.Object, expectedType string) (bool, error) {
	if expectedType == "" || expectedType == "object" {
		// 型指定がないか、任意の型を許可する場合はチェックしない
		return true, nil
	}

	// 期待される型のObjectTypeを取得
	expectedObjType, ok := typeMapping[expectedType]
	if !ok {
		return false, fmt.Errorf("未知の型定義: %s", expectedType)
	}

	// 実際の型をチェック
	actualType := result.Type()
	if actualType != expectedObjType {
		return false, fmt.Errorf("💩の型が不正です: 期待=%s, 実際=%s", expectedType, mapObjectTypeToName(actualType))
	}

	return true, nil
}

// mapObjectTypeToName はobject.ObjectTypeから読みやすい型名に変換する
func mapObjectTypeToName(objType object.ObjectType) string {
	switch objType {
	case object.INTEGER_OBJ:
		return "int"
	case object.FLOAT_OBJ:
		return "float"
	case object.BOOLEAN_OBJ:
		return "bool"
	case object.STRING_OBJ:
		return "str"
	case object.NULL_OBJ:
		return "null"
	case object.ARRAY_OBJ:
		return "array"
	case object.HASH_OBJ:
		return "hash"
	case object.CLASS_OBJ:
		return "class"
	case object.INSTANCE_OBJ:
		return "instance"
	case object.FUNCTION_OBJ:
		return "function"
	case object.BUILTIN_OBJ:
		return "builtin"
	default:
		return string(objType)
	}
}
