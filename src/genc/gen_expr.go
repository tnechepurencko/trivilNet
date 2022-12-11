package genc

import (
	"fmt"
	"trivil/ast"
)

var _ = fmt.Printf

func (genc *genContext) genExpr(expr ast.Expr) string {

	switch x := expr.(type) {
	case *ast.IdentExpr:
		return x.Name

	case *ast.CallExpr:
		var left = genc.genExpr(x.X)
		return left + "()"

	default:
		panic(fmt.Sprintf("gen expression: ni %T", expr))
	}

}
