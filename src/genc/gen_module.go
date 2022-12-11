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

	genc.genStatementSeq(entry.Seq)

	genc.c("  return 0;")
	genc.c("}")
}

//====

func (genc *genContext) genStatementSeq(seq *ast.StatementSeq) {

	for _, s := range seq.Statements {

		switch x := s.(type) {
		case *ast.ExprStatement:
			s := genc.genExpr(x.X)
			genc.c(s + ";")

		default:
			panic(fmt.Sprintf("gen statement: ni %T", s))

		}
	}

}
