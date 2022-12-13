package ast

import (
	"fmt"
	//	"trivil/env"
)

var _ = fmt.Printf

//=== модуль

type Module struct {
	//source *env.Source
	Name  string
	Inner *Scope
	Decls []Decl
	Entry *EntryFn
}

func NewModule() *Module {
	return &Module{
		Inner: NewScope(topScope),
		Decls: make([]Decl, 0),
	}
}

//=== вход

type EntryFn struct {
	Pos int
	Seq *StatementSeq
}

func (n *EntryFn) GetPos() int {
	return n.Pos
}
