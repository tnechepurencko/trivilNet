package genc

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"trivil/ast"
	"trivil/lexer"
)

var _ = fmt.Printf

func (genc *genContext) genExpr(expr ast.Expr) string {

	switch x := expr.(type) {
	case *ast.IdentExpr:
		return genc.genIdent(x)
	case *ast.LiteralExpr:
		return genc.genLiteral(x)
	case *ast.UnaryExpr:
		return fmt.Sprintf("%s(%s)", unaryOp(x.Op), genc.genExpr(x.X))
	case *ast.BinaryExpr:
		return genc.genBinaryExpr(x)
	case *ast.SelectorExpr:
		return genc.genSelector(x)
	case *ast.CallExpr:
		return genc.genCall(x)
	case *ast.ConversionExpr:
		if x.Caution {
			return genc.genCautionCast(x)
		} else {
			return genc.genConversion(x)
		}
	case *ast.GeneralBracketExpr:
		return genc.genBracketExpr(x)
	case *ast.ClassCompositeExpr:
		return genc.genClassComposite(x)

	default:
		panic(fmt.Sprintf("gen expression: ni %T", expr))
	}
}

func (genc *genContext) genIdent(id *ast.IdentExpr) string {

	c, ok := id.Obj.(*ast.ConstDecl)
	if ok {
		blit, ok := c.Value.(*ast.BoolLiteral)
		if ok {
			if blit.Value {
				return "true"
			} else {
				return "false"
			}
		}
	}

	return genc.declName(id.Obj.(ast.Decl))
}

//==== literals

func (genc *genContext) genLiteral(li *ast.LiteralExpr) string {
	switch li.Kind {
	case ast.Lit_Byte:
		return li.Lit
	case ast.Lit_Int:
		return li.Lit
	case ast.Lit_Float:
		return li.Lit
	case ast.Lit_Symbol:
		r, _ := utf8.DecodeRuneInString(li.Lit)
		return fmt.Sprintf("0x%x", r)
	case ast.Lit_String:
		return genc.genStringLiteral(li)
	default:
		panic("ni")
	}
}

func (genc *genContext) genStringLiteral(li *ast.LiteralExpr) string {

	if len(li.Lit) == 0 {
		return fmt.Sprintf("%s()", rt_emptyString)
	}

	var name = genc.localName(nm_stringLiteral)
	genc.g("static TString %s = NULL;", name)

	return fmt.Sprintf("%s(&%s, %d, %d, \"%s\")", rt_newLiteralString, name, len(li.Lit), utf8.RuneCountInString(li.Lit), li.Lit)
}

//==== унарные операции

func unaryOp(op lexer.Token) string {
	switch op {
	case lexer.SUB:
		return "-"
	case lexer.NOT:
		return "!"

	default:
		panic("ni unary" + op.String())
	}
}

//==== бинарные операции

func binaryOp(op lexer.Token) string {
	switch op {
	case lexer.OR:
		return "||"
	case lexer.AND:
		return "&&"
	case lexer.EQ:
		return "=="
	case lexer.NEQ:
		return "!="
	case lexer.LSS:
		return "<"
	case lexer.LEQ:
		return "<="
	case lexer.GTR:
		return ">"
	case lexer.GEQ:
		return ">="
	case lexer.ADD:
		return "+"
	case lexer.SUB:
		return "-"
	case lexer.BITOR:
		return "|"
	case lexer.MUL:
		return "*"
	case lexer.QUO:
		return "/"
	case lexer.REM:
		return "%"
	case lexer.BITAND:
		return "&"

	default:
		panic("ni binary" + op.String())
	}
}

func (genc *genContext) genBinaryExpr(x *ast.BinaryExpr) string {

	if ast.IsStringType(x.X.GetType()) {
		var not = ""
		if x.Op == lexer.NEQ {
			not = "!"
		}

		return fmt.Sprintf("%s%s(%s, %s)", not, rt_equalStrings, genc.genExpr(x.X), genc.genExpr(x.Y))
	}

	return fmt.Sprintf("(%s %s %s)", genc.genExpr(x.X), binaryOp(x.Op), genc.genExpr(x.Y))
}

