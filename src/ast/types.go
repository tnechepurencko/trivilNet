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
}

//==== function type

type FuncType struct {
	TypeBase
	Params    []*Param
	ReturnTyp Type
}

type Param struct {
	DeclBase
}

//==== variadic type

type VariadicType struct {
	TypeBase
	ElementTyp Type
}

//==== predicates

func IsIntegerType(t Type) bool {
	if tr, ok := t.(*TypeRef); ok {
		t = tr.Typ
	}

	return t == Int64 || t == Byte
}

func IsInt64(t Type) bool {
	if tr, ok := t.(*TypeRef); ok {
		t = tr.Typ
	}

	return t == Int64
}

func IsFloatType(t Type) bool {
	if tr, ok := t.(*TypeRef); ok {
		t = tr.Typ
	}

	return t == Float64
}

func IsBoolType(t Type) bool {
	if tr, ok := t.(*TypeRef); ok {
		t = tr.Typ
	}

	return t == Bool
}

func IsStringType(t Type) bool {
	if tr, ok := t.(*TypeRef); ok {
		t = tr.Typ
	}

	return t == String
}

func IsIndexableType(t Type) bool {
	if tr, ok := t.(*TypeRef); ok {
		t = tr.Typ
	}

	_, ok := t.(*VectorType)
	return ok
}

func ElementType(t Type) Type {
	if tr, ok := t.(*TypeRef); ok {
		t = tr.Typ
	}

	v, ok := t.(*VectorType)
	if !ok {
		panic("assert - должен быть индексируемый тип")
	}
	return v.ElementTyp
}

func IsVariadicType(t Type) bool {
	if tr, ok := t.(*TypeRef); ok {
		t = tr.Typ
	}

	_, ok := t.(*VariadicType)
	return ok
}

func IsClassType(t Type) bool {
	if tr, ok := t.(*TypeRef); ok {
		t = tr.Typ
	}

	_, ok := t.(*ClassType)

	return ok
}

func UnderType(t Type) Type {
	if tr, ok := t.(*TypeRef); ok {
		t = tr.Typ
	}

	return t
}

//==== invalid type

func IsInvalidType(t Type) bool {
	_, ok := UnderType(t).(*InvalidType)
	return ok
}

func MakeInvalidType(pos int) *InvalidType {
	return &InvalidType{TypeBase: TypeBase{Pos: pos}}
}

//==== for error messages

func TypeString(t Type) string {

	if tr, ok := t.(*TypeRef); ok {
		t = tr.Typ
	}

	switch x := t.(type) {
	case nil:
		return "*nil*"
	case *InvalidType:
		return "*invalid*"
	case *PredefinedType:
		return x.Name
	case *VectorType:
		return "[]" + TypeString(x.ElementTyp)
	default:
		return fmt.Sprintf("TypeString ni: %T", t)
	}
}

func TypeName(t Type) string {

	if tr, ok := t.(*TypeRef); ok {
		if tr.ModuleName != "" {
			return tr.ModuleName + "." + tr.TypeName
		} else {
			return tr.TypeName
		}
	} else {
		return TypeString(t)
	}
}
