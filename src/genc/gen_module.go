package genc

import (
	"fmt"
	"trivil/ast"
)

var _ = fmt.Printf

func (genc *genContext) genModule() {
	genc.genEntry(genc.module.Entry, true)
}

func (genc *genContext) genEntry(entry *ast.EntryFn, main bool) {

	if !main {
		panic("ni")
	}

	genc.c("int main() {")

	genc.c("  return 0;")
	genc.c("}")
}
