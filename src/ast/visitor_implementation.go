package ast

import (
	"fmt"
)

type TreePrinter struct {
	jsonString string
}

func (t TreePrinter) VisitModule(m *Module) {
	fmt.Println(m.Pos)
	fmt.Println(m.Name)
	for _, importName := range m.Imports {
		importName.Accept(t)
	}
	for _, declName := range m.Decls {
		declName.Accept(t)
	}
	m.Entry.Accept(t)
	m.Inner.Accept(t)
	if m.Setting != nil {
		m.Setting.Accept(t)
	}
}

func (t TreePrinter) VisitDeclaration(m *Decl) {
}

func (t TreePrinter) VisitImport(m *Import) {
	fmt.Println(m.Pos)
	fmt.Println(m.Path)
	m.Mod.Accept(t)
	fmt.Println(m.Sources)

}

func (t TreePrinter) VisitSetting(m *Setting) {
	fmt.Println(m.Pos)
	fmt.Println(m.Path)
}

func (t TreePrinter) VisitEntryFn(m *EntryFn) {
	fmt.Println(m.Pos)
	m.Seq.Accept(t)
}

func (t TreePrinter) VisitDeclBase(m *DeclBase) {
	fmt.Println(m.Pos)
	fmt.Println(m.Name)
	m.Typ.Accept(t)
	fmt.Println(m.Exported)
	if m.Exported {
		m.Host.Accept(t)
	}
}

func (t TreePrinter) VisitInvalidDecl(m *InvalidDecl) {
	m.DeclBase.Accept(t)
}

func (t TreePrinter) VisitFunction(m *Function) {
	m.DeclBase.Accept(t)
	m.Recv.Accept(t)
	m.Seq.Accept(t)
	fmt.Println(m.External)
	m.Mod.Accept(t)
}

func (t TreePrinter) VisitVarDecl(m *VarDecl) {
	m.DeclBase.Accept(t)
	m.Init.Accept(t)
	fmt.Println(m.Later)
	fmt.Println(m.AssignOnce)
	fmt.Println(m.OutParam)
}

func (t TreePrinter) VisitConstDecl(m *ConstDecl) {
	m.DeclBase.Accept(t)
	m.Value.Accept(t)
}

func (t TreePrinter) VisitTypeDecl(m *TypeDecl) {
	m.DeclBase.Accept(t)
}

func (t TreePrinter) VisitStdFunction(m *StdFunction) {
	m.DeclBase.Accept(t)
	fmt.Println(m.Method)
}

func (t TreePrinter) VisitExprBase(m *ExprBase) {
	fmt.Println(m.Pos)
	if m.Typ != nil {
		m.Typ.Accept(t)
	}
	fmt.Println(m.ReadOnly)
}

func (t TreePrinter) VisitInvalidExpr(m *InvalidExpr) {
	m.ExprBase.Accept(t)
}

func (t TreePrinter) VisitBinaryExpr(m *BinaryExpr) {
	m.ExprBase.Accept(t)
	m.X.Accept(t)
	fmt.Println(m.Op)
	m.Y.Accept(t)
}

func (t TreePrinter) VisitUnaryExpr(m *UnaryExpr) {
	m.ExprBase.Accept(t)
	m.X.Accept(t)
	fmt.Println(m.Op)
}

func (t TreePrinter) VisitOfTypeExpr(m *OfTypeExpr) {
	m.ExprBase.Accept(t)
	m.X.Accept(t)
	m.TargetTyp.Accept(t)
}

func (t TreePrinter) VisitLiteralKind(m *LiteralKind) {
	fmt.Println(m)
}

func (t TreePrinter) VisitLiteralExpr(m *LiteralExpr) {
	m.ExprBase.Accept(t)
	fmt.Println(m.Kind)
	fmt.Println(m.IntVal)
	fmt.Println(m.WordVal)
	fmt.Println(m.FloatStr)
	fmt.Println(m.StrVal)
}

func (t TreePrinter) VisitBoolLiteral(m *BoolLiteral) {
	m.ExprBase.Accept(t)
	fmt.Println(m.Value)
}

func (t TreePrinter) VisitIdentExpr(m *IdentExpr) {
	m.ExprBase.Accept(t)
	fmt.Println(m.Name)
	m.Obj.Accept(t)
}

