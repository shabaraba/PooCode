package runtime

import (
	"fmt"
	"os"

	"github.com/uncode/ast"
	"github.com/uncode/config"
	"github.com/uncode/evaluator"
	"github.com/uncode/lexer"
	"github.com/uncode/logger"
	"github.com/uncode/object"
	"github.com/uncode/parser"
	"github.com/uncode/token"
)

// SourceCodeResult は処理結果を表す構造体
type SourceCodeResult struct {
	Tokens   []token.Token
	Program  *ast.Program
	Result   object.Object
	ExitCode int
}