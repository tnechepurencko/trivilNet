package ast

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

type TreePrinter struct {
	JsonString string
}

func (t TreePrinter) VisitModule(m *Module) TreePrinter {
	var result = "{"
	result = result + "\"Pos\": " + strconv.Itoa(m.Pos) + ","
	result = result + "\"Name\": \"" + m.Name + "\","
	result = result + "\"Imports\": ["
	for _, importName := range m.Imports {
		result = result + importName.Accept(t).JsonString + ","
	}
	result = strings.TrimSuffix(result, ",")
	result = result + "],"
	result = result + "\"Decl\": ["
	for _, declName := range m.Decls {
		result = result + declName.Accept(t).JsonString + ","
	}
	result = strings.TrimSuffix(result, ",")
	result = result + "],"
	result = result + "\"Entry\": "
	result = result + m.Entry.Accept(t).JsonString + ","
	result = result + "\"Inner\": "
	result = result + m.Inner.Accept(t).JsonString + ","
	result = result + "\"Setting\": "
	if m.Setting != nil {
		result = result + m.Setting.Accept(t).JsonString
	} else {
		result = result + " null"
	}
	result = result + "}"
	t.JsonString = result
	return t
}

func (t TreePrinter) VisitDeclaration(m *Decl) TreePrinter {
	return t
}

func (t TreePrinter) VisitImport(m *Import) TreePrinter {
	var result = "{"
	result = result + "\"Pos\": " + strconv.Itoa(m.Pos) + ","
	result = result + "\"Path\": \"" + m.Path + "\","
	result = result + "\"Mod\": "
	result = result + m.Mod.Accept(t).JsonString + ","
	var IDs []string
	for _, i := range m.Sources {
		IDs = append(IDs, strconv.Itoa(i))
	}
	result = result + "\"Sources\": [" + strings.TrimSuffix(strings.Join(IDs, ", "), ",") + "],\n"
	result = result + "}"
	t.JsonString = result
	return t
}

func (t TreePrinter) VisitSetting(m *Setting) TreePrinter {
	var result = "{"
	result = result + "\"Pos\": " + strconv.Itoa(m.Pos) + ","
	result = result + "\"Path\": \"" + m.Path + "\","
	result = result + "}"
	t.JsonString = result
	return t
}

func (t TreePrinter) VisitEntryFn(m *EntryFn) TreePrinter {
	var result = "{"
	result = result + "\"Pos\": " + strconv.Itoa(m.Pos) + ","
	result = result + "\"Seq\": "
	result = result + m.Seq.Accept(t).JsonString + ","
	result = result + "}"
	t.JsonString = result
	return t
}

func (t TreePrinter) VisitDeclBase(m *DeclBase) TreePrinter {
	var result = "{"
	result = result + "\"Pos\": " + strconv.Itoa(m.Pos) + ","
	result = result + "\"Name\": \"" + m.Name + "\","
	result = result + "\"Typ\": "
	if m.Typ != nil {
		result = result + m.Typ.Accept(t).JsonString + ","
	} else {
		result = result + " null,"
	}
	result = result + "\"Exported\": " + strconv.FormatBool(m.Exported) + ","
	result = result + "\"Host\": "
	if m.Exported {
		result = result + m.Host.Accept(t).JsonString
	} else {
		result = result + " null"
	}
	result = result + "}"
	t.JsonString = result
	return t
}

func (t TreePrinter) VisitInvalidDecl(m *InvalidDecl) TreePrinter {
	//m.DeclBase.Accept(t)
	return t
}

func (t TreePrinter) VisitFunction(m *Function) TreePrinter {
	var result = "{"
	result = result + "\"DeclBase\": "
	result = result + m.DeclBase.Accept(t).JsonString + ","
	result = result + "\"Recv\": "
	result = result + m.Recv.Accept(t).JsonString + ","
	result = result + "\"Seq\": "
	result = result + m.Seq.Accept(t).JsonString + ","
	result = result + "\"External\": " + strconv.FormatBool(m.External) + ","
	result = result + "\"Mod\": "
	result = result + m.Mod.Accept(t).JsonString
	result = result + "}"
	t.JsonString = result
	return t
}

