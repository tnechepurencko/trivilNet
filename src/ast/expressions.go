package ast

import (
	"fmt"
	"trivil/lexer"
)

var _ = fmt.Printf

//====

type ExprBase struct {
	Pos      int
	Typ      Type
	ReadOnly bool
}

func (n *ExprBase) ExprNode() {}

func (n *ExprBase) GetPos() int {
	return n.Pos
}

func (n *ExprBase) GetType() Type {
	return n.Typ
}

func (n *ExprBase) IsReadOnly() bool {
	return n.ReadOnly
}

//====

type InvalidExpr struct {
	ExprBase
}

//====

type BinaryExpr struct {
	ExprBase
	X  Expr
	Op lexer.Token
	Y  Expr
}

type UnaryExpr struct {
	ExprBase
	Op lexer.Token
	X  Expr
}

type LiteralExpr struct {
	ExprBase
	Kind lexer.Token // если STRING - см. Typ - может быть Символ или Байт
	Lit  string
}

type BoolLiteral struct {
	ExprBase
	Value bool
}

type IdentExpr struct {
	ExprBase
	Name string
	Obj  Node // Decl: Var, Const or Function or TypeRef
}

type SelectorExpr struct {
	ExprBase
	X    Expr // == nil, если импортированный объект
	Name string
	Obj  Node // импортированный объект или поле или метод
}

type CallExpr struct {
	ExprBase
	X       Expr
	Args    []Expr
	StdFunc *StdFunction
}

type ConversionExpr struct {
	ExprBase
	X         Expr
	TargetTyp Type
	Done      bool // X уже преобразован к целевому типу
}

//==== index

type GeneralBracketExpr struct {
	ExprBase
	X         Expr
	Index     Expr // indexation if != nil, otherwise composite
	Composite *ArrayCompositeExpr
}

type ArrayCompositeExpr struct {
	ExprBase
	Elements []ElementPair
	Keys     bool // both L and R are used: L - are indexes, R - values
}

type ElementPair struct {
	Key   Expr
	Value Expr
}

//=== class composite

type ClassCompositeExpr struct {
	ExprBase
	X      Expr
	Values []ValuePair
}

type ValuePair struct {
	Pos   int
	Name  string
	Value Expr
}
