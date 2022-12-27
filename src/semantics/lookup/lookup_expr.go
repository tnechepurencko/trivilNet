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

	default:
		panic(fmt.Sprintf("expression: ni %T", expr))

	}
}
