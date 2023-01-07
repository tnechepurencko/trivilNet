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
	//main *ast.Module
}

func compile(src *env.Source) {

	var cc = &compileContext{}

	var m = cc.parse(src)

	if env.ErrorCount() != 0 {
		return
	}

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

	if *env.ShowAST >= 1 {
		fmt.Println(ast.SExpr(m))
	}

	for _, i := range m.Imports {
		cc.importModule(m, i)
	}

	return m
}

func (cc *compileContext) importModule(m *ast.Module, i *ast.Import) {

	src := env.AddSource(i.Path)
	if src.Err != nil {
		fmt.Printf("Ошибка чтения исходного файла '%s': %s\n", i.Path, src.Err.Error())
		//TODO: error
		return
	}

	i.Mod = cc.parse(src)
	//TODO добавить в область видимости
}
