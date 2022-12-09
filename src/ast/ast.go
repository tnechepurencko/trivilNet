package ast

import (
	"fmt"
)

var _ = fmt.Printf

//==== Interfaces

type Node interface {
	GetPos() int
}

type Type interface {
	Node
	TypeNode()
}

type Decl interface {
	Node
	DeclNode()
}

type Expr interface {
	Node
	ExprNode()
}

type Statement interface {
	Node
	StatementNode()
}

//==== init

func init() {
	initScopes()
}
