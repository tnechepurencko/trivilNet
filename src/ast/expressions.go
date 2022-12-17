package ast

import (
	"fmt"
	"trivil/lexer"
)

var _ = fmt.Printf

//====

type ExprBase struct {
	Pos int
}

func (n *ExprBase) GetPos() int {
	return n.Pos
}
func (n *ExprBase) ExprNode() {}

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

type IdentExpr struct {
	ExprBase
	Name string
	Obj  Decl
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