func (t TreePrinter) VisitVarDecl(m *VarDecl) TreePrinter {
	var result = "{"
	result = result + "\"DeclBase\": "
	result = result + m.DeclBase.Accept(t).JsonString + ","
	result = result + "\"Init\": "
	result = result + m.Init.Accept(t).JsonString + ","
	result = result + "\"Later\": " + strconv.FormatBool(m.Later) + ","
	result = result + "\"AssignOnce\": " + strconv.FormatBool(m.AssignOnce) + ","
	result = result + "\"OutParam\": " + strconv.FormatBool(m.OutParam)
	result = result + "}"
	t.JsonString = result
	return t
}

func (t TreePrinter) VisitConstDecl(m *ConstDecl) TreePrinter {
	var result = "{"
	result = result + "\"DeclBase\": "
	result = result + m.DeclBase.Accept(t).JsonString + ","
	result = result + "\"Value\": "
	result = result + m.Value.Accept(t).JsonString
	result = result + "}"
	t.JsonString = result
	return t
}

func (t TreePrinter) VisitTypeDecl(m *TypeDecl) TreePrinter {
	var result = "{"
	result = result + "\"DeclBase\": "
	result = result + m.DeclBase.Accept(t).JsonString
	result = result + "}"
	t.JsonString = result
	return t
}

func (t TreePrinter) VisitStdFunction(m *StdFunction) TreePrinter {
	var result = "{"
	result = result + "\"DeclBase\": "
	result = result + m.DeclBase.Accept(t).JsonString + ","
	result = result + "\"Method\": " + strconv.FormatBool(m.Method)
	result = result + "}"
	t.JsonString = result
	return t
}

func (t TreePrinter) VisitExprBase(m *ExprBase) TreePrinter {
	var result = "{"
	result = result + "\"Pos\": " + strconv.Itoa(m.Pos) + ","
	result = result + "\"Typ\": "
	if m.Typ != nil {
		result = result + m.Typ.Accept(t).JsonString + ","
	} else {
		result = result + " null,"
	}
	result = result + "\"ReadOnly\": " + strconv.FormatBool(m.ReadOnly)
	result = result + "}"
	t.JsonString = result
	return t
}

func (t TreePrinter) VisitInvalidExpr(m *InvalidExpr) TreePrinter {
	//m.ExprBase.Accept(t)
	return t
}

func (t TreePrinter) VisitBinaryExpr(m *BinaryExpr) TreePrinter {
	var result = "{"
	result = result + "\"ExprBase\": "
	result = result + m.ExprBase.Accept(t).JsonString + ","
	result = result + "\"X\": "
	result = result + m.X.Accept(t).JsonString + ","
	result = result + "\"Op\": " + strconv.Itoa(int(m.Op)) + ","
	result = result + "\"Y\": "
	result = result + m.Y.Accept(t).JsonString
	result = result + "}"
	t.JsonString = result
	return t
}

func (t TreePrinter) VisitUnaryExpr(m *UnaryExpr) TreePrinter {
	var result = "{"
	result = result + "\"ExprBase\": "
	result = result + m.ExprBase.Accept(t).JsonString + ","
	result = result + "\"X\": "
	result = result + m.X.Accept(t).JsonString + ","
	result = result + "\"Op\": " + strconv.Itoa(int(m.Op))
	result = result + "}"
	t.JsonString = result
	return t
}

func (t TreePrinter) VisitOfTypeExpr(m *OfTypeExpr) TreePrinter {
	var result = "{"
	result = result + "\"ExprBase\": "
	result = result + m.ExprBase.Accept(t).JsonString + ","
	result = result + "\"X\": "
	result = result + m.X.Accept(t).JsonString + ","
	result = result + "\"TargetTyp\": "
	result = result + m.TargetTyp.Accept(t).JsonString
	result = result + "}"
	t.JsonString = result
	return t
}

