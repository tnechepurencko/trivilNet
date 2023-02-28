package check

import (
	"fmt"
	"trivil/ast"
	"trivil/env"
)

var _ = fmt.Printf

func (cc *checkContext) typeExpr(expr ast.Expr) ast.Type {

	switch x := expr.(type) {
	case *ast.IdentExpr:
		if tr, ok := x.Obj.(*ast.TypeRef); ok {
			return tr
		} else {
			return nil
		}
	case *ast.SelectorExpr:
		if tr, ok := x.Obj.(*ast.TypeRef); ok {
			return tr
		} else {
			return nil
		}
	}

	return nil
}

//==== конструктор вектора

func (cc *checkContext) arrayComposite(c *ast.ArrayCompositeExpr, t ast.Type) {

	var elemT ast.Type = nil

	if t == nil {
		env.AddError(c.Pos, "СЕМ-КОМПОЗИТ-НЕТ-ТИПА")
	} else if !ast.IsIndexableType(t) {
		env.AddError(c.Pos, "СЕМ-МАССИВ-КОМПОЗИТ-ОШ-ТИП")
	} else {
		c.Typ = t
		elemT = ast.ElementType(t)
	}

	if c.Length != nil {
		cc.expr(c.Length)
		cc.checkAssignable(ast.Int64, c.Length)
	}

	if c.Capacity != nil {
		cc.expr(c.Capacity)
		cc.checkAssignable(ast.Int64, c.Capacity)
	}

	if c.Default != nil {
		cc.expr(c.Default)
		if elemT != nil {
			cc.checkAssignable(elemT, c.Default)
		}
	}

	for _, inx := range c.Indexes {
		cc.expr(inx)
		cc.checkAssignable(ast.Int64, inx)
		cc.checkConstExpr(inx)
	}

	for _, val := range c.Values {
		cc.expr(val)
		if elemT != nil {
			cc.checkAssignable(elemT, val)
		}
	}

	// TODO: проверить отсутствие дупликатов
	// TODO: проверить индексы в [0..длина-1]

	cc.arrayCompositeIndexes(c)

	// TODO: добавить тесты
}

func (cc *checkContext) arrayCompositeIndexes(c *ast.ArrayCompositeExpr) {

	// если были ошибки, не пытаюсь проверить индексы и длину
	if env.ErrorCount() > 0 {
		return
	}

	if len(c.Indexes) == 0 {
		return
	}

	/*
		var max = 0
		for _, inx := range c.Indexes {

		}
	*/
}

//==== конструктор класса

func (cc *checkContext) classComposite(c *ast.ClassCompositeExpr) {

	var t = cc.typeExpr(c.X)

	if t == nil {
		env.AddError(c.Pos, "СЕМ-КОМПОЗИТ-НЕТ-ТИПА")
		c.Typ = ast.MakeInvalidType(c.X.GetPos())
		return
	}

	cl, ok := ast.UnderType(t).(*ast.ClassType)
	if !ok {
		env.AddError(c.Pos, "СЕМ-КЛАСС-КОМПОЗИТ-ОШ-ТИП")
		c.Typ = ast.MakeInvalidType(c.X.GetPos())
	} else {
		c.Typ = t
	}

	for _, vp := range c.Values {
		cc.expr(vp.Value)
	}

	if cl == nil {
		return
	}

	// проверяю поля и типы
	var vals = make(map[string]bool)
	for _, vp := range c.Values {
		d, ok := cl.Members[vp.Name]
		if !ok {
			env.AddError(vp.Pos, "СЕМ-КЛАСС-КОМПОЗИТ-НЕТ-ПОЛЯ", vp.Name)
		} else {
			f, ok := d.(*ast.Field)
			if !ok {
				env.AddError(vp.Pos, "СЕМ-КЛАСС-КОМПОЗИТ-НЕ-ПОЛE")
			} else if f.Host != cc.module && !f.Exported {
				env.AddError(vp.Pos, "СЕМ-НЕ-ЭКСПОРТИРОВАН", f.Name, f.Host.Name)
			} else {
				vals[vp.Name] = true
				cc.checkAssignable(f.Typ, vp.Value)
			}
		}
	}
	// проверяю позднюю инициализацию
	for name, d := range cl.Members {
		if f, ok := d.(*ast.Field); ok && f.Later {
			_, ok := vals[name]
			if !ok {
				env.AddError(c.Pos, "СЕМ-НЕТ-ПОЗЖЕ-ПОЛЯ", name)
			}
		}
	}
}
