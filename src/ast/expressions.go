package ast

import (
	"fmt"
	"trivil/lexer"
)

var _ = fmt.Printf

//====

type ExprBase struct {
	Pos      int
	Typ      Type
	ReadOnly bool
}

func (m *ExprBase) Accept(visitor Visitor) {
	visitor.VisitExprBase(m)
}

func (n *ExprBase) ExprNode() {}

func (n *ExprBase) GetPos() int {
	return n.Pos
}

func (n *ExprBase) GetType() Type {
	return n.Typ
}

func (n *ExprBase) SetType(t Type) {
	n.Typ = t
}

func (n *ExprBase) IsReadOnly() bool {
	return n.ReadOnly
}

//====

type InvalidExpr struct {
	ExprBase
}

func (m *InvalidExpr) Accept(visitor Visitor) {
	visitor.VisitInvalidExpr(m)
}

//====

type BinaryExpr struct {
	ExprBase
	X  Expr
	Op lexer.Token
	Y  Expr
}

func (m *BinaryExpr) Accept(visitor Visitor) {
	visitor.VisitBinaryExpr(m)
}

type UnaryExpr struct {
	ExprBase
	Op lexer.Token
	X  Expr
}

func (m *UnaryExpr) Accept(visitor Visitor) {
	visitor.VisitUnaryExpr(m)
}

type OfTypeExpr struct {
	ExprBase
	X         Expr
	TargetTyp Type
}

func (m *OfTypeExpr) Accept(visitor Visitor) {
	visitor.VisitOfTypeExpr(m)
}

type LiteralKind int

const (
	Lit_Int = iota
	Lit_Word
	Lit_Float
	Lit_Symbol
	Lit_String
	Lit_Null
)

func (m *LiteralKind) Accept(visitor Visitor) {
	visitor.VisitLiteralKind(m)
}

type LiteralExpr struct {
	ExprBase
	Kind     LiteralKind
	IntVal   int64  // Цел
	WordVal  uint64 // Байт, Слово, Символ
	FloatStr string // Вещ, чтобы не терять точность
	StrVal   []rune // Строка
}

func (m *LiteralExpr) Accept(visitor Visitor) {
	visitor.VisitLiteralExpr(m)
}

type BoolLiteral struct {
	ExprBase
	Value bool
}

func (m *BoolLiteral) Accept(visitor Visitor) {
	visitor.VisitBoolLiteral(m)
}

type IdentExpr struct {
	ExprBase
	Name string
	Obj  Node // Decl: Var, Const or Function or TypeRef
}

func (m *IdentExpr) Accept(visitor Visitor) {
	visitor.VisitIdentExpr(m)
}

type SelectorExpr struct {
	ExprBase
	X         Expr // == nil, если импортированный объект
	Name      string
	Obj       Node // импортированный объект или поле или метод
	StdMethod *StdFunction
}

func (m *SelectorExpr) Accept(visitor Visitor) {
	visitor.VisitSelectorExpr(m)
}

type CallExpr struct {
	ExprBase
	X       Expr
	Args    []Expr
	StdFunc *StdFunction
}

func (m *CallExpr) Accept(visitor Visitor) {
	visitor.VisitCallExpr(m)
}

type UnfoldExpr struct {
	ExprBase
	X Expr
}

func (m *UnfoldExpr) Accept(visitor Visitor) {
	visitor.VisitUnfoldExpr(m)
}

type ConversionExpr struct {
	ExprBase
	X         Expr
	TargetTyp Type
	Caution   bool
	Done      bool // X уже преобразован к целевому типу
}

func (m *ConversionExpr) Accept(visitor Visitor) {
	visitor.VisitConversionExpr(m)
}

// Если тип передается, как параметр, например, в функции 'тег'
type TypeExpr struct {
	ExprBase
}

func (m *TypeExpr) Accept(visitor Visitor) {
	visitor.VisitTypeExpr(m)
}

type NotNilExpr struct {
	ExprBase
	X Expr
}

func (m *NotNilExpr) Accept(visitor Visitor) {
	visitor.VisitNotNilExpr(m)
}

//==== index

type GeneralBracketExpr struct {
	ExprBase
	X         Expr
	Index     Expr // indexation if != nil, otherwise composite
	Composite *ArrayCompositeExpr
}

func (m *GeneralBracketExpr) Accept(visitor Visitor) {
	visitor.VisitGeneralBracketExpr(m)
}

type ArrayCompositeExpr struct {
	ExprBase
	LenExpr  Expr  // если константое выражение, то значение сохраняется в Length
	Length   int64 // если вычислено, или -1
	CapExpr  Expr
	Default  Expr // default value
	Indexes  []Expr
	MaxIndex int64 // из значение индекса, -1, если нет индексов
	Values   []Expr
}

func (m *ArrayCompositeExpr) Accept(visitor Visitor) {
	visitor.VisitArrayCompositeExpr(m)
}

//=== class composite

type ClassCompositeExpr struct {
	ExprBase
	X      Expr
	Values []ValuePair
}

func (m *ClassCompositeExpr) Accept(visitor Visitor) {
	visitor.VisitClassCompositeExpr(m)
}

type ValuePair struct {
	Pos   int
	Name  string
	Field *Field
	Value Expr
}

func (m *ValuePair) Accept(visitor Visitor) {
	visitor.VisitValuePair(m)
}
