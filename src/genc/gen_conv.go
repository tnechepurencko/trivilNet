package genc

import (
	"fmt"

	"trivil/ast"
)

var _ = fmt.Printf

func (genc *genContext) genConversion(x *ast.ConversionExpr) string {

	var expr = genc.genExpr(x.X)
	if x.Done {
		return expr
	}

	var to = ast.UnderType(x.TargetTyp)

	var from = ast.UnderType(x.X.GetType())
	fromPred, _ := from.(*ast.PredefinedType)

	switch to {
	case ast.Byte:
		return genc.convertPredefined(expr, fromPred, ast.Byte)
	case ast.Int64:
		if fromPred == ast.Byte || fromPred == ast.Symbol {
			return genc.castPredefined(expr, ast.Int64)
		} else {
			return genc.convertPredefined(expr, fromPred, ast.Int64)
		}
	case ast.Float64:
		return genc.castPredefined(expr, ast.Float64)
	case ast.Symbol:
		return genc.convertPredefined(expr, fromPred, ast.Symbol)
	case ast.String:
		return genc.convertToString(expr, ast.UnderType(x.X.GetType()))
	}

	switch xt := to.(type) {
	case *ast.VectorType:
		return genc.convertToVector(expr, from, xt)
	case *ast.ClassType:
		return genc.convertToClass(expr, x.TargetTyp)
	default:
		panic(fmt.Sprintf("ni %T '%s'", to, ast.TypeString(to)))
	}
}

func (genc *genContext) convertPredefined(expr string, from, to *ast.PredefinedType) string {
	return fmt.Sprintf("%s%s_to_%s(%s)", rt_convert, predefinedTypeName(from.Name), predefinedTypeName(to.Name), expr)
}

func (genc *genContext) castPredefined(expr string, to *ast.PredefinedType) string {
	return fmt.Sprintf("(%s)(%s)", predefinedTypeName(to.Name), expr)
}

func (genc *genContext) convertToString(expr string, from ast.Type) string {

	if from == ast.Symbol {
		return genc.convertPredefined(expr, ast.Symbol, ast.String)
	}

	vt, ok := from.(*ast.VectorType)
	if !ok {
		panic("ni")
	}

	var et = ast.UnderType(vt.ElementTyp)
	if et == ast.Byte {
		return fmt.Sprintf("%s%s_to_%s(%s)", rt_convert, "Bytes", predefinedTypeName(ast.String.Name), expr)
	} else if et == ast.Symbol {
		return fmt.Sprintf("%s%s_to_%s(%s)", rt_convert, "Symbols", predefinedTypeName(ast.String.Name), expr)
	} else {
		panic("ni")
	}

}

func (genc *genContext) convertToVector(expr string, from ast.Type, to *ast.VectorType) string {

	if from != ast.String {
		panic("ni")
	}

	var et = ast.UnderType(to.ElementTyp)
	if et == ast.Byte {
		return fmt.Sprintf("%s%s_to_%s(%s)", rt_convert, predefinedTypeName(ast.String.Name), "Bytes", expr)
	} else if et == ast.Symbol {
		return fmt.Sprintf("%s%s_to_%s(%s)", rt_convert, predefinedTypeName(ast.String.Name), "Symbols", expr)
	} else {
		panic("ni")
	}
}

func (genc *genContext) convertToClass(expr string, target ast.Type) string {
	var tname = genc.typeRef(target)

	return fmt.Sprintf("((%s)%s(%s, %s))", tname, rt_checkClassType, expr, tname+nm_class_info_ptr_suffix)
}
