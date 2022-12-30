package check

import (
	"fmt"
	"trivil/ast"
	"trivil/env"
	//	"trivil/lexer"
)

var _ = fmt.Printf

/*
  по целевому типу:
	Байт: Цел64, Символ (0..255), Строковый литерал (из 1-го символа)
	Цел64: Байт, Вещ64, Символ, Строковый литерал (из 1-го символа)
	Вещ64: Цел64
	Лог: -
	Символ: Цел64, Строка
	Строка: Символ, []Символ, []Байт
	[]Байт: Строка
	[]Символ: Строка
	Класс: Класс (вверх или вниз)
*/
func (cc *checkContext) conversion(x *ast.ConversionExpr) {

	cc.expr(x.X)

	var target = ast.UnderType(x.TargetTyp)

	switch target {
	case ast.Byte:
		cc.conversionToByte(x)
		return
	case ast.Int64:
		cc.conversionToInt64(x)
		return

	case ast.Bool:
		env.AddError(x.Pos, "СЕМ-ОШ-ПРИВЕДЕНИЯ-ТИПА", ast.TypeString(x.X.GetType()), ast.Bool.Name)
		x.Typ = invalidType(x.Pos)
		return

	}

	switch target.(type) {

	default:
		panic(fmt.Sprintf("ni %T '%s'", target, ast.TypeString(target)))
	}
}

func (cc *checkContext) conversionToByte(x *ast.ConversionExpr) {
	panic("ni")
}

func (cc *checkContext) conversionToInt64(x *ast.ConversionExpr) {

	var t = ast.UnderType(x.X.GetType())

	if t == ast.Int64 {
		env.AddError(x.Pos, "СЕМ-ПРИВЕДЕНИЕ-ТИПА-К-СЕБЕ", ast.TypeString(x.X.GetType()))
		x.Typ = ast.Int64
		return
	}

	if x.Typ == nil {
		env.AddError(x.Pos, "СЕМ-ОШ-ПРИВЕДЕНИЯ-ТИПА", ast.TypeString(x.X.GetType()), ast.Int64.Name)
		x.Typ = invalidType(x.Pos)
	}
}
