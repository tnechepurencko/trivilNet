package check

import (
	"fmt"
	"trivil/ast"
	"trivil/env"
)

var _ = fmt.Printf

type checkContext struct {
	checkedTypes map[string]struct{}
}

func Process(m *ast.Module) {
	var cc = &checkContext{
		checkedTypes: make(map[string]struct{}),
	}

	for _, d := range m.Decls {
		switch x := d.(type) {
		case *ast.TypeDecl:
			cc.typeDecl(x)
		case *ast.VarDecl:
		case *ast.ConstDecl:
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
/*
func (cc *checkContext) lookVarDecl(v *ast.VarDecl) {
	cc.TypeRef(v.Typ)
}

func (cc *checkContext) lookConstDecl(v *ast.ConstDecl) {
	cc.TypeRef(v.Typ)

}
*/

//==== functions

func (cc *checkContext) function(f *ast.Function) {

	if f.Seq != nil {
		cc.statements(f.Seq)
	}
}

/*
func (cc *checkContext) addMethodToType(f *ast.Function) {

	var rt = f.Recv.Typ.(*ast.TypeRef)

	cl, ok := rt.Typ.(*ast.ClassType)
	if !ok {
		env.AddError(f.Recv.Pos, "СЕМ-ПОЛУЧАТЕЛЬ-КЛАСС")
		return
	}

	cl.Methods = append(cl.Methods, f)

}

func (cc *checkContext) addVarForParameter(p *ast.Param) {
	var v = &ast.VarDecl{
		Typ: p.Typ,
	}
	v.Name = p.Name
	addToScope(v.Name, v, lc.scope)
}
*/

func (cc *checkContext) entry(e *ast.EntryFn) {
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
		//TODO: проверить, что есть вызов, иначе ошибка
	case *ast.DeclStatement:
		//cc.localDecl(seq, x.D)
	case *ast.AssignStatement:
		cc.expr(x.L)
		cc.expr(x.R)
	case *ast.IncStatement:
		cc.expr(x.L)
	case *ast.DecStatement:
		cc.expr(x.L)
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
		cc.statements(x.Seq)
	case *ast.Return:
		if x.X != nil {
			cc.expr(x.X)
		}

	default:
		panic(fmt.Sprintf("statement: ni %T", s))
	}
}

/*
func (cc *checkContext) localDecl(seq *ast.StatementSeq, decl ast.Decl) {

	if lc.scope != seq.Inner {
		seq.Inner = ast.NewScope(lc.scope)
		lc.scope = seq.Inner
	}
	switch x := decl.(type) {
	case *ast.VarDecl:
		addToScope(x.Name, x, lc.scope)
		cc.VarDecl(x)
	default:
		panic(fmt.Sprintf("local decl: ni %T", decl))
	}
}
*/

//====
