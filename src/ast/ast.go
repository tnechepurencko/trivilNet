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
	IsExported() bool
}

type Expr interface {
	Node
	ExprNode()
	GetType() Type
	IsReadOnly() bool
}

type Statement interface {
	Node
	StatementNode()
}

//==== init

func init() {
	initScopes()
}
