package ast

import (
	"fmt"
	//	"trivil/env"
)

var _ = fmt.Printf

//====

type DeclBase struct {
	Pos      int
	Exported bool
}

func (n *DeclBase) DeclNode() {}

func (n *DeclBase) GetPos() int {
	return n.Pos
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
	Name     string
	Typ      Type
	Inner    *Scope
	Seq      *StatementSeq
	External bool
}

type VarDecl struct {
	DeclBase
	Name string
	Typ  Type
}

type ConstDecl struct {
	DeclBase
	Name string
	Typ  Type
}

type TypeDecl struct {
	DeclBase
	Name string
	Typ  Type
}
