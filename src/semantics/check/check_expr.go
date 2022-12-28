package check

import (
	"fmt"
	"trivil/ast"
	"trivil/env"
	"trivil/lexer"
)

var _ = fmt.Printf

func (cc *checkContext) expr(expr ast.Expr) {
	switch x := expr.(type) {
	case *ast.IdentExpr:
		if _, ok := x.Obj.(*ast.TypeRef); ok {
			env.AddError(x.Pos, "СЕМ-ТИП-В-ВЫРАЖЕНИИ")
			x.Typ = &ast.InvalidType{TypeBase: ast.TypeBase{Pos: x.Pos}}
			return
		}

		x.Typ = x.Obj.(ast.Decl).GetType()

		_, isVar := x.Obj.(*ast.VarDecl)
		x.ReadOnly = !isVar
		//fmt.Printf("ident %v %v\n", x.Obj, x.Typ)

	case *ast.UnaryExpr:
		cc.expr(x.X)
		cc.unaryExpr(x)

	case *ast.BinaryExpr:
		cc.expr(x.X)
		cc.expr(x.Y)
		cc.binaryExpr(x)

	case *ast.SelectorExpr:
		cc.selector(x)

	case *ast.CallExpr:
		cc.expr(x.X)
		for _, a := range x.Args {
			cc.expr(a)
		}
		cc.call(x)

	case *ast.GeneralBracketExpr:
		cc.generalBracketExpr(x)

	case *ast.ClassCompositeExpr:
		cc.classComposite(x)
	case *ast.LiteralExpr:
		switch x.Kind {
		case lexer.INT:
			x.Typ = ast.Int64
		case lexer.FLOAT:
			x.Typ = ast.Float64
		case lexer.STRING:
			x.Typ = ast.String
		default:
			panic(fmt.Sprintf("LiteralExpr - bad kind: ni %v", x))
		}
		x.ReadOnly = true
	case *ast.BoolLiteral:
		x.Typ = ast.Bool
		x.ReadOnly = true
	default:
		panic(fmt.Sprintf("expression: ni %T", expr))
	}

}

func (cc *checkContext) selector(x *ast.SelectorExpr) {
	if x.Obj != nil {
		// imported object
		if _, ok := x.Obj.(*ast.TypeRef); ok {
			env.AddError(x.Pos, "СЕМ-ТИП-В-ВЫРАЖЕНИИ")
			x.Typ = &ast.InvalidType{TypeBase: ast.TypeBase{Pos: x.Pos}}
		}
		return
	}
	cc.expr(x.X)
	var t = x.X.GetType()

	var cl = getClassType(t)
	if cl == nil {
		env.AddError(x.X.GetPos(), "СЕМ-ОЖИДАЛСЯ-ТИП-КЛАССА", ast.TypeString(t))
		x.Typ = &ast.InvalidType{TypeBase: ast.TypeBase{Pos: x.X.GetPos()}}
		return
	}

	d, ok := cl.Members[x.Name]
	if !ok {
		//TODO: проверить экспорт
		env.AddError(x.Pos, "СЕМ-ОЖИДАЛОСЬ-ПОЛЕ-ИЛИ-МЕТОД", x.Name)
	} else {
		x.Typ = d.GetType()
		x.Obj = d
	}

	if x.Typ == nil {
		x.Typ = &ast.InvalidType{TypeBase: ast.TypeBase{Pos: x.X.GetPos()}}
	}
}

func (cc *checkContext) call(x *ast.CallExpr) {

	ft, ok := x.X.GetType().(*ast.FuncType)
	if !ok {
		env.AddError(x.X.GetPos(), "СЕМ-ВЫЗОВ-НЕ_ФУНКТИП")
		return
	}

	if ft.ReturnTyp == nil {
		x.Typ = ast.Void
	} else {
		x.Typ = ft.ReturnTyp
	}

	if len(x.Args) != len(ft.Params) {
		env.AddError(x.X.GetPos(), "СЕМ-ЧИСЛО-АРГУМЕНТОВ", len(x.Args), len(ft.Params))
		return
	}

	for i, p := range ft.Params {
		cc.checkAssignable(p.Typ, x.Args[i])
	}
}

