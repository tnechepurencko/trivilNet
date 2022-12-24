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
		/*

			case *ast.UnaryExpr:
				cc.expr(x.X)

			case *ast.BinaryExpr:
				cc.expr(x.X)
				cc.expr(x.Y)

			case *ast.SelectorExpr:
				cc.expr(x.X)
				panic("ni")
		*/
	case *ast.CallExpr:
		cc.expr(x.X)
		for _, a := range x.Args {
			cc.expr(a)
		}
		cc.call(x)

		/*
			case *ast.IndexExpr:
				cc.expr(x.X)
				if x.Index != nil {
					cc.expr(x.Index)
				}

				for _, e := range x.Elements {
					cc.expr(e.L)
					if e.R != nil {
						cc.expr(e.R)
					}
				}
			case *ast.CompositeExpr:
				cc.expr(x.X)

				for _, vp := range x.Values {
					cc.expr(vp.V)

				}
		*/
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
		//cc.expr(x.X)
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
		res := assignable(p.Typ, x.Args[i])
		if res != "" {

			if res == "без уточнения" {
				res = ""
			}

			env.AddError(x.Args[i].GetPos(), "СЕМ-НЕСОВМЕСТИМЫй-АРГУМЕНТ", res,
				ast.TypeString(p.Typ), ast.TypeString(x.Args[i].GetType()))
		}
	}
}
