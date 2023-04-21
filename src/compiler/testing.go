package compiler

import (
	"fmt"
	/*
		"trivil/ast"
		"trivil/env"
		"trivil/genc"
		"trivil/semantics"
	*/)

var _ = fmt.Printf

func TestOne(spath string) {
	fmt.Printf("тест %v\n", spath)

	/*
		var list = env.GetSources(spath)
		var src = list[0]
		if src.Err != nil {
			env.FatalError("ОКР-ОШ-ЧТЕНИЕ-ИСХОДНОГО", spath, src.Err.Error())
			return
		}

		var cc = &compileContext{
			imported: make(map[string]*ast.Module),
		}

		cc.main = cc.parseModule(true, list)

		if env.ErrorCount() != 0 {
			return
		}

		cc.orderedList()

		for _, m := range cc.list {

			if env.ErrorCount() != 0 {
				break
			}

			if *env.TraceCompile {
				fmt.Printf("Анализ и генерация: '%s'\n", m.Name)
				//fmt.Printf("Анализ и генерация: '%s' %p\n", m.Name, m)
			}

			cc.process(m)
		}

		if env.ErrorCount() == 0 && *env.DoGen && *env.BuildExe {
			genc.BuildExe(cc.list)
		}
	*/
}
