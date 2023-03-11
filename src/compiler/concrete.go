package compiler

import (
	"fmt"

	"trivil/ast"
	"trivil/env"
)

var _ = fmt.Printf

func (cc *compileContext) concretize(spec *ast.Module) []*ast.Module {

	var con = spec.Concrete

	var err = env.CheckFolder(con.Path)
	if err != nil {
		env.AddError(con.Pos, "ОКР-ИМПОРТ-НЕ-ПАПКА", con.Path, err.Error())
		return nil
	}

	var list = env.GetFolderSources(con.Path)

	if len(list) == 0 {
		// TODO: изменить ошибку
		env.AddError(con.Pos, "ОКР-ИМПОРТ-ПУСТАЯ-ПАПКА", con.Path)
		return nil
	}

	if len(list) == 1 && list[0].Err != nil {
		env.AddError(con.Pos, "ОКР-ОШ-ЧТЕНИЕ-ИСХОДНОГО", list[0].Path, list[0].Err.Error())
		return nil
	}

	var src = list[0]

	for key, val := range con.Attrs {
		var s = fmt.Sprintf("тип %s = %s\n", key, val)
		src.Bytes = append(src.Bytes, []byte(s)...)
	}

	var mods = cc.parseList(false, list)

	for _, m := range mods {
		m.Name = spec.Name
	}

	return mods
}