func (t TreePrinter) VisitLiteralKind(m *LiteralKind) TreePrinter {
	return t
}

func (t TreePrinter) VisitLiteralExpr(m *LiteralExpr) TreePrinter {
	var result = "{"
	result = result + "\"ExprBase\": "
	result = result + m.ExprBase.Accept(t).JsonString + ","
	result = result + "\"Kind\": " + strconv.Itoa(int(m.Kind)) + ","
	result = result + "\"IntVal\": " + strconv.FormatInt(m.IntVal, 10) + ","
	result = result + "\"WordVal\": " + strconv.FormatInt(int64(m.WordVal), 10) + ","
	result = result + "\"Float.Str\": \"" + m.FloatStr + "\","
	var IDs string
	for _, i := range m.StrVal {
		IDs = IDs + string(i)
	}
	result = result + "\"StrVal\": \"" + IDs + "\""
	result = result + "}"
	t.JsonString = result
	return t
}

func (t TreePrinter) VisitBoolLiteral(m *BoolLiteral) TreePrinter {
	var result = "{"
	result = result + "\"ExprBase\": "
	result = result + m.ExprBase.Accept(t).JsonString + ","
	result = result + "\"Value\": " + strconv.FormatBool(m.Value)
	result = result + "}"
	t.JsonString = result
	return t
}

func (t TreePrinter) VisitIdentExpr(m *IdentExpr) TreePrinter {
	var result = "{"
	result = result + "\"ExprBase\": "
	result = result + m.ExprBase.Accept(t).JsonString + ","
	result = result + "\"Name\": \"" + m.Name + "\","
	result = result + "\"Obj\": "
	result = result + m.Obj.Accept(t).JsonString
	result = result + "}"
	t.JsonString = result
	return t
}

func (t TreePrinter) VisitSelectorExpr(m *SelectorExpr) TreePrinter {
	var result = "{"
	result = result + "\"ExprBase\": "
	result = result + m.ExprBase.Accept(t).JsonString + ","
	result = result + "\"Name\": \"" + m.Name + "\","
	result = result + "\"X\": "
	if m.X != nil {
		result = result + m.X.Accept(t).JsonString + ","
	} else {
		result = result + "null,"
	}
	result = result + "\"Obj\": "
	result = result + m.Obj.Accept(t).JsonString + ","
	result = result + "\"StdMethod\": "
	result = result + m.StdMethod.Accept(t).JsonString
	result = result + "}"
	t.JsonString = result
	return t
}

func (t TreePrinter) VisitCallExpr(m *CallExpr) TreePrinter {
	var result = "{"
	result = result + "\"ExprBase\": "
	result = result + m.ExprBase.Accept(t).JsonString + ","
	result = result + "\"X\": "
	if m.X != nil {
		result = result + m.X.Accept(t).JsonString + ","
	} else {
		result = result + "null,"
	}
	result = result + "\"Args\": ["
	for _, argument := range m.Args {
		result = result + argument.Accept(t).JsonString + ","
	}
	result = strings.TrimSuffix(result, ",")
	result = result + "],"
	result = result + "\"StdFunc\": "
	if m.StdFunc != nil {
		result = result + m.StdFunc.Accept(t).JsonString
	} else {
		result = result + "null"
	}
	result = result + "}"
	t.JsonString = result
	return t
}

func (t TreePrinter) VisitUnfoldExpr(m *UnfoldExpr) TreePrinter {
	var result = "{"
	result = result + "\"ExprBase\": "
	result = result + m.ExprBase.Accept(t).JsonString + ","
	result = result + "\"X\": " + m.X.Accept(t).JsonString
	result = result + "}"
	t.JsonString = result
	return t
}

