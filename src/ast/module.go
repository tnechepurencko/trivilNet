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

func (n *Module) DeclNode() {}

func (n *Module) GetPos() int {
	panic("assert")
}

func (n *Module) GetName() string {
	return n.Name
}

func (n *Module) GetType() Type {
	panic("assert")
}

func (n *Module) IsExported() bool {
	return false
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
	Mod  *Module
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
