package genc

import (
	"fmt"
	"strings"

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

// genClassType - forward VT
// genClassDesc - before entry
func (genc *genContext) genClassType(td *ast.TypeDecl, x *ast.ClassType) {

	var tname = genc.outTypeName(td.Name)
	var tname_st = tname + nm_class_struct_suffix
	var vt_type = tname + nm_VT_suffix

	var fields = make([]string, len(x.Fields))
	for i, f := range x.Fields {
		fields[i] = fmt.Sprintf("%s %s;",
			genc.typeRef(f.Typ),
			genc.outName(f.Name))
	}

	genc.c("typedef struct %s {", tname_st)
	if x.BaseTyp != nil {
		genc.c("%s%s _B;", genc.typeRef(x.BaseTyp), nm_class_struct_suffix)
	}
	genc.code = append(genc.code, fields...)
	genc.c("} %s;", tname_st)

	genc.c("struct %s;", vt_type)

	genc.c("typedef struct { struct %s* %s; %s %s;} *%s;", vt_type, nm_VT_field, tname_st, nm_class_fields, tname)

	genc.c("")
}

func (genc *genContext) genClassDesc(td *ast.TypeDecl, x *ast.ClassType) {

	var tname = genc.outTypeName(td.Name)
	var tname_st = tname + "_ST"
	var meta_type = tname + nm_meta_suffix
	var vt_type = tname + nm_VT_suffix

	genc.genMeta(x, meta_type)

	var vtable = collectVTable(x)

	genc.genVTable(x, vtable, tname, meta_type, vt_type)
	genc.genClassInit(x, vtable, tname, tname_st, meta_type, vt_type)
}

func (genc *genContext) genMeta(x *ast.ClassType, meta_type string) {

	genc.c("typedef struct %s {", meta_type)
	genc.c("size_t object_size;")
	genc.c("void* base;")

	genc.c("} %s;", meta_type)
}

func (genc *genContext) genVTable(x *ast.ClassType, vtable []*ast.Function, tname, meta_type, vt_type string) {

	genc.c("typedef struct %s {", vt_type)
	genc.c("size_t self_size;")

	for _, f := range vtable {
		genc.c("%s", genc.genMethodField(f, tname))
	}

	genc.c("} %s;", vt_type)
}

func (genc *genContext) genMethodField(f *ast.Function, tname string) string {
	var ft = f.Typ.(*ast.FuncType)

	var ps = make([]string, len(ft.Params)+1)

	ps[0] = genc.typeRef(f.Recv.Typ)
	for i, p := range ft.Params {
		ps[i+1] = genc.typeRef(p.Typ)
	}

	return fmt.Sprintf("%s (*%s)(%s);",
		genc.returnType(ft),
		genc.outName(f.Name),
		strings.Join(ps, ", "))
}

func (genc *genContext) genClassInit(x *ast.ClassType, vtable []*ast.Function, tname, tname_st, meta_type, vt_type string) {

	var desc_var = tname + nm_desc_var_suffix
	genc.c("struct { %s vt; %s meta; } %s;", vt_type, meta_type, desc_var)

	var meta_init_fn = tname + "_init"

	genc.c("void %s() {", meta_init_fn)

	//-- Meta
	var base = "NULL"
	if x.BaseTyp != nil {
		base = "&" + genc.typeRef(x.BaseTyp) + nm_desc_var_suffix
	}
	genc.c("%s.meta.object_size = sizeof(%s);", desc_var, tname_st)
	genc.c("%s.meta.base = %s;", desc_var, base)

	//-- VT
	genc.c("%s.vt.self_size = sizeof(%s);", desc_var, vt_type)

	for _, f := range vtable {
		genc.c("%s.vt.%s = &%s;", desc_var, genc.outName(f.Name), genc.outFnName(f))
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

	//fmt.Printf("! len = %d\n", len(vtable))

	return vtable
}

func addMethodsToVT(vtable []*ast.Function, cl, sub *ast.ClassType) []*ast.Function {

	for _, m := range sub.Methods {

		d, ok := cl.Members[m.Name]
		if !ok {
			panic("assert")
		}
		vtable = append(vtable, d.(*ast.Function))
		//fmt.Printf("! add %s\n", d.GetName())
	}
	return vtable
}
