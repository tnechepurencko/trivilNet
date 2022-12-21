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

	if cl.BaseTyp != nil {
		cc.classBaseType(cl, members)
	}

	for _, f := range cl.Fields {
		prev, ok := members[f.Name]
		if ok {
			env.AddError(f.Pos, "СЕМ-ДУБЛЬ-В-КЛАССЕ", f.Name, env.PosString(prev.(ast.Node).GetPos()))
		} else {
			members[f.Name] = f
		}
	}

	for _, m := range cl.Methods {
		prev, ok := members[m.Name]
		if ok {
			env.AddError(m.Pos, "СЕМ-ДУБЛЬ-В-КЛАССЕ", m.Name, env.PosString(prev.(ast.Node).GetPos()))
		} else {
			members[m.Name] = m
		}
	}

}

func (cc *checkContext) classBaseType(cl *ast.ClassType, members map[string]interface{}) {

	var tr = cl.BaseTyp.(*ast.TypeRef)

	baseClass, ok := tr.Typ.(*ast.ClassType)
	if !ok {
		env.AddError(tr.Pos, "СЕМ-БАЗА-НЕ-КЛАСС")
		return
	}

	if !cc.isCheckedType(tr.TypeDecl) {
		cc.classType(tr.TypeDecl, baseClass)
	}

	if baseClass.BaseTyp != nil {
		cc.classBaseType(baseClass, members)
	}

	for _, f := range baseClass.Fields {
		members[f.Name] = f
	}
	for _, m := range baseClass.Methods {
		members[m.Name] = m
	}

}
