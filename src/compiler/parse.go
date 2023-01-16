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

func (cc *compileContext) parseList(list []*env.Source) *ast.Module {

	var moduleName = ""
	var mods = make([]*ast.Module, len(list))
	for n, src := range list {

		var m = cc.parseFile(src)
		mods[n] = m

		if env.ErrorCount() == 0 {
			if n == 0 {
				moduleName = m.Name

				if cc.main != nil && m.Name != src.FolderName {
					// не проверяю соответствие имени папки для головного модуля
					env.AddError(m.Pos, "ОКР-ОШ-ИМЯ-МОДУЛЯ", m.Name, src.FolderName)
				}
			} else if moduleName != m.Name {
				env.AddError(m.Pos, "ОКР-ОШ-МОДУЛИ-В-ПАПКЕ", moduleName, m.Name, src.FolderPath)
			}
		}

	}

	if env.ErrorCount() == 0 && len(list) > 1 {
		mergeModules(mods)
	}

	if env.ErrorCount() == 0 && *env.ShowAST >= 1 {
		fmt.Println(ast.SExpr(mods[0]))
	}

	return mods[0]
}
