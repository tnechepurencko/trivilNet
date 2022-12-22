package ast

import (
	"fmt"
	"trivil/lexer"
)

var _ = fmt.Printf

//====

type ExprBase struct {
	Pos int
	Typ Type
}

func (n *ExprBase) ExprNode() {}

func (n *ExprBase) GetPos() int {
	return n.Pos
}

func (n *ExprBase) GetType() Type {
	return n.Typ
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
	Kind lexer.Token
	Lit  string
}

type BoolLiteral struct {
	ExprBase
	Value bool
}

type IdentExpr struct {
	ExprBase
	Name   string
	Obj    Decl     // var, const or function, never type
	TypRef *TypeRef // != if using type name in expression, like in composite
}

type SelectorExpr struct {
	ExprBase
	X    Expr
	Name string
}

type CallExpr struct {
	ExprBase
	X    Expr
	Args []Expr
}

type ConversionExpr struct {
	ExprBase
	X         Expr
	TargetTyp Type
}

//==== index

type IndexExpr struct {
	ExprBase
	X     Expr
	Index Expr // indexation if != nil
	// composite:
	Elements []ElementPair
	Pairs    bool // both L and R are used: L - are indexes, R - values
}

type ElementPair struct {
	L Expr
	R Expr
}

//=== composite

type CompositeExpr struct {
	ExprBase
	X      Expr
	Values []ValuePair
}

type ValuePair struct {
	Pos  int
	Name string
	V    Expr
}
