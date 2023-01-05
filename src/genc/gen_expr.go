package genc

import (
	"fmt"
	"strings"

	"trivil/ast"
	"trivil/lexer"
	"unicode/utf8"
)

var _ = fmt.Printf

func (genc *genContext) genExpr(expr ast.Expr) string {

	switch x := expr.(type) {
	case *ast.IdentExpr:
		return genc.genIdent(x)
	case *ast.LiteralExpr:
		return genc.genLiteral(x)
	case *ast.UnaryExpr:
		return fmt.Sprintf("%s(%s)", x.Op.String(), genc.genExpr(x.X))
	case *ast.BinaryExpr:
		return fmt.Sprintf("(%s %s %s)", genc.genExpr(x.X), x.Op.String(), genc.genExpr(x.Y))
	case *ast.SelectorExpr:
		if x.X == nil {
			panic("ni")
		} else {
			return fmt.Sprintf("%s->%s", genc.genExpr(x.X), genc.outName(x.Name))
		}
	case *ast.CallExpr:
		return genc.genCall(x)
	case *ast.GeneralBracketExpr:
		return genc.genBracketExpr(x)
	case *ast.ClassCompositeExpr:
		return genc.genClassComposite(x)

	default:
		panic(fmt.Sprintf("gen expression: ni %T", expr))
	}
}

func (genc *genContext) genIdent(x *ast.IdentExpr) string {
	return genc.outName(x.Name)
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

	} else if _, ok := t.(*ast.VectorType); ok {
		return fmt.Sprintf("%s(%s)", rt_lenVector, genc.genExpr(a))
	} else {
		panic("ni")
	}
}

func (genc *genContext) genBracketExpr(x *ast.GeneralBracketExpr) string {

	if x.Index != nil {
		var name string
		if id, ok := x.X.(*ast.IdentExpr); ok {
			name = genc.genIdent(id)
		} else {
			name = genc.localName("loc")

			genc.c("%s %s = %s;",
				genc.typeRef(x.X.GetType()),
				name,
				genc.genExpr(x.X))
		}
		return fmt.Sprintf("%s->body[%s(%s, %s)]",
			name,
			rt_vcheck,
			name,
			genc.genExpr(x.Index))
	}

	return genc.genArrayComposite(x.Composite)

}

func (genc *genContext) genArrayComposite(x *ast.ArrayCompositeExpr) string {
	var name = genc.localName("loc")

	var vt = ast.UnderType(x.Typ).(*ast.VectorType)
	var s = fmt.Sprintf("%s %s = %s(sizeof(%s), %d);",
		genc.typeRef(x.Typ),
		name,
		rt_newVector,
		genc.typeRef(vt.ElementTyp),
		len(x.Elements))

	var list = make([]string, len(x.Elements))
	for i, e := range x.Elements {
		var inx string
		if e.Key == nil {
			inx = fmt.Sprintf("%d", i)
		} else {
			inx = genc.genExpr(e.Key)
		}
		list[i] = fmt.Sprintf("%s->body[%s] = %s;", name, inx, genc.genExpr(e.Value))
	}
	s += strings.Join(list, " ")

	genc.c(s)

	return name
}

func (genc *genContext) genClassComposite(x *ast.ClassCompositeExpr) string {
	var name = genc.localName("loc")

	//	var ct = ast.UnderType(x.Typ).(*ast.ClassType)
	var tname = genc.typeRef(x.Typ)
	var s = fmt.Sprintf("%s %s = %s(&%s);",
		tname,
		name,
		rt_newObject,
		tname+nm_meta_var_suffix)

	var list = make([]string, len(x.Values))
	for i, v := range x.Values {
		list[i] = fmt.Sprintf("%s->%s = %s;",
			name, genc.outName(v.Name), genc.genExpr(v.Value))
	}
	s += strings.Join(list, " ")

	genc.c(s)
	return name
}
