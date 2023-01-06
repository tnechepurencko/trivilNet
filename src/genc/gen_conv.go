package genc

import (
	"fmt"

	"trivil/ast"
)

var _ = fmt.Printf

func (genc *genContext) genConversion(x *ast.ConversionExpr) string {

	var e = genc.genExpr(x.X)
	if x.Done {
		return e
	}

	var target = ast.UnderType(x.TargetTyp)
	/*

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
		case ast.String:
			cc.conversionToString(x)
			return

		}
	*/
	switch /*xt :=*/ target.(type) {
	/*
		case *ast.VectorType:
			cc.conversionToVector(x, xt)
		case *ast.ClassType:
			cc.conversionToClass(x, xt)
	*/
	default:
		panic(fmt.Sprintf("ni %T '%s'", target, ast.TypeString(target)))
	}
}

/*
func (cc *checkContext) conversionToByte(x *ast.ConversionExpr) {

	var t = ast.UnderType(x.X.GetType())

	switch t {
	case ast.Byte:
		env.AddError(x.Pos, "СЕМ-ПРИВЕДЕНИЕ-ТИПА-К-СЕБЕ", ast.TypeString(x.X.GetType()))
		x.Typ = ast.Byte
		return

	case ast.Int64,
		ast.Symbol:
		var li = literal(x.X)
		if li != nil {
			i, err := strconv.ParseInt(li.Lit, 0, 64)
			if err != nil || i < 0 || i > 255 {
				env.AddError(x.Pos, "СЕМ-ЗНАЧЕНИЕ-НЕ-В_ДИАПАЗОНЕ", ast.Byte.Name)
			} else {
				x.Done = true
				li.Typ = ast.Byte
			}
		}
		x.Typ = ast.Byte
		return
	case ast.String:
		var li = literal(x.X)
		if li != nil {
			if utf8.RuneCountInString(li.Lit) == 1 {
				r, _ := utf8.DecodeRuneInString(li.Lit)
				if r < 0 || r > 255 {
					env.AddError(x.Pos, "СЕМ-ЗНАЧЕНИЕ-НЕ-В_ДИАПАЗОНЕ", ast.Byte.Name)
				} else {
					x.Done = true
					li.Typ = ast.Byte
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
	case ast.Byte,
		ast.Symbol:
		var li = literal(x.X)
		if li != nil {
			li.Typ = ast.Byte
			x.Done = true
		}
		x.Typ = ast.Int64
		return
	case ast.Float64:
		// пока не работаю с литералами
		x.Typ = ast.Int64
		return
	case ast.String:
		var li = oneSymbolString(x.X)
		if li != nil {
			r, _ := utf8.DecodeRuneInString(li.Lit)
			li.Kind = lexer.INT
			li.Lit = fmt.Sprintf("0x%x", r)
			x.Typ = ast.Int64
			x.Done = true
			return
		}
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
		var li = literal(x.X)
		if li != nil {
			i, err := strconv.ParseInt(li.Lit, 0, 64)
			if err != nil || i < 0 || i > unicode.MaxRune {
				env.AddError(x.Pos, "СЕМ-ЗНАЧЕНИЕ-НЕ-В_ДИАПАЗОНЕ", ast.Symbol.Name)
			} else {
				x.Done = true
				li.Typ = ast.Symbol
			}
		}
		x.Typ = ast.Symbol
		return
	case ast.String:
		var lit = oneSymbolString(x.X)
		if lit != nil {
			x.Typ = ast.Symbol
			x.Done = true
			return
		}
	}

	env.AddError(x.Pos, "СЕМ-ОШ-ПРИВЕДЕНИЯ-ТИПА", ast.TypeString(x.X.GetType()), ast.Symbol.Name)
	x.Typ = invalidType(x.Pos)

}

func (cc *checkContext) conversionToString(x *ast.ConversionExpr) {

	var t = ast.UnderType(x.X.GetType())

	switch t {
	case ast.String:
		env.AddError(x.Pos, "СЕМ-ПРИВЕДЕНИЕ-ТИПА-К-СЕБЕ", ast.TypeString(x.X.GetType()))
		x.Typ = ast.String
		return
	case ast.Symbol:
		var li = literal(x.X)
		if li != nil {
			li.Typ = ast.String
			x.Done = true
		}
		x.Typ = ast.String
		return
	}

	vt, ok := t.(*ast.VectorType)
	if ok {

		var et = ast.UnderType(vt.ElementTyp)

		if et == ast.Byte || et == ast.Symbol {
			x.Typ = ast.String
			return
		}
	}

	env.AddError(x.Pos, "СЕМ-ОШ-ПРИВЕДЕНИЯ-ТИПА", ast.TypeString(x.X.GetType()), ast.String.Name)
	x.Typ = invalidType(x.Pos)
}

func (cc *checkContext) conversionToVector(x *ast.ConversionExpr, target *ast.VectorType) {

	var t = ast.UnderType(x.X.GetType())

	if t == ast.String {

		var et = ast.UnderType(target.ElementTyp)

		if et == ast.Byte || et == ast.Symbol {
			x.Typ = x.TargetTyp
			return
		}
	}
	env.AddError(x.Pos, "СЕМ-ОШ-ПРИВЕДЕНИЯ-ТИПА",
		ast.TypeString(x.X.GetType()), ast.TypeString(x.TargetTyp))
	x.Typ = invalidType(x.Pos)

}

func (cc *checkContext) conversionToClass(x *ast.ConversionExpr, target *ast.ClassType) {

	var t = ast.UnderType(x.X.GetType())

	if t == target {
		env.AddError(x.Pos, "СЕМ-ПРИВЕДЕНИЕ-ТИПА-К-СЕБЕ", ast.TypeString(target))
		x.Typ = x.TargetTyp
		return
	}

	tClass, ok := t.(*ast.ClassType)
	if ok {
		if !isDerivedClass(tClass, target) {
			env.AddError(x.Pos, "СЕМ-ДОЛЖЕН-БЫТЬ-НАСЛЕДНИКОМ", ast.TypeName(x.X.GetType()), ast.TypeName(x.TargetTyp))
		}
		return
	}

	env.AddError(x.Pos, "СЕМ-ОШ-ПРИВЕДЕНИЯ-ТИПА",
		ast.TypeName(x.X.GetType()), ast.TypeName(x.TargetTyp))
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
	var li = literal(expr)
	if li == nil {
		return nil
	}
	if li.Kind != lexer.STRING {
		return nil
	}

	if utf8.RuneCountInString(li.Lit) != 1 {
		return nil
	}
	return li
}

func isDerivedClass(base, derived *ast.ClassType) bool {

	var c = derived

	for c.BaseTyp != nil {
		var t = ast.UnderType(c.BaseTyp)
		if t == base {
			return true
		}
		c = t.(*ast.ClassType)
	}
	return false
}
*/
