package main

import (
	"fmt"

	"trivil/ast"
	"trivil/env"
	"trivil/genc"
	"trivil/parser"
	"trivil/semantics"
)

type compileContext struct {
	modules []*ast.Module
}

func compile(src *env.Source) {

	var cc = &compileContext{
		modules: make([]*ast.Module, 0),
	}

	//var main =
	cc.parse(src)

	if env.ErrorCount() != 0 {
		return
	}

	//TODO: reorder and check cycles

	for i := len(cc.modules) - 1; i >= 0; i-- {

		if env.ErrorCount() != 0 {
			break
		}

		var m = cc.modules[i]
		fmt.Printf("> %s\n", m.Name)

		cc.process(m)

		fmt.Printf("< %s\n", m.Name)
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
		genc.Generate(m)
	}
}

func (cc *compileContext) parse(src *env.Source) *ast.Module {
	var m = parser.Parse(src)
	if env.ErrorCount() != 0 {
		return m
	}

	cc.modules = append(cc.modules, m)

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
