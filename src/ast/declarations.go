package ast

import (
	"fmt"
	//"trivil/env"
)

var _ = fmt.Printf

//====

type DeclBase struct {
	Pos      int
	Name     string
	Typ      Type
	Host     *Module
	Exported bool
}

func (m *DeclBase) Accept(visitor Visitor) TreePrinter {
	return visitor.VisitDeclBase(m)
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

func (n *DeclBase) GetHost() *Module {
	return n.Host
}

func (n *DeclBase) SetHost(host *Module) {
	n.Host = host
}

func (n *DeclBase) IsExported() bool {
	return n.Exported
}

//====

type InvalidDecl struct {
	DeclBase
}

func (m *InvalidDecl) Accept(visitor Visitor) TreePrinter {
	return visitor.VisitInvalidDecl(m)
}

//=== описания

type Function struct {
	DeclBase
	Recv     *Param
	Seq      *StatementSeq
	External bool
	Mod      *Modifier
}

func (m *Function) Accept(visitor Visitor) TreePrinter {
	return visitor.VisitFunction(m)
}

type VarDecl struct {
	DeclBase
	Init       Expr
	Later      bool
	AssignOnce bool
	OutParam   bool // если это выходной параметр
}

func (m *VarDecl) Accept(visitor Visitor) TreePrinter {
	return visitor.VisitVarDecl(m)
}

type ConstDecl struct {
	DeclBase
	Value Expr
}

func (m *ConstDecl) Accept(visitor Visitor) TreePrinter {
	return visitor.VisitConstDecl(m)
}

type TypeDecl struct {
	DeclBase
}

func (m *TypeDecl) Accept(visitor Visitor) TreePrinter {
	return visitor.VisitTypeDecl(m)
}

//====

type StdFunction struct {
	DeclBase
	Method bool
}

func (m *StdFunction) Accept(visitor Visitor) TreePrinter {
	return visitor.VisitStdFunction(m)
}
