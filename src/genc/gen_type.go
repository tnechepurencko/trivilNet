package genc

import (
	"fmt"

	"trivil/ast"
)

var _ = fmt.Printf

func (genc *genContext) typeRef(t ast.Type) string {
	if tr, ok := t.(*ast.TypeRef); ok {
		t = tr.Typ
	}

	switch x := t.(type) {
	case *ast.PredefinedType:
		return predefinedTypeName(x.Name)
	default:
		panic(fmt.Sprintf("genTypeRef: ni %T", t))
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
