package check

import (
	"fmt"
	"trivil/ast"
	//	"trivil/env"
)

var _ = fmt.Printf

func (cc *checkContext) typeDecl(v *ast.TypeDecl) {

	switch x := v.Typ.(type) {
	case *ast.InvalidType:
		// nothing
	case *ast.VectorType:
		// есть ли что проверять?
	case *ast.ClassType:
		cc.classType(x)
	default:
		panic(fmt.Sprintf("check typeDecl: ni %T", v.Typ))
	}
}

func (cc *checkContext) classType(v *ast.ClassType) {

	/*
		if x.BaseTyp != nil {
			lc.lookTypeRef(x.BaseTyp)
		}
		for _, f := range x.Fields {
			lc.lookTypeRef(f.Typ)
		}
	*/

}
