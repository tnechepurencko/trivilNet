package semantics

import (
	"fmt"
	"trivil/ast"
	"trivil/semantics/lookup"
)

var _ = fmt.Printf

func Analyse(m *ast.Module) {
	lookup.Process(m)

}
