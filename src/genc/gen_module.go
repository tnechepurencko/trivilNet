package genc

import (
	"fmt"
	"strings"

	"trivil/ast"
)

var _ = fmt.Printf

func (genc *genContext) genModule(main bool) {

	//=== import
	for _, i := range genc.module.Imports {
		genc.h("#include \"%s\"", genc.declName(i.Mod)+".h")
	}
	genc.h("")

	//=== gen types
	genc.genTypes = true
	for _, d := range genc.module.Decls {
		d, ok := d.(*ast.TypeDecl)
		if ok {
			genc.genTypeDecl(d)
		}
	}
	genc.genTypes = false

	//=== gen vars, consts
	for _, d := range genc.module.Decls {
		switch x := d.(type) {
		case *ast.ConstDecl:
			genc.genGlobalConst(x)
		case *ast.VarDecl:
			genc.genGlobalVar(x)
		}
	}

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

	genc.genEntry(genc.module.Entry, main)
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
			genc.declName(f.Recv))
		if len(ft.Params) > 0 {
			receiver += ", "
		}
	}

	var fn_header = fmt.Sprintf("%s %s(%s%s)",
		genc.returnType(ft),
		genc.functionName(f),
		receiver,
		genc.params(ft))

	genc.h("%s;", fn_header)
	genc.c("%s {", fn_header)

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

	var list = make([]string, len(ft.Params))

	for i, p := range ft.Params {

		if ast.IsVariadicType(p.Typ) {
			list[i] = fmt.Sprintf("TInt64 %s%s", genc.declName(p), nm_variadic_len_suffic)
			list = append(list, fmt.Sprintf("void* %s", genc.declName(p)))
		} else {
			list[i] = fmt.Sprintf("%s %s", genc.typeRef(p.Typ), genc.declName(p))
		}
	}
	return strings.Join(list, ", ")
}

//==== entry

const (
	init_fn  = "init"
	init_var = "init_done"
)

func (genc *genContext) genEntry(entry *ast.EntryFn, main bool) {

	if main {
		genc.c("int main(int argc, char *argv[]) {")
		genc.c("%s(argc, argv);", rt_init)
	} else {
		var init_header = fmt.Sprintf("void %s__%s()", genc.outname, init_fn)

		genc.h("%s;", init_header)

		genc.c("static TBool %s = false;", init_var)
		genc.c("%s {", init_header)
		genc.c("if (%s) return;", init_var)
		genc.c("%s = true;", init_var)
	}

	for _, i := range genc.module.Imports {
		genc.c("%s__%s();", genc.declName(i.Mod), init_fn)
	}

	genc.code = append(genc.code, genc.init...)

	if entry != nil {
		genc.genStatementSeq(entry.Seq)
	}

	if main {
		genc.c("  return 0;")
	}
	genc.c("}")
}

//===

func (genc *genContext) genGlobalConst(x *ast.ConstDecl) {

	var name = genc.declName(x)
	var def = fmt.Sprintf("%s %s", genc.typeRef(x.Typ), name)
	var val = genc.genExpr(x.Value)

	if initializeInPlace(x.Typ) {
		def = "const " + def
		genc.c("%s = %s;", def, val)
	} else {
		genc.c("%s;", def)
		genc.init = append(genc.init, fmt.Sprintf("%s = %s;", name, val))
	}

	if x.Exported {
		genc.h("extern %s;", def)
	}
}

func initializeInPlace(t ast.Type) bool {

	t = ast.UnderType(t)
	switch t {
	case ast.Byte, ast.Int64, ast.Float64, ast.Bool, ast.Symbol:
		return true
	}
	return false
}

func (genc *genContext) genGlobalVar(x *ast.VarDecl) {
	panic("ni")
}

func (genc *genContext) genLocalDecl(d ast.Decl) string {
	switch x := d.(type) {
	case *ast.VarDecl:

		return fmt.Sprintf("%s %s = %s%s;",
			genc.typeRef(x.Typ),
			genc.declName(x),
			genc.assignCast(x.Typ, x.Init.GetType()),
			genc.genExpr(x.Init))
	default:
		panic(fmt.Sprintf("genDecl: ni %T", d))
	}
}
