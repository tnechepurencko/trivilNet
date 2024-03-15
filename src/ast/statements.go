package ast

import (
	"fmt"
)

var _ = fmt.Printf

//====

type StatementBase struct {
	Pos int
}

func (m *StatementBase) Accept(visitor Visitor) TreePrinter {
	return visitor.VisitStatementBase(m)
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

func (m *StatementSeq) Accept(visitor Visitor) TreePrinter {
	return visitor.VisitStatementSeq(m)
}

type ExprStatement struct {
	StatementBase
	X Expr
}

func (m *ExprStatement) Accept(visitor Visitor) TreePrinter {
	return visitor.VisitExprStatement(m)
}

type DeclStatement struct {
	StatementBase
	D Decl
}

func (m *DeclStatement) Accept(visitor Visitor) TreePrinter {
	return visitor.VisitDeclStatement(m)
}

type AssignStatement struct {
	StatementBase
	L Expr
	R Expr
}

func (m *AssignStatement) Accept(visitor Visitor) TreePrinter {
	return visitor.VisitAssignStatement(m)
}

type IncStatement struct {
	StatementBase
	L Expr
}

func (m *IncStatement) Accept(visitor Visitor) TreePrinter {
	return visitor.VisitIncStatement(m)
}

type DecStatement struct {
	StatementBase
	L Expr
}

func (m *DecStatement) Accept(visitor Visitor) TreePrinter {
	return visitor.VisitDecStatement(m)
}

//==== управление

type If struct {
	StatementBase
	Cond Expr
	Then *StatementSeq
	Else Statement
}

func (m *If) Accept(visitor Visitor) TreePrinter {
	return visitor.VisitIf(m)
}

type Guard struct {
	StatementBase
	Cond Expr
	Else Statement
}

func (m *Guard) Accept(visitor Visitor) TreePrinter {
	return visitor.VisitGuard(m)
}

type Select struct {
	StatementBase
	X     Expr // = nil, если предикатный оператор
	Cases []*Case
	Else  *StatementSeq
}

func (m *Select) Accept(visitor Visitor) TreePrinter {
	return visitor.VisitSelect(m)
}

type Case struct {
	StatementBase
	Exprs []Expr
	Seq   *StatementSeq
}

func (m *Case) Accept(visitor Visitor) TreePrinter {
	return visitor.VisitCase(m)
}

type SelectType struct {
	StatementBase
	VarIdent string
	X        Expr
	Cases    []*CaseType
	Else     *StatementSeq
}

func (m *SelectType) Accept(visitor Visitor) TreePrinter {
	return visitor.VisitSelectType(m)
}

type CaseType struct {
	StatementBase
	Types []Type
	Var   *VarDecl // nil, если переменная не задана
	Seq   *StatementSeq
}

func (m *CaseType) Accept(visitor Visitor) TreePrinter {
	return visitor.VisitCaseType(m)
}

//==== циклы

type While struct {
	StatementBase
	Cond Expr
	Seq  *StatementSeq
}

func (m *While) Accept(visitor Visitor) TreePrinter {
	return visitor.VisitWhile(m)
}

type Cycle struct {
	StatementBase
	IndexVar   *VarDecl
	ElementVar *VarDecl
	Expr       Expr
	Seq        *StatementSeq
}

func (m *Cycle) Accept(visitor Visitor) TreePrinter {
	return visitor.VisitCycle(m)
}

//==== завершающие

type Crash struct {
	StatementBase
	X Expr //
}

func (m *Crash) Accept(visitor Visitor) TreePrinter {
	return visitor.VisitCrash(m)
}

type Return struct {
	StatementBase
	ReturnTyp Type
	X         Expr
}

func (m *Return) Accept(visitor Visitor) TreePrinter {
	return visitor.VisitReturn(m)
}

type Break struct {
	StatementBase
}

func (m *Break) Accept(visitor Visitor) TreePrinter {
	return visitor.VisitBreak(m)
}
