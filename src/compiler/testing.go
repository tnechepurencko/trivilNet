package compiler

import (
	"fmt"

	//	"trivil/ast"
	"trivil/env"
	/*
		"trivil/genc"
		"trivil/semantics"
	*/)

var _ = fmt.Printf

const testsFolderName = "_тест_"

func TestOne(tpath string) {
	fmt.Printf("тест %v\n", tpath)

	env.Normalizer.Process(tpath)
	if env.Normalizer.Err != nil {
		env.FatalError("ОКР-ОШ-ПУТЬ-ТЕСТ", tpath, env.Normalizer.Err.Error())
		return
	}

	var npath = env.Normalizer.NPath

	var err = env.EnsureFolder(npath)
	if err != nil {
		env.FatalError("ОКР-ТЕСТ-НЕ-ПАПКА", tpath, err.Error())
		return
	}

	var list = env.GetFolderSources(tpath, npath)

	if len(list) == 0 {
		env.FatalError("ОКР-ТЕСТ-ПУСТАЯ-ПАПКА", tpath, npath)
		return
	}

	if len(list) == 1 && list[0].Err != nil {
		env.FatalError("ОКР-ОШ-ЧТЕНИЕ-ИСХОДНОГО", list[0].FilePath, list[0].Err.Error())
		return
	}

	/*

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
