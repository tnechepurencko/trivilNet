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
	Setting *Setting
}

func (m *Module) Accept(visitor Visitor) TreePrinter {
	return visitor.VisitModule(m)
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

func (mod *Module) SetDeclsHost() {
	for _, d := range mod.Decls {
		d.SetHost(mod)

		if td, ok := d.(*TypeDecl); ok {
			if cl, ok := td.Typ.(*ClassType); ok {
				for _, f := range cl.Fields {
					f.SetHost(mod)
				}
			}
		}
	}
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

func (m *Setting) Accept(visitor Visitor) TreePrinter {
	return visitor.VisitSetting(m)
}

//=== импорт

type Import struct {
	Pos  int
	Path string
	Mod  *Module
	// Собирается на слиянии модулей:
	Sources []int // список номеров исходных файлов, в которых есть импорт
}

func (m *Import) Accept(visitor Visitor) TreePrinter {
	return visitor.VisitImport(m)
}

func (n *Import) GetPos() int {
	return n.Pos
}

//=== вход

type EntryFn struct {
	Pos int
	Seq *StatementSeq
}

func (m *EntryFn) Accept(visitor Visitor) TreePrinter {
	return visitor.VisitEntryFn(m)
}

func (n *EntryFn) GetPos() int {
	return n.Pos
}