func (t TreePrinter) VisitConversionExpr(m *ConversionExpr) TreePrinter {
	var result = "{"
	result = result + "\"ExprBase\": "
	result = result + m.ExprBase.Accept(t).JsonString + ","
	result = result + "\"X\": " + m.X.Accept(t).JsonString + ","
	result = result + "\"TargetTyp\": " + m.TargetTyp.Accept(t).JsonString + ","
	result = result + "\"Caution\": " + strconv.FormatBool(m.Caution) + ","
	result = result + "\"Done\": " + strconv.FormatBool(m.Done)
	result = result + "}"
	t.JsonString = result
	return t
}

func (t TreePrinter) VisitTypeExpr(m *TypeExpr) TreePrinter {
	var result = "{"
	result = result + "\"ExprBase\": "
	result = result + m.ExprBase.Accept(t).JsonString
	result = result + "}"
	t.JsonString = result
	return t
}

func (t TreePrinter) VisitNotNilExpr(m *NotNilExpr) TreePrinter {
	var result = "{"
	result = result + "\"ExprBase\": "
	result = result + m.ExprBase.Accept(t).JsonString + ","
	result = result + "\"X\": "
	result = result + m.X.Accept(t).JsonString + ","
	result = result + "\"ReadOnly\": " + strconv.FormatBool(m.ReadOnly)
	result = result + "}"
	t.JsonString = result
	return t
}

func (t TreePrinter) VisitGeneralBracketExpr(m *GeneralBracketExpr) TreePrinter {
	var result = "{"
	result = result + "\"ExprBase\": "
	result = result + m.ExprBase.Accept(t).JsonString + ","
	result = result + "\"X\": "
	if m.X != nil {
		result = result + m.X.Accept(t).JsonString + ","
	} else {
		result = result + "null,"
	}
	result = result + "\"Index\": "
	if m.Index != nil {
		result = result + m.Index.Accept(t).JsonString + ","
	} else {
		result = result + "null,"
	}
	result = result + "\"Composite\": "
	if m.Composite != nil {
		result = result + m.Composite.Accept(t).JsonString + ","
	} else {
		result = result + "null,"
	}
	result = result + "}"
	t.JsonString = result
	return t
}

func (t TreePrinter) VisitArrayCompositeExpr(m *ArrayCompositeExpr) TreePrinter {
	var result = "{"
	result = result + "\"ExprBase\": "
	result = result + m.ExprBase.Accept(t).JsonString + ","
	result = result + "\"LenExpr\": "
	if m.LenExpr != nil {
		result = result + m.LenExpr.Accept(t).JsonString + ","
	} else {
		result = result + "null,"
	}
	result = result + "\"Length\": " + strconv.FormatInt(m.Length, 10) + ","
	result = result + "\"CapExpr\": "
	if m.CapExpr != nil {
		result = result + m.CapExpr.Accept(t).JsonString + ","
	} else {
		result = result + "null,"
	}
	result = result + "\"Default\": "
	if m.Default != nil {
		result = result + m.Default.Accept(t).JsonString + ","
	} else {
		result = result + "null,"
	}
	result = result + "\"Indexes\": ["
	for _, indexName := range m.Indexes {
		result = result + indexName.Accept(t).JsonString + ","
	}
	result = strings.TrimSuffix(result, ",")
	result = result + "],"
	result = result + "\"MaxIndex\": " + strconv.FormatInt(m.MaxIndex, 10) + ","
	result = result + "\"Values\": ["
	for _, value := range m.Values {
		result = result + value.Accept(t).JsonString + ","
	}
	result = strings.TrimSuffix(result, ",")
	result = result + "]"
	result = result + "}"
	t.JsonString = result
	return t

}

func (t TreePrinter) VisitClassCompositeExpr(m *ClassCompositeExpr) TreePrinter {
	var result = "{"
	result = result + "\"ExprBase\": "
	result = result + m.ExprBase.Accept(t).JsonString + ","
	result = result + "\"X\": "
	result = result + m.X.Accept(t).JsonString + ","
	result = result + "\"Values\": ["
	for _, value := range m.Values {
		result = result + value.Accept(t).JsonString + ","
	}
	result = strings.TrimSuffix(result, ",")
	result = result + "]"
	result = result + "}"
	t.JsonString = result
	return t
}

