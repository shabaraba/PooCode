package object

import (
	"bytes"
	"fmt"
	"strings"

	"github.com/uncode/ast"
)

const (
	FUNCTION_OBJ = "FUNCTION"
)

// Function represents a function object
type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Type() ObjectType { return FUNCTION_OBJ }
func (f *Function) Inspect() string {
	var out bytes.Buffer

	params := []string{}
	for _, p := range f.Parameters {
		params = append(params, p.String())
	}

	out.WriteString("fn")
	out.WriteString("(")
	out.WriteString(strings.Join(params, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}

// Closure represents a closure (a function with captured environment)
type Closure struct {
	Function *Function
	Env      *Environment
}

func (c *Closure) Type() ObjectType { return FUNCTION_OBJ }
func (c *Closure) Inspect() string {
	return fmt.Sprintf("Closure: %s", c.Function.Inspect())
}
