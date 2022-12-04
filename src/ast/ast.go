package ast

import (
	"fmt"
)

var _ = fmt.Printf

//==== Interfaces

type Node interface {
	Pos() int
}

type Type interface {
	Node
}

type Decl interface {
	Node
}

type Expr interface {
	Node
}

type Stmt interface {
	Node
}

//==== declarations

type VarDecl struct {
	DeclBase
}
