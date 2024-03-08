package ast

type Visitor interface {
	VisitModule(m *Module)
	VisitDeclaration(m *Decl)
	VisitImport(m *Import)
	VisitSetting(m *Setting)
	VisitEntryFn(m *EntryFn)
	VisitDeclBase(m *DeclBase)
	VisitInvalidDecl(m *InvalidDecl)
	VisitFunction(m *Function)
	VisitVarDecl(m *VarDecl)
	VisitConstDecl(m *ConstDecl)
	VisitTypeDecl(m *TypeDecl)
	VisitStdFunction(m *StdFunction)
	VisitExprBase(m *ExprBase)
	VisitInvalidExpr(m *InvalidExpr)
	VisitBinaryExpr(m *BinaryExpr)
	VisitUnaryExpr(m *UnaryExpr)
	VisitOfTypeExpr(m *OfTypeExpr)
	VisitLiteralKind(m *LiteralKind)
	VisitLiteralExpr(m *LiteralExpr)
	VisitBoolLiteral(m *BoolLiteral)
	VisitIdentExpr(m *IdentExpr)
	VisitSelectorExpr(m *SelectorExpr)
	VisitCallExpr(m *CallExpr)
	VisitUnfoldExpr(m *UnfoldExpr)
	VisitConversionExpr(m *ConversionExpr)
	VisitTypeExpr(m *TypeExpr)
	VisitNotNilExpr(m *NotNilExpr)
	VisitGeneralBracketExpr(m *GeneralBracketExpr)
	VisitArrayCompositeExpr(m *ArrayCompositeExpr)
	VisitClassCompositeExpr(m *ClassCompositeExpr)
	VisitValuePair(m *ValuePair)
	VisitTypeBase(m *TypeBase)
	VisitPredefinedType(m *PredefinedType)
	VisitInvalidType(m *InvalidType)
	VisitTypeRef(m *TypeRef)
	VisitVectorType(m *VectorType)
	VisitClassType(m *ClassType)
	VisitField(m *Field)
	VisitFuncType(m *FuncType)
	VisitParam(m *Param)
	VisitVariadicType(m *VariadicType)
	VisitMayBeType(m *MayBeType)
	VisitStatementBase(m *StatementBase)
	VisitStatementSeq(m *StatementSeq)
	VisitExprStatement(m *ExprStatement)
	VisitDeclStatement(m *DeclStatement)
	VisitAssignStatement(m *AssignStatement)
	VisitIncStatement(m *IncStatement)
	VisitDecStatement(m *DecStatement)
	VisitIf(m *If)
	VisitGuard(m *Guard)
	VisitSelect(m *Select)
	VisitCase(m *Case)
	VisitSelectType(m *SelectType)
	VisitCaseType(m *CaseType)
	VisitWhile(m *While)
	VisitCycle(m *Cycle)
	VisitCrash(m *Crash)
	VisitReturn(m *Return)
	VisitBreak(m *Break)
	VisitScope(m *Scope)
	VisitModifier(m *Modifier)
}
