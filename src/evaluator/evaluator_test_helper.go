package evaluator

import (
	"testing"

	"github.com/uncode/lexer"
	"github.com/uncode/object"
	"github.com/uncode/parser"
)

// testEval は入力文字列を評価し、評価結果を返す
func testEval(input string) object.Object {
	l := lexer.NewLexer(input)
	tokens, _ := l.Tokenize()
	p := parser.NewParser(tokens)
	program, _ := p.ParseProgram()
	env := object.NewEnvironment()

	return Eval(program, env)
}

// testIntegerObject は評価結果が期待する整数値であるかチェックする
func testIntegerObject(t *testing.T, obj object.Object, expected int64) bool {
	result, ok := obj.(*object.Integer)
	if !ok {
		t.Errorf("object is not Integer. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%d, want=%d",
			result.Value, expected)
		return false
	}
	return true
}

// testBooleanObject は評価結果が期待する真偽値であるかチェックする
func testBooleanObject(t *testing.T, obj object.Object, expected bool) bool {
	result, ok := obj.(*object.Boolean)
	if !ok {
		t.Errorf("object is not Boolean. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%t, want=%t",
			result.Value, expected)
		return false
	}
	return true
}

// testStringObject は評価結果が期待する文字列であるかチェックする
func testStringObject(t *testing.T, obj object.Object, expected string) bool {
	result, ok := obj.(*object.String)
	if !ok {
		t.Errorf("object is not String. got=%T (%+v)", obj, obj)
		return false
	}
	if result.Value != expected {
		t.Errorf("object has wrong value. got=%q, want=%q",
			result.Value, expected)
		return false
	}
	return true
}

// testNullObject は評価結果がNULLオブジェクトであるかチェックする
func testNullObject(t *testing.T, obj object.Object) bool {
	if obj != NullObj {
		t.Errorf("object is not NULL. got=%T (%+v)", obj, obj)
		return false
	}
	return true
}
