package lookup

import (
	"fmt"
	"trivil/ast"
	"trivil/env"
)

var _ = fmt.Printf

//==== ссылка на тип

func (lc *lookContext) lookTypeRef(t ast.Type) {

	var tr, ok = t.(*ast.TypeRef)
	if !ok {
		if t == nil {
			panic("assert")
		}
		vTyp, ok := t.(*ast.VariadicType)
		if ok {
			lc.lookTypeRef(vTyp.ElementTyp)
		}
		return
	}

	if tr.Typ != nil {
		return // уже сделано
	}

	var td *ast.TypeDecl

	if tr.ModuleName != "" {
		td = lc.lookTypeDeclInModule(tr.ModuleName, tr.TypeName, tr.Pos)
	} else {
		td = lc.lookTypeDeclInScopes(tr.TypeName, tr.Pos)
	}

	tr.TypeDecl = td
	tr.Typ = tr.TypeDecl.Typ

	if tr.Typ == nil {
		panic("not resolved")
	}

	//fmt.Printf("! %v %T\n", tr.TypeDecl, tr.Typ)
}

func (lc *lookContext) lookTypeDeclInScopes(name string, pos int) *ast.TypeDecl {

	var d = findInScopes(lc.scope, name, pos)

	if d == nil {
		env.AddError(pos, "СЕМ-НЕ-НАЙДЕНО", name)
		var td = lc.makeTypeDecl(name, pos)
		addToScope(name, td, lc.scope)
		return td
	}

	td, ok := d.(*ast.TypeDecl)
	if !ok {
		env.AddError(pos, "СЕМ-ДОЛЖЕН-БЫТЬ-ТИП", name)
		return lc.makeTypeDecl(name, pos)
	}

	return td
}

func (lc *lookContext) lookTypeDeclInModule(moduleName, name string, pos int) *ast.TypeDecl {
	var d = findInScopes(lc.scope, moduleName, pos)

	if d == nil {
		env.AddError(pos, "СЕМ-НЕ-НАЙДЕН-МОДУЛЬ", moduleName)
		return lc.makeTypeDecl(name, pos)
	}

	m, ok := d.(*ast.Module)
	if !ok {
		env.AddError(pos, "СЕМ-ДОЛЖЕН-БЫТЬ-МОДУЛЬ", moduleName)
		return lc.makeTypeDecl(name, pos)
	}

	d, ok = m.Inner.Names[name]
	if !ok {
		env.AddError(pos, "СЕМ-НЕ-НАЙДЕНО-В-МОДУЛЕ", m.Name, name)
		var td = lc.makeTypeDecl(name, pos)
		addToScope(name, td, lc.scope)
		return td
	}

	td, ok := d.(*ast.TypeDecl)
	if !ok {
		env.AddError(pos, "СЕМ-ДОЛЖЕН-БЫТЬ-ТИП", name)
		return lc.makeTypeDecl(name, pos)
	}

	if !d.IsExported() {
		env.AddError(pos, "СЕМ-НЕ-ЭКСПОРТИРОВАН", name, m.Name)
	}

	return td
}

func (lc *lookContext) makeTypeDecl(name string, pos int) *ast.TypeDecl {
	var td = &ast.TypeDecl{
		DeclBase: ast.DeclBase{
			Pos:      pos,
			Name:     name,
			Typ:      ast.MakeInvalidType(pos),
			Host:     lc.module,
			Exported: true,
		},
	}
	return td
}

//==== типы

func (lc *lookContext) lookTypeDecl(v *ast.TypeDecl) {

	switch x := v.Typ.(type) {
	case *ast.VectorType:
		lc.lookTypeRef(x.ElementTyp)
	case *ast.ClassType:
		if x.BaseTyp != nil {
			lc.lookTypeRef(x.BaseTyp)
		}
		for _, f := range x.Fields {
			if f.Typ != nil {
				lc.lookTypeRef(f.Typ)
			}
			lc.lookExpr(f.Init)
		}

	default:
		panic(fmt.Sprintf("lookTypeDecl: ni %T", v.Typ))
	}
}