//==== selector

func (genc *genContext) genSelector(x *ast.SelectorExpr) string {
	if x.X == nil {
		return genc.declName(x.Obj.(ast.Decl))
	}

	var cl = ast.UnderType(x.X.GetType()).(*ast.ClassType)
	return fmt.Sprintf("%s->%s.%s%s",
		genc.genExpr(x.X),
		nm_class_fields,
		pathToField(cl, x.Name),
		genc.outName(x.Name))
}

func pathToField(cl *ast.ClassType, name string) string {
	var path = ""
	for {
		if cl.BaseTyp == nil {
			break
		}
		cl = ast.UnderType(cl.BaseTyp).(*ast.ClassType)

		_, ok := cl.Members[name]
		if !ok {
			break
		}
		path += nm_base_fields + "."
	}
	return path
}

//==== индексация и композит массива

func (genc *genContext) genBracketExpr(x *ast.GeneralBracketExpr) string {

	if x.Index != nil {

		switch xt := ast.UnderType(x.X.GetType()).(type) {
		case *ast.VectorType:
			return genc.genVectorIndex(x.X, x.Index)
		case *ast.VariadicType:
			return genc.genVariadicIndex(xt, x.X, x.Index)
		default:
			panic("assert")
		}
	}

	return genc.genArrayComposite(x.Composite)
}

func (genc *genContext) genVectorIndex(x, inx ast.Expr) string {
	var name string
	if id, ok := x.(*ast.IdentExpr); ok {
		name = genc.genIdent(id)
	} else {
		name = genc.localName("")

		genc.c("%s %s = %s;",
			genc.typeRef(x.GetType()),
			name,
			genc.genExpr(x))
	}
	return fmt.Sprintf("%s->body[%s(%s, %s->len)]",
		name,
		rt_indexcheck,
		genc.genExpr(inx),
		name)

}

func (genc *genContext) genVariadicIndex(vt *ast.VariadicType, x, inx ast.Expr) string {

	var vPar = genc.genExpr(x)

	if ast.IsTagPairType(vt.ElementTyp) {
		/*
			return fmt.Sprintf("((%s*)(%s + sizeof(TInt64)))[%s(%s, *(TInt64 *)%s) << 1]",
				predefinedTypeName(ast.Int64.Name),
				vPar,
				rt_indexcheck,
				genc.genExpr(inx),
				vPar)
		*/
		panic("assert")
	} else {

		return fmt.Sprintf("((%s*)%s)[%s(%s, %s%s)]",
			genc.typeRef(vt.ElementTyp),
			vPar,
			rt_indexcheck,
			genc.genExpr(inx),
			vPar,
			nm_variadic_len_suffic)
	}
}

func (genc *genContext) genArrayComposite(x *ast.ArrayCompositeExpr) string {
	var name = genc.localName("")

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

	genc.c("%s", s)

	return name
}

func (genc *genContext) genClassComposite(x *ast.ClassCompositeExpr) string {
	var name = genc.localName("")

	var tname = genc.typeRef(x.Typ)
	var s = fmt.Sprintf("%s %s = %s(%s);",
		tname,
		name,
		rt_newObject,
		tname+nm_class_info_ptr_suffix)

	var cl = ast.UnderType(x.Typ).(*ast.ClassType)

	var list = make([]string, len(x.Values))
	for i, v := range x.Values {
		list[i] = fmt.Sprintf("%s->%s.%s%s = %s;",
			name, nm_class_fields, pathToField(cl, v.Name), genc.outName(v.Name), genc.genExpr(v.Value))
	}
	s += strings.Join(list, " ")

	genc.c("%s", s)
	return name
}
