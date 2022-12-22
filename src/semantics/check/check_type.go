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

	if cl.BaseTyp != nil {
		cc.classBaseType(cl, cl.Members)
	}

	for _, f := range cl.Fields {
		prev, ok := cl.Members[f.Name]
		if ok {
			env.AddError(f.Pos, "СЕМ-ДУБЛЬ-В-КЛАССЕ", f.Name, env.PosString(prev.(ast.Node).GetPos()))
		} else {
			cl.Members[f.Name] = f
		}
	}

	for _, m := range cl.Methods {
		prev, ok := cl.Members[m.Name]
		if ok {
			prevM, ok := prev.(*ast.Function)
			if ok && prevM.Recv.Typ != m.Recv.Typ {
				// сигнатуры при переопределении должны совпадать
				var res = cc.compareFuncTypes(m.Typ, prevM.Typ)
				if res != "" {
					env.AddError(m.Pos, "СЕМ-РАЗНЫЕ-ТИПЫ-МЕТОДОВ", m.Name, res)
				}
			} else {
				env.AddError(m.Pos, "СЕМ-ДУБЛЬ-В-КЛАССЕ", m.Name, env.PosString(prev.(ast.Node).GetPos()))
			}
		} else {
			cl.Members[m.Name] = m
		}
	}
}

func (cc *checkContext) classBaseType(cl *ast.ClassType, members map[string]ast.Decl) {

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

	//TODO: Учесть экспорт для типов из другого модуля!
	for _, f := range baseClass.Fields {
		members[f.Name] = f
	}
	for _, m := range baseClass.Methods {
		members[m.Name] = m
	}
}

// Возвращает "", если равны или причину, если разные
func (cc *checkContext) compareFuncTypes(t1, t2 ast.Type) string {
	ft1, ok1 := t1.(*ast.FuncType)
	ft2, ok2 := t2.(*ast.FuncType)
	if !ok1 || !ok2 {
		return "" // а вдруг где-то Invalid type
	}

	if len(ft1.Params) != len(ft2.Params) {
		return "разное число параметров"
	}

	for i, p := range ft1.Params {
		if p.Typ != ft2.Params[i].Typ {
			return fmt.Sprintf("не совпадает тип у параметра '%s'", p.Name)
		}
	}

	if ft1.ReturnTyp != ft2.ReturnTyp {
		return "разные типы результата"
	}

	return ""
}