func (t TreePrinter) VisitValuePair(m *ValuePair) TreePrinter {
	var result = "{"
	result = result + "\"Pos\": " + strconv.Itoa(m.Pos) + ","
	result = result + "\"Name\": \"" + m.Name + "\","
	result = result + "\"Field\": " + m.Field.Accept(t).JsonString + ","
	result = result + "\"Value\": " + m.Value.Accept(t).JsonString
	result = result + "}"
	t.JsonString = result
	return t
}

func (t TreePrinter) VisitTypeBase(m *TypeBase) TreePrinter {
	t.JsonString = "{\"Pos\": " + strconv.Itoa(m.Pos) + "}"
	return t
}

func (t TreePrinter) VisitPredefinedType(m *PredefinedType) TreePrinter {
	var result = "{"
	result = result + "\"TypeBase\": "
	result = result + m.TypeBase.Accept(t).JsonString + ","
	result = result + "\"Name\": \"" + m.Name + "\""
	result = result + "}"
	t.JsonString = result
	return t
}

func (t TreePrinter) VisitInvalidType(m *InvalidType) TreePrinter {
	var result = "{"
	result = result + "\"TypeBase\": "
	result = result + m.TypeBase.Accept(t).JsonString
	result = result + "}"
	t.JsonString = result
	return t
}

func (t TreePrinter) VisitTypeRef(m *TypeRef) TreePrinter {
	var result = "{"
	result = result + "\"TypeBase\": "
	result = result + m.TypeBase.Accept(t).JsonString + ","
	result = result + "\"TypeName\": \"" + m.TypeName + "\","
	result = result + "\"ModuleName\": \"" + m.ModuleName + "\","
	result = result + "\"TypeDecl\": "
	result = result + m.TypeDecl.Accept(t).JsonString + ","
	result = result + "\"Typ\": "
	result = result + m.Typ.Accept(t).JsonString
	result = result + "}"
	t.JsonString = result
	return t
}

func (t TreePrinter) VisitVectorType(m *VectorType) TreePrinter {
	var result = "{"
	result = result + "\"TypeBase\": "
	result = result + m.TypeBase.Accept(t).JsonString + ","
	result = result + "\"ElementTyp\": "
	result = result + m.ElementTyp.Accept(t).JsonString
	result = result + "}"
	t.JsonString = result
	return t
}

func (t TreePrinter) VisitClassType(m *ClassType) TreePrinter {
	var result = "{"
	result = result + "\"TypeBase\": "
	result = result + m.TypeBase.Accept(t).JsonString + ","
	result = result + "\"BaseTyp\": "
	result = result + m.BaseTyp.Accept(t).JsonString + ","
	result = result + "\"Fields\": ["
	for _, field := range m.Fields {
		result = result + field.Accept(t).JsonString + ","
	}
	result = strings.TrimSuffix(result, ",")
	result = result + "],"
	result = result + "\"Methods\": ["
	for _, method := range m.Methods {
		result = result + method.Accept(t).JsonString + ","
	}
	result = strings.TrimSuffix(result, ",")
	result = result + "],"
	b := new(bytes.Buffer)
	for key, value := range m.Members {
		fmt.Fprintf(b, "\"%s\":%s,", key, value.Accept(t).JsonString)
	}
	var s = strings.TrimSuffix(b.String(), ",")
	result = result + "\"Members\": [" + s + "]"
	result = result + "}"
	t.JsonString = result
	return t
}

func (t TreePrinter) VisitField(m *Field) TreePrinter {
	var result = "{"
	result = result + "\"DeclBase\": "
	result = result + m.DeclBase.Accept(t).JsonString + ","
	result = result + "\"Init\": "
	result = result + m.Init.Accept(t).JsonString + ","
	result = result + "\"AssignOnce\": " + strconv.FormatBool(m.AssignOnce) + ","
	result = result + "\"Later\": " + strconv.FormatBool(m.Later)
	result = result + "}"
	t.JsonString = result
	return t
}

