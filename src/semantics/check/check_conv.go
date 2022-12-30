package check

import (
	"fmt"
	"trivil/ast"
	//	"trivil/env"
	//	"trivil/lexer"
)

var _ = fmt.Printf

/*
  по целевому типу:
	Байт: Цел64, Символ (0..255)
	Цел64: Байт, Вещ64, Символ
	Вещ64: Цел64
	Лог: -
	Символ: Цел64, Строка(из 1-го символа)
	Строка: Символ, []Символ, []Байт
	[]Байт: Строка
	[]Символ: Строка
	Класс: Класс (вверх или вниз)
*/
func (cc *checkContext) conversion(x *ast.ConversionExpr) {
}
