package check

import (
	"fmt"
	"trivil/ast"
	"trivil/env"
	"trivil/lexer"
)

var _ = fmt.Printf

type checkContext struct {
	module       *ast.Module
	checkedTypes map[string]struct{}
	returnTyp    ast.Type
	errorHint    string
	loopCount    int
}

func Process(m *ast.Module) {
	var cc = &checkContext{
		module:       m,
		checkedTypes: make(map[string]struct{}),
	}

	for _, d := range m.Decls {
		switch x := d.(type) {
		case *ast.TypeDecl:
			cc.typeDecl(x)
		case *ast.VarDecl:
			cc.varDecl(x)
		case *ast.ConstDecl:
			cc.constDecl(x)
		case *ast.Function:
			cc.function(x)
		default:
			panic(fmt.Sprintf("check: ni %T", d))
		}
	}

	if m.Entry != nil {
		cc.entry(m.Entry)
	}
}

//==== константы и переменные

func (cc *checkContext) varDecl(v *ast.VarDecl) {
	cc.expr(v.Init)

	if v.Typ != nil {
		cc.checkAssignable(v.Typ, v.Init)
	} else {
		v.Typ = v.Init.GetType()
		if v.Typ == nil {
			panic("assert - не задан тип переменной")
		}
	}
}

func (cc *checkContext) constDecl(v *ast.ConstDecl) {
	cc.expr(v.Value)

	if v.Typ != nil {
		cc.checkAssignable(v.Typ, v.Value)
	} else {
		v.Typ = v.Value.GetType()
		if v.Typ == nil {
			panic("assert - не задан тип константы")
		}
	}

}

//==== functions

func (cc *checkContext) function(f *ast.Function) {

	if f.Seq != nil {
		cc.returnTyp = f.Typ.(*ast.FuncType).ReturnTyp

		cc.loopCount = 0
		cc.statements(f.Seq)

		cc.returnTyp = nil
	}
}

func (cc *checkContext) entry(e *ast.EntryFn) {
	cc.loopCount = 0
	cc.statements(e.Seq)
}

//==== statements

func (cc *checkContext) statements(seq *ast.StatementSeq) {

	for _, s := range seq.Statements {
		cc.statement(s)
	}
}

func (cc *checkContext) statement(s ast.Statement) {
	switch x := s.(type) {
	case *ast.StatementSeq:
		cc.statements(x) // из else
	case *ast.ExprStatement:
		cc.expr(x.X)
		if _, isCall := x.X.(*ast.CallExpr); !isCall {
			env.AddError(x.Pos, "СЕМ-ЗНАЧЕНИЕ-НЕ-ИСПОЛЬЗУЕТСЯ")
		}
	case *ast.DeclStatement:
		cc.localDecl(x.D)
	case *ast.AssignStatement:
		cc.expr(x.L)
		cc.expr(x.R)
		cc.checkAssignable(x.L.GetType(), x.R)
		cc.checkLValue(x.L)
	case *ast.IncStatement:
		cc.expr(x.L)
		if !ast.IsIntegerType(x.L.GetType()) {
			env.AddError(x.GetPos(), "СЕМ-ОШ-УНАРНАЯ-ТИП",
				ast.TypeString(x.L.GetType()), lexer.INC.String())
		}
		cc.checkLValue(x.L)
	case *ast.DecStatement:
		cc.expr(x.L)
		if !ast.IsIntegerType(x.L.GetType()) {
			env.AddError(x.GetPos(), "СЕМ-ОШ-УНАРНАЯ-ТИП",
				ast.TypeString(x.L.GetType()), lexer.DEC.String())
		}
		cc.checkLValue(x.L)
	case *ast.If:
		cc.expr(x.Cond)
		if !ast.IsBoolType(x.Cond.GetType()) {
			env.AddError(x.Cond.GetPos(), "СЕМ-ТИП-ВЫРАЖЕНИЯ", ast.Bool.Name)
		}

		cc.statements(x.Then)
		if x.Else != nil {
			cc.statement(x.Else)
		}
	case *ast.While:
		cc.expr(x.Cond)
		if !ast.IsBoolType(x.Cond.GetType()) {
			env.AddError(x.Cond.GetPos(), "СЕМ-ТИП-ВЫРАЖЕНИЯ", ast.Bool.Name)
		}
		cc.loopCount++
		cc.statements(x.Seq)
		cc.loopCount--

	case *ast.Guard:
		cc.expr(x.Cond)
		if !ast.IsBoolType(x.Cond.GetType()) {
			env.AddError(x.Cond.GetPos(), "СЕМ-ТИП-ВЫРАЖЕНИЯ", ast.Bool.Name)
		}
		cc.statement(x.Else)
		if seq, ok := x.Else.(*ast.StatementSeq); ok {
			if !isTerminating(seq) {
				env.AddError(x.Else.GetPos(), "СЕМ-НЕ-ЗАВЕРШАЮЩИЙ")
			}
		}

	case *ast.Return:
		if x.X != nil {
			cc.expr(x.X)

			if cc.returnTyp == nil {
				env.AddError(x.Pos, "СЕМ-ОШ-ВЕРНУТЬ-ЛИШНЕЕ")
			} else {
				cc.checkAssignable(cc.returnTyp, x.X)
			}
		} else if cc.returnTyp != nil {
			env.AddError(x.Pos, "СЕМ-ОШ-ВЕРНУТЬ-НУЖНО")
		}
	case *ast.Break:
		if cc.loopCount == 0 {
			env.AddError(x.Pos, "СЕМ-ПРЕРВАТЬ-ВНЕ-ЦИКЛА")
		}
	case *ast.Crash:
		cc.expr(x.X)
		if !ast.IsStringType(x.X.GetType()) {
			env.AddError(x.X.GetPos(), "СЕМ-ТИП-ВЫРАЖЕНИЯ", ast.String.Name)
		}

	default:
		panic(fmt.Sprintf("statement: ni %T", s))
	}
}

func (cc *checkContext) localDecl(decl ast.Decl) {

	switch x := decl.(type) {
	case *ast.VarDecl:
		cc.varDecl(x)
	default:
		panic(fmt.Sprintf("local decl: ni %T", decl))
	}
}

//====

func isTerminating(seq *ast.StatementSeq) bool {
	if len(seq.Statements) == 0 {
		return false
	}
	var st = seq.Statements[len(seq.Statements)-1]
	switch st.(type) {
	case *ast.Return, *ast.Break, *ast.Crash:
		return true
	default:
		return false
	}

}