func (t TreePrinter) VisitFuncType(m *FuncType) TreePrinter {
	var result = "{"
	result = result + "\"TypeBase\": "
	result = result + m.TypeBase.Accept(t).JsonString + ","
	result = result + "\"Params\": ["
	for _, param := range m.Params {
		result = result + param.Accept(t).JsonString + ","
	}
	result = strings.TrimSuffix(result, ",")
	result = result + "],"
	result = result + "\"ReturnTyp\": "
	result = result + m.ReturnTyp.Accept(t).JsonString
	result = result + "}"
	t.JsonString = result
	return t
}

func (t TreePrinter) VisitParam(m *Param) TreePrinter {
	var result = "{"
	result = result + "\"DeclBase\": "
	result = result + m.DeclBase.Accept(t).JsonString + ","
	result = result + "\"Out\": " + strconv.FormatBool(m.Out)
	result = result + "}"
	t.JsonString = result
	return t
}

func (t TreePrinter) VisitVariadicType(m *VariadicType) TreePrinter {
	var result = "{"
	result = result + "\"TypeBase\": "
	result = result + m.TypeBase.Accept(t).JsonString + ","
	result = result + "\"ElementTyp\": "
	result = result + m.ElementTyp.Accept(t).JsonString
	result = result + "}"
	t.JsonString = result
	return t
}

func (t TreePrinter) VisitMayBeType(m *MayBeType) TreePrinter {
	var result = "{"
	result = result + "\"TypeBase\": "
	result = result + m.TypeBase.Accept(t).JsonString + ","
	result = result + "\"Typ\": "
	result = result + m.Typ.Accept(t).JsonString
	result = result + "}"
	t.JsonString = result
	return t
}

func (t TreePrinter) VisitStatementBase(m *StatementBase) TreePrinter {
	t.JsonString = "{\"Pos\": " + strconv.Itoa(m.Pos) + "}"
	return t
}

func (t TreePrinter) VisitStatementSeq(m *StatementSeq) TreePrinter {
	var result = "{"
	result = result + "\"StatementBase\": "
	result = result + m.StatementBase.Accept(t).JsonString + ","
	result = result + "\"Statements\": ["
	for _, statement := range m.Statements {
		result = result + statement.Accept(t).JsonString + ","
	}
	result = strings.TrimSuffix(result, ",")
	result = result + "],"
	result = result + "\"Inner\": "
	if m.Inner != nil {
		result = result + m.Inner.Accept(t).JsonString
	} else {
		result = result + "null"
	}

	result = result + "}"
	t.JsonString = result
	return t
}

func (t TreePrinter) VisitExprStatement(m *ExprStatement) TreePrinter {
	var result = "{"
	result = result + "\"StatementBase\": "
	result = result + m.StatementBase.Accept(t).JsonString + ","
	result = result + "\"X\": "
	result = result + m.X.Accept(t).JsonString
	result = result + "}"
	t.JsonString = result
	return t
}

func (t TreePrinter) VisitDeclStatement(m *DeclStatement) TreePrinter {
	var result = "{"
	result = result + "\"StatementBase\": "
	result = result + m.StatementBase.Accept(t).JsonString + ","
	result = result + "\"D\": "
	result = result + m.D.Accept(t).JsonString
	result = result + "}"
	t.JsonString = result
	return t
}

func (t TreePrinter) VisitAssignStatement(m *AssignStatement) TreePrinter {
	var result = "{"
	result = result + "\"StatementBase\": "
	result = result + m.StatementBase.Accept(t).JsonString + ","
	result = result + "\"L\": "
	result = result + m.L.Accept(t).JsonString + ","
	result = result + "\"R\": "
	result = result + m.R.Accept(t).JsonString
	result = result + "}"
	t.JsonString = result
	return t
}

