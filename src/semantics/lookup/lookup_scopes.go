package lookup

import (
	"fmt"
	"trivil/ast"
	"trivil/env"
)

var _ = fmt.Printf

func addToScope(name string, d ast.Decl, scope *ast.Scope) {
	_, ok := scope.Names[name]
	if ok {
		env.AddError(d.GetPos(), "СЕМ-УЖЕ-ОПИСАНО", name)
		return
	}
	scope.Names[name] = d

	//fmt.Printf("scope: %v\n", scope.Names)
}

func findInScopes(scope *ast.Scope, name string, pos int) ast.Decl {

	var cur = scope

	for {
		if cur == nil {

			//ast.ShowScopes("not found "+name, scope)

			env.AddError(pos, "СЕМ-НЕ-НАЙДЕНО", name)
			var inv = &ast.InvalidDecl{
				Name: name,
			}
			addToScope(name, inv, scope)
			return inv
		}

		d, ok := cur.Names[name]
		if ok {
			return d
		}

		cur = cur.Outer
	}

}
