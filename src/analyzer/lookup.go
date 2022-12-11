package analyzer

import (
	"fmt"
	"trivil/ast"
	//"trivil/env"
)

var _ = fmt.Printf

func lookup(m *ast.Module) {

	// добавление имен
	for _, d := range m.Decls {
		switch x := d.(type) {
		case *ast.Function:
			//			fmt.Printf("Function %v\n", x.Name)
			addToScope(x.Name, x, m.Inner)
		default:
			panic(fmt.Sprintf("lookup: ni %T", d))
		}
	}

	// process decls

	if m.Entry != nil {
		processEntry(m.Inner, m.Entry)
	}

}

//====

func processEntry(scope *ast.Scope, e *ast.EntryFn) {
	processStatements(scope, e.Seq)
}

func processStatements(scope *ast.Scope, seq *ast.StatementSeq) {

	for _, s := range seq.Statements {

		switch x := s.(type) {
		case *ast.ExprStatement:
			processExpr(scope, x.X)

		default:
			panic(fmt.Sprintf("statement: ni %T", s))

		}
	}
}

//====

func processExpr(scope *ast.Scope, expr ast.Expr) {

	switch x := expr.(type) {
	case *ast.IdentExpr:
		x.Obj = findInScopes(scope, x.Name, x.Pos)
		//fmt.Printf("found %v => %v\n", x.Name, x.Obj)

	case *ast.CallExpr:
		processExpr(scope, x.X)
		//TODO: args

	default:
		panic(fmt.Sprintf("expression: ni %T", expr))

	}

}
