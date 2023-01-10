package genc

import (
	"fmt"

	"trivil/ast"
)

var _ = fmt.Printf

func (genc *genContext) genStatementSeq(seq *ast.StatementSeq) {

	for _, s := range seq.Statements {
		genc.genStatement(s)
	}
}

func (genc *genContext) genStatement(s ast.Statement) {
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

		var cast = genc.assignCast(x.L.GetType(), x.R.GetType())
		genc.c("%s = %s%s;", l, cast, r)
	case *ast.IncStatement:
		l := genc.genExpr(x.L)
		genc.c(l + "++;")
	case *ast.DecStatement:
		l := genc.genExpr(x.L)
		genc.c(l + "--;")
	case *ast.If:
		genc.genIf(x, "")
	case *ast.While:
		genc.genWhile(x)
	case *ast.Return:
		r := ""
		if x.X != nil {
			r = " " + genc.genExpr(x.X)
		}
		genc.c("return" + r + ";")

	case *ast.Break:
		genc.c("break;")

	default:
		panic(fmt.Sprintf("gen statement: ni %T", s))

	}
}

func (genc *genContext) assignCast(lt, rt ast.Type) string {
	if ast.UnderType(lt) != ast.UnderType(rt) {
		return "(" + genc.typeRef(lt) + ")"
	}
	return ""
}

func (genc *genContext) genIf(x *ast.If, prefix string) {
	genc.c("%sif (%s) {", prefix, genc.genExpr(x.Cond))
	genc.genStatementSeq(x.Then)
	genc.c("}")
	if x.Else != nil {

		elsif, ok := x.Else.(*ast.If)
		if ok {
			genc.genIf(elsif, "else ")
		} else {
			genc.c("else {")
			genc.genStatementSeq(x.Else.(*ast.StatementSeq))
			genc.c("}")
		}
	}
}

func (genc *genContext) genWhile(x *ast.While) {
	genc.c("while (%s) {", genc.genExpr(x.Cond))
	genc.genStatementSeq(x.Seq)
	genc.c("}")
}
