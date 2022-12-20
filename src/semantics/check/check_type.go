package check

import (
	"fmt"
	"trivil/ast"
	"trivil/env"
)

var _ = fmt.Printf

func (cc *checkContext) isCheckedType(v *ast.TypeDecl) bool {
	// if other module return true
	_, ok := cc.checkedTypes[v.Name]
	return ok
}

func (cc *checkContext) typeDecl(td *ast.TypeDecl) {

	switch x := td.Typ.(type) {
	case *ast.InvalidType:
		// nothing
	case *ast.VectorType:
		// есть ли что проверять?
	case *ast.ClassType:
		cc.classType(td, x)
	default:
		panic(fmt.Sprintf("check typeDecl: ni %T", td.Typ))
	}
}

func (cc *checkContext) classType(td *ast.TypeDecl, cl *ast.ClassType) {

	if cc.isCheckedType(td) {
		return
	}
	cc.checkedTypes[td.Name] = struct{}{}

	var members = make(map[string]interface{})

	cc.classBaseType(cl, members)

}

func (cc *checkContext) classBaseType(cl *ast.ClassType, members map[string]interface{}) {
	if cl.BaseTyp == nil {
		return
	}

	var tr = cl.BaseTyp.(*ast.TypeRef)

	base, ok := tr.Typ.(*ast.ClassType)
	if !ok {
		env.AddError(tr.Pos, "СЕМ-БАЗА-НЕ-КЛАСС")
		return
	}

	if base != nil {
	}

}
