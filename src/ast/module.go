package ast

import (
	"fmt"
	//	"trivil/env"
)

var _ = fmt.Printf

//=== модуль

type Module struct {
	//source *env.Source
	Name    string
	Imports []*Import
	Decls   []Decl
	Entry   *EntryFn
	Inner   *Scope
}

func NewModule() *Module {
	return &Module{
		Inner: NewScope(topScope),
		Decls: make([]Decl, 0),
	}
}

//=== импорт

type Import struct {
	Pos  int
	Path string
	Mod  Module
}

func (n *Import) GetPos() int {
	return n.Pos
}

//=== вход

type EntryFn struct {
	Pos int
	Seq *StatementSeq
}

func (n *EntryFn) GetPos() int {
	return n.Pos
}
