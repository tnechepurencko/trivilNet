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

	from, _ /*is_predefined*/ := ast.UnderType(x.X.GetType()).(*ast.PredefinedType)

	switch to {
	case ast.Byte:
		return genc.convertPredefined(expr, from, ast.Byte)
	case ast.Int64:
		if from == ast.Byte || from == ast.Symbol {
			return genc.castPredefined(expr, ast.Int64)
		} else {
			return genc.convertPredefined(expr, from, ast.Int64)
		}
	case ast.Float64:
		return genc.castPredefined(expr, ast.Float64)
	case ast.Symbol:
		return genc.convertPredefined(expr, from, ast.Symbol)
	case ast.String:
		return genc.convertToString(expr, ast.UnderType(x.X.GetType()))
	}

	switch /*xt :=*/ to.(type) {
	/*
		case *ast.VectorType:
			cc.conversionToVector(x, xt)
		case *ast.ClassType:
			cc.conversionToClass(x, xt)
	*/
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
