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

	addType("цел")
}

func addType(name string) {
	var td = &TypeDecl{
		Name: name,
		Typ: &PredefinedType{
			Name: name,
		},
	}
	topScope.Names[name] = td
}
