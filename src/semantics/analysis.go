package semantics

import (
	_ "encoding/json" //JSON package from the standard library
	"fmt"
	_ "os"
	"trivil/ast"
	"trivil/env"
	"trivil/semantics/check"
	"trivil/semantics/lookup"
)

var _ = fmt.Printf

func Analyse(m *ast.Module) {
	lookup.Process(m)

	if env.ErrorCount() > 0 {
		return
	}

	check.Process(m)
}
