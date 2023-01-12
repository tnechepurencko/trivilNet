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
		return genc.declName(x.TypeDecl)

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
		var tname = genc.declName(td)
		var desc = tname + "Desc"
		var et = genc.typeRef(x.ElementTyp)
		//TODO meta
		genc.h("typedef struct %s { TInt64 len; %s* body; } %s;", desc, et, desc)
		genc.h("typedef %s* %s;", desc, tname)
	case *ast.ClassType:
		genc.genClassType(td, x)
	default:
		panic(fmt.Sprintf("getTypeDecl: ni %T", td.Typ))
	}
}

// genClassType - forward VT
// genClassDesc - before entry
func (genc *genContext) genClassType(td *ast.TypeDecl, x *ast.ClassType) {

	var tname = genc.declName(td)
	var tname_st = tname + nm_class_struct_suffix
	var vt_type = tname + nm_VT_suffix

	var fields = make([]string, len(x.Fields))
	for i, f := range x.Fields {
		fields[i] = fmt.Sprintf("%s %s;",
			genc.typeRef(f.Typ),
			genc.declName(f))
	}

	genc.h("typedef struct %s {", tname_st)
	if x.BaseTyp != nil {
		genc.h("%s%s _B;", genc.typeRef(x.BaseTyp), nm_class_struct_suffix)
	}
	genc.header = append(genc.header, fields...)
	genc.h("} %s;", tname_st)

	genc.h("struct %s;", vt_type)

	genc.h("typedef struct %s { struct %s* %s; %s %s;} *%s;", tname, vt_type, nm_VT_field, tname_st, nm_class_fields, tname)

	genc.h("")
}

func (genc *genContext) genClassDesc(td *ast.TypeDecl, x *ast.ClassType) {

	var tname = genc.declName(td)
	var tname_st = tname + "_ST"
	var meta_type = tname + nm_meta_suffix
	var vt_type = tname + nm_VT_suffix

	genc.genMeta(x, meta_type)

	var col collector

	col.collectVTable(x)

	genc.genVTable(x, col.vtable, tname, meta_type, vt_type)
	genc.genClassInit(x, col.vtable, tname, tname_st, meta_type, vt_type)
}

func (genc *genContext) genMeta(x *ast.ClassType, meta_type string) {

	genc.c("typedef struct %s {", meta_type)
	genc.c("size_t object_size;")
	genc.c("void* base;")

	genc.c("} %s;", meta_type)
}

func (genc *genContext) genVTable(x *ast.ClassType, vtable []*ast.Function, tname, meta_type, vt_type string) {

	genc.h("typedef struct %s {", vt_type)
	genc.h("size_t self_size;")

	for _, f := range vtable {
		genc.h("%s", genc.genMethodField(f, tname))
	}

	genc.h("} %s;", vt_type)
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
		genc.outName(f.Name), // только имя, без префикса модуля
		strings.Join(ps, ", "))
}

func (genc *genContext) genClassInit(x *ast.ClassType, vtable []*ast.Function, tname, tname_st, meta_type, vt_type string) {

	var desc_var = tname + nm_class_info_suffix
	genc.c("struct { %s vt; %s meta; } %s;", vt_type, meta_type, desc_var)

	var ptr = fmt.Sprintf("void * %s;", tname+nm_class_info_ptr_suffix)
	genc.h("extern %s", ptr)
	genc.c("%s", ptr)

	var meta_init_fn = tname + "_init"

	genc.c("void %s() {", meta_init_fn)

	//-- Meta
	var base = "NULL"
	if x.BaseTyp != nil {
		base = genc.typeRef(x.BaseTyp) + nm_class_info_ptr_suffix
	}
	genc.c("%s.meta.object_size = sizeof(struct %s);", desc_var, tname)
	genc.c("%s.meta.base = %s;", desc_var, base)

	//-- VT
	genc.c("%s.vt.self_size = sizeof(%s);", desc_var, vt_type)

	for _, f := range vtable {
		genc.c("%s.vt.%s = &%s;", desc_var, genc.outName(f.Name), genc.functionName(f))
	}

	genc.c("%s = &%s;", tname+nm_class_info_ptr_suffix, desc_var)
	genc.c("}")

	genc.init = append(genc.init, fmt.Sprintf("%s();", meta_init_fn))
}

//=== collect methods

type collector struct {
	cl     *ast.ClassType
	vtable []*ast.Function
	done   map[string]struct{}
}

func (col *collector) collectVTable(x *ast.ClassType) {

	col.cl = x
	col.vtable = make([]*ast.Function, 0)
	col.done = make(map[string]struct{})

	if x.BaseTyp != nil {
		col.addMethods(ast.UnderType(x.BaseTyp).(*ast.ClassType))
	}
	col.addMethods(x)

	//fmt.Printf("! len = %d\n", len(col.vtable))
}

func (col *collector) addMethods(sub *ast.ClassType) {

	for _, m := range sub.Methods {

		d, ok := col.cl.Members[m.Name]
		if !ok {
			panic("assert")
		}

		f := d.(*ast.Function)

		_, ok = col.done[f.Name]
		if !ok {
			col.vtable = append(col.vtable, d.(*ast.Function))
			col.done[f.Name] = struct{}{}
			//fmt.Printf("! add %s\n", f.Name)
		}
	}

}
