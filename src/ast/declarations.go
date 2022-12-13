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
	Exported bool
}

func (n *DeclBase) DeclNode() {}

func (n *DeclBase) GetPos() int {
	return n.Pos
}

func (n *DeclBase) GetName() string {
	return n.Name
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
	Typ      Type
	Inner    *Scope
	Seq      *StatementSeq
	External bool
}

type VarDecl struct {
	DeclBase
	Typ Type
}

type ConstDecl struct {
	DeclBase
	Typ   Type
	Value Expr
}

type TypeDecl struct {
	DeclBase
	Typ Type
}
