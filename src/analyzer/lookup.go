package analyzer

import (
	"fmt"
	"trivil/ast"
)

var _ = fmt.Printf

func lookup(m *ast.Module) {
	fmt.Printf("lookup\n")
}
