package ast

import (
	"fmt"
)

var _ = fmt.Printf

type Scope struct {
	Outer *Scope
	Names map[string]Decl
}

var top *Scope

func initScopes() {
	top = &Scope{
		Names: make(map[string]Decl),
	}
}
