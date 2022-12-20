package check

import (
	"fmt"
	"trivil/ast"
	//	"trivil/env"
)

var _ = fmt.Printf

type checkContext struct {
}

func Process(m *ast.Module) {
	var cc = &checkContext{}

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

	/*
		if m.Entry != nil {
			lc.lookEntry(m.Entry)
		}
	*/
}

//==== константы и переменные
/*
func (lc *lookContext) lookVarDecl(v *ast.VarDecl) {
	lc.lookTypeRef(v.Typ)
}

func (lc *lookContext) lookConstDecl(v *ast.ConstDecl) {
	lc.lookTypeRef(v.Typ)

}

//==== functions

func (lc *lookContext) lookFunction(f *ast.Function) {

	f.Inner = ast.NewScope(lc.scope)
	lc.scope = f.Inner

	if f.Recv.Typ != nil {
		lc.lookTypeRef(f.Recv.Typ)

		lc.addMethodToType(f)
	}

	var ft = f.Typ.(*ast.FuncType)

	for _, p := range ft.Params {
		lc.lookTypeRef(p.Typ)
		if !f.External {
			lc.addVarForParameter(p)
		}
	}

	if ft.ReturnTyp != nil {
		lc.lookTypeRef(ft.ReturnTyp)
	}

	if !f.External {
		lc.lookStatements(f.Seq)
	}

	lc.scope = lc.scope.Outer
}

func (lc *lookContext) addMethodToType(f *ast.Function) {

	var rt = f.Recv.Typ.(*ast.TypeRef)

	cl, ok := rt.Typ.(*ast.ClassType)
	if !ok {
		env.AddError(f.Recv.Pos, "СЕМ-ПОЛУЧАТЕЛЬ-КЛАСС")
		return
	}

	cl.Methods = append(cl.Methods, f)

}

func (lc *lookContext) addVarForParameter(p *ast.Param) {
	var v = &ast.VarDecl{
		Typ: p.Typ,
	}
	v.Name = p.Name
	addToScope(v.Name, v, lc.scope)
}

func (lc *lookContext) lookEntry(e *ast.EntryFn) {
	lc.lookStatements(e.Seq)
}

//==== statements

func (lc *lookContext) lookStatements(seq *ast.StatementSeq) {

	for _, s := range seq.Statements {
		lc.lookStatement(seq, s)
	}

	if lc.scope == seq.Inner {
		lc.scope = seq.Inner.Outer
	}
}

func (lc *lookContext) lookStatement(seq *ast.StatementSeq, s ast.Statement) {
	switch x := s.(type) {
	case *ast.StatementSeq:
		lc.lookStatements(x)
	case *ast.ExprStatement:
		lc.lookExpr(x.X)
	case *ast.DeclStatement:
		lc.lookLocalDecl(seq, x.D)
	case *ast.AssignStatement:
		lc.lookExpr(x.L)
		lc.lookExpr(x.R)
	case *ast.IncStatement:
		lc.lookExpr(x.L)
	case *ast.DecStatement:
		lc.lookExpr(x.L)
	case *ast.If:
		lc.lookExpr(x.Cond)
		lc.lookStatements(x.Then)
		if x.Else != nil {
			lc.lookStatement(nil, x.Else)
		}
	case *ast.While:
		lc.lookExpr(x.Cond)
		lc.lookStatements(x.Seq)
	case *ast.Return:
		if x.X != nil {
			lc.lookExpr(x.X)
		}

	default:
		panic(fmt.Sprintf("statement: ni %T", s))

	}
}

func (lc *lookContext) lookLocalDecl(seq *ast.StatementSeq, decl ast.Decl) {
	if lc.scope != seq.Inner {
		seq.Inner = ast.NewScope(lc.scope)
		lc.scope = seq.Inner
	}
	switch x := decl.(type) {
	case *ast.VarDecl:
		addToScope(x.Name, x, lc.scope)
		lc.lookVarDecl(x)
	default:
		panic(fmt.Sprintf("local decl: ni %T", decl))
	}
	//ast.ShowScopes("", lc.scope)
}

//====

func (lc *lookContext) lookExpr(expr ast.Expr) {

	switch x := expr.(type) {
	case *ast.IdentExpr:
		var d = findInScopes(lc.scope, x.Name, x.Pos)
		if td, ok := d.(*ast.TypeDecl); ok {
			x.TypRef = &ast.TypeRef{
				TypeName: td.Name,
				//ModuleName: ?
				TypeDecl: td,
				Typ:      td.Typ,
			}
			x.TypRef.Pos = x.Pos
		} else {
			x.Obj = d
		}

		//fmt.Printf("found %v => %v\n", x.Name, x.Obj)

	case *ast.LiteralExpr:
		//lc.lookExpr(x.X)

	case *ast.UnaryExpr:
		lc.lookExpr(x.X)

	case *ast.BinaryExpr:
		lc.lookExpr(x.X)
		lc.lookExpr(x.Y)

	case *ast.SelectorExpr:
		lc.lookExpr(x.X)
		panic("ni")

	case *ast.CallExpr:
		lc.lookExpr(x.X)
		for _, a := range x.Args {
			lc.lookExpr(a)
		}

	case *ast.IndexExpr:
		lc.lookExpr(x.X)
		if x.Index != nil {
			lc.lookExpr(x.Index)
		}

		for _, e := range x.Elements {
			lc.lookExpr(e.L)
			if e.R != nil {
				lc.lookExpr(e.R)
			}
		}
	case *ast.CompositeExpr:
		lc.lookExpr(x.X)

		for _, vp := range x.Values {
			lc.lookExpr(vp.V)

		}

	default:
		panic(fmt.Sprintf("expression: ni %T", expr))

	}
}
*/
