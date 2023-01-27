package check

import (
	"fmt"
	"strconv"

	"trivil/ast"
	"trivil/env"
)

var _ = fmt.Printf

func (cc *checkContext) isCheckedType(v *ast.TypeDecl) bool {

	if v.Host != cc.module {
		return true
	}
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

		if f.Later {
			if f.Typ == nil {
				env.AddError(f.Pos, "СЕМ-ДЛЯ-ПОЗЖЕ-НУЖЕН-ТИП")
			}
		} else {
			cc.expr(f.Init)

			if f.Typ != nil {
				cc.checkAssignable(f.Typ, f.Init)
			} else {
				f.Typ = f.Init.GetType()
				if f.Typ == nil {
					panic("assert - не задан тип поля")
				}
			}
		}

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
				} else {
					cl.Members[m.Name] = m
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
		if !equalTypes(p.Typ, ft2.Params[i].Typ) {
			return fmt.Sprintf("не совпадает тип у параметра '%s'", p.Name)
		}
	}

	if !equalTypes(ft1.ReturnTyp, ft2.ReturnTyp) {
		return "разные типы результата"
	}

	return ""
}

// Возвращает "", если равны или причину ошибки
func (cc *checkContext) assignable(lt ast.Type, r ast.Expr) bool {

	cc.errorHint = ""

	if equalTypes(lt, r.GetType()) {
		return true
	}

	var t = ast.UnderType(lt)

	switch t {
	case ast.Byte:
		var li = literal(r)
		if li != nil && li.Kind == ast.Lit_Int {
			i, err := strconv.ParseInt(li.Lit, 0, 64)
			if err == nil && i >= 0 || i <= 255 {
				li.Typ = ast.Byte
				return true
			}
		}
	case ast.TagPair:
		return ast.HasTag(r.GetType())
	}

	switch xt := t.(type) {
	case *ast.ClassType:
		rcl, ok := ast.UnderType(r.GetType()).(*ast.ClassType)
		if ok && isDerivedClass(xt, rcl) {
			return true
		}
	}

	// TODO: function types, целые литералы?, ...
	return false
}

func (cc *checkContext) checkAssignable(lt ast.Type, r ast.Expr) {
	if cc.assignable(lt, r) {
		return
	}
	if ast.IsInvalidType(lt) || ast.IsInvalidType(r.GetType()) {
		return
	}

	env.AddError(r.GetPos(), "СЕМ-НЕСОВМЕСТИМО-ПРИСВ", cc.errorHint,
		ast.TypeName(lt), ast.TypeName(r.GetType()))
}

func equalTypes(t1, t2 ast.Type) bool {

	if tr, ok := t1.(*ast.TypeRef); ok {
		t1 = tr.Typ
	}

	if tr, ok := t2.(*ast.TypeRef); ok {
		t2 = tr.Typ
	}

	return t1 == t2
}
