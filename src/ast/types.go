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
	Init       Expr
	Later      bool
	AssignOnce bool
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

func VariadicParam(ft *FuncType) *Param {
	if len(ft.Params) == 0 {
		return nil
	}
	var last = ft.Params[len(ft.Params)-1]
	if IsVariadicType(last.Typ) {
		return last
	}
	return nil
}

//==== variadic type

type VariadicType struct {
	TypeBase
	ElementTyp Type
}

//==== мб тип

type MayBeType struct {
	TypeBase
	Typ Type
}

//====

// Снимает все TypeRef, может быть два в контексте тип А = Б
func UnderType(t Type) Type {
	for {
		if tr, ok := t.(*TypeRef); ok {
			t = tr.Typ
		} else {
			return t
		}
	}
}

//==== predicates

func IsIntegerType(t Type) bool {
	t = UnderType(t)
	return t == Int64 || t == Byte || t == Word64
}

func IsInt64(t Type) bool {
	return UnderType(t) == Int64
}

func IsFloatType(t Type) bool {
	return UnderType(t) == Float64
}

func IsWord64(t Type) bool {
	return UnderType(t) == Word64
}

func IsBoolType(t Type) bool {
	return UnderType(t) == Bool
}

func IsStringType(t Type) bool {
	return UnderType(t) == String
}

func IsVoidType(t Type) bool {
	return UnderType(t) == VoidType
}

func IsIndexableType(t Type) bool {
	t = UnderType(t)

	switch t.(type) {
	case *VectorType, *VariadicType:
		return true
	default:
		return t == String8
	}
}

func ElementType(t Type) Type {
	t = UnderType(t)

	switch x := t.(type) {
	case *VectorType:
		return x.ElementTyp
	case *VariadicType:
		return x.ElementTyp
	default:
		if t == String8 {
			return Byte
		}
		panic("assert - должен быть индексируемый тип")
	}
}

func IsVectorType(t Type) bool {
	_, ok := UnderType(t).(*VectorType)
	return ok
}

func IsVariadicType(t Type) bool {
	_, ok := UnderType(t).(*VariadicType)
	return ok
}

func IsClassType(t Type) bool {
	_, ok := UnderType(t).(*ClassType)

	return ok
}

func IsMayBeType(t Type) bool {
	_, ok := UnderType(t).(*MayBeType)

	return ok
}

func IsTagPairType(t Type) bool {
	t = UnderType(t)
	return t == TagPairType
}

func IsReferenceType(t Type) bool {
	t = UnderType(t)
	switch t.(type) {
	case *VectorType, *ClassType:
		return true
	case *InvalidType:
		return true
	default:
		return t == String
	}
}

//==== tags

// Объекты каких типов имеют тег
func HasTag(t Type) bool {
	return !IsVariadicType(t) //пока так, можно разрешить для всех, но потом

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

	t = UnderType(t)

	switch x := t.(type) {
	case nil:
		return "*nil*"
	case *InvalidType:
		return "*invalid*"
	case *PredefinedType:
		return x.Name
	case *VectorType:
		return "[]" + TypeName(x.ElementTyp)
	case *VariadicType:
		return "..." + TypeName(x.ElementTyp)
	case *MayBeType:
		return "мб " + TypeName(x.Typ)
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
