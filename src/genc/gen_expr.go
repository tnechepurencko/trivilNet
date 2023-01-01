package genc

import (
	"fmt"
	"trivil/ast"
	"trivil/lexer"
	"unicode/utf8"
)

var _ = fmt.Printf

func (genc *genContext) genExpr(expr ast.Expr) string {

	switch x := expr.(type) {
	case *ast.IdentExpr:
		return genc.outName(x.Name)
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
		return genc.genStringLiteral(li)
	default:
		panic("ni")
	}
}

func (genc *genContext) genStringLiteral(li *ast.LiteralExpr) string {

	var name = genc.localName(nm_stringLiteral)
	genc.g("TString %s = NULL;", name)

	return fmt.Sprintf("%s(&%s, %d, %d, %s)",
		rt_newLiteralString,
		name, len(li.Lit), utf8.RuneCountInString(li.Lit), "\""+li.Lit+"\"")
}

func (genc *genContext) genCall(call *ast.CallExpr) string {

	if call.StdFunc != nil {
		return genc.genStdFuncCall(call)
	}

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

func (genc *genContext) genStdFuncCall(call *ast.CallExpr) string {

	switch call.StdFunc.Name {
	case "длина":
		return genc.genStdLen(call)

	default:
		panic("assert: не реализована стандартная функция " + call.StdFunc.Name)
	}
}

func (genc *genContext) genStdLen(call *ast.CallExpr) string {
	var a = call.Args[0]

	var t = ast.UnderType(a.GetType())
	if t == ast.String {

		return fmt.Sprintf("%s(%s)", rt_lenString, genc.genExpr(a))

	} else {
		panic("ni")
	}
}
