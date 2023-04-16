package lookup

import (
	"fmt"
	//"strings"

	"trivil/ast"
	"trivil/env"
)

var _ = fmt.Printf

type lookContext struct {
	module    *ast.Module
	scope     *ast.Scope
	processed map[ast.Decl]bool
	decls     []ast.Decl // в правильном порядке
}

func Process(m *ast.Module) {

	var lc = &lookContext{
		module:    m,
		scope:     m.Inner,
		processed: make(map[ast.Decl]bool),
		decls:     make([]ast.Decl, 0),
	}

	// добавление импорта
	for _, i := range m.Imports {
		addToScope(i.Mod.Name, i.Mod, m.Inner)
	}

	// добавление имен
	for _, d := range m.Decls {
		switch x := d.(type) {
		case *ast.TypeDecl:
			addToScope(x.Name, x, m.Inner)
		case *ast.VarDecl:
			addToScope(x.Name, x, m.Inner)
		case *ast.ConstDecl:
			addToScope(x.Name, x, m.Inner)
		case *ast.Function:
			if x.Recv == nil {
				addToScope(x.Name, x, m.Inner)
			}
		default:
			panic(fmt.Sprintf("lookup 1: ni %T", d))
		}
	}

	if lc.scope != m.Inner {
		panic("assert - should be module scope")
	}

	/*
		// TODO обойти типы
		for _, d := range m.Decls {
			td, ok := d.(*ast.TypeDecl)
			if ok {
				lc.lookTypeDecl(td)
			}
		}
	*/

	// обойти описания, кроме функций
	for _, d := range m.Decls {
		lc.lookDecl(d)
	}

	// обойти функции
	for _, d := range m.Decls {
		f, ok := d.(*ast.Function)
		if ok {
			lc.lookFunction(f)
			lc.decls = append(lc.decls, d)
		}
	}

	if m.Entry != nil {
		lc.lookEntry(m.Entry)
	}
	//show(m.Decls)
	//show(lc.decls)

	// Меняем порядок описаний - определение до использования
	m.Decls = lc.decls
}

/* Отладочное, временно оставляю
func show(decls []ast.Decl) {
	var s = make([]string, len(decls))
	for i, d := range decls {
		s[i] = d.GetName()
	}
	fmt.Printf("%v\n", strings.Join(s, ","))
}
*/

// Обрабатывает описания, кроме функций
// Проверяет рекурсивные описания, задает порядок описаний
func (lc *lookContext) lookDecl(d ast.Decl) {

	_, ok := d.(*ast.Function)
	if ok {
		return
	}

	completed, exist := lc.processed[d]
	//fmt.Printf("! %v %v %v\n", d.GetName(), completed, exist)

	if exist {
		if !completed {
			env.AddError(d.GetPos(), "СЕМ-РЕКУРСИВНОЕ-ОПРЕДЕЛЕНИЕ", d.GetName())
		}
		return
	}

	lc.processed[d] = false

	switch x := d.(type) {
	case *ast.TypeDecl:
		lc.lookTypeDecl(x)
	case *ast.ConstDecl:
		lc.lookConstDecl(x)
	case *ast.VarDecl:
		lc.lookVarDecl(x)
	case *ast.Function:
		return
	default:
		panic(fmt.Sprintf("lookup 3: ni %T", d))
	}

	lc.processed[d] = true
	lc.decls = append(lc.decls, d)
}

//==== константы и переменные

func (lc *lookContext) lookVarDecl(v *ast.VarDecl) {
	if v.Typ != nil {
		lc.lookTypeRef(v.Typ)
	}
	if !v.Later {
		lc.lookExpr(v.Init)
	}
}

func (lc *lookContext) lookConstDecl(v *ast.ConstDecl) {
	if v.Typ != nil {
		lc.lookTypeRef(v.Typ)
	}
	lc.lookExpr(v.Value)
}

//==== functions

