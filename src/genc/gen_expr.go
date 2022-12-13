package genc

import (
	"fmt"
	"trivil/ast"
	"trivil/lexer"
)

var _ = fmt.Printf

func (genc *genContext) genExpr(expr ast.Expr) string {

	switch x := expr.(type) {
	case *ast.IdentExpr:
		return x.Name
	case *ast.LiteralExpr:
		return genc.genLiteral(x)
	case *ast.UnaryExpr:
		return x.Op.String() + genc.genExpr(x.X)
	case *ast.BinaryExpr:
		return fmt.Sprintf("(%s %s %s)", genc.genExpr(x.X), x.Op.String(), genc.genExpr(x.Y))
	case *ast.CallExpr:
		return genc.genCall(x)

	default:
		panic(fmt.Sprintf("gen expression: ni %T", expr))
	}
}

func (genc *genContext) genLiteral(li *ast.LiteralExpr) string {
	switch li.Kind {
	case lexer.INT, lexer.FLOAT:
		return li.Lit
	case lexer.STRING:
		return "\"" + li.Lit + "\""
	default:
		panic("ni")
	}
}

func (genc *genContext) genCall(call *ast.CallExpr) string {

	var left = genc.genExpr(call.X)

	var cargs = ""
	for i, a := range call.Args {
		var ca = genc.genExpr(a)

		cargs += ca

		if i < len(call.Args)-1 {
			cargs += ", "
		}
	}

	return left + "(" + cargs + ")"

}