func (t TreePrinter) VisitIncStatement(m *IncStatement) TreePrinter {
	var result = "{"
	result = result + "\"StatementBase\": "
	result = result + m.StatementBase.Accept(t).JsonString + ","
	result = result + "\"LInc\": "
	result = result + m.L.Accept(t).JsonString
	result = result + "}"
	t.JsonString = result
	return t
}

func (t TreePrinter) VisitDecStatement(m *DecStatement) TreePrinter {
	var result = "{"
	result = result + "\"StatementBase\": "
	result = result + m.StatementBase.Accept(t).JsonString + ","
	result = result + "\"LDec\": "
	result = result + m.L.Accept(t).JsonString
	result = result + "}"
	t.JsonString = result
	return t
}

func (t TreePrinter) VisitIf(m *If) TreePrinter {
	var result = "{"
	result = result + "\"StatementBase\": "
	result = result + m.StatementBase.Accept(t).JsonString + ","
	result = result + "\"Cond\": "
	result = result + m.Cond.Accept(t).JsonString + ","
	result = result + "\"Then\": "
	result = result + m.Then.Accept(t).JsonString + ","
	result = result + "\"Else\": "
	result = result + m.Else.Accept(t).JsonString
	result = result + "}"
	t.JsonString = result
	return t
}

func (t TreePrinter) VisitGuard(m *Guard) TreePrinter {
	var result = "{"
	result = result + "\"StatementBase\": "
	result = result + m.StatementBase.Accept(t).JsonString + ","
	result = result + "\"Cond\": "
	result = result + m.Cond.Accept(t).JsonString + ","
	result = result + "\"Else\": "
	result = result + m.Else.Accept(t).JsonString
	result = result + "}"
	t.JsonString = result
	return t
}

func (t TreePrinter) VisitSelect(m *Select) TreePrinter {
	var result = "{"
	result = result + "\"StatementBase\": "
	result = result + m.StatementBase.Accept(t).JsonString + ","
	result = result + "\"X\": "
	if m.X != nil {
		result = result + m.X.Accept(t).JsonString + ","
	} else {
		result = result + "null,"
	}
	result = result + "\"Cases\": ["
	for _, caseName := range m.Cases {
		result = result + caseName.Accept(t).JsonString + ","
	}
	result = strings.TrimSuffix(result, ",")
	result = result + "],"
	result = result + "\"Else\": "
	result = result + m.Else.Accept(t).JsonString
	result = result + "}"
	t.JsonString = result
	return t
}

func (t TreePrinter) VisitCase(m *Case) TreePrinter {
	var result = "{"
	result = result + "\"StatementBase\": "
	result = result + m.StatementBase.Accept(t).JsonString + ","
	result = result + "\"Exprs\": "
	for _, expr := range m.Exprs {
		result = result + expr.Accept(t).JsonString + ","
	}
	result = strings.TrimSuffix(result, ",")
	result = result + "\"Seq\": "
	result = result + m.Seq.Accept(t).JsonString
	result = result + "}"
	t.JsonString = result
	return t
}

func (t TreePrinter) VisitSelectType(m *SelectType) TreePrinter {
	var result = "{"
	result = result + "\"StatementBase\": "
	result = result + m.StatementBase.Accept(t).JsonString + ","
	result = result + "\"VarIdent\": \"" + m.VarIdent + "\","
	result = result + "\"X\": "
	result = result + m.X.Accept(t).JsonString + ","
	result = result + "\"Cases\": ["
	for _, caseName := range m.Cases {
		result = result + caseName.Accept(t).JsonString + ","
	}
	result = strings.TrimSuffix(result, ",")
	result = result + "],"
	result = result + "\"Else\": "
	result = result + m.Else.Accept(t).JsonString
	result = result + "}"
	t.JsonString = result
	return t
}