func (cc *checkContext) generalBracketExpr(x *ast.GeneralBracketExpr) {

	var t = cc.typeName(x.X)

	if t != nil || len(x.Composite.Elements) != 1 || x.Composite.Keys { // composite
		cc.arrayComposite(x.Composite, t)

		if t == nil {
			t = &ast.InvalidType{TypeBase: ast.TypeBase{Pos: x.X.GetPos()}}
		}
		x.Typ = t
		x.X = nil
		return
	}

	// это индексация
	cc.expr(x.X)

	t = x.X.GetType()

	if !ast.IsIndexableType(t) {
		env.AddError(x.X.GetPos(), "СЕМ-ОЖИДАЛСЯ-ТИП-МАССИВА", ast.TypeString(t))
		x.Typ = &ast.InvalidType{TypeBase: ast.TypeBase{Pos: x.Pos}}
	} else {
		x.Index = x.Composite.Elements[0].Value
		cc.expr(x.Index)
		if !ast.IsIntegerType(x.Index.GetType()) {
			env.AddError(x.Index.GetPos(), "СЕМ-ОШ-ТИП-ИНДЕКСА", ast.TypeString(x.Index.GetType()))
		}
		x.Typ = ast.ElementType(t)
	}
	x.Composite = nil

	if x.X.IsReadOnly() {
		x.ReadOnly = true
	}
}

func (cc *checkContext) typeName(expr ast.Expr) ast.Type {

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

	for _, p := range c.Elements {

		if p.Key != nil {
			cc.expr(p.Key)
			if !ast.IsIntegerType(p.Key.GetType()) {
				env.AddError(c.Pos, "СЕМ-МАССИВ-КОМПОЗИТ-ТИП-КЛЮЧА")
			}
			cc.checkConstExpr(p.Key)
		}

		cc.expr(p.Value)
		if elemT != nil {
			cc.checkAssignable(elemT, p.Value)
		}
	}
}

func getClassType(t ast.Type) *ast.ClassType {
	if tr, ok := t.(*ast.TypeRef); ok {
		t = tr.Typ
	}

	cl, _ := t.(*ast.ClassType)
	return cl
}

