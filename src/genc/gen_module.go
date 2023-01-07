package genc

import (
	"fmt"
	"strings"

	"trivil/ast"
)

var _ = fmt.Printf

func (genc *genContext) genModule() {

	//=== gen types
	for _, d := range genc.module.Decls {
		d, ok := d.(*ast.TypeDecl)
		if ok {
			genc.genTypeDecl(d)
		}
	}

	//=== gen vars, consts

	//=== gen functions
	for _, d := range genc.module.Decls {
		f, ok := d.(*ast.Function)
		if ok {
			genc.genFunction(f)
		}
	}

	//=== gen class desc
	for _, d := range genc.module.Decls {
		if td, ok := d.(*ast.TypeDecl); ok {
			if cl, ok := ast.UnderType(td.GetType()).(*ast.ClassType); ok {
				genc.genClassDesc(td, cl)
			}
		}
	}

	genc.genEntry(genc.module.Entry, true)
}

//=== functions

func (genc *genContext) genFunction(f *ast.Function) {

	if f.External {
		return
	}

	var ft = f.Typ.(*ast.FuncType)

	var receiver string
	if f.Recv != nil {
		receiver = fmt.Sprintf("%s %s",
			genc.typeRef(f.Recv.Typ),
			genc.outName(f.Recv.Name))
		if len(ft.Params) > 0 {
			receiver += ", "
		}
	}

	genc.c("%s %s(%s%s) {",
		genc.returnType(ft),
		genc.outFnName(f),
		receiver,
		genc.params(ft))

	genc.genStatementSeq(f.Seq)

	genc.c("}")
}

func (genc *genContext) returnType(ft *ast.FuncType) string {
	if ft.ReturnTyp == nil {
		return "void"
	} else {
		return genc.typeRef(ft.ReturnTyp)
	}
}

func (genc *genContext) params(ft *ast.FuncType) string {

	var b strings.Builder

	for i, p := range ft.Params {

		b.WriteString(fmt.Sprintf("%s %s", genc.typeRef(p.Typ), genc.outName(p.Name)))
		if i < len(ft.Params)-1 {
			b.WriteRune(',')
		}
	}

	return b.String()
}

func (genc *genContext) outFnName(f *ast.Function) string {

	var name = genc.outName(f.Name)

	if f.Recv != nil {
		name = genc.typeRef(f.Recv.Typ) + "_" + name
	}
	return name
}

//==== entry

func (genc *genContext) genEntry(entry *ast.EntryFn, main bool) {

	if !main {
		panic("ni")
	}

	genc.c("int main() {")

	genc.code = append(genc.code, genc.init...)

	if entry != nil {
		genc.genStatementSeq(entry.Seq)
	}

	genc.c("  return 0;")
	genc.c("}")
}

//===

func (genc *genContext) genLocalDecl(d ast.Decl) string {
	switch x := d.(type) {
	case *ast.VarDecl:

		return fmt.Sprintf("%s %s = %s%s;",
			genc.typeRef(x.Typ),
			genc.outName(x.Name),
			genc.assignCast(x.Typ, x.Init.GetType()),
			genc.genExpr(x.Init))
	default:
		panic(fmt.Sprintf("genDecl: ni %T", d))
	}
}
