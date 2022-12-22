package check

import (
	"fmt"
	"trivil/ast"
	"trivil/env"
	"trivil/lexer"
)

var _ = fmt.Printf

type checkContext struct {
	checkedTypes map[string]struct{}
}

func Process(m *ast.Module) {
	var cc = &checkContext{
		checkedTypes: make(map[string]struct{}),
	}

	for _, d := range m.Decls {
		switch x := d.(type) {
		case *ast.TypeDecl:
			cc.typeDecl(x)
		case *ast.VarDecl:
		case *ast.ConstDecl:
		case *ast.Function:
		default:
			panic(fmt.Sprintf("check: ni %T", d))
		}
	}

	if m.Entry != nil {
		cc.entry(m.Entry)
	}
}

//==== константы и переменные
/*
func (cc *checkContext) lookVarDecl(v *ast.VarDecl) {
	cc.TypeRef(v.Typ)
}

func (cc *checkContext) lookConstDecl(v *ast.ConstDecl) {
	cc.TypeRef(v.Typ)

}

//==== functions

func (cc *checkContext) lookFunction(f *ast.Function) {

	f.Inner = ast.NewScope(lc.scope)
	lc.scope = f.Inner

	if f.Recv.Typ != nil {
		cc.TypeRef(f.Recv.Typ)

		lc.addMethodToType(f)
	}

	var ft = f.Typ.(*ast.FuncType)

	for _, p := range ft.Params {
		cc.TypeRef(p.Typ)
		if !f.External {
			lc.addVarForParameter(p)
		}
	}

	if ft.ReturnTyp != nil {
		cc.TypeRef(ft.ReturnTyp)
	}

	if !f.External {
		cc.Statements(f.Seq)
	}

	lc.scope = lc.scope.Outer
}

func (cc *checkContext) addMethodToType(f *ast.Function) {

	var rt = f.Recv.Typ.(*ast.TypeRef)

	cl, ok := rt.Typ.(*ast.ClassType)
	if !ok {
		env.AddError(f.Recv.Pos, "СЕМ-ПОЛУЧАТЕЛЬ-КЛАСС")
		return
	}

	cl.Methods = append(cl.Methods, f)

}

func (cc *checkContext) addVarForParameter(p *ast.Param) {
	var v = &ast.VarDecl{
		Typ: p.Typ,
	}
	v.Name = p.Name
	addToScope(v.Name, v, lc.scope)
}
*/

func (cc *checkContext) entry(e *ast.EntryFn) {
	cc.statements(e.Seq)
}

//==== statements

func (cc *checkContext) statements(seq *ast.StatementSeq) {

	for _, s := range seq.Statements {
		cc.statement(s)
	}
}

func (cc *checkContext) statement(s ast.Statement) {
	switch x := s.(type) {
	case *ast.StatementSeq:
		cc.statements(x) // из else
	case *ast.ExprStatement:
		cc.expr(x.X)
		//TODO: проверить, что есть вызов, иначе ошибка
	case *ast.DeclStatement:
		//cc.localDecl(seq, x.D)
	case *ast.AssignStatement:
		cc.expr(x.L)
		cc.expr(x.R)
	case *ast.IncStatement:
		cc.expr(x.L)
	case *ast.DecStatement:
		cc.expr(x.L)
	case *ast.If:
		cc.expr(x.Cond)
		if x.Cond.GetType() != ast.Bool {
			env.AddError(x.Cond.GetPos(), "СЕМ-ТИП-ВЫРАЖЕНИЯ", ast.Bool.Name)
		}

		cc.statements(x.Then)
		if x.Else != nil {
			cc.statement(x.Else)
		}
	case *ast.While:
		cc.expr(x.Cond)
		cc.statements(x.Seq)
	case *ast.Return:
		if x.X != nil {
			cc.expr(x.X)
		}

	default:
		panic(fmt.Sprintf("statement: ni %T", s))
	}
}

/*
func (cc *checkContext) localDecl(seq *ast.StatementSeq, decl ast.Decl) {

	if lc.scope != seq.Inner {
		seq.Inner = ast.NewScope(lc.scope)
		lc.scope = seq.Inner
	}
	switch x := decl.(type) {
	case *ast.VarDecl:
		addToScope(x.Name, x, lc.scope)
		cc.VarDecl(x)
	default:
		panic(fmt.Sprintf("local decl: ni %T", decl))
	}
}
*/

//====

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
			//TODO: нужный тип
			x.Typ = ast.Int
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

			env.AddError(x.Args[i].GetPos(), "СЕМ-НЕСОВМЕСТИМЫй-АРГУМЕНТ", res)
		}
	}
}
