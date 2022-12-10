package analyzer

import (
	"fmt"
	"trivil/ast"
)

var _ = fmt.Printf

func Analyse(m *ast.Module) {
	lookup(m)

}