func (t TreePrinter) VisitCaseType(m *CaseType) TreePrinter {
	var result = "{"
	result = result + "\"StatementBase\": "
	result = result + m.StatementBase.Accept(t).JsonString + ","
	result = result + "\"Var\": "
	if m.Var != nil {
		result = result + m.Var.Accept(t).JsonString + ","
	} else {
		result = result + "null,"
	}
	result = result + "\"Types\": ["
	for _, typeName := range m.Types {
		result = result + typeName.Accept(t).JsonString + ","
	}
	result = strings.TrimSuffix(result, ",")
	result = result + "],"
	result = result + "\"Seq\": "
	result = result + m.Seq.Accept(t).JsonString
	result = result + "}"
	t.JsonString = result
	return t
}

func (t TreePrinter) VisitWhile(m *While) TreePrinter {
	var result = "{"
	result = result + "\"StatementBase\": "
	result = result + m.StatementBase.Accept(t).JsonString + ","
	result = result + "\"Cond\": "
	result = result + m.Cond.Accept(t).JsonString + ","
	result = result + "\"Seq\": "
	result = result + m.Seq.Accept(t).JsonString
	result = result + "}"
	t.JsonString = result
	return t
}

func (t TreePrinter) VisitCycle(m *Cycle) TreePrinter {
	var result = "{"
	result = result + "\"StatementBase\": "
	result = result + m.StatementBase.Accept(t).JsonString + ","
	result = result + "\"IndexVar\": "
	result = result + m.IndexVar.Accept(t).JsonString + ","
	result = result + "\"ElementVar\": "
	result = result + m.ElementVar.Accept(t).JsonString + ","
	result = result + "\"Expr\": "
	result = result + m.Expr.Accept(t).JsonString + ","
	result = result + "\"Seq\": "
	result = result + m.Seq.Accept(t).JsonString
	result = result + "}"
	t.JsonString = result
	return t
}

func (t TreePrinter) VisitCrash(m *Crash) TreePrinter {
	var result = "{"
	result = result + "\"StatementBase\": "
	result = result + m.StatementBase.Accept(t).JsonString + ","
	result = result + "\"X\": "
	result = result + m.X.Accept(t).JsonString
	result = result + "}"
	t.JsonString = result
	return t
}

func (t TreePrinter) VisitReturn(m *Return) TreePrinter {
	var result = "{"
	result = result + "\"StatementBase\": "
	result = result + m.StatementBase.Accept(t).JsonString + ","
	result = result + "\"X\": "
	result = result + m.X.Accept(t).JsonString + ","
	result = result + "\"ReturnTyp\": "
	result = result + m.ReturnTyp.Accept(t).JsonString
	result = result + "}"
	t.JsonString = result
	return t
}

func (t TreePrinter) VisitBreak(m *Break) TreePrinter {
	var result = "{"
	result = result + "\"StatementBase\": "
	result = result + m.StatementBase.Accept(t).JsonString
	result = result + "}"
	t.JsonString = result
	return t
}

func (t TreePrinter) VisitScope(m *Scope) TreePrinter {
	var result = "{"
	result = result + "\"Outer\": "
	if m.Outer != nil {
		result = result + m.Outer.Accept(t).JsonString + ","
	} else {
		result = result + "null,"
	}
	b := new(bytes.Buffer)
	for key, value := range m.Names {
		fmt.Fprintf(b, "\"%s\":%s,", key, value.Accept(t).JsonString)
	}
	var s = strings.TrimSuffix(b.String(), ",")
	result = result + "\"Names\": [" + s + "]"
	result = result + "}"
	t.JsonString = result
	return t
}

func (t TreePrinter) VisitModifier(m *Modifier) TreePrinter {
	var result = "{"
	result = result + "\"Name\": \"" + m.Name + "\","
	b := new(bytes.Buffer)
	for key, value := range m.Attrs {
		fmt.Fprintf(b, "%s:\"%s\",", key, value)
	}
	var s = strings.TrimSuffix(b.String(), ",")
	result = result + "\"Attrs\": [" + s + "]"
	result = result + "}"
	t.JsonString = result
	return t
}
