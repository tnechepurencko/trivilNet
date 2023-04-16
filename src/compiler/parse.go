package compiler

import (
	"fmt"
	"strings"

	"trivil/ast"
	"trivil/env"
	"trivil/parser"
)

var _ = fmt.Printf

func (cc *compileContext) parseFile(src *env.Source) *ast.Module {

	if *env.TraceCompile {
		fmt.Printf("Синтаксис: '%s'\n", src.FilePath)
	}

	var m = parser.Parse(src)

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

			if !isMain && m.Name != src.FolderName() {
				// не проверяю соответствие имени папки для головного модуля
				env.AddError(m.Pos, "ОКР-ОШ-ИМЯ-МОДУЛЯ", m.Name, src.FolderName())
			}
		} else if moduleName != m.Name {
			env.AddError(m.Pos, "ОКР-ОШ-МОДУЛИ-В-ПАПКЕ", moduleName, m.Name, src.FolderPath)
		}

		if m.Setting != nil {
			mods = append(mods, cc.setup(m)...)
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

	var m = mods[0]

	if env.ErrorCount() != 0 {
		return m
	}

	for _, i := range m.Imports {
		cc.importModule(m, i)
	}

	return m
}

func mergeModules(mods []*ast.Module) {

	if *env.TraceCompile {
		var list = make([]string, len(mods))
		for i, m := range mods {
			source, _, _ := env.SourcePos(m.Pos)
			list[i] = source.FileName
		}
		fmt.Printf("Слияние: %s\n", strings.Join(list, " + "))
	}

	var combined = mods[0]

	// соединить импорт
	var allImport = make(map[string]struct{}, len(combined.Imports))
	for _, i := range combined.Imports {
		allImport[i.Path] = struct{}{}
	}

	for n := 1; n < len(mods); n++ {
		m := mods[n]
		for _, i := range m.Imports {
			_, ok := allImport[i.Path]
			if !ok {
				allImport[i.Path] = struct{}{}
				combined.Imports = append(combined.Imports, i)
			}
		}
	}

	// соединить описания
	for n := 1; n < len(mods); n++ {

		var m = mods[n]

		setHost(combined, m.Decls)
		combined.Decls = append(combined.Decls, m.Decls...)

		if m.Entry != nil {
			if combined.Entry != nil {
				env.AddError(combined.Entry.Pos, "ПАР-ДУБЛЬ-ВХОД", env.PosString(m.Entry.Pos))
			} else {
				combined.Entry = m.Entry
			}
		}
	}
}

func setHost(combined *ast.Module, decls []ast.Decl) {

	for _, d := range decls {
		d.SetHost(combined)

		if td, ok := d.(*ast.TypeDecl); ok {
			if cl, ok := td.Typ.(*ast.ClassType); ok {
				for _, f := range cl.Fields {
					f.SetHost(combined)
				}
			}
		}

	}

}
