package compiler

import (
	"fmt"

	"trivil/ast"
	"trivil/env"
	"trivil/parser"
)

var _ = fmt.Printf

func (cc *compileContext) parseFile(src *env.Source) *ast.Module {

	if *env.TraceCompile {
		fmt.Printf("Синтаксис: '%s'\n", src.Path)
	}

	m := parser.Parse(src)

	if env.ErrorCount() != 0 {
		return m
	}

	for _, i := range m.Imports {
		cc.importModule(m, i)
	}

	return m
}

func (cc *compileContext) parseList(isMain bool, list []*env.Source) []*ast.Module {

	var mods = make([]*ast.Module, 0)
	var moduleName = ""

	for _, src := range list {

		var m = cc.parseFile(src)
		mods = append(mods, m)

		if len(mods) == 1 {
			moduleName = m.Name

			if !isMain && m.Name != src.FolderName {
				// не проверяю соответствие имени папки для головного модуля
				env.AddError(m.Pos, "ОКР-ОШ-ИМЯ-МОДУЛЯ", m.Name, src.FolderName)
			}
		} else if moduleName != m.Name {
			env.AddError(m.Pos, "ОКР-ОШ-МОДУЛИ-В-ПАПКЕ", moduleName, m.Name, src.FolderPath)
		}

		if m.Concrete != nil {
			mods = append(mods[:len(mods)-1], cc.concretize(m)...)
		}
	}

	return mods
}

func (cc *compileContext) parseModule(isMain bool, list []*env.Source) *ast.Module {

	var mods = cc.parseList(isMain, list)

	if env.ErrorCount() > 0 {
		return mods[0]
	}

	if env.ErrorCount() == 0 && len(mods) > 1 {
		mergeModules(mods)
	}

	if env.ErrorCount() == 0 && *env.ShowAST >= 1 {
		fmt.Println(ast.SExpr(mods[0]))
	}

	return mods[0]
}
