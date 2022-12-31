package check

import (
	"fmt"
	"strconv"
	"unicode"
	"unicode/utf8"

	"trivil/ast"
	"trivil/env"
	"trivil/lexer"
)

var _ = fmt.Printf

/*
  по целевому типу:
	Байт: Цел64, Символ (0..255), Строковый литерал (из 1-го символа)
	Цел64: Байт, Вещ64, Символ, Строковый литерал (из 1-го символа)
	Вещ64: Цел64
	Лог: -
	Символ: Цел64, Строковый литерал
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
	case ast.Float64:
		cc.conversionToFloat64(x)
		return
	case ast.Bool:
		env.AddError(x.Pos, "СЕМ-ОШ-ПРИВЕДЕНИЯ-ТИПА", ast.TypeString(x.X.GetType()), ast.Bool.Name)
		x.Typ = invalidType(x.Pos)
		return
	case ast.Symbol:
		cc.conversionToSymbol(x)
		return

	}

	switch target.(type) {

	default:
		panic(fmt.Sprintf("ni %T '%s'", target, ast.TypeString(target)))
	}
}

func (cc *checkContext) conversionToByte(x *ast.ConversionExpr) {

	var t = ast.UnderType(x.X.GetType())

	switch t {
	case ast.Byte:
		env.AddError(x.Pos, "СЕМ-ПРИВЕДЕНИЕ-ТИПА-К-СЕБЕ", ast.TypeString(x.X.GetType()))
		x.Typ = ast.Byte
		return

	case ast.Int64, ast.Symbol:
		var lit = literal(x.X)
		if lit != nil {
			i, err := strconv.ParseInt(lit.Lit, 0, 64)
			if err != nil || i < 0 || i > 255 {
				env.AddError(x.Pos, "СЕМ-ЗНАЧЕНИЕ-НЕ-В_ДИАПАЗОНЕ", ast.Byte.Name)
			} else {
				x.Done = true
				lit.Typ = ast.Byte
			}
		}
		x.Typ = ast.Byte
		return
	case ast.String:
		var lit = literal(x.X)
		if lit != nil {
			if utf8.RuneCountInString(lit.Lit) == 1 {
				r, _ := utf8.DecodeRuneInString(lit.Lit)
				if r < 0 || r > 255 {
					env.AddError(x.Pos, "СЕМ-ЗНАЧЕНИЕ-НЕ-В_ДИАПАЗОНЕ", ast.Byte.Name)
				} else {
					x.Done = true
					lit.Typ = ast.Byte
				}

			} else {
				env.AddError(x.Pos, "СЕМ-ДЛИНА-СТРОКИ-НЕ-1")
			}

		}
		x.Typ = ast.Byte
		return
	}

	env.AddError(x.Pos, "СЕМ-ОШ-ПРИВЕДЕНИЯ-ТИПА", ast.TypeString(x.X.GetType()), ast.Byte.Name)
	x.Typ = invalidType(x.Pos)

}

func (cc *checkContext) conversionToInt64(x *ast.ConversionExpr) {

	var t = ast.UnderType(x.X.GetType())

	switch t {
	case ast.Int64:
		env.AddError(x.Pos, "СЕМ-ПРИВЕДЕНИЕ-ТИПА-К-СЕБЕ", ast.TypeString(x.X.GetType()))
		x.Typ = ast.Int64
		return
	case ast.Byte, ast.Symbol:
		var lit = literal(x.X)
		if lit != nil {
			lit.Typ = ast.Byte
			x.Done = true
		}
		x.Typ = ast.Int64
		return
	case ast.Float64:
		// пока не работаю с литералами
		x.Typ = ast.Int64
		return
	case ast.String:
		panic("ni")
	}

	env.AddError(x.Pos, "СЕМ-ОШ-ПРИВЕДЕНИЯ-ТИПА", ast.TypeString(x.X.GetType()), ast.Int64.Name)
	x.Typ = invalidType(x.Pos)
}

func (cc *checkContext) conversionToFloat64(x *ast.ConversionExpr) {

	var t = ast.UnderType(x.X.GetType())

	switch t {
	case ast.Float64:
		env.AddError(x.Pos, "СЕМ-ПРИВЕДЕНИЕ-ТИПА-К-СЕБЕ", ast.TypeString(x.X.GetType()))
		x.Typ = ast.Float64
		return
	case ast.Int64:
		// пока не работаю с литералами
		x.Typ = ast.Float64
		return
	}

	env.AddError(x.Pos, "СЕМ-ОШ-ПРИВЕДЕНИЯ-ТИПА", ast.TypeString(x.X.GetType()), ast.Float64.Name)
	x.Typ = invalidType(x.Pos)

}

func (cc *checkContext) conversionToSymbol(x *ast.ConversionExpr) {

	var t = ast.UnderType(x.X.GetType())

	switch t {
	case ast.Symbol:
		env.AddError(x.Pos, "СЕМ-ПРИВЕДЕНИЕ-ТИПА-К-СЕБЕ", ast.TypeString(x.X.GetType()))
		x.Typ = ast.Symbol
		return

	case ast.Int64:
		var lit = literal(x.X)
		if lit != nil {
			i, err := strconv.ParseInt(lit.Lit, 0, 64)
			if err != nil || i < 0 || i > unicode.MaxRune {
				env.AddError(x.Pos, "СЕМ-ЗНАЧЕНИЕ-НЕ-В_ДИАПАЗОНЕ", ast.Symbol.Name)
			} else {
				x.Done = true
				lit.Typ = ast.Symbol
			}
		}
		x.Typ = ast.Symbol
		return
	case ast.String:
		panic("ni")
	}

	env.AddError(x.Pos, "СЕМ-ОШ-ПРИВЕДЕНИЯ-ТИПА", ast.TypeString(x.X.GetType()), ast.Symbol.Name)
	x.Typ = invalidType(x.Pos)

}

//====

func literal(expr ast.Expr) *ast.LiteralExpr {

	switch x := expr.(type) {
	case *ast.LiteralExpr:
		return x
	case *ast.ConversionExpr:
		if x.Done {
			return literal(x.X)
		}
	}
	return nil
}

func oneSymbolString(expr ast.Expr) *ast.LiteralExpr {
	var lit = literal(expr)
	if lit == nil {
		return nil
	}
	if lit.Kind != lexer.STRING {
		return nil
	}

	if utf8.RuneCountInString(lit.Lit) != 1 {
		return nil
	}
	return lit
}
