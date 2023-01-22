package genc

import (
	"fmt"
	"strings"

	"trivil/ast"
)

var _ = fmt.Printf

//==== call

func (genc *genContext) genCall(call *ast.CallExpr) string {

	if call.StdFunc != nil {
		return genc.genStdFuncCall(call)
	}

	if isMethodCall(call.X) {
		return genc.genMethodCall(call)
	}

	var left = genc.genExpr(call.X)

	var cargs = genc.genArgs(call)

	return left + "(" + cargs + ")"
}

func (genc *genContext) genArgs(call *ast.CallExpr) string {

	var ft = call.X.GetType().(*ast.FuncType)

	var vPar = ast.VariadicParam(ft)

	if vPar == nil {
		var cargs = make([]string, len(call.Args))
		for i, a := range call.Args {
			cargs[i] = genc.genExpr(a)
		}
		return strings.Join(cargs, ", ")
	} else {
		var cargs = make([]string, len(ft.Params))
		var normCount = len(ft.Params) - 1

		for i := 0; i < normCount; i++ {
			cargs[i] = genc.genExpr(call.Args[i])
		}

		var vTyp = vPar.Typ.(*ast.VariadicType)

		if ast.IsTagPairType(vTyp.ElementTyp) {
			cargs[normCount] = genc.genVariadicTaggedArgs(call, vPar, normCount)
		} else {
			cargs[normCount] = genc.genVariadicArgs(call, vPar, vTyp, normCount)
		}

		return strings.Join(cargs, ", ")
	}
}

func (genc *genContext) genVariadicArgs(call *ast.CallExpr, vPar *ast.Param, vTyp *ast.VariadicType, normCount int) string {

	var loc = genc.localName("loc")
	var et = genc.typeRef(vTyp.ElementTyp)
	var vLen = len(call.Args) - normCount

	genc.c("struct { TInt64 len; %s body[%d]; } %s;", et, vLen, loc)

	//TODO: нужно ли выдержать какой-то порядок вычисления аргументов?
	var cargs = make([]string, vLen)
	var n = 0
	for i := normCount; i < len(call.Args); i++ {
		cargs[n] = fmt.Sprintf("%s.body[%d]=%s;", loc, n, genc.genExpr(call.Args[i]))
		n++
	}
	genc.c("%s.len=%d;%s", loc, vLen, strings.Join(cargs, ""))

	return "&" + loc
}

func (genc *genContext) genVariadicTaggedArgs(call *ast.CallExpr, vPar *ast.Param, normCount int) string {

	var loc = genc.localName("loc")
	var vLen = len(call.Args) - normCount

	genc.c("struct { TInt64 len; TInt64 body[%d]; } %s;", vLen*2, loc)

	var cargs = make([]string, vLen*2)
	var n = 0
	for i := normCount; i < len(call.Args); i++ {
		cargs[n] = fmt.Sprintf("%s.body[%d]=%s;", loc, n, genc.genTypeTag(call.Args[i].GetType()))
		cargs[n+1] = fmt.Sprintf("%s.body[%d]=(%s)%s;",
			loc, n+1, predefinedTypeName(ast.Int64.Name), genc.genExpr(call.Args[i]))
		n += 2
	}
	genc.c("%s.len=%d;%s", loc, vLen, strings.Join(cargs, ""))

	return "&" + loc
}

func isMethodCall(left ast.Expr) bool {
	sel, ok := left.(*ast.SelectorExpr)
	if !ok {
		return false
	}

	f, ok := sel.Obj.(*ast.Function)
	if !ok {
		return false
	}

	return f.Recv != nil
}

func (genc *genContext) genMethodCall(call *ast.CallExpr) string {

	sel := call.X.(*ast.SelectorExpr)
	f := sel.Obj.(*ast.Function)

	var name string
	if id, ok := sel.X.(*ast.IdentExpr); ok {
		name = genc.genIdent(id)
	} else {
		name = genc.localName("loc")

		genc.c("%s %s = %s;",
			genc.typeRef(sel.X.GetType()),
			name,
			genc.genExpr(sel.X))
	}

	//TODO - можно убрать каст, если лишний
	var args = fmt.Sprintf("(%s)%s", genc.typeRef(f.Recv.Typ), name)

	if len(call.Args) > 0 {
		args += ", " + genc.genArgs(call)
	}

	return fmt.Sprintf("%s->%s->%s(%s)", name, nm_VT_field, genc.outName(f.Name), args)
}

