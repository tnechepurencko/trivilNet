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
		genc.c("typedef struct %s { TInt64 len; %s* body; } %s;", desc, et, desc)
		genc.c("typedef %s* %s;", desc, tname)
	case *ast.ClassType:
		genc.genClassType(td, x)
	default:
		panic(fmt.Sprintf("getTypeDecl: ni %T", td.Typ))
	}
}

func (genc *genContext) genClassType(td *ast.TypeDecl, x *ast.ClassType) {

	var vtableType = genc.genVTable(td.Name, x)

	var tname = genc.outTypeName(td.Name)
	var st = tname + "Struct"

	var fields = make([]string, len(x.Fields)+1)
	fields[0] = fmt.Sprintf("%s %s;", vtableType, nm_VT_field)

	for i, f := range x.Fields {
		fields[i+1] = fmt.Sprintf("%s %s;",
			genc.typeRef(f.Typ),
			genc.outName(f.Name))
	}
	genc.c("typedef struct %s {", st)

	genc.code = append(genc.code, fields...)

	genc.c("} %s;", st)
	genc.c("typedef %s* %s;", st, tname)
}

func (genc *genContext) genVTable(name string, x *ast.ClassType) string {
	var tname = genc.outTypeName(name)
	var vt_name = tname + nm_VT_suffix
	var vt_desc = vt_name + "Desc"
	var meta_name = tname + nm_meta_suffix
	var meta_desc = meta_name + "Desc"

	genc.c("typedef struct %s { size_t sz; } %s;", meta_desc, meta_desc)
	genc.c("typedef %s* %s;", meta_desc, meta_name)

	//init meta

	genc.c("typedef struct %s { %s %s; } %s;",
		vt_desc, meta_name, nm_meta_field, vt_desc)
	genc.c("typedef %s* %s;", vt_desc, vt_name)

	//init vtable

	return vt_name
}
