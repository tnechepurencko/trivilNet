package genc

import (
	"fmt"
	"strings"

	"trivil/ast"
	"trivil/env"
)

var _ = fmt.Printf

func (genc *genContext) genStatementSeq(seq *ast.StatementSeq) {

	for _, s := range seq.Statements {
		genc.genStatement(s)
	}
}

func (genc *genContext) genStatement(s ast.Statement) {
	switch x := s.(type) {
	case *ast.DeclStatement:
		s := genc.genLocalDecl(x.D)
		genc.c("%s", s)
	case *ast.ExprStatement:
		s := genc.genExpr(x.X)
		genc.c("%s;", s)
	case *ast.AssignStatement:
		l := genc.genExpr(x.L)
		r := genc.genExpr(x.R)

		var cast = genc.assignCast(x.L.GetType(), x.R.GetType())
		genc.c("%s = %s%s;", l, cast, r)
	case *ast.IncStatement:
		l := genc.genExpr(x.L)
		genc.c("%s++;", l)
	case *ast.DecStatement:
		l := genc.genExpr(x.L)
		genc.c("%s--;", l)
	case *ast.If:
		genc.genIf(x, "")
	case *ast.While:
		genc.genWhile(x)
	case *ast.Guard:
		genc.genGuard(x)
	case *ast.When:
		if canWhenAsSwitch(x) {
			genc.genWhenAsSwitch(x)
		} else {
			genc.genWhenAsIfs(x)
		}
	case *ast.Return:
		r := ""
		if x.X != nil {
			r = " " + genc.genExpr(x.X)
		}
		genc.c("return %s;", r)

	case *ast.Break:
		genc.c("break;")
	case *ast.Crash:
		genc.genCrash(x)

	default:
		panic(fmt.Sprintf("gen statement: ni %T", s))
	}
}

func (genc *genContext) assignCast(lt, rt ast.Type) string {
	if ast.UnderType(lt) != ast.UnderType(rt) {
		return "(" + genc.typeRef(lt) + ")"
	}
	return ""
}

func (genc *genContext) genIf(x *ast.If, prefix string) {
	genc.c("%sif (%s) {", prefix, removeExtraPars(genc.genExpr(x.Cond)))
	genc.genStatementSeq(x.Then)
	genc.c("}")
	if x.Else != nil {

		elsif, ok := x.Else.(*ast.If)
		if ok {
			genc.genIf(elsif, "else ")
		} else {
			genc.c("else {")
			genc.genStatementSeq(x.Else.(*ast.StatementSeq))
			genc.c("}")
		}
	}
}

func removeExtraPars(s string) string {
	if len(s) == 0 {
		return s
	}
	if s[0] == '(' && s[len(s)-1] == ')' {
		return s[1 : len(s)-1]
	}
	return s
}

func (genc *genContext) genWhile(x *ast.While) {
	genc.c("while (%s) {", removeExtraPars(genc.genExpr(x.Cond)))
	genc.genStatementSeq(x.Seq)
	genc.c("}")
}

func (genc *genContext) genGuard(x *ast.Guard) {
	genc.c("if (!(%s)) {", removeExtraPars(genc.genExpr(x.Cond)))
	seq, ok := x.Else.(*ast.StatementSeq)
	if ok {
		genc.genStatementSeq(seq)
	} else {
		genc.genStatement(x.Else)
	}
	genc.c("}")
}

func (genc *genContext) genCrash(x *ast.Crash) {

	var expr string
	var li = literal(x.X)
	if li != nil {
		expr = "\"" + encodeLiteralString(li.StrVal) + "\""
	} else {
		expr = genc.genExpr(x.X) + "->body"
	}

	genc.c("%s(%s,%s);", rt_crash, expr, genPos(x.Pos))
}

func genPos(pos int) string {
	src, line, col := env.SourcePos(pos)
	return fmt.Sprintf("\"%s:%d:%d\"", src.Path, line, col)
}

func literal(expr ast.Expr) *ast.LiteralExpr {

	switch x := expr.(type) {
	case *ast.LiteralExpr:
		return x
	case *ast.ConversionExpr:
		if x.Done {
			return literal(x.X)
		}
	}
	return nil
}

//==== когда

func canWhenAsSwitch(x *ast.When) bool {
	var t = ast.UnderType(x.X.GetType())
	switch t {
	case ast.Byte, ast.Int64, ast.Word64, ast.Symbol:
	default:
		return false
	}

	for _, c := range x.Cases {
		for _, e := range c.Exprs {
			if _, ok := e.(*ast.LiteralExpr); !ok {
				return false
			}
		}
	}
	return true
}

func (genc *genContext) genWhenAsSwitch(x *ast.When) {
	genc.c("switch (%s) {", genc.genExpr(x.X))

	for _, c := range x.Cases {
		for _, e := range c.Exprs {
			genc.c("case %s: ", genc.genExpr(e))
		}
		genc.genStatementSeq(c.Seq)
		genc.c("break;")
	}

	if x.Else != nil {
		genc.c("default:")
		genc.genStatementSeq(x.Else)
	}

	genc.c("}")
}

func (genc *genContext) genWhenAsIfs(x *ast.When) {

	var strCompare = ast.IsStringType(x.X.GetType())

	var loc = genc.localName("")
	genc.c("%s %s = %s;", genc.typeRef(x.X.GetType()), loc, genc.genExpr(x.X))

	var els = ""
	for _, c := range x.Cases {

		var conds = make([]string, 0)
		for _, e := range c.Exprs {
			if strCompare {
				conds = append(conds, fmt.Sprintf("%s(%s, %s)", rt_equalStrings, loc, genc.genExpr(e)))
			} else {
				conds = append(conds, fmt.Sprintf("%s == %s", loc, genc.genExpr(e)))
			}
		}
		genc.c("%sif (%s) {", els, strings.Join(conds, " || "))
		els = "else "
		genc.genStatementSeq(c.Seq)
		genc.c("}")
	}

	if x.Else != nil {
		genc.c("else {")
		genc.genStatementSeq(x.Else)
		genc.c("}")
	}
}