func (genc *genContext) genStdFuncCall(call *ast.CallExpr) string {

	switch call.StdFunc.Name {
	case ast.StdLen:
		return genc.genStdLen(call)
	case ast.StdTag:
		return genc.genStdTag(call)
	case ast.StdSomething:
		return genc.genStdSomething(call)

	case ast.VectorAppend:
		return genc.genVectorAppend(call)

	default:
		panic("assert: не реализована стандартная функция " + call.StdFunc.Name)
	}
}

func (genc *genContext) genStdLen(call *ast.CallExpr) string {
	var a = call.Args[0]

	var t = ast.UnderType(a.GetType())
	if t == ast.String {
		return fmt.Sprintf("%s(%s)", rt_lenString, genc.genExpr(a))
	}

	switch t.(type) {
	case *ast.VectorType:
		return fmt.Sprintf("%s(%s)", rt_lenVector, genc.genExpr(a))
	case *ast.VariadicType:
		return fmt.Sprintf("(*(TInt64 *)%s)", genc.genExpr(a))
	default:
		panic("ni")
	}
}

func (genc *genContext) genStdTag(call *ast.CallExpr) string {

	var a = call.Args[0]

	if tExpr, ok := a.(*ast.TypeExpr); ok {
		return genc.genTypeTag(tExpr.Typ)
	} else {
		var t = a.GetType()
		if ast.IsTagPairType(t) {
			return genc.genExprTag(a)
		} else {
			return genc.genTypeTag(t)
		}
	}
}

func (genc *genContext) genTypeTag(t ast.Type) string {
	t = ast.UnderType(t)
	switch x := t.(type) {
	case *ast.PredefinedType:
		return fmt.Sprintf("%s%s()", rt_tag, predefinedTypeName(x.Name))
	default:
		panic(("ni"))
	}
}

func (genc *genContext) genExprTag(e ast.Expr) string {
	switch x := e.(type) {
	case *ast.GeneralBracketExpr:
		if x.Index == nil {
			panic("assert")
		}

		var left = genc.genExpr(x.X)

		return fmt.Sprintf("((%s*)(%s + sizeof(TInt64)))[%s(%s, *(TInt64 *)%s) << 1]",
			predefinedTypeName(ast.Word64.Name),
			left,
			rt_indexcheck,
			genc.genExpr(x.Index),
			left)

	case *ast.IdentExpr:
		panic("ni")
	}
	panic("assert")
}

func (genc *genContext) genStdSomething(call *ast.CallExpr) string {

	var a = call.Args[0]

	var t = a.GetType()
	if !ast.IsTagPairType(t) {
		panic("assert")
	}

	switch x := a.(type) {
	case *ast.GeneralBracketExpr:
		if x.Index == nil {
			panic("assert")
		}

		var left = genc.genExpr(x.X)

		return fmt.Sprintf("((%s*)(%s + sizeof(TInt64)))[(%s(%s, *(TInt64 *)%s) << 1)+1]",
			predefinedTypeName(ast.Word64.Name),
			left,
			rt_indexcheck,
			genc.genExpr(x.Index),
			left)

	case *ast.IdentExpr:
		panic("ni")
	}
	panic("assert")
}

//==== векторные

func (genc *genContext) genVectorAppend(call *ast.CallExpr) string {

	var vt = ast.UnderType(call.X.GetType()).(*ast.VectorType)

	var loc = genc.localName("loc")
	var et = genc.typeRef(vt.ElementTyp)

	var cargs = make([]string, len(call.Args))
	for i, a := range call.Args {
		cargs[i] = genc.genExpr(a)
	}

	genc.c("%s %s[%d] = {%s};", et, loc, len(call.Args), strings.Join(cargs, ", "))

	return fmt.Sprintf("%s(%s, sizeof(%s), %d, %s)", rt_vector_append, genc.genExpr(call.X), et, len(call.Args), loc)
}
