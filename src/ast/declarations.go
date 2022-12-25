package ast

import (
	"fmt"
	//	"trivil/env"
)

var _ = fmt.Printf

//====

type DeclBase struct {
	Pos      int
	Name     string
	Typ      Type
	Exported bool
}

func (n *DeclBase) DeclNode() {}

func (n *DeclBase) GetPos() int {
	return n.Pos
}

func (n *DeclBase) GetName() string {
	return n.Name
}

func (n *DeclBase) GetType() Type {
	return n.Typ
}

func (n *DeclBase) IsExported() bool {
	return n.Exported
}

func (n *DeclBase) SetExported() {
	n.Exported = true
}

//====

type InvalidDecl struct {
	DeclBase
	Name string
}

//=== описания

type Function struct {
	DeclBase
	Recv     *Param
	Inner    *Scope
	Seq      *StatementSeq
	External bool
}

type VarDecl struct {
	DeclBase
	Init Expr
}

type ConstDecl struct {
	DeclBase
	Value Expr
}

type TypeDecl struct {
	DeclBase
}
