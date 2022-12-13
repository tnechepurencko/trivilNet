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
			lc.processVarDecl(x)
		default:
			panic(fmt.Sprintf("lookup 2: ni %T", d))
		}
	}

	if m.Entry != nil {
		lc.processEntry(m.Entry)
	}

}

//==== описания

func (lc *lookContext) processVarDecl(v *ast.VarDecl) {
	lc.processTypeRef(v.Typ)

}

//==== ссылка на тип

func (lc *lookContext) processTypeRef(t ast.Type) {
	var tr = t.(*ast.TypeRef)

	if tr.ModuleName != "" {
		panic("ni")
	}

	tr.TypeDecl = findInScopes(lc.scope, tr.TypeName, tr.Pos).(*ast.TypeDecl)
	tr.Typ = tr.TypeDecl.Typ

	fmt.Printf("! %v %T\n", tr.TypeDecl, tr.Typ)

}

//====

func (lc *lookContext) processEntry(e *ast.EntryFn) {
	lc.processStatements(e.Seq)
}

func (lc *lookContext) processStatements(seq *ast.StatementSeq) {

	for _, s := range seq.Statements {

		switch x := s.(type) {
		case *ast.ExprStatement:
			lc.processExpr(x.X)
		case *ast.DeclStatement:
			lc.processLocalDecl(seq, x.D)

		default:
			panic(fmt.Sprintf("statement: ni %T", s))

		}
	}

	if lc.scope == seq.Inner {
		lc.scope = seq.Inner.Outer
	}
}

func (lc *lookContext) processLocalDecl(seq *ast.StatementSeq, decl ast.Decl) {
	if lc.scope != seq.Inner {
		seq.Inner = ast.NewScope(lc.scope)
		lc.scope = seq.Inner
	}
	switch x := decl.(type) {
	case *ast.VarDecl:
		addToScope(x.Name, x, lc.scope)
		lc.processVarDecl(x)
	default:
		panic(fmt.Sprintf("local decl: ni %T", decl))
	}
	ast.ShowScopes("", lc.scope)
}

//====

func (lc *lookContext) processExpr(expr ast.Expr) {

	switch x := expr.(type) {
	case *ast.IdentExpr:
		x.Obj = findInScopes(lc.scope, x.Name, x.Pos)
		//fmt.Printf("found %v => %v\n", x.Name, x.Obj)

	case *ast.CallExpr:
		lc.processExpr(x.X)
		//TODO: args

	default:
		panic(fmt.Sprintf("expression: ni %T", expr))

	}
}
