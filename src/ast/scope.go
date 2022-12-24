package ast

import (
	"fmt"
)

var _ = fmt.Printf

var (
	Byte    *PredefinedType
	Int     *PredefinedType
	Int32   *PredefinedType
	Int64   *PredefinedType
	Float64 *PredefinedType
	Bool    *PredefinedType
	Symbol  *PredefinedType
	String  *PredefinedType

	Void *PredefinedType // только для вызова функции без результата
)

type Scope struct {
	Outer *Scope
	Names map[string]Decl
}

var topScope *Scope

func initScopes() {
	topScope = &Scope{
		Names: make(map[string]Decl),
	}

	Byte = addType("Байт")
	//	Int32 = addType("Цел32")
	Int64 = addType("Цел64")
	Float64 = addType("Вещ64")
	Bool = addType("Лог")
	Symbol = addType("Символ")
	String = addType("Строка")

	addBoolConst("истина", true)
	addBoolConst("ложь", false)

	Void = &PredefinedType{Name: "нет результата"}
	//	ShowScopes("top", topScope)
}

func addType(name string) *PredefinedType {
	var pt = &PredefinedType{
		Name: name,
	}

	var td = &TypeDecl{}
	td.Typ = pt
	td.Name = name
	topScope.Names[name] = td

	return pt
}

func addBoolConst(name string, val bool) {
	var c = &ConstDecl{
		Value: &BoolLiteral{Value: val},
	}
	c.Typ = Bool
	c.Name = name

	topScope.Names[name] = c
}

func NewScope(outer *Scope) *Scope {
	return &Scope{
		Outer: outer,
		Names: make(map[string]Decl),
	}
}

func ShowScopes(label string, cur *Scope) {
	if label != "" {
		fmt.Println(label)
	}

	var n = 0
	for cur != nil {
		n++
		fmt.Printf("<--- scope %d\n", n)
		for _, d := range cur.Names {
			fmt.Printf("%s ", d.GetName())
		}
		if len(cur.Names) > 0 {
			fmt.Println()
		}
		cur = cur.Outer
	}
	fmt.Printf("--- end scopes\n")
}
