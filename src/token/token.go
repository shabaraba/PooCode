package token

import (
	"fmt"
)

// TokenType は字句解析で識別されるトークンの種類を表す
type TokenType string

// トークンの種類を定義
const (
	// 特殊トークン
	ILLEGAL = "ILLEGAL" // 不正なトークン
	EOF     = "EOF"     // ファイル終端

	// 識別子・リテラル
	IDENT   = "IDENT"   // 識別子
	INT     = "INT"     // 整数リテラル
	FLOAT   = "FLOAT"   // 浮動小数点リテラル
	STRING  = "STRING"  // 文字列リテラル
	BOOLEAN = "BOOLEAN" // 真偽値リテラル

	// 演算子
	ASSIGN   = ">>" // 代入演算子
	EQUAL    = "="  // = 演算子
	PLUS     = "+"  // 加算
	MINUS    = "-"  // 減算
	ASTERISK = "*"  // 乗算
	SLASH    = "/"  // 除算
	MODULO   = "%"  // 剰余

	// 比較演算子
	EQ     = "==" // 等価
	NOT_EQ = "!=" // 非等価
	LT     = "<"  // 小なり
	GT     = ">"  // 大なり
	LE     = "<=" // 以下
	GE     = ">=" // 以上

	// 論理演算子
	AND   = "&&"  // 論理積
	OR    = "||"  // 論理和
	NOT   = "not" // 論理否定
	BANG  = "!"   // ! 論理否定

	// デリミタ
	COMMA     = "," // カンマ
	SEMICOLON = ";" // セミコロン
	COLON     = ":" // コロン
	LPAREN    = "(" // 左括弧
	RPAREN    = ")" // 右括弧
	LBRACE    = "{" // 左中括弧
	RBRACE    = "}" // 右中括弧
	LBRACKET  = "[" // 左角括弧
	RBRACKET  = "]" // 右角括弧
	DOT       = "." // ドット
	DOTDOT    = ".." // 範囲演算子

	// パイプライン
	PIPE     = "|>" // パイプライン
	PIPE_PAR = "|"  // 並列パイプライン
	
	// 特殊パイプ演算子
	MAP_PIPE = "MAP_PIPE" // map関数をパイプとして使用
	FILTER_PIPE = "FILTER_PIPE" // filter関数をパイプとして使用

	// キーワード
	FUNCTION = "def"     // 関数定義
	CLASS    = "class"   // クラス定義
	IF       = "if"      // 条件分岐
	ELSE     = "else"    // 条件分岐（その他）
	PUBLIC   = "public"  // 公開アクセス修飾子
	PRIVATE  = "private" // 非公開アクセス修飾子
	GLOBAL   = "global"  // グローバル変数
	ENUM     = "enum"    // 列挙型
	EXTENDS  = "extends" // 継承

	// 特殊変数
	PIZZA = "🍕" // 入力値
	POO   = "💩" // 出力値

	// アクセス演算子
	APOSTROPHE_S = "'s" // プロパティアクセス
)

// Token は字句解析で生成されるトークンを表す
type Token struct {
	Type    TokenType // トークンの種類
	Literal string    // トークンのリテラル値
	Line    int       // トークンの行番号
	Column  int       // トークンの列番号
}

// キーワードマップ
var keywords = map[string]TokenType{
	"def":     FUNCTION,
	"class":   CLASS,
	"if":      IF,
	"else":    ELSE,
	"public":  PUBLIC,
	"private": PRIVATE,
	"global":  GLOBAL,
	"enum":    ENUM,
	"extends": EXTENDS,
	"not":     NOT,
	"true":    BOOLEAN,
	"false":   BOOLEAN,
	"eq":      EQ,
	// "add":     PLUS, // 加算を表す関数ですが、直接+に変換するとパイプラインでエラーになるため
	"add":     IDENT, // 関数として扱う
	"print":   IDENT, // print関数を明示的に追加
	"show":    IDENT, // 代替としてshowも追加
	"map":     MAP_PIPE, // mapをパイプとして追加
	"filter":  FILTER_PIPE, // filterをパイプとして追加
}

// LookupIdent は識別子がキーワードかどうかを判定する
func LookupIdent(ident string) TokenType {
	if tok, ok := keywords[ident]; ok {
		return tok
	}
	return IDENT
}

// String はトークンの文字列表現を返す
func (t Token) String() string {
	return fmt.Sprintf("Token{Type: %s, Literal: %q, Line: %d, Column: %d}", 
		t.Type, t.Literal, t.Line, t.Column)
}