func (t TreePrinter) VisitSelectorExpr(m *SelectorExpr) {
	m.ExprBase.Accept(t)
	fmt.Println(m.Name)
	if m.X != nil {
		m.X.Accept(t)
	}
	m.Obj.Accept(t)
	m.StdMethod.Accept(t)
}

func (t TreePrinter) VisitCallExpr(m *CallExpr) {
	m.ExprBase.Accept(t)
	m.X.Accept(t)
	for _, argument := range m.Args {
		argument.Accept(t)
	}
	m.StdFunc.Accept(t)
}

func (t TreePrinter) VisitUnfoldExpr(m *UnfoldExpr) {
	m.ExprBase.Accept(t)
	m.X.Accept(t)
}

func (t TreePrinter) VisitConversionExpr(m *ConversionExpr) {
	m.ExprBase.Accept(t)
	m.X.Accept(t)
	m.TargetTyp.Accept(t)
	fmt.Println(m.Caution)
	fmt.Println(m.Done)
}

func (t TreePrinter) VisitTypeExpr(m *TypeExpr) {
	m.ExprBase.Accept(t)
}

func (t TreePrinter) VisitNotNilExpr(m *NotNilExpr) {
	m.ExprBase.Accept(t)
	m.X.Accept(t)
	fmt.Println(m.ReadOnly)
}

func (t TreePrinter) VisitGeneralBracketExpr(m *GeneralBracketExpr) {
	m.ExprBase.Accept(t)
	if m.X != nil {
		m.X.Accept(t)
	}
	if m.Index != nil {
		m.Index.Accept(t)
	}
	if m.Composite != nil {
		m.Composite.Accept(t)
	}
}

func (t TreePrinter) VisitArrayCompositeExpr(m *ArrayCompositeExpr) {
	m.ExprBase.Accept(t)
	if m.LenExpr != nil {
		m.LenExpr.Accept(t)
	}
	fmt.Println(m.Length)
	if m.CapExpr != nil {
		m.CapExpr.Accept(t)
	}
	if m.Default != nil {
		m.Default.Accept(t)
	}
	for _, indexName := range m.Indexes {
		indexName.Accept(t)
	}
	fmt.Println(m.MaxIndex)
	for _, value := range m.Values {
		value.Accept(t)
	}

}

func (t TreePrinter) VisitClassCompositeExpr(m *ClassCompositeExpr) {
	m.ExprBase.Accept(t)
	m.X.Accept(t)
	for _, value := range m.Values {
		value.Accept(t)
	}
}

func (t TreePrinter) VisitValuePair(m *ValuePair) {
	fmt.Println(m.Pos)
	fmt.Println(m.Name)
	m.Field.Accept(t)
	m.Value.Accept(t)
}

func (t TreePrinter) VisitTypeBase(m *TypeBase) {
	fmt.Println(m.Pos)
}

func (t TreePrinter) VisitPredefinedType(m *PredefinedType) {
	m.TypeBase.Accept(t)
	fmt.Println(m.Name)
}

func (t TreePrinter) VisitInvalidType(m *InvalidType) {
	m.TypeBase.Accept(t)
}

func (t TreePrinter) VisitTypeRef(m *TypeRef) {
	m.TypeBase.Accept(t)
	fmt.Println(m.TypeName)
	fmt.Println(m.ModuleName)
	m.TypeDecl.Accept(t)
	m.Typ.Accept(t)
}

func (t TreePrinter) VisitVectorType(m *VectorType) {
	m.TypeBase.Accept(t)
	m.ElementTyp.Accept(t)
}

func (t TreePrinter) VisitClassType(m *ClassType) {
	m.TypeBase.Accept(t)
	m.BaseTyp.Accept(t)
	for _, field := range m.Fields {
		field.Accept(t)
	}
	for _, method := range m.Methods {
		method.Accept(t)
	}
	for _, member := range m.Members {
		member.Accept(t)
	}
}

func (t TreePrinter) VisitField(m *Field) {
	m.DeclBase.Accept(t)
	m.Init.Accept(t)
	fmt.Println(m.AssignOnce)
	fmt.Println(m.Later)
}

func (t TreePrinter) VisitFuncType(m *FuncType) {
	m.TypeBase.Accept(t)
	for _, param := range m.Params {
		param.Accept(t)
	}
	m.ReturnTyp.Accept(t)
}

