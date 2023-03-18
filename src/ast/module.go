package ast

import (
	"fmt"
	//	"trivil/env"
)

var _ = fmt.Printf

//=== модуль

type Module struct {
	//source *env.Source
	Pos     int
	Name    string
	Imports []*Import
	Decls   []Decl
	Entry   *EntryFn
	Inner   *Scope
	Caution bool
	Setting *Setting
}

func (n *Module) DeclNode() {}

func (n *Module) GetPos() int {
	return n.Pos
}

func (n *Module) GetName() string {
	return n.Name
}

func (n *Module) GetType() Type {
	panic("assert")
}

func (n *Module) GetHost() *Module {
	return nil
}

func (n *Module) SetHost(host *Module) {
}

func (n *Module) IsExported() bool {
	return false
}

func NewModule() *Module {
	return &Module{
		Inner:   NewScope(topScope),
		Decls:   make([]Decl, 0),
		Imports: make([]*Import, 0),
	}
}

//=== конкретизация

type Setting struct {
	Pos  int
	Path string
	//Attrs map[string]string
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
