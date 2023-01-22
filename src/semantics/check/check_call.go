package check

import (
	"fmt"
	"trivil/ast"
	"trivil/env"
)

var _ = fmt.Printf

func (cc *checkContext) call(x *ast.CallExpr) {

	cc.expr(x.X)

	sel, ok := x.X.(*ast.SelectorExpr)
	if ok && sel.StdMethod != nil {
		x.StdFunc = sel.StdMethod
		x.X = sel.X // убрал лишний селектор
		cc.callStdFunction(x)
		return
	}

	ft, ok := x.X.GetType().(*ast.FuncType)
	if !ok {

		env.AddError(x.X.GetPos(), "СЕМ-ВЫЗОВ-НЕ_ФУНКТИП", ast.TypeName(x.X.GetType()))
		return
	}

	if ft.ReturnTyp == nil {
		x.Typ = ast.Void
	} else {
		x.Typ = ft.ReturnTyp
	}

	var vPar = ast.VariadicParam(ft)

	if vPar == nil {
		if len(x.Args) != len(ft.Params) {
			env.AddError(x.X.GetPos(), "СЕМ-ЧИСЛО-АРГУМЕНТОВ", len(x.Args), len(ft.Params))
			return
		}

		for i, p := range ft.Params {
			cc.expr(x.Args[i])
			cc.checkAssignable(p.Typ, x.Args[i])
		}
	} else {
		var normCount = len(ft.Params) - 1
		if len(x.Args) < normCount {
			env.AddError(x.X.GetPos(), "СЕМ-ВАРИАДИК-ЧИСЛО-АРГУМЕНТОВ", normCount, len(x.Args))
			return
		}

		for i := 0; i < normCount; i++ {
			cc.expr(x.Args[i])
			cc.checkAssignable(ft.Params[i].Typ, x.Args[i])
		}

		var vTyp = vPar.Typ.(*ast.VariadicType)
		if cc.checkUnfold(x.Args, normCount, vTyp.ElementTyp) {
			// проверено
		} else {
			for i := normCount; i < len(x.Args); i++ {
				cc.expr(x.Args[i])
				cc.checkAssignable(vTyp.ElementTyp, x.Args[i])
			}
		}

	}
}

func (cc *checkContext) checkUnfold(args []ast.Expr, start int, elementTyp ast.Type) bool {
	for i := start; i < len(args); i++ {
		if u, ok := args[i].(*ast.UnfoldExpr); ok {

			cc.expr(u.X)

			if i != start || len(args)-start > 1 {
				env.AddError(args[i].GetPos(), "СЕМ-ОДНО-РАЗВОРАЧИВАНИЕ")
			}
			var t = u.X.GetType()
			switch xt := ast.UnderType(t).(type) {
			case *ast.VectorType:
				if !equalTypes(elementTyp, xt.ElementTyp) {
					env.AddError(args[i].GetPos(), "СЕМ-ТИПЫ-ЭЛЕМЕНТОВ-НЕ-СОВПАДАЮТ",
						ast.TypeName(elementTyp), ast.TypeName(xt.ElementTyp))
				}
			default:
				env.AddError(args[i].GetPos(), "СЕМ-ОЖИДАЛСЯ-ТИП-ВЕКТОРА", ast.TypeName(t))
			}

			return true
		}
	}
	return false
}

//=== стд. функции

func (cc *checkContext) callStdFunction(x *ast.CallExpr) {

	switch x.StdFunc.Name {
	case "":
		return
	case ast.StdLen:
		cc.callStdLen(x)
	case ast.StdTag:
		cc.callStdTag(x)
	case ast.StdSomething:
		cc.callStdSomething(x)

	case ast.VectorAppend:
		cc.callVectorAppend(x)

	default:
		panic("assert: не реализована стандартная функция " + x.StdFunc.Name)
	}
}

func (cc *checkContext) callStdLen(x *ast.CallExpr) {
	x.Typ = ast.Int64

	if len(x.Args) != 1 {
		env.AddError(x.Pos, "СЕМ-СТДФУНК-ОШ-ЧИСЛО-АРГ", x.StdFunc.Name, "1")
		return
	}

	cc.expr(x.Args[0])

	var t = ast.UnderType(x.Args[0].GetType())

	if ast.IsIndexableType(t) || t == ast.String {
		// ok
	} else {
		env.AddError(x.Pos, "СЕМ-СТД-ДЛИНА-ОШ-ТИП-АРГ", x.StdFunc.Name)
	}
}

func (cc *checkContext) callStdTag(x *ast.CallExpr) {
	x.Typ = ast.Word64

	if len(x.Args) != 1 {
		env.AddError(x.Pos, "СЕМ-СТДФУНК-ОШ-ЧИСЛО-АРГ", x.StdFunc.Name, "1")
		return
	}

	var t = cc.typeExpr(x.Args[0])
	if t != nil {
		var prev = x.Args[0]
		x.Args[0] = &ast.TypeExpr{
			ExprBase: ast.ExprBase{Pos: prev.GetPos(), Typ: t, ReadOnly: true},
		}
	} else {
		cc.expr(x.Args[0])

		if !ast.HasTag(x.Args[0].GetType()) {
			env.AddError(x.Pos, "СЕМ-СТД-ТЕГ-ОШ-АРГ")
			return
		}
	}
}

func (cc *checkContext) callStdSomething(x *ast.CallExpr) {
	x.Typ = ast.Word64

	if len(x.Args) != 1 {
		env.AddError(x.Pos, "СЕМ-СТДФУНК-ОШ-ЧИСЛО-АРГ", x.StdFunc.Name, "1")
		return
	}

	cc.expr(x.Args[0])

	if !ast.IsTagPairType(x.Args[0].GetType()) {
		env.AddError(x.Pos, "СЕМ-СТД-НЕЧТО-ОШ-АРГ")
		return
	}
}

//==== векторные

func (cc *checkContext) callVectorAppend(x *ast.CallExpr) {

	var vt = ast.UnderType(x.X.GetType()).(*ast.VectorType)

	if cc.checkUnfold(x.Args, 0, vt.ElementTyp) {
		// проверено
	} else {
		for _, a := range x.Args {
			cc.expr(a)
			cc.checkAssignable(vt.ElementTyp, a)
		}
	}
	x.Typ = ast.Void
}