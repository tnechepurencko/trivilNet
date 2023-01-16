package compiler

import (
	"fmt"

	"trivil/ast"
	"trivil/env"
	"trivil/genc"
	"trivil/parser"
	"trivil/semantics"
)

var _ = fmt.Printf

type compileContext struct {
	main *ast.Module
	//modules map[string]*ast.Module
	imported map[string]*ast.Module // map[folder]

	// упорядоченный список для обработки
	// головной модуль - в конце
	list   []*ast.Module
	status map[*ast.Module]int
}

func Compile(spath string) {

	var files = env.AddSource(spath)
	var src = files[0]
	if src.Err != nil {
		env.FatalError("ОКР-ОШ-ЧТЕНИЕ-ИСХОДНОГО", spath, src.Err.Error())
		return
	}

	var cc = &compileContext{
		imported: make(map[string]*ast.Module),
	}

	cc.main = cc.parse(src)
	//TODO: добавить главный в imported[folder]
	//cc.modules[src.Original] = cc.main

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

	m := parser.Parse(src)
	//	cc.modules[src.Path] = m

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

//=== process

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

//=== traverse

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
