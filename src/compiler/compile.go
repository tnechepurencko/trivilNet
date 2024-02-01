package compiler

import (
	"fmt"

	"trivil/ast"
	"trivil/env"
	"trivil/genc"
	"trivil/semantics"
)

var _ = fmt.Printf

type CompileContext struct {
	Main     *ast.Module
	imported map[string]*ast.Module // map[folder]

	testModulePath string // импорт путь для тестируемого модуля

	// упорядоченный список для обработки
	// головной модуль - в конце
	list   []*ast.Module
	status map[*ast.Module]int

	// Путь к папке для модуля, только для создания интерфейса модуля
	folders map[*ast.Module]string
}

func Compile(spath string) {

	var list = env.GetSources(spath)
	var src = list[0]
	if src.Err != nil {
		env.FatalError("ОКР-ОШ-ЧТЕНИЕ-ИСХОДНОГО", spath, src.Err.Error())
		return
	}

	var cc = &CompileContext{
		imported: make(map[string]*ast.Module),
		folders:  make(map[*ast.Module]string),
	}

	cc.Main = cc.parseModule(true, list)

	if env.ErrorCount() != 0 {
		return
	}

	cc.build()
}

func (cc *CompileContext) build() {
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
}

//=== process

func (cc *CompileContext) process(m *ast.Module) {

	ast.CurHost = m
	semantics.Analyse(m)
	ast.CurHost = nil

	if env.ErrorCount() != 0 {
		return
	}

	if *env.ShowAST >= 2 {
		fmt.Println(ast.SExpr(m))
	}

	if *env.MakeDef && m != cc.Main {
		makeDef(m, cc.folders[m])
	}

	if *env.DoGen {
		genc.Generate(m, m == cc.Main)
	}
}

func (cc *CompileContext) orderedList() {

	cc.list = make([]*ast.Module, 0)
	cc.status = make(map[*ast.Module]int)

	cc.traverse(cc.Main, cc.Main.Pos)
}

const (
	processing = 1
	processed  = 2
)

//=== traverse

func (cc *CompileContext) traverse(m *ast.Module, pos int) {

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
