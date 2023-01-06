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
		/*
			case ast.Int64:
				cc.conversionToInt64(x)
				return
			case ast.Float64:
				cc.conversionToFloat64(x)
				return
			case ast.Bool:
				env.AddError(x.Pos, "СЕМ-ОШ-ПРИВЕДЕНИЯ-ТИПА", ast.TypeString(x.X.GetType()), ast.Bool.Name)
				x.Typ = invalidType(x.Pos)
				return
			case ast.Symbol:
				cc.conversionToSymbol(x)
				return
			case ast.String:
				cc.conversionToString(x)
				return
		*/
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

//cast
