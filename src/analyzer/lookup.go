package analyzer

import (
	"fmt"
	"trivil/ast"
	//"trivil/env"
)

var _ = fmt.Printf

type lookContext struct {
	scope *ast.Scope
}

func lookup(m *ast.Module) {

	var lc = &lookContext{
		scope: m.Inner,
	}

	// добавление имен
	for _, d := range m.Decls {
		switch x := d.(type) {
		case *ast.Function:
			//			fmt.Printf("Function %v\n", x.Name)
			addToScope(x.Name, x, m.Inner)
		case *ast.VarDecl:
			//			fmt.Printf("Function %v\n", x.Name)
			addToScope(x.Name, x, m.Inner)
		default:
			panic(fmt.Sprintf("lookup: ni %T", d))
		}
	}

	if lc.scope != m.Inner {
		panic("assert - should be module scope")
	}

	// TODO обойти типы

	// обойти описания
	for _, d := range m.Decls {
		switch x := d.(type) {
		case *ast.Function:
		case *ast.VarDecl:
			lc.lookVarDecl(x)
		default:
			panic(fmt.Sprintf("lookup 2: ni %T", d))
		}
	}

	if m.Entry != nil {
		lc.lookEntry(m.Entry)
	}

}

//==== описания

func (lc *lookContext) lookVarDecl(v *ast.VarDecl) {
	lc.lookTypeRef(v.Typ)

}

//==== ссылка на тип

func (lc *lookContext) lookTypeRef(t ast.Type) {
	var tr = t.(*ast.TypeRef)

	if tr.ModuleName != "" {
		panic("ni")
	}

	var x = findInScopes(lc.scope, tr.TypeName, tr.Pos)
	td, ok := x.(*ast.TypeDecl)
	if !ok {
		return
	}

	tr.TypeDecl = td
	tr.Typ = tr.TypeDecl.Typ

	//fmt.Printf("! %v %T\n", tr.TypeDecl, tr.Typ)
}

//====

func (lc *lookContext) lookEntry(e *ast.EntryFn) {
	lc.lookStatements(e.Seq)
}

func (lc *lookContext) lookStatements(seq *ast.StatementSeq) {

	for _, s := range seq.Statements {

		switch x := s.(type) {
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

	if lc.scope == seq.Inner {
		lc.scope = seq.Inner.Outer
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
		x.Obj = findInScopes(lc.scope, x.Name, x.Pos)
		//fmt.Printf("found %v => %v\n", x.Name, x.Obj)

	case *ast.LiteralExpr:
		//lc.lookExpr(x.X)

	case *ast.UnaryExpr:
		lc.lookExpr(x.X)

	case *ast.BinaryExpr:
		lc.lookExpr(x.X)
		lc.lookExpr(x.Y)

	case *ast.CallExpr:
		lc.lookExpr(x.X)
		//TODO: args

	default:
		panic(fmt.Sprintf("expression: ni %T", expr))

	}
}