func (lc *lookContext) lookFunction(f *ast.Function) {

	f.Inner = ast.NewScope(lc.scope)
	lc.scope = f.Inner

	if f.Recv != nil {
		lc.lookTypeRef(f.Recv.Typ)

		lc.addMethodToType(f)

		lc.addVarForParameter(f.Recv)
	}

	var ft = f.Typ.(*ast.FuncType)

	for _, p := range ft.Params {
		lc.lookTypeRef(p.Typ)
		if !f.External {
			lc.addVarForParameter(p)
		}
	}

	if ft.ReturnTyp != nil {
		lc.lookTypeRef(ft.ReturnTyp)
	}

	if !f.External {
		lc.lookStatements(f.Seq)
	}

	lc.scope = lc.scope.Outer
}

func (lc *lookContext) addMethodToType(f *ast.Function) {

	var rt = f.Recv.Typ.(*ast.TypeRef)

	cl, ok := rt.Typ.(*ast.ClassType)
	if !ok {
		env.AddError(f.Recv.Pos, "СЕМ-ПОЛУЧАТЕЛЬ-НЕ-КЛАСС")
		return
	}

	cl.Methods = append(cl.Methods, f)

}

func (lc *lookContext) addVarForParameter(p *ast.Param) {
	var v = &ast.VarDecl{}
	v.Typ = p.Typ
	v.Name = p.Name
	addToScope(v.Name, v, lc.scope)
}

func (lc *lookContext) lookEntry(e *ast.EntryFn) {
	lc.lookStatements(e.Seq)
}

//==== statements

func (lc *lookContext) lookStatements(seq *ast.StatementSeq) {

	for _, s := range seq.Statements {
		lc.lookStatement(seq, s)
	}

	if lc.scope == seq.Inner {
		lc.scope = seq.Inner.Outer
	}
}

func (lc *lookContext) lookStatement(seq *ast.StatementSeq, s ast.Statement) {
	switch x := s.(type) {
	case *ast.StatementSeq:
		lc.lookStatements(x)
	case *ast.ExprStatement:
		lc.lookExpr(x.X)
	case *ast.DeclStatement:
		lc.lookLocalDecl(seq, x.D)
	case *ast.AssignStatement:
		lc.lookExpr(x.L)
		lc.lookExpr(x.R)
	case *ast.IncStatement:
		lc.lookExpr(x.L)
	case *ast.DecStatement:
		lc.lookExpr(x.L)
	case *ast.If:
		lc.lookExpr(x.Cond)
		lc.lookStatements(x.Then)
		if x.Else != nil {
			lc.lookStatement(nil, x.Else)
		}
	case *ast.While:
		lc.lookExpr(x.Cond)
		lc.lookStatements(x.Seq)
	case *ast.Guard:
		lc.lookExpr(x.Cond)
		lc.lookStatement(nil, x.Else)
	case *ast.When:
		lc.lookWhen(x)
	case *ast.Return:
		if x.X != nil {
			lc.lookExpr(x.X)
		}
	case *ast.Break:
		//nothing
	case *ast.Crash:
		lc.lookExpr(x.X)

	default:
		panic(fmt.Sprintf("statement: ni %T", s))

	}
}

func (lc *lookContext) lookLocalDecl(seq *ast.StatementSeq, decl ast.Decl) {
	if lc.scope != seq.Inner {
		seq.Inner = ast.NewScope(lc.scope)
		lc.scope = seq.Inner
	}
	switch x := decl.(type) {
	case *ast.VarDecl:
		lc.lookVarDecl(x)
		addToScope(x.Name, x, lc.scope)
	default:
		panic(fmt.Sprintf("local decl: ni %T", decl))
	}
	//ast.ShowScopes("", lc.scope)
}

func (lc *lookContext) lookWhen(x *ast.When) {
	lc.lookExpr(x.X)

	for _, c := range x.Cases {
		for _, e := range c.Exprs {
			lc.lookExpr(e)
		}
		lc.lookStatements(c.Seq)
	}
	if x.Else != nil {
		lc.lookStatements(x.Else)
	}
}
