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

type Stmt interface {
	Node
	StmtNode()
}

//==== declarations

type DeclBase struct {
	Pos int
}

func (n *DeclBase) GetPos() int {
	return n.Pos
}
func (n *DeclBase) DeclNode() {}

type VarDecl struct {
	DeclBase
	Name string
	Typ  Type
}

func init() {
}
