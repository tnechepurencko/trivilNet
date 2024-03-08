package ast

import (
	"fmt"
)

var _ = fmt.Printf

//==== Interfaces

type Node interface {
	GetPos() int
	Accept(visitor Visitor)
}

type Type interface {
	Node
	TypeNode()
	Accept(visitor Visitor)
}

type Decl interface {
	Node
	DeclNode()
	GetName() string
	GetType() Type
	GetHost() *Module // только для объектов уровня модуля, для остальных - nil
	SetHost(host *Module)
	IsExported() bool
	Accept(visitor Visitor)
}

type Expr interface {
	Node
	ExprNode()
	GetType() Type
	SetType(t Type)
	IsReadOnly() bool
	Accept(visitor Visitor)
}

type Statement interface {
	Node
	StatementNode()
	Accept(visitor Visitor)
}

//==== init

func init() {
	initScopes()
}
