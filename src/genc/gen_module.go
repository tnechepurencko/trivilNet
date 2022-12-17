package genc

import (
	"fmt"
	"strings"

	"trivil/ast"
)

var _ = fmt.Printf

func (genc *genContext) genModule() {

	//=== gen types

	//=== gen decls

	//=== gen functions
	for _, d := range genc.module.Decls {
		f, ok := d.(*ast.Function)
		if ok {
			genc.genFunction(f)
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

	genc.c("%s %s(%s) {", genc.returnType(ft), genc.outName(f.Name), genc.params(ft))

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

func (genc *genContext) genEntry(entry *ast.EntryFn, main bool) {

	if !main {
		panic("ni")
	}

	genc.c("int main() {")

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
		return fmt.Sprintf("%s %s;", genc.typeRef(x.Typ), genc.outName(x.Name))
	default:
		panic(fmt.Sprintf("genDecl: ni %T", d))
	}
}

func (genc *genContext) typeRef(t ast.Type) string {
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
