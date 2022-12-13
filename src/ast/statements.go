package ast

import (
	"fmt"
)

var _ = fmt.Printf

//====

type StatementBase struct {
	Pos int
}

func (n *StatementBase) GetPos() int {
	return n.Pos
}
func (n *StatementBase) StatementNode() {}

//====

type StatementSeq struct {
	StatementBase
	Statements []Statement
	Inner      *Scope
}

type ExprStatement struct {
	StatementBase
	X Expr
}

type DeclStatement struct {
	StatementBase
	D Decl
}

type AssignStatement struct {
	StatementBase
	L Expr
	R Expr
}

type IncStatement struct {
	StatementBase
	L Expr
}

type DecStatement struct {
	StatementBase
	L Expr
}
