package compiler

import (
	"fmt"

	"trivil/ast"
	"trivil/env"
	"trivil/genc"
	"trivil/parser"
	"trivil/semantics"
)

type compileContext struct {
	main    *ast.Module
	modules map[string]*ast.Module

	list   []*ast.Module
	status map[*ast.Module]int
}

func Compile(src *env.Source) {

	var cc = &compileContext{
		modules: make(map[string]*ast.Module),
	}

	cc.main = cc.parse(src)
	cc.modules[src.Original] = cc.main

	if env.ErrorCount() != 0 {
		return
	}

	cc.orderedList()

	/*
		for i, m := range cc.list {
			fmt.Printf("%d: %v %p\n", i, m.Name, m)
		}
	*/

	for _, m := range cc.list {

		if env.ErrorCount() != 0 {
			break
		}

		if *env.TraceCompile {
			fmt.Printf("-->анализ и генерация модуля '%s'\n", m.Name)
		}

		cc.process(m)

		if *env.TraceCompile {
			fmt.Printf("<-- '%s'\n", m.Name)
		}
	}

	if env.ErrorCount() == 0 && *env.DoGen && *env.BuildExe {
		genc.BuildExe(cc.list)
	}
}

func (cc *compileContext) parse(src *env.Source) *ast.Module {

	m, ok := cc.modules[src.Path]
	if ok {
		return m
	}

	m = parser.Parse(src)
	cc.modules[src.Path] = m

	if env.ErrorCount() != 0 {
		return m
	}

	if *env.ShowAST >= 1 {
		fmt.Println(ast.SExpr(m))
	}

	for _, i := range m.Imports {
		cc.importModule(m, i)
	}

	return m
}

func (cc *compileContext) importModule(m *ast.Module, i *ast.Import) {

	//check already imported

	src := env.AddSource(i.Path)
	if src.Err != nil {
		env.AddError(i.Pos, "ОКР-ОШ-ЧТЕНИЕ-ИСХОДНОГО", src.Path, src.Err.Error())
		return
	}

	i.Mod = cc.parse(src)

	if i.Mod.Name != src.LastName {
		env.AddError(i.Pos, "ОКР-ОШ-ИМЯ-МОДУЛЯ", i.Mod.Name, src.LastName)
	}
}

func (cc *compileContext) process(m *ast.Module) {
	semantics.Analyse(m)

	if env.ErrorCount() != 0 {
		return
	}

	if *env.ShowAST >= 2 {
		fmt.Println(ast.SExpr(m))
	}

	if *env.DoGen {
		genc.Generate(m, m == cc.main)
	}
}

func (cc *compileContext) orderedList() {

	cc.list = make([]*ast.Module, 0)
	cc.status = make(map[*ast.Module]int)

	cc.traverse(cc.main, cc.main.Pos)
}

const (
	processing = 1
	processed  = 2
)

func (cc *compileContext) traverse(m *ast.Module, pos int) {

	s, ok := cc.status[m]
	if ok {
		if s == processing {
			env.AddError(pos, "СЕМ-ЦИКЛ-ИМПОРТА", m.Name)
			cc.status[m] = processed
		}
		return
	}

	cc.status[m] = processing

	for _, i := range m.Imports {
		cc.traverse(i.Mod, i.Pos)
	}

	cc.status[m] = processed
	cc.list = append(cc.list, m)
}
