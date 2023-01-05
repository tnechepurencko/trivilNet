package genc

import (
	"fmt"

	"trivil/ast"
)

var _ = fmt.Printf

func (genc *genContext) typeRef(t ast.Type) string {
	switch x := t.(type) {
	case *ast.PredefinedType:
		return predefinedTypeName(x.Name)
	case *ast.TypeRef:
		if pt, ok := x.Typ.(*ast.PredefinedType); ok {
			return predefinedTypeName(pt.Name)
		}

		if x.ModuleName != "" {
			panic("ni")
		}

		return genc.outName(typeNamePrefix + x.TypeName)

	default:
		panic(fmt.Sprintf("assert: %T", t))
	}

}

func predefinedTypeName(name string) string {
	switch name {
	case "Байт":
		return "TByte"
		//	case "Цел":
		//		return "int"
	case "Цел64":
		return "TInt64"
	case "Вещ64":
		return "TFloat64"
	case "Лог":
		return "TBool"
	case "Символ":
		return "TSymbol"
	case "Строка":
		return "TString"
	default:
		panic(fmt.Sprintf("predefinedTypeName: ni %s", name))
	}
}

func (genc *genContext) genTypeDecl(td *ast.TypeDecl) {
	switch x := td.Typ.(type) {
	case *ast.VectorType:
		var tname = genc.outTypeName(td.Name)
		var desc = tname + "Desc"
		var et = genc.typeRef(x.ElementTyp)
		//TODO meta
		genc.c("typedef struct %s { TInt64 len; %s* body; } %s;", desc, et, desc)
		genc.c("typedef %s* %s;", desc, tname)
	case *ast.ClassType:
		genc.genClassType(td, x)
	default:
		panic(fmt.Sprintf("getTypeDecl: ni %T", td.Typ))
	}
}

func (genc *genContext) genClassType(td *ast.TypeDecl, x *ast.ClassType) {

	var tname = genc.outTypeName(td.Name)
	var tname_st = tname + "_ST"
	var meta_type = tname + nm_meta_suffix
	var vt_type = tname + nm_VT_suffix

	var vtable = collectVTable(x)

	genc.genVTable(x, vtable, meta_type, vt_type)

	var fields = make([]string, len(x.Fields)+1)
	fields[0] = fmt.Sprintf("%s* %s;", vt_type, nm_VT_field)

	for i, f := range x.Fields {
		fields[i+1] = fmt.Sprintf("%s %s;",
			genc.typeRef(f.Typ),
			genc.outName(f.Name))
	}
	genc.c("typedef struct %s {", tname_st)

	genc.code = append(genc.code, fields...)

	genc.c("} %s;", tname_st)
	genc.c("typedef %s* %s;", tname_st, tname)
	genc.c("")

	genc.genClassInit(x, vtable, tname, tname_st, meta_type, vt_type)
}

func (genc *genContext) genVTable(x *ast.ClassType, vtable []*ast.Function, meta_type, vt_type string) {

	genc.c("typedef struct %s { size_t object_size; } %s;", meta_type, meta_type)

	genc.c("typedef struct %s {", vt_type)
	genc.c("size_t self_size;")

	for _, f := range vtable {
		genc.c("%s %s;", "ftype", genc.outName(f.Name))
	}

	genc.c("} %s;", vt_type)
}

func (genc *genContext) genClassInit(x *ast.ClassType, vtable []*ast.Function, tname, tname_st, meta_type, vt_type string) {

	var meta_var = tname + nm_meta_var_suffix
	genc.c("struct { %s vt; %s meta; } %s;", vt_type, meta_type, meta_var)

	var meta_init_fn = tname + "_init"

	genc.c("void %s() {", meta_init_fn)
	genc.c("%s.vt.self_size = sizeof(%s);", meta_var, vt_type)
	genc.c("%s.meta.object_size = sizeof(%s);", meta_var, tname_st)

	for _, f := range vtable {
		genc.c("%s.vt.%s = %s;", meta_var, genc.outName(f.Name), genc.outFnName(f))
	}

	genc.c("}")

	genc.init = append(genc.init, fmt.Sprintf("%s();", meta_init_fn))
}

func collectVTable(x *ast.ClassType) []*ast.Function {

	var vtable = make([]*ast.Function, 0)

	if x.BaseTyp != nil {
		vtable = addMethodsToVT(vtable, x, ast.UnderType(x.BaseTyp).(*ast.ClassType))
	}
	vtable = addMethodsToVT(vtable, x, x)

	fmt.Printf("! len = %d\n", len(vtable))

	return vtable
}

func addMethodsToVT(vtable []*ast.Function, cl, sub *ast.ClassType) []*ast.Function {

	for _, m := range sub.Methods {

		d, ok := cl.Members[m.Name]
		if !ok {
			panic("assert")
		}
		vtable = append(vtable, d.(*ast.Function))
		fmt.Printf("! add %s\n", d.GetName())
	}
	return vtable
}
