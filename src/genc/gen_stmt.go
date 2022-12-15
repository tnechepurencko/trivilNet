package genc

import (
	"fmt"

	"trivil/ast"
)

var _ = fmt.Printf

func (genc *genContext) genStatementSeq(seq *ast.StatementSeq) {

	for _, s := range seq.Statements {

		switch x := s.(type) {
		case *ast.DeclStatement:
			s := genc.genLocalDecl(x.D)
			genc.c(s)
		case *ast.ExprStatement:
			s := genc.genExpr(x.X)
			genc.c(s + ";")
		case *ast.AssignStatement:
			l := genc.genExpr(x.L)
			r := genc.genExpr(x.R)
			genc.c(l + "=" + r + ";")
		case *ast.IncStatement:
			l := genc.genExpr(x.L)
			genc.c(l + "++;")
		case *ast.DecStatement:
			l := genc.genExpr(x.L)
			genc.c(l + "--;")
		case *ast.While:
			genc.genWhile(x)
		case *ast.Return:
			r := ""
			if x.X != nil {
				r = " " + genc.genExpr(x.X)
			}
			genc.c("return" + r + ";")

		default:
			panic(fmt.Sprintf("gen statement: ni %T", s))

		}
	}
}

func (genc *genContext) genWhile(x *ast.While) {
	genc.c("while (%s) {", genc.genExpr(x.Cond))
	genc.genStatementSeq(x.Seq)
	genc.c("}")
}
