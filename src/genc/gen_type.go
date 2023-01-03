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
		genc.genClassType(genc.outTypeName(td.Name), x)
	default:
		panic(fmt.Sprintf("getTypeDecl: ni %T", td.Typ))
	}
}

func (genc *genContext) genClassType(tname string, x *ast.ClassType) {
	var st = tname + "Struct"

	var fields = make([]string, len(x.Fields))
	for i, f := range x.Fields {
		fields[i] = fmt.Sprintf("%s %s;",
			genc.typeRef(f.Typ),
			genc.outName(f.Name))
	}
	genc.c("typedef struct %s {", st)

	genc.code = append(genc.code, fields...)

	genc.c("} %s;", st)
	genc.c("typedef %s* %s;", st, tname)
}
