package ast

import (
	"fmt"
)

var _ = fmt.Printf

//==== Interfaces

type Node interface {
	GetPos() int
	Accept(visitor Visitor) TreePrinter
}

type Type interface {
	Node
	TypeNode()
	Accept(visitor Visitor) TreePrinter
}

type Decl interface {
	Node
	DeclNode()
	GetName() string
	GetType() Type
	GetHost() *Module // только для объектов уровня модуля, для остальных - nil
	SetHost(host *Module)
	IsExported() bool
	Accept(visitor Visitor) TreePrinter
}

type Expr interface {
	Node
	ExprNode()
	GetType() Type
	SetType(t Type)
	IsReadOnly() bool
	Accept(visitor Visitor) TreePrinter
}

type Statement interface {
	Node
	StatementNode()
	Accept(visitor Visitor) TreePrinter
}

//==== init

func init() {
	initScopes()
}
