package ast

import (
	"fmt"
)

var _ = fmt.Printf

//====

type TypeBase struct {
	Pos int
}

func (n *TypeBase) GetPos() int {
	return n.Pos
}
func (n *TypeBase) TypeNode() {}

//==== predefined types

type PredefinedType struct {
	TypeBase
	Name string
}

type InvalidType struct {
	TypeBase
}

//=== type ref

type TypeRef struct {
	TypeBase
	TypeName   string
	ModuleName string
	TypeDecl   *TypeDecl
	Typ        Type
}

//==== vector type

type VectorType struct {
	TypeBase
	ElementTyp Type
}

//==== class type

type ClassType struct {
	TypeBase
	BaseTyp Type
	Fields  []*Field        // поля самого класса
	Methods []*Function     // методы самого класса
	Members map[string]Decl // включая поля и методы базовых типов
}

type Field struct {
	DeclBase
	Typ Type
}

//==== function type

type FuncType struct {
	TypeBase
	Params    []*Param
	ReturnTyp Type
}

type Param struct {
	TypeBase
	Name string
	Typ  Type
}
