package lookup

import (
	"fmt"
	"trivil/ast"
)

var _ = fmt.Printf

func (lc *lookContext) lookExpr(expr ast.Expr) {

	switch x := expr.(type) {
	case *ast.IdentExpr:
		var d = findInScopes(lc.scope, x.Name, x.Pos)
		x.Obj = lc.considerTypeRef(d, x.Pos)

		//fmt.Printf("found %v => %v\n", x.Name, x.Obj)

	case *ast.UnaryExpr:
		lc.lookExpr(x.X)

	case *ast.BinaryExpr:
		lc.lookExpr(x.X)
		lc.lookExpr(x.Y)

	case *ast.SelectorExpr:
		lc.lookExpr(x.X)
		lc.lookImported(x)

	case *ast.CallExpr:
		lc.lookExpr(x.X)
		for _, a := range x.Args {
			lc.lookExpr(a)
		}
		lc.lookStdFunction(x)

	case *ast.GeneralBracketExpr:
		lc.lookExpr(x.X)
		if x.Index != nil {
			lc.lookExpr(x.Index)
		}

		for _, e := range x.Composite.Elements {
			if e.Key != nil {
				lc.lookExpr(e.Key)
			}
			lc.lookExpr(e.Value)
		}
	case *ast.ClassCompositeExpr:
		lc.lookExpr(x.X)

		for _, vp := range x.Values {
			lc.lookExpr(vp.Value)

		}
	case *ast.LiteralExpr:
		//nothing

	default:
		panic(fmt.Sprintf("expression: ni %T", expr))

	}
}

// Возврашает TypeRef для TypeDecl, или сам объект
func (lc *lookContext) considerTypeRef(d ast.Decl, pos int) ast.Node {

	if td, ok := d.(*ast.TypeDecl); ok {
		return &ast.TypeRef{
			TypeBase: ast.TypeBase{Pos: pos},
			TypeName: td.Name,
			//ModuleName: ?
			TypeDecl: td,
			Typ:      td.Typ,
		}
	}

	return d

}

func (lc *lookContext) lookImported(x *ast.SelectorExpr) {

	ident, ok := x.X.(*ast.IdentExpr)
	if !ok {
		return
	}

	m, ok := ident.Obj.(*ast.Module)
	if !ok {
		return
	}

	if d, ok := m.Inner.Names[x.Name]; ok {
		x.Obj = lc.considerTypeRef(d, x.Pos)
	} else {
		var inv = &ast.InvalidDecl{
			DeclBase: ast.DeclBase{Pos: x.Pos, Name: x.Name},
		}
		x.Obj = inv
		m.Inner.Names[x.Name] = inv
		panic("add and test error")
	}
	x.X = nil
}

func (lc *lookContext) lookStdFunction(x *ast.CallExpr) {

	ident, ok := x.X.(*ast.IdentExpr)
	if !ok {
		return
	}

	if ident.Obj == nil {
		return
	}

	stdf, ok := ident.Obj.(*ast.StdFunction)
	if !ok {
		return
	}

	x.StdFunc = stdf
	x.X = nil
}
