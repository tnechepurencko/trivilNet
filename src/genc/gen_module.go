package genc

import (
	"fmt"
	"trivil/ast"
)

var _ = fmt.Printf

func (genc *genContext) genModule() {
	genc.genEntry(genc.module.Entry, true)
}

//===

func (genc *genContext) genDecl(d ast.Decl) string {
	switch x := d.(type) {
	case *ast.VarDecl:
		return fmt.Sprintf("%s %s;", genc.genTypeRef(x.Typ), x.Name)
	default:
		panic(fmt.Sprintf("genDecl: ni %T", d))
	}
}

func (genc *genContext) genTypeRef(t ast.Type) string {
	var tr = t.(*ast.TypeRef)

	switch x := tr.Typ.(type) {
	case *ast.PredefinedType:
		return predefinedTypeName(x.Name)
	default:
		panic(fmt.Sprintf("genTypeRef: ni %T", tr.Typ))
	}
}

func predefinedTypeName(name string) string {
	switch name {
	case "Байт":
		return "int8_t"
	case "Цел":
		return "int"
	case "Цел64":
		return "int64_t"
	case "Вещ64":
		return "double"
	case "Лог":
		return "_Bool"
	default:
		panic(fmt.Sprintf("predefinedTypeName: ni %s", name))
	}
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
		case *ast.DeclStatement:
			s := genc.genDecl(x.D)
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
