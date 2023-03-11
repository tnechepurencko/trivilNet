package compiler

import (
	"fmt"

	"trivil/ast"
	"trivil/env"
)

var _ = fmt.Printf

func (cc *compileContext) importModule(m *ast.Module, i *ast.Import) {

	var npath = env.NormalizeFolderPath(i.Path)
	m, ok := cc.imported[npath]
	if ok {
		// Модуль уже был импортирован
		i.Mod = m
		//fmt.Printf("already imported %s\n", i.Path)
		return
	}

	var err = env.CheckFolder(i.Path)
	if err != nil {
		env.AddError(i.Pos, "ОКР-ИМПОРТ-НЕ-ПАПКА", i.Path, err.Error())
		return
	}

	var list = env.GetFolderSources(i.Path)

	if len(list) == 0 {
		env.AddError(i.Pos, "ОКР-ИМПОРТ-ПУСТАЯ-ПАПКА", i.Path)
		return
	}

	if len(list) == 1 && list[0].Err != nil {
		env.AddError(i.Pos, "ОКР-ОШ-ЧТЕНИЕ-ИСХОДНОГО", list[0].Path, list[0].Err.Error())
		return
	}

	i.Mod = cc.parseModule(false, list)
	cc.imported[npath] = i.Mod
}

func mergeModules(mods []*ast.Module) {
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
		m := mods[n]

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