func (cc *checkContext) classComposite(c *ast.ClassCompositeExpr) {

	var t = cc.typeName(c.X)

	if t == nil {
		env.AddError(c.Pos, "СЕМ-КОМПОЗИТ-НЕТ-ТИПА")
		c.Typ = &ast.InvalidType{TypeBase: ast.TypeBase{Pos: c.X.GetPos()}}
		return
	}

	var cl = getClassType(t)
	if cl == nil {
		env.AddError(c.Pos, "СЕМ-КЛАСС-КОМПОЗИТ-ОШ-ТИП")
		c.Typ = &ast.InvalidType{TypeBase: ast.TypeBase{Pos: c.X.GetPos()}}
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
	for _, vp := range c.Values {
		d, ok := cl.Members[vp.Name]
		if !ok {
			//TODO: проверить экспорт
			env.AddError(vp.Pos, "СЕМ-КЛАСС-КОМПОЗИТ-НЕТ-ПОЛЯ", vp.Name)
		} else {
			f, ok := d.(*ast.Field)
			if !ok {
				env.AddError(vp.Pos, "СЕМ-КЛАСС-КОМПОЗИТ-НЕ-ПОЛE")
			} else {
				cc.checkAssignable(f.Typ, vp.Value)
			}
		}
	}
	//TODO: проверить обязательные значения (без умолчания)
}

func (cc *checkContext) unaryExpr(x *ast.UnaryExpr) {
	switch x.Op {
	case lexer.SUB:
		panic("ni")
	case lexer.NOT:
		if !ast.IsBoolType(x.X.GetType()) {
			env.AddError(x.X.GetPos(), "СЕМ-ОШ-УНАРНАЯ-ТИП",
				ast.TypeString(x.X.GetType()), x.Op.String())
		}
		x.Typ = ast.Bool
	default:
		panic(fmt.Sprintf("unary expr ni: %T op=%s", x, x.Op.String()))
	}
}

func (cc *checkContext) binaryExpr(x *ast.BinaryExpr) {

	switch x.Op {
	case lexer.ADD, lexer.SUB, lexer.MUL, lexer.REM, lexer.QUO:
		if ast.IsIntegerType(x.X.GetType()) || ast.IsFloatType(x.X.GetType()) {
			checkOperandTypes(x)
		} else {
			env.AddError(x.X.GetPos(), "СЕМ-ОШ-ТИП-ОПЕРАНДА",
				ast.TypeString(x.X.GetType()), x.Op.String())
		}
		x.Typ = x.X.GetType()
	case lexer.AND, lexer.OR:
		if !ast.IsBoolType(x.X.GetType()) {
			env.AddError(x.X.GetPos(), "СЕМ-ОШ-ТИП-ОПЕРАНДА",
				ast.TypeString(x.X.GetType()), x.Op.String())
		} else if !ast.IsBoolType(x.Y.GetType()) {
			env.AddError(x.Y.GetPos(), "СЕМ-ОШ-ТИП-ОПЕРАНДА",
				ast.TypeString(x.Y.GetType()), x.Op.String())
		}
		x.Typ = ast.Bool

	//case lexer.BITAND, lexer.BITOR:
	case lexer.EQ, lexer.NEQ:
		if ast.IsIntegerType(x.X.GetType()) || ast.IsFloatType(x.X.GetType()) {
			checkOperandTypes(x)

			//TODO: add other
		} else {
			env.AddError(x.Pos, "СЕМ-ОШ-ТИП-ОПЕРАНДА",
				ast.TypeString(x.X.GetType()), x.Op.String())
		}

		x.Typ = ast.Bool
	case lexer.LSS, lexer.LEQ, lexer.GTR, lexer.GEQ:
		if ast.IsIntegerType(x.X.GetType()) || ast.IsFloatType(x.X.GetType()) {
			checkOperandTypes(x)
		} else {
			env.AddError(x.Pos, "СЕМ-ОШ-ТИП-ОПЕРАНДА",
				ast.TypeString(x.X.GetType()), x.Op.String())
		}
		x.Typ = ast.Bool

	default:
		panic(fmt.Sprintf("binary expr ni: %T op=%s", x, x.Op.String()))
	}
}

func checkOperandTypes(x *ast.BinaryExpr) {
	if equalTypes(x.X.GetType(), x.Y.GetType()) {
		return
	}
	env.AddError(x.Pos, "СЕМ-ОПЕРАНДЫ-НЕ-СОВМЕСТИМЫ",
		ast.TypeString(x.X.GetType()), x.Op.String(), ast.TypeString(x.Y.GetType()))

}

func (cc *checkContext) checkConstExpr(expr ast.Expr) {
	switch x := expr.(type) {
	case *ast.LiteralExpr:
		return
	case *ast.IdentExpr:
		if x.Obj != nil {
			if _, ok := x.Obj.(*ast.ConstDecl); ok {
				return
			}
		}
	}

	env.AddError(expr.GetPos(), "СЕМ-ОШ-КОНСТ-ВЫРАЖЕНИЕ")
}

func isLValue(expr ast.Expr) bool {

	if expr.IsReadOnly() {
		return false
	}

	switch x := expr.(type) {
	case *ast.IdentExpr:
		return true
	case *ast.GeneralBracketExpr:
		return x.Index != nil
	case *ast.SelectorExpr:
		return true
	case *ast.ConversionExpr:
		return isLValue(x.X)
	default:
		return false
	}
}

func (cc *checkContext) checkLValue(expr ast.Expr) {
	if isLValue(expr) {
		return
	}
	env.AddError(expr.GetPos(), "СЕМ-НЕ-ПРИСВОИТЬ")
}
