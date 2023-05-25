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

//====

type If struct {
	StatementBase
	Cond Expr
	Then *StatementSeq
	Else Statement
}

type While struct {
	StatementBase
	Cond Expr
	Seq  *StatementSeq
}

type Guard struct {
	StatementBase
	Cond Expr
	Else Statement
}

type Select struct {
	StatementBase
	X     Expr // = nil, если предикатный оператор
	Cases []*Case
	Else  *StatementSeq
}

type Case struct {
	StatementBase
	Exprs []Expr
	Seq   *StatementSeq
}

type SelectType struct {
	StatementBase
	Ident string
	X     Expr
	Cases []*CaseType
	Else  *StatementSeq
}

type CaseType struct {
	StatementBase
	Types []Type
	Seq   *StatementSeq
}

type Crash struct {
	StatementBase
	X Expr //
}

type Return struct {
	StatementBase
	X Expr
}

type Break struct {
	StatementBase
}
