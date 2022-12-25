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
		//TODO: check not type
		x.Typ = x.Obj.GetType()

		//fmt.Printf("ident %v %v\n", x.Obj, x.Typ)

	case *ast.UnaryExpr:
		cc.expr(x.X)
		cc.unaryExpr(x)

	case *ast.BinaryExpr:
		cc.expr(x.X)
		cc.expr(x.Y)
		cc.binaryExpr(x)

	case *ast.SelectorExpr:
		cc.expr(x.X)
		panic("ni")

	case *ast.CallExpr:
		cc.expr(x.X)
		for _, a := range x.Args {
			cc.expr(a)
		}
		cc.call(x)

	case *ast.CompositeExpr:
		cc.expr(x.X)

		for _, vp := range x.Values {
			cc.expr(vp.V)
		}
		panic("ni")
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
	case *ast.BoolLiteral:
		x.Typ = ast.Bool
	default:
		panic(fmt.Sprintf("expression: ni %T", expr))
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
		if !cc.assignable(p.Typ, x.Args[i]) {

			env.AddError(x.Args[i].GetPos(), "СЕМ-НЕСОВМЕСТИМО-ПРИСВ", cc.errorHint,
				ast.TypeString(p.Typ), ast.TypeString(x.Args[i].GetType()))
		}
	}
}

func (cc *checkContext) unaryExpr(x *ast.UnaryExpr) {
	switch x.Op {
	case lexer.SUB:
		panic("ni")
	case lexer.NOT:
		if !ast.IsBoolType(x.X.GetType()) {
			env.AddError(x.X.GetPos(), "СЕМ-ОШ-УНАРНАЯ-ТИП",
				x.Op.String(), ast.TypeString(x.X.GetType()))
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
