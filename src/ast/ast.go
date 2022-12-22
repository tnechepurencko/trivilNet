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
	GetName() string
	GetType() Type
}

type Expr interface {
	Node
	ExprNode()
	GetType() Type
}

type Statement interface {
	Node
	StatementNode()
}

//==== init

func init() {
	initScopes()
}
