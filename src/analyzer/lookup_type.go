package analyzer

import (
	"fmt"
	"trivil/ast"
	//"trivil/env"
)

var _ = fmt.Printf

//==== ссылка на тип

func (lc *lookContext) lookTypeRef(t ast.Type) {
	var tr = t.(*ast.TypeRef)
	if tr.Typ != nil {
		return // уже сделано
	}

	if tr.ModuleName != "" {
		panic("ni")
	}

	var x = findInScopes(lc.scope, tr.TypeName, tr.Pos)
	td, ok := x.(*ast.TypeDecl)
	if !ok {
		return
	}

	tr.TypeDecl = td
	tr.Typ = tr.TypeDecl.Typ

	if tr.Typ == nil {
		panic("not resolved")
	}

	//fmt.Printf("! %v %T\n", tr.TypeDecl, tr.Typ)
}

//==== типы

func (lc *lookContext) lookTypeDecl(v *ast.TypeDecl) {

	switch x := v.Typ.(type) {
	case *ast.DynArrayType:
		lc.lookTypeRef(x.ElementTyp)
	case *ast.ClassType:
		if x.BaseTyp != nil {
			lc.lookTypeRef(x.BaseTyp)
		}
		for _, f := range x.Fields {
			lc.lookTypeRef(f.Typ)
		}

	default:
		panic(fmt.Sprintf("lookTypeDecl: ni %T", v.Typ))
	}
}