func (t TreePrinter) VisitParam(m *Param) {
	m.DeclBase.Accept(t)
	fmt.Println(m.Out)
}

func (t TreePrinter) VisitVariadicType(m *VariadicType) {
	m.TypeBase.Accept(t)
	m.ElementTyp.Accept(t)
}

func (t TreePrinter) VisitMayBeType(m *MayBeType) {
	m.TypeBase.Accept(t)
	m.Typ.Accept(t)
}

func (t TreePrinter) VisitStatementBase(m *StatementBase) {
	fmt.Println(m.Pos)
}

func (t TreePrinter) VisitStatementSeq(m *StatementSeq) {
	m.StatementBase.Accept(t)
	for _, statement := range m.Statements {
		statement.Accept(t)
	}
	m.Inner.Accept(t)
}

func (t TreePrinter) VisitExprStatement(m *ExprStatement) {
	m.StatementBase.Accept(t)
	m.X.Accept(t)
}

func (t TreePrinter) VisitDeclStatement(m *DeclStatement) {
	m.StatementBase.Accept(t)
	m.D.Accept(t)
}

func (t TreePrinter) VisitAssignStatement(m *AssignStatement) {
	m.StatementBase.Accept(t)
	m.L.Accept(t)
	m.R.Accept(t)
}

func (t TreePrinter) VisitIncStatement(m *IncStatement) {
	m.StatementBase.Accept(t)
	m.L.Accept(t)
}

func (t TreePrinter) VisitDecStatement(m *DecStatement) {
	m.StatementBase.Accept(t)
	m.L.Accept(t)
}

func (t TreePrinter) VisitIf(m *If) {
	m.StatementBase.Accept(t)
	m.Cond.Accept(t)
	m.Then.Accept(t)
	m.Else.Accept(t)
}

func (t TreePrinter) VisitGuard(m *Guard) {
	m.StatementBase.Accept(t)
	m.Cond.Accept(t)
	m.Else.Accept(t)
}

func (t TreePrinter) VisitSelect(m *Select) {
	m.StatementBase.Accept(t)
	if m.X != nil {
		m.X.Accept(t)
	}
	for _, caseName := range m.Cases {
		caseName.Accept(t)
	}
	m.Else.Accept(t)
}

func (t TreePrinter) VisitCase(m *Case) {
	m.StatementBase.Accept(t)
	for _, expr := range m.Exprs {
		expr.Accept(t)
	}
	m.Seq.Accept(t)
}

func (t TreePrinter) VisitSelectType(m *SelectType) {
	m.StatementBase.Accept(t)
	fmt.Println(m.VarIdent)
	m.X.Accept(t)
	for _, caseName := range m.Cases {
		caseName.Accept(t)
	}
	m.Else.Accept(t)
}

func (t TreePrinter) VisitCaseType(m *CaseType) {
	m.StatementBase.Accept(t)
	if m.Var != nil {
		m.Var.Accept(t)
	}
	for _, typeName := range m.Types {
		typeName.Accept(t)
	}
	m.Seq.Accept(t)
}

func (t TreePrinter) VisitWhile(m *While) {
	m.StatementBase.Accept(t)
	m.Cond.Accept(t)
	m.Seq.Accept(t)
}

func (t TreePrinter) VisitCycle(m *Cycle) {
	m.StatementBase.Accept(t)
	m.IndexVar.Accept(t)
	m.ElementVar.Accept(t)
	m.Expr.Accept(t)
	m.Seq.Accept(t)
}

func (t TreePrinter) VisitCrash(m *Crash) {
	m.StatementBase.Accept(t)
	m.X.Accept(t)
}

func (t TreePrinter) VisitReturn(m *Return) {
	m.StatementBase.Accept(t)
	m.X.Accept(t)
	m.ReturnTyp.Accept(t)
}

func (t TreePrinter) VisitBreak(m *Break) {
	m.StatementBase.Accept(t)
}

func (t TreePrinter) VisitScope(m *Scope) {
	if m.Outer != nil {
		m.Outer.Accept(t)
	}
	for _, decl := range m.Names {
		decl.Accept(t)
	}
}

func (t TreePrinter) VisitModifier(m *Modifier) {
	fmt.Println(m.Name)
	fmt.Println(m.Attrs)
}
