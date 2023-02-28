package check

import (
	"fmt"
	"trivil/ast"
	"trivil/env"
)

var _ = fmt.Printf

func (cc *checkContext) isConstExpr(expr ast.Expr) bool {
	switch x := expr.(type) {
	case *ast.LiteralExpr:
		return true
	case *ast.IdentExpr:
		if x.Obj != nil {
			if _, ok := x.Obj.(*ast.ConstDecl); ok {
				return true
			}
		}
	}
	return false
}

func (cc *checkContext) checkConstExpr(expr ast.Expr) {
	if cc.isConstExpr(expr) {
		return
	}

	env.AddError(expr.GetPos(), "СЕМ-ОШ-КОНСТ-ВЫРАЖЕНИЕ")
}

func (cc *checkContext) calculateIntConstExpr(expr ast.Expr) int64 {
	switch x := expr.(type) {
	case *ast.LiteralExpr:
		if x.Kind == ast.Lit_Int {
			return x.IntVal
		} else if x.Kind == ast.Lit_Word {
			return int64(x.WordVal)
		}
		panic("wrong literal kind")
	case *ast.IdentExpr:
		c := x.Obj.(*ast.ConstDecl)

		return cc.calculateIntConstExpr(c.Value)
	}
	panic("wrong const expression")
}
