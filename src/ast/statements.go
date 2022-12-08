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
}
