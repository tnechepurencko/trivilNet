package ast

import (
	"fmt"
)

var _ = fmt.Printf

var (
	Byte *PredefinedType
	//Int32   *PredefinedType
	Int64   *PredefinedType
	Float64 *PredefinedType
	Bool    *PredefinedType
	Symbol  *PredefinedType
	String  *PredefinedType
	Word64  *PredefinedType

	Void    *PredefinedType // только для вызова функции без результата
	TagPair *PredefinedType // только в типе параметров и в типе элемента variadic
)

type Scope struct {
	Outer *Scope
	Names map[string]Decl
}

const (
	StdLen       = "длина"
	StdTag       = "тег"
	StdSomething = "нечто"
)

var topScope *Scope

func initScopes() {
	topScope = &Scope{
		Names: make(map[string]Decl),
	}

	Byte = addType("Байт")
	//	Int32 = addType("Цел32")
	Int64 = addType("Цел64")
	Float64 = addType("Вещ64")
	Word64 = addType("Слово64")
	Bool = addType("Лог")
	Symbol = addType("Символ")
	String = addType("Строка")

	addBoolConst("истина", true)
	addBoolConst("ложь", false)

	Void = &PredefinedType{Name: "нет результата"}
	TagPair = &PredefinedType{Name: "ТегСлово"}

	addStdFunction(StdLen)

	addStdFunction(StdTag)
	addStdFunction(StdSomething)

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

func addStdFunction(name string) {
	var f = &StdFunction{}
	f.Typ = Void
	f.Name = name

	topScope.Names[name] = f
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
