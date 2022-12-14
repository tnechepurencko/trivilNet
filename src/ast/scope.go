package ast

import (
	"fmt"
)

var _ = fmt.Printf

type Scope struct {
	Outer *Scope
	Names map[string]Decl
}

var topScope *Scope

func initScopes() {
	topScope = &Scope{
		Names: make(map[string]Decl),
	}

	addType("Байт")
	addType("Цел")
	addType("Цел64")
	addType("Вещ64")
	addType("Лог")
	addType("Строка")
	//ShowScopes("top", topScope)

}

func addType(name string) {
	var td = &TypeDecl{
		Typ: &PredefinedType{
			Name: name,
		},
	}
	td.Name = name
	topScope.Names[name] = td
}

func NewScope(outer *Scope) *Scope {
	return &Scope{
		Outer: outer,
		Names: make(map[string]Decl),
	}
}

func ShowScopes(label string, cur *Scope) {
	if label != "" {
		fmt.Println(label)
	}

	var n = 0
	for cur != nil {
		n++
		fmt.Printf("<--- scope %d\n", n)
		for _, d := range cur.Names {
			fmt.Printf("%s ", d.GetName())
		}
		if len(cur.Names) > 0 {
			fmt.Println()
		}
		cur = cur.Outer
	}
	fmt.Printf("--- end scopes\n")
}